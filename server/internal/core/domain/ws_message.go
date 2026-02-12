package domain

import (
	"encoding/json"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WSMessageType string

const (
	// Client -> Server
	WSMsgPlayerAction WSMessageType = "player_action"
	WSMsgJoinTable    WSMessageType = "join_table"
	WSMsgLeaveTable   WSMessageType = "leave_table"
	WSMsgChat         WSMessageType = "chat"

	// Server -> Client
	WSMsgTableState    WSMessageType = "table_state"
	WSMsgCardsDealt    WSMessageType = "cards_dealt"
	WSMsgPlayerActed   WSMessageType = "player_acted"
	WSMsgCommunity     WSMessageType = "community_cards"
	WSMsgHandResult    WSMessageType = "hand_result"
	WSMsgPlayerJoined  WSMessageType = "player_joined"
	WSMsgPlayerLeft    WSMessageType = "player_left"
	WSMsgTurnChanged   WSMessageType = "turn_changed"
	WSMsgNewHand       WSMessageType = "new_hand"
	WSMsgError         WSMessageType = "error"
	WSMsgPotUpdated    WSMessageType = "pot_updated"
)

type WSMessage struct {
	Type    WSMessageType   `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type WSPlayerAction struct {
	TableID  uuid.UUID       `json:"table_id"`
	Action   ActionType      `json:"action"`
	Amount   decimal.Decimal `json:"amount,omitempty"`
}

type WSTableState struct {
	TableID        uuid.UUID       `json:"table_id"`
	Name           string          `json:"name"`
	SmallBlind     decimal.Decimal `json:"small_blind"`
	BigBlind       decimal.Decimal `json:"big_blind"`
	Pot            decimal.Decimal `json:"pot"`
	CommunityCards []Card          `json:"community_cards"`
	Stage          GameStage       `json:"stage"`
	DealerSeat     int             `json:"dealer_seat"`
	CurrentTurn    *uuid.UUID      `json:"current_turn,omitempty"`
	Players        []WSPlayerInfo  `json:"players"`
}

type WSPlayerInfo struct {
	UserID     uuid.UUID       `json:"user_id"`
	Username   string          `json:"username"`
	Stack      decimal.Decimal `json:"stack"`
	SeatNumber int             `json:"seat_number"`
	Status     PlayerStatus    `json:"status"`
	BetAmount  decimal.Decimal `json:"bet_amount"`
	IsDealer   bool            `json:"is_dealer"`
}

type WSCardsDealt struct {
	HoleCards []Card `json:"hole_cards"`
	HandID    uuid.UUID `json:"hand_id"`
}
