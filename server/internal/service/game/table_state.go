package game

import (
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type TableState struct {
	Table      domain.PokerTable
	Players    map[int]*domain.PokerPlayer
	Hand       *HandState
	DealerSeat int
	HandCount  int
}

type HandState struct {
	Hand           domain.PokerHand
	FSM            GameFSM
	Deck           []domain.Card
	CommunityCards []domain.Card
	PlayerHands    map[uuid.UUID][]domain.Card
	Betting        BettingState
	CumulativeBets map[uuid.UUID]decimal.Decimal
	Pots           []domain.Pot
	ServerSeed     string
	SeedHash       string
	ClientSeed     string
	Nonce          int
	ActionOrder    int
}

func NewTableState(table domain.PokerTable) TableState {
	return TableState{
		Table:      table,
		Players:    make(map[int]*domain.PokerPlayer),
		DealerSeat: 0,
		HandCount:  0,
	}
}

func (ts TableState) ActivePlayerCount() int {
	count := 0
	for _, p := range ts.Players {
		if p.Status == domain.PlayerStatusActive || p.Status == domain.PlayerStatusAllIn {
			count++
		}
	}
	return count
}

func (ts TableState) PlayerList() []domain.PokerPlayer {
	players := make([]domain.PokerPlayer, 0, len(ts.Players))
	for _, p := range ts.Players {
		players = append(players, *p)
	}
	return players
}

func (ts TableState) FindPlayerBySeat(seat int) *domain.PokerPlayer {
	p, ok := ts.Players[seat]
	if !ok {
		return nil
	}
	return p
}

func (ts TableState) FindPlayerByUserID(userID uuid.UUID) *domain.PokerPlayer {
	for _, p := range ts.Players {
		if p.UserID == userID {
			return p
		}
	}
	return nil
}

func (ts TableState) OccupiedSeats() []int {
	seats := make([]int, 0, len(ts.Players))
	for seat := range ts.Players {
		seats = append(seats, seat)
	}
	return seats
}

func (ts TableState) ToWSTableState() domain.WSTableState {
	var communityCards []domain.Card
	var pot decimal.Decimal
	var stage domain.GameStage = domain.StageWaiting

	if ts.Hand != nil {
		communityCards = ts.Hand.CommunityCards
		pot = ts.Hand.Hand.Pot
		stage = ts.Hand.FSM.Stage()
	}

	players := make([]domain.WSPlayerInfo, 0, len(ts.Players))
	for _, p := range ts.Players {
		betAmount := decimal.Zero
		if ts.Hand != nil {
			for _, bp := range ts.Hand.Betting.Players {
				if bp.PlayerID == p.ID {
					betAmount = bp.BetThisRound
					break
				}
			}
		}

		players = append(players, domain.WSPlayerInfo{
			UserID:     p.UserID,
			Username:   p.Username,
			Stack:      p.Stack,
			SeatNumber: p.SeatNumber,
			Status:     p.Status,
			BetAmount:  betAmount,
			IsDealer:   p.SeatNumber == ts.DealerSeat,
		})
	}

	return domain.WSTableState{
		TableID:        ts.Table.ID,
		Name:           ts.Table.Name,
		SmallBlind:     ts.Table.SmallBlind,
		BigBlind:       ts.Table.BigBlind,
		Pot:            pot,
		CommunityCards: communityCards,
		Stage:          stage,
		DealerSeat:     ts.DealerSeat,
		Players:        players,
	}
}
