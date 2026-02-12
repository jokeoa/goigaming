package game

import (
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type BettingPlayer struct {
	PlayerID     uuid.UUID
	Stack        decimal.Decimal
	BetThisRound decimal.Decimal
	HasActed     bool
	IsAllIn      bool
	IsFolded     bool
}

type BettingState struct {
	Players    []BettingPlayer
	CurrentBet decimal.Decimal
	MinRaise   decimal.Decimal
	PotSize    decimal.Decimal
	CurrentIdx int
	BigBlind   decimal.Decimal
}

func NewBettingState(players []BettingPlayer, bigBlind decimal.Decimal, potSize decimal.Decimal) BettingState {
	return BettingState{
		Players:    players,
		CurrentBet: decimal.Zero,
		MinRaise:   bigBlind,
		PotSize:    potSize,
		CurrentIdx: 0,
		BigBlind:   bigBlind,
	}
}

func ValidateAction(state BettingState, playerID uuid.UUID, action domain.ActionType, amount decimal.Decimal) (BettingState, error) {
	idx := -1
	for i, p := range state.Players {
		if p.PlayerID == playerID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return state, domain.ErrPlayerNotFound
	}

	if state.CurrentIdx != idx {
		return state, domain.ErrNotPlayerTurn
	}

	player := state.Players[idx]
	if player.IsFolded || player.IsAllIn {
		return state, domain.ErrInvalidAction
	}

	switch action {
	case domain.ActionFold:
		return applyFold(state, idx), nil
	case domain.ActionCheck:
		return applyCheck(state, idx)
	case domain.ActionCall:
		return applyCall(state, idx)
	case domain.ActionBet:
		return applyBet(state, idx, amount)
	case domain.ActionRaise:
		return applyRaise(state, idx, amount)
	case domain.ActionAllIn:
		return applyAllIn(state, idx)
	default:
		return state, domain.ErrInvalidAction
	}
}

func applyFold(state BettingState, idx int) BettingState {
	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     newPlayers[idx].PlayerID,
		Stack:        newPlayers[idx].Stack,
		BetThisRound: newPlayers[idx].BetThisRound,
		HasActed:     true,
		IsAllIn:      newPlayers[idx].IsAllIn,
		IsFolded:     true,
	}

	return BettingState{
		Players:    newPlayers,
		CurrentBet: state.CurrentBet,
		MinRaise:   state.MinRaise,
		PotSize:    state.PotSize,
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}
}

func applyCheck(state BettingState, idx int) (BettingState, error) {
	if state.CurrentBet.GreaterThan(state.Players[idx].BetThisRound) {
		return state, domain.ErrInvalidAction
	}

	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     newPlayers[idx].PlayerID,
		Stack:        newPlayers[idx].Stack,
		BetThisRound: newPlayers[idx].BetThisRound,
		HasActed:     true,
		IsAllIn:      newPlayers[idx].IsAllIn,
		IsFolded:     newPlayers[idx].IsFolded,
	}

	return BettingState{
		Players:    newPlayers,
		CurrentBet: state.CurrentBet,
		MinRaise:   state.MinRaise,
		PotSize:    state.PotSize,
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}, nil
}

func applyCall(state BettingState, idx int) (BettingState, error) {
	player := state.Players[idx]
	callAmount := state.CurrentBet.Sub(player.BetThisRound)

	if callAmount.LessThanOrEqual(decimal.Zero) {
		return state, domain.ErrInvalidAction
	}

	isAllIn := false
	if callAmount.GreaterThanOrEqual(player.Stack) {
		callAmount = player.Stack
		isAllIn = true
	}

	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     player.PlayerID,
		Stack:        player.Stack.Sub(callAmount),
		BetThisRound: player.BetThisRound.Add(callAmount),
		HasActed:     true,
		IsAllIn:      isAllIn,
		IsFolded:     false,
	}

	return BettingState{
		Players:    newPlayers,
		CurrentBet: state.CurrentBet,
		MinRaise:   state.MinRaise,
		PotSize:    state.PotSize.Add(callAmount),
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}, nil
}

func applyBet(state BettingState, idx int, amount decimal.Decimal) (BettingState, error) {
	if state.CurrentBet.GreaterThan(decimal.Zero) {
		return state, domain.ErrInvalidAction
	}

	if amount.LessThan(state.BigBlind) {
		return state, domain.ErrInvalidBetAmount
	}

	player := state.Players[idx]
	if amount.GreaterThan(player.Stack) {
		return state, domain.ErrInsufficientStack
	}

	isAllIn := amount.Equal(player.Stack)

	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     player.PlayerID,
		Stack:        player.Stack.Sub(amount),
		BetThisRound: player.BetThisRound.Add(amount),
		HasActed:     true,
		IsAllIn:      isAllIn,
		IsFolded:     false,
	}

	resetHasActed(newPlayers, idx)

	return BettingState{
		Players:    newPlayers,
		CurrentBet: amount,
		MinRaise:   amount,
		PotSize:    state.PotSize.Add(amount),
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}, nil
}

func applyRaise(state BettingState, idx int, totalAmount decimal.Decimal) (BettingState, error) {
	if state.CurrentBet.IsZero() {
		return state, domain.ErrInvalidAction
	}

	raiseBy := totalAmount.Sub(state.CurrentBet)
	if raiseBy.LessThan(state.MinRaise) {
		return state, domain.ErrInvalidBetAmount
	}

	player := state.Players[idx]
	toCall := totalAmount.Sub(player.BetThisRound)
	if toCall.GreaterThan(player.Stack) {
		return state, domain.ErrInsufficientStack
	}

	isAllIn := toCall.Equal(player.Stack)

	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     player.PlayerID,
		Stack:        player.Stack.Sub(toCall),
		BetThisRound: player.BetThisRound.Add(toCall),
		HasActed:     true,
		IsAllIn:      isAllIn,
		IsFolded:     false,
	}

	resetHasActed(newPlayers, idx)

	return BettingState{
		Players:    newPlayers,
		CurrentBet: totalAmount,
		MinRaise:   raiseBy,
		PotSize:    state.PotSize.Add(toCall),
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}, nil
}

func applyAllIn(state BettingState, idx int) (BettingState, error) {
	player := state.Players[idx]
	allInAmount := player.Stack
	newTotalBet := player.BetThisRound.Add(allInAmount)

	newPlayers := copyPlayers(state.Players)
	newPlayers[idx] = BettingPlayer{
		PlayerID:     player.PlayerID,
		Stack:        decimal.Zero,
		BetThisRound: newTotalBet,
		HasActed:     true,
		IsAllIn:      true,
		IsFolded:     false,
	}

	currentBet := state.CurrentBet
	minRaise := state.MinRaise
	if newTotalBet.GreaterThan(state.CurrentBet) {
		raiseBy := newTotalBet.Sub(state.CurrentBet)
		if raiseBy.GreaterThanOrEqual(state.MinRaise) {
			resetHasActed(newPlayers, idx)
			minRaise = raiseBy
		}
		currentBet = newTotalBet
	}

	return BettingState{
		Players:    newPlayers,
		CurrentBet: currentBet,
		MinRaise:   minRaise,
		PotSize:    state.PotSize.Add(allInAmount),
		CurrentIdx: NextPlayer(newPlayers, idx),
		BigBlind:   state.BigBlind,
	}, nil
}

func IsBettingComplete(state BettingState) bool {
	active := 0
	for _, p := range state.Players {
		if !p.IsFolded && !p.IsAllIn {
			active++
		}
	}

	if active <= 1 {
		return true
	}

	for _, p := range state.Players {
		if p.IsFolded || p.IsAllIn {
			continue
		}
		if !p.HasActed {
			return false
		}
		if p.BetThisRound.LessThan(state.CurrentBet) {
			return false
		}
	}

	return true
}

func ActivePlayerCount(state BettingState) int {
	count := 0
	for _, p := range state.Players {
		if !p.IsFolded {
			count++
		}
	}
	return count
}

func NextPlayer(players []BettingPlayer, currentIdx int) int {
	n := len(players)
	for i := 1; i < n; i++ {
		next := (currentIdx + i) % n
		if !players[next].IsFolded && !players[next].IsAllIn {
			return next
		}
	}
	return currentIdx
}

func copyPlayers(players []BettingPlayer) []BettingPlayer {
	cp := make([]BettingPlayer, len(players))
	copy(cp, players)
	return cp
}

// resetHasActed mutates the provided slice in place.
// Callers must ensure the slice is a fresh copy (via copyPlayers).
func resetHasActed(players []BettingPlayer, exceptIdx int) {
	for i := range players {
		if i != exceptIdx && !players[i].IsFolded && !players[i].IsAllIn {
			players[i].HasActed = false
		}
	}
}
