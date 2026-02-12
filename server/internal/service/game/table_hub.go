package game

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/jokeoa/goigaming/internal/core/ports"
	"github.com/shopspring/decimal"
)

var ErrHubClosed = errors.New("hub is closed")

type TableHub struct {
	state       TableState
	eventCh     chan HubEvent
	stateCh     chan chan TableState
	turnTimer   *time.Timer
	turnTimeout time.Duration
	broadcaster ports.Broadcaster
	walletSvc   ports.WalletService
	rngSvc      ports.RNGService
	handRepo    ports.PokerHandRepository
	playerRepo  ports.PokerPlayerRepository
	logger      *slog.Logger
	done        chan struct{}
}

func NewTableHub(
	table domain.PokerTable,
	turnTimeout time.Duration,
	broadcaster ports.Broadcaster,
	walletSvc ports.WalletService,
	rngSvc ports.RNGService,
	handRepo ports.PokerHandRepository,
	playerRepo ports.PokerPlayerRepository,
	logger *slog.Logger,
) *TableHub {
	return &TableHub{
		state:       NewTableState(table),
		eventCh:     make(chan HubEvent, 64),
		stateCh:     make(chan chan TableState, 8),
		turnTimeout: turnTimeout,
		broadcaster: broadcaster,
		walletSvc:   walletSvc,
		rngSvc:      rngSvc,
		handRepo:    handRepo,
		playerRepo:  playerRepo,
		logger:      logger.With("table_id", table.ID),
		done:        make(chan struct{}),
	}
}

func (h *TableHub) Send(event HubEvent) error {
	select {
	case h.eventCh <- event:
		return nil
	case <-h.done:
		if event.ResultCh != nil {
			event.ResultCh <- HubResult{Err: ErrHubClosed}
		}
		return ErrHubClosed
	}
}

func (h *TableHub) Done() <-chan struct{} {
	return h.done
}

func (h *TableHub) State() TableState {
	replyCh := make(chan TableState, 1)
	select {
	case h.stateCh <- replyCh:
		return <-replyCh
	case <-h.done:
		return h.state
	}
}

func (h *TableHub) Run(ctx context.Context) {
	defer close(h.done)
	h.logger.Info("table hub started")

	timerCh := make(<-chan time.Time)

	for {
		if h.turnTimer != nil {
			timerCh = h.turnTimer.C
		}

		select {
		case <-ctx.Done():
			h.logger.Info("table hub stopping: context cancelled")
			h.stopTimer()
			return

		case event := <-h.eventCh:
			if event.Type == EventShutdown {
				h.stopTimer()
				if event.ResultCh != nil {
					event.ResultCh <- HubResult{}
				}
				return
			}
			h.handleEvent(ctx, event)

		case replyCh := <-h.stateCh:
			replyCh <- h.copyState()

		case <-timerCh:
			h.handleTimeout(ctx)
			timerCh = make(<-chan time.Time)
		}
	}
}

func (h *TableHub) copyState() TableState {
	cp := TableState{
		Table:      h.state.Table,
		Players:    make(map[int]*domain.PokerPlayer, len(h.state.Players)),
		Hand:       h.state.Hand,
		DealerSeat: h.state.DealerSeat,
		HandCount:  h.state.HandCount,
	}
	for seat, p := range h.state.Players {
		playerCopy := *p
		cp.Players[seat] = &playerCopy
	}
	return cp
}

func (h *TableHub) handleEvent(ctx context.Context, event HubEvent) {
	var result HubResult

	switch event.Type {
	case EventPlayerAction:
		result.Err = h.handlePlayerAction(ctx, event.UserID, event.Action, event.Amount)
	case EventPlayerJoin:
		result.Err = h.handlePlayerJoin(ctx, event)
	case EventPlayerLeave:
		stack, err := h.handlePlayerLeave(ctx, event.UserID)
		result.Err = err
		result.Stack = stack
	case EventStartHand:
		result.Err = h.tryStartHand(ctx)
	}

	if event.ResultCh != nil {
		event.ResultCh <- result
	}
}

func (h *TableHub) handlePlayerJoin(ctx context.Context, event HubEvent) error {
	if event.PlayerID == uuid.Nil {
		return fmt.Errorf("missing player id for join event")
	}
	if _, exists := h.state.Players[event.SeatNum]; exists {
		return domain.ErrSeatTaken
	}

	if len(h.state.Players) >= h.state.Table.MaxPlayers {
		return domain.ErrTableFull
	}

	player := &domain.PokerPlayer{
		ID:         event.PlayerID,
		TableID:    h.state.Table.ID,
		UserID:     event.UserID,
		Username:   event.Username,
		Stack:      event.BuyIn,
		SeatNumber: event.SeatNum,
		Status:     domain.PlayerStatusActive,
		JoinedAt:   time.Now(),
	}

	h.state.Players[event.SeatNum] = player

	h.broadcastTableState()

	if h.state.Hand == nil && len(h.state.Players) >= 2 {
		return h.tryStartHand(ctx)
	}

	return nil
}

func (h *TableHub) handlePlayerLeave(ctx context.Context, userID uuid.UUID) (*decimal.Decimal, error) {
	player := h.state.FindPlayerByUserID(userID)
	if player == nil {
		return nil, domain.ErrPlayerNotFound
	}

	if h.state.Hand != nil {
		for i, bp := range h.state.Hand.Betting.Players {
			if bp.PlayerID == player.ID && !bp.IsFolded {
				newState := applyFold(h.state.Hand.Betting, i)
				h.state.Hand.Betting = newState
				break
			}
		}
	}

	stack := player.Stack
	delete(h.state.Players, player.SeatNumber)

	h.broadcastTableState()

	if h.state.Hand != nil && ActivePlayerCount(h.state.Hand.Betting) <= 1 {
		h.completeHand(ctx)
	}

	return &stack, nil
}

func (h *TableHub) handlePlayerAction(ctx context.Context, userID uuid.UUID, action domain.ActionType, amount decimal.Decimal) error {
	if h.state.Hand == nil {
		return domain.ErrGameNotStarted
	}

	player := h.state.FindPlayerByUserID(userID)
	if player == nil {
		return domain.ErrPlayerNotFound
	}

	newBetting, err := ValidateAction(h.state.Hand.Betting, player.ID, action, amount)
	if err != nil {
		return err
	}

	h.state.Hand.Betting = newBetting
	h.state.Hand.ActionOrder++

	for _, bp := range newBetting.Players {
		for seat, p := range h.state.Players {
			if p.ID == bp.PlayerID {
				updated := *p
				updated.Stack = bp.Stack
				if bp.IsFolded {
					updated.Status = domain.PlayerStatusFolded
				} else if bp.IsAllIn {
					updated.Status = domain.PlayerStatusAllIn
				}
				h.state.Players[seat] = &updated
			}
		}
	}

	h.state.Hand.Hand.Pot = newBetting.PotSize

	h.broadcastPlayerActed(userID, action, amount)
	h.broadcastPotUpdate()

	if ActivePlayerCount(newBetting) <= 1 {
		h.completeHand(ctx)
		return nil
	}

	if IsBettingComplete(newBetting) {
		h.advanceStage(ctx)
		return nil
	}

	h.resetTimer()
	h.broadcastTurnChanged()
	return nil
}

func (h *TableHub) handleTimeout(ctx context.Context) {
	if h.state.Hand == nil {
		return
	}

	currentIdx := h.state.Hand.Betting.CurrentIdx
	if currentIdx >= len(h.state.Hand.Betting.Players) {
		return
	}

	currentPlayer := h.state.Hand.Betting.Players[currentIdx]
	var userID uuid.UUID
	for _, p := range h.state.Players {
		if p.ID == currentPlayer.PlayerID {
			userID = p.UserID
			break
		}
	}

	if err := h.handlePlayerAction(ctx, userID, domain.ActionFold, decimal.Zero); err != nil {
		h.logger.Error("auto-fold on timeout failed", "error", err)
	}
}

func (h *TableHub) tryStartHand(ctx context.Context) error {
	if h.state.Hand != nil {
		return domain.ErrGameAlreadyStarted
	}

	if len(h.state.Players) < 2 {
		return domain.ErrMinPlayersRequired
	}

	return h.startHand(ctx)
}

func (h *TableHub) startHand(ctx context.Context) error {
	serverSeed, err := h.rngSvc.GenerateServerSeed()
	if err != nil {
		return fmt.Errorf("generate server seed: %w", err)
	}

	h.state.HandCount++
	h.state.DealerSeat = h.nextOccupiedSeat(h.state.DealerSeat)

	deck := h.rngSvc.ShuffleDeck(serverSeed, "default", h.state.HandCount)

	hand := domain.PokerHand{
		ID:         uuid.New(),
		TableID:    h.state.Table.ID,
		HandNumber: h.state.HandCount,
		Pot:        decimal.Zero,
		Stage:      domain.StagePreflop,
	}

	seats := h.sortedSeats()
	bettingPlayers := make([]BettingPlayer, 0, len(seats))
	playerHands := make(map[uuid.UUID][]domain.Card)
	cumulativeBets := make(map[uuid.UUID]decimal.Decimal)
	deckIdx := 0

	for _, seat := range seats {
		p := h.state.Players[seat]
		holeCards := []domain.Card{deck[deckIdx], deck[deckIdx+1]}
		deckIdx += 2
		playerHands[p.ID] = holeCards
		cumulativeBets[p.ID] = decimal.Zero

		bettingPlayers = append(bettingPlayers, BettingPlayer{
			PlayerID:     p.ID,
			Stack:        p.Stack,
			BetThisRound: decimal.Zero,
			HasActed:     false,
			IsAllIn:      false,
			IsFolded:     false,
		})
	}

	betting := NewBettingState(bettingPlayers, h.state.Table.BigBlind, decimal.Zero)

	sbIdx, bbIdx := h.blindPositions(seats)
	betting = h.postBlinds(betting, sbIdx, bbIdx)

	for _, bp := range betting.Players {
		cumulativeBets[bp.PlayerID] = bp.BetThisRound
	}

	firstToAct := (bbIdx + 1) % len(bettingPlayers)
	betting.CurrentIdx = firstToAct

	fsm, err := NewGameFSM().Transition(domain.StagePreflop)
	if err != nil {
		return fmt.Errorf("FSM transition to preflop: %w", err)
	}

	h.state.Hand = &HandState{
		Hand:           hand,
		FSM:            fsm,
		Deck:           deck[deckIdx:],
		CommunityCards: nil,
		PlayerHands:    playerHands,
		Betting:        betting,
		CumulativeBets: cumulativeBets,
		ServerSeed:     serverSeed,
		SeedHash:       h.rngSvc.HashSeed(serverSeed),
		ClientSeed:     "default",
		Nonce:          h.state.HandCount,
		ActionOrder:    0,
	}

	if _, err := h.handRepo.Create(ctx, hand); err != nil {
		h.logger.Error("failed to persist hand", "error", err)
	}

	for id, cards := range playerHands {
		var userID uuid.UUID
		for _, p := range h.state.Players {
			if p.ID == id {
				userID = p.UserID
				break
			}
		}
		h.sendCardsDealt(userID, hand.ID, cards)
	}

	h.broadcastNewHand()
	h.broadcastPotUpdate()
	h.resetTimer()
	h.broadcastTurnChanged()

	return nil
}

func (h *TableHub) advanceStage(ctx context.Context) {
	handState := h.state.Hand
	nextStage := handState.FSM.NextStage()

	if nextStage == domain.StageShowdown {
		h.doShowdown(ctx)
		return
	}

	newFSM, err := handState.FSM.Transition(nextStage)
	if err != nil {
		h.logger.Error("FSM transition failed", "error", err, "to", nextStage)
		return
	}

	switch nextStage {
	case domain.StageFlop:
		h.dealCommunity(3)
	case domain.StageTurn, domain.StageRiver:
		h.dealCommunity(1)
	}

	for _, p := range handState.Betting.Players {
		h.state.Hand.CumulativeBets[p.PlayerID] = h.state.Hand.CumulativeBets[p.PlayerID].Add(p.BetThisRound)
	}

	newPlayers := make([]BettingPlayer, len(handState.Betting.Players))
	for i, p := range handState.Betting.Players {
		newPlayers[i] = BettingPlayer{
			PlayerID:     p.PlayerID,
			Stack:        p.Stack,
			BetThisRound: decimal.Zero,
			HasActed:     p.IsFolded || p.IsAllIn,
			IsAllIn:      p.IsAllIn,
			IsFolded:     p.IsFolded,
		}
	}

	firstActive := 0
	for i, p := range newPlayers {
		if !p.IsFolded && !p.IsAllIn {
			firstActive = i
			break
		}
	}

	h.state.Hand.FSM = newFSM
	h.state.Hand.Hand.Stage = nextStage
	h.state.Hand.Betting = BettingState{
		Players:    newPlayers,
		CurrentBet: decimal.Zero,
		MinRaise:   h.state.Table.BigBlind,
		PotSize:    handState.Betting.PotSize,
		CurrentIdx: firstActive,
		BigBlind:   h.state.Table.BigBlind,
	}

	h.broadcastCommunityCards()

	allActed := true
	for _, p := range newPlayers {
		if !p.IsFolded && !p.IsAllIn {
			allActed = false
			break
		}
	}

	if allActed {
		h.advanceStage(ctx)
		return
	}

	h.resetTimer()
	h.broadcastTurnChanged()
}

func (h *TableHub) doShowdown(ctx context.Context) {
	handState := h.state.Hand

	for _, p := range handState.Betting.Players {
		h.state.Hand.CumulativeBets[p.PlayerID] = h.state.Hand.CumulativeBets[p.PlayerID].Add(p.BetThisRound)
	}

	for handState.FSM.Stage() != domain.StageRiver && handState.FSM.Stage() != domain.StageShowdown {
		nextStage := handState.FSM.NextStage()
		if nextStage == domain.StageShowdown {
			break
		}
		switch nextStage {
		case domain.StageFlop:
			h.dealCommunity(3)
		case domain.StageTurn, domain.StageRiver:
			h.dealCommunity(1)
		}
		newFSM, err := handState.FSM.Transition(nextStage)
		if err != nil {
			break
		}
		handState.FSM = newFSM
	}

	newFSM, _ := handState.FSM.Transition(domain.StageShowdown)
	h.state.Hand.FSM = newFSM

	contributions := make([]PotContribution, 0, len(handState.Betting.Players))
	for _, bp := range handState.Betting.Players {
		contributions = append(contributions, PotContribution{
			PlayerID: bp.PlayerID,
			TotalBet: h.state.Hand.CumulativeBets[bp.PlayerID],
			IsAllIn:  bp.IsAllIn,
			IsFolded: bp.IsFolded,
		})
	}

	pots := CalculateSidePots(contributions)
	if len(pots) == 0 {
		pots = []domain.Pot{{Amount: handState.Betting.PotSize, EligibleIDs: h.activePlayerIDs()}}
	}

	var handPlayers []HandPlayerCards
	for playerID, cards := range handState.PlayerHands {
		isFolded := false
		for _, bp := range handState.Betting.Players {
			if bp.PlayerID == playerID {
				isFolded = bp.IsFolded
				break
			}
		}
		if !isFolded {
			handPlayers = append(handPlayers, HandPlayerCards{
				PlayerID:  playerID,
				HoleCards: cards,
			})
		}
	}

	result := DetermineWinners(handPlayers, handState.CommunityCards, pots)
	result.HandID = handState.Hand.ID

	h.completePayout(ctx, result)
	h.broadcastHandResult(result)
	h.cleanupHand(ctx)
}

func (h *TableHub) completeHand(ctx context.Context) {
	if h.state.Hand == nil {
		return
	}

	winnerIDs := h.activePlayerIDs()
	if len(winnerIDs) == 1 {
		winnerID := winnerIDs[0]
		result := domain.HandResult{
			HandID: h.state.Hand.Hand.ID,
			Winners: []domain.WinnerInfo{{
				PlayerID: winnerID,
				Amount:   h.state.Hand.Betting.PotSize,
				HandRank: "last player standing",
			}},
			Pots: []domain.Pot{{
				Amount:      h.state.Hand.Betting.PotSize,
				EligibleIDs: winnerIDs,
			}},
		}

		h.completePayout(ctx, result)
		h.broadcastHandResult(result)
	}

	h.cleanupHand(ctx)
}

func (h *TableHub) completePayout(ctx context.Context, result domain.HandResult) {
	for _, winner := range result.Winners {
		for seat, p := range h.state.Players {
			if p.ID == winner.PlayerID {
				if _, err := h.walletSvc.Deposit(ctx, p.UserID, winner.Amount.String()); err != nil {
					h.logger.Error("CRITICAL: payout deposit failed",
						"user_id", p.UserID, "amount", winner.Amount,
						"hand_id", result.HandID, "error", err)
					continue
				}

				updated := *p
				updated.Stack = updated.Stack.Add(winner.Amount)
				updated.Status = domain.PlayerStatusActive
				h.state.Players[seat] = &updated
				break
			}
		}
	}
}

func (h *TableHub) cleanupHand(ctx context.Context) {
	h.stopTimer()

	if h.state.Hand != nil {
		now := time.Now()
		h.state.Hand.Hand.Stage = domain.StageComplete
		h.state.Hand.Hand.EndedAt = &now

		if _, err := h.handRepo.Update(ctx, h.state.Hand.Hand); err != nil {
			h.logger.Error("failed to update completed hand", "error", err)
		}
	}

	h.state.Hand = nil

	for seat, p := range h.state.Players {
		if p.Stack.IsZero() {
			delete(h.state.Players, seat)
			continue
		}
		updated := *p
		updated.Status = domain.PlayerStatusActive
		h.state.Players[seat] = &updated
	}

	h.broadcastTableState()

	if len(h.state.Players) >= 2 {
		time.AfterFunc(3*time.Second, func() {
			h.Send(HubEvent{Type: EventStartHand})
		})
	}
}

func (h *TableHub) dealCommunity(count int) {
	if h.state.Hand == nil || len(h.state.Hand.Deck) < count {
		return
	}

	cards := h.state.Hand.Deck[:count]
	h.state.Hand.Deck = h.state.Hand.Deck[count:]
	h.state.Hand.CommunityCards = append(h.state.Hand.CommunityCards, cards...)
	h.state.Hand.Hand.CommunityCards = domain.CardsToString(h.state.Hand.CommunityCards)
}

func (h *TableHub) postBlinds(betting BettingState, sbIdx, bbIdx int) BettingState {
	sb := h.state.Table.SmallBlind
	bb := h.state.Table.BigBlind

	sbPlayer := betting.Players[sbIdx]
	sbAmount := decimal.Min(sb, sbPlayer.Stack)
	betting.Players[sbIdx] = BettingPlayer{
		PlayerID:     sbPlayer.PlayerID,
		Stack:        sbPlayer.Stack.Sub(sbAmount),
		BetThisRound: sbAmount,
		HasActed:     false,
		IsAllIn:      sbAmount.Equal(sbPlayer.Stack),
		IsFolded:     false,
	}

	bbPlayer := betting.Players[bbIdx]
	bbAmount := decimal.Min(bb, bbPlayer.Stack)
	betting.Players[bbIdx] = BettingPlayer{
		PlayerID:     bbPlayer.PlayerID,
		Stack:        bbPlayer.Stack.Sub(bbAmount),
		BetThisRound: bbAmount,
		HasActed:     false,
		IsAllIn:      bbAmount.Equal(bbPlayer.Stack),
		IsFolded:     false,
	}

	betting.CurrentBet = bbAmount
	betting.PotSize = sbAmount.Add(bbAmount)

	for seat, p := range h.state.Players {
		if p.ID == sbPlayer.PlayerID {
			updated := *p
			updated.Stack = sbPlayer.Stack.Sub(sbAmount)
			h.state.Players[seat] = &updated
		}
		if p.ID == bbPlayer.PlayerID {
			updated := *p
			updated.Stack = bbPlayer.Stack.Sub(bbAmount)
			h.state.Players[seat] = &updated
		}
	}

	return betting
}

func (h *TableHub) blindPositions(seats []int) (sbIdx, bbIdx int) {
	n := len(seats)
	if n == 2 {
		dealerSeatIdx := 0
		for i, s := range seats {
			if s == h.state.DealerSeat {
				dealerSeatIdx = i
				break
			}
		}
		return dealerSeatIdx, (dealerSeatIdx + 1) % n
	}

	dealerSeatIdx := 0
	for i, s := range seats {
		if s == h.state.DealerSeat {
			dealerSeatIdx = i
			break
		}
	}
	return (dealerSeatIdx + 1) % n, (dealerSeatIdx + 2) % n
}

func (h *TableHub) sortedSeats() []int {
	seats := h.state.OccupiedSeats()
	sort.Ints(seats)
	return seats
}

func (h *TableHub) nextOccupiedSeat(current int) int {
	seats := h.sortedSeats()
	if len(seats) == 0 {
		return 0
	}

	for _, s := range seats {
		if s > current {
			return s
		}
	}
	return seats[0]
}

func (h *TableHub) activePlayerIDs() []uuid.UUID {
	if h.state.Hand == nil {
		return nil
	}

	var ids []uuid.UUID
	for _, bp := range h.state.Hand.Betting.Players {
		if !bp.IsFolded {
			ids = append(ids, bp.PlayerID)
		}
	}
	return ids
}

func (h *TableHub) resetTimer() {
	h.stopTimer()
	h.turnTimer = time.NewTimer(h.turnTimeout)
}

func (h *TableHub) stopTimer() {
	if h.turnTimer != nil {
		h.turnTimer.Stop()
		h.turnTimer = nil
	}
}

func (h *TableHub) broadcastTableState() {
	msg := h.buildMessage(domain.WSMsgTableState, h.state.ToWSTableState())
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastPlayerActed(userID uuid.UUID, action domain.ActionType, amount decimal.Decimal) {
	payload := map[string]any{
		"user_id": userID,
		"action":  action,
		"amount":  amount,
	}
	msg := h.buildMessage(domain.WSMsgPlayerActed, payload)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastCommunityCards() {
	if h.state.Hand == nil {
		return
	}
	payload := map[string]any{
		"cards": h.state.Hand.CommunityCards,
		"stage": h.state.Hand.FSM.Stage(),
	}
	msg := h.buildMessage(domain.WSMsgCommunity, payload)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastHandResult(result domain.HandResult) {
	msg := h.buildMessage(domain.WSMsgHandResult, result)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastNewHand() {
	if h.state.Hand == nil {
		return
	}
	payload := map[string]any{
		"hand_id":     h.state.Hand.Hand.ID,
		"hand_number": h.state.Hand.Hand.HandNumber,
		"dealer_seat": h.state.DealerSeat,
		"seed_hash":   h.state.Hand.SeedHash,
	}
	msg := h.buildMessage(domain.WSMsgNewHand, payload)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastTurnChanged() {
	if h.state.Hand == nil {
		return
	}
	idx := h.state.Hand.Betting.CurrentIdx
	if idx >= len(h.state.Hand.Betting.Players) {
		return
	}
	currentPlayer := h.state.Hand.Betting.Players[idx]

	var userID uuid.UUID
	for _, p := range h.state.Players {
		if p.ID == currentPlayer.PlayerID {
			userID = p.UserID
			break
		}
	}

	payload := map[string]any{
		"user_id": userID,
		"timeout": h.turnTimeout.Seconds(),
	}
	msg := h.buildMessage(domain.WSMsgTurnChanged, payload)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) broadcastPotUpdate() {
	if h.state.Hand == nil {
		return
	}
	payload := map[string]any{"pot": h.state.Hand.Betting.PotSize}
	msg := h.buildMessage(domain.WSMsgPotUpdated, payload)
	h.broadcaster.BroadcastToTable(h.state.Table.ID, msg)
}

func (h *TableHub) sendCardsDealt(userID, handID uuid.UUID, cards []domain.Card) {
	payload := domain.WSCardsDealt{
		HoleCards: cards,
		HandID:    handID,
	}
	msg := h.buildMessage(domain.WSMsgCardsDealt, payload)
	h.broadcaster.SendToPlayer(h.state.Table.ID, userID, msg)
}

func (h *TableHub) buildMessage(msgType domain.WSMessageType, payload any) domain.WSMessage {
	data, err := json.Marshal(payload)
	if err != nil {
		h.logger.Error("failed to marshal message payload", "type", msgType, "error", err)
		return domain.WSMessage{Type: msgType}
	}
	return domain.WSMessage{
		Type:    msgType,
		Payload: data,
	}
}
