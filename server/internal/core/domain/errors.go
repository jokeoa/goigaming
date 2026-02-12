package domain

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidToken      = errors.New("invalid token")
	ErrWalletNotFound    = errors.New("wallet not found")
	ErrInsufficientFunds = errors.New("insufficient funds")
	ErrOptimisticLock    = errors.New("optimistic lock conflict")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrForbidden         = errors.New("forbidden")

	ErrRoundNotFound      = errors.New("round not found")
	ErrBettingClosed      = errors.New("betting is closed")
	ErrInvalidBetType     = errors.New("invalid bet type")
	ErrBetAmountOutOfRange = errors.New("bet amount out of range")
	ErrTableNotActive     = errors.New("table is not active")
)
