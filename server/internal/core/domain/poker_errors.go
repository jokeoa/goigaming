package domain

import "errors"

var (
	ErrTableNotFound      = errors.New("table not found")
	ErrTableFull          = errors.New("table is full")
	ErrPlayerNotFound     = errors.New("player not found")
	ErrPlayerAlreadySeated = errors.New("player already seated at this table")
	ErrHandNotFound       = errors.New("hand not found")
	ErrNotPlayerTurn      = errors.New("not player's turn")
	ErrInvalidAction      = errors.New("invalid action")
	ErrInvalidBetAmount   = errors.New("invalid bet amount")
	ErrInsufficientStack  = errors.New("insufficient stack")
	ErrGameNotStarted     = errors.New("game not started")
	ErrGameAlreadyStarted = errors.New("game already started")
	ErrInvalidBuyIn       = errors.New("invalid buy-in amount")
	ErrSeatTaken          = errors.New("seat is taken")
	ErrMinPlayersRequired = errors.New("minimum 2 players required")
	ErrInvalidTransition  = errors.New("invalid stage transition")
)
