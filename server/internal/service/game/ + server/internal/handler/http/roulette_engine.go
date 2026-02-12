package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RouletteBetType represents different bet types in European Roulette
type RouletteBetType string

const (
	BetStraight    RouletteBetType = "straight"     // Single number (0-36)
	BetSplit       RouletteBetType = "split"        // Two adjacent numbers
	BetStreet      RouletteBetType = "street"       // Three numbers in a row
	BetCorner      RouletteBetType = "corner"       // Four numbers in a square
	BetSixLine     RouletteBetType = "sixline"      // Six numbers (two rows)
	BetRed         RouletteBetType = "red"          // All red numbers
	BetBlack       RouletteBetType = "black"        // All black numbers
	BetEven        RouletteBetType = "even"         // All even numbers
	BetOdd         RouletteBetType = "odd"          // All odd numbers
	BetLow         RouletteBetType = "low"          // 1-18
	BetHigh        RouletteBetType = "high"         // 19-36
	BetDozen1      RouletteBetType = "dozen1"       // 1-12
	BetDozen2      RouletteBetType = "dozen2"       // 13-24
	BetDozen3      RouletteBetType = "dozen3"       // 25-36
	BetColumn1     RouletteBetType = "column1"      // 1,4,7...34
	BetColumn2     RouletteBetType = "column2"      // 2,5,8...35
	BetColumn3     RouletteBetType = "column3"      // 3,6,9...36
)

// RouletteBet represents a single bet on the roulette table
type RouletteBet struct {
	ID       uuid.UUID       `json:"id"`
	UserID   uuid.UUID       `json:"user_id"`
	Type     RouletteBetType `json:"type"`
	Numbers  []int           `json:"numbers"`
	Amount   float64         `json:"amount"`
	PlacedAt time.Time       `json:"placed_at"`
}

// RouletteResult represents the outcome of a spin
type RouletteResult struct {
	SpinID      uuid.UUID       `json:"spin_id"`
	Number      int             `json:"number"`
	Color       string          `json:"color"`
	IsEven      bool            `json:"is_even"`
	IsLow       bool            `json:"is_low"`
	Dozen       int             `json:"dozen"`
	Column      int             `json:"column"`
	SpinHash    string          `json:"spin_hash"`
	ServerSeed  string          `json:"server_seed"`
	Timestamp   time.Time       `json:"timestamp"`
	Bets        []RouletteBet   `json:"bets"`
	Winnings    map[uuid.UUID]float64 `json:"winnings"`
	TotalWagered float64        `json:"total_wagered"`
	TotalPaidOut float64        `json:"total_paid_out"`
}

// RouletteTable manages a single roulette table instance
type RouletteTable struct {
	ID            uuid.UUID
	Status        string // betting, spinning, completed
	CurrentBets   []RouletteBet
	SpinHistory   []*RouletteResult
	BettingTimer  *time.Timer
	mu            sync.RWMutex
}

// RouletteEngine handles roulette game logic
type RouletteEngine struct {
	tables map[uuid.UUID]*RouletteTable
	mu     sync.RWMutex
	rng    *rand.Rand
}

// NewRouletteEngine creates a new roulette engine
func NewRouletteEngine() *RouletteEngine {
	return &RouletteEngine{
		tables: make(map[uuid.UUID]*RouletteTable),
		rng:    rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// CreateTable creates a new roulette table
func (e *RouletteEngine) CreateTable(ctx context.Context) (*RouletteTable, error) {
	table := &RouletteTable{
		ID:          uuid.New(),
		Status:      "betting",
		CurrentBets: []RouletteBet{},
		SpinHistory: []*RouletteResult{},
	}

	e.mu.Lock()
	e.tables[table.ID] = table
	e.mu.Unlock()

	return table, nil
}

// PlaceBet places a bet on a table
func (e *RouletteEngine) PlaceBet(ctx context.Context, tableID, userID uuid.UUID, betType RouletteBetType, numbers []int, amount float64) (*RouletteBet, error) {
	e.mu.RLock()
	table, exists := e.tables[tableID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("table not found")
	}

	table.mu.Lock()
	defer table.mu.Unlock()

	if table.Status != "betting" {
		return nil, fmt.Errorf("betting is closed")
	}

	// Validate bet
	if err := e.validateBet(betType, numbers, amount); err != nil {
		return nil, err
	}

	bet := RouletteBet{
		ID:       uuid.New(),
		UserID:   userID,
		Type:     betType,
		Numbers:  numbers,
		Amount:   amount,
		PlacedAt: time.Now(),
	}

	table.CurrentBets = append(table.CurrentBets, bet)

	return &bet, nil
}

// Spin executes a roulette spin and determines winners
func (e *RouletteEngine) Spin(ctx context.Context, tableID uuid.UUID) (*RouletteResult, error) {
	e.mu.RLock()
	table, exists := e.tables[tableID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("table not found")
	}

	table.mu.Lock()
	defer table.mu.Unlock()

	if table.Status != "betting" {
		return nil, fmt.Errorf("table not in betting state")
	}

	if len(table.CurrentBets) == 0 {
		return nil, fmt.Errorf("no bets placed")
	}

	// Change status to spinning
	table.Status = "spinning"

	// Generate provably fair result
	serverSeed := e.generateSeed()
	number := e.generateNumber(serverSeed)
	spinHash := e.generateHash(serverSeed, number)

	result := &RouletteResult{
		SpinID:     uuid.New(),
		Number:     number,
		Color:      e.getNumberColor(number),
		IsEven:     number > 0 && number%2 == 0,
		IsLow:      number >= 1 && number <= 18,
		Dozen:      e.getDozen(number),
		Column:     e.getColumn(number),
		SpinHash:   spinHash,
		ServerSeed: serverSeed,
		Timestamp:  time.Now(),
		Bets:       table.CurrentBets,
		Winnings:   make(map[uuid.UUID]float64),
	}

	// Calculate winnings for each bet
	totalWagered := 0.0
	totalPaidOut := 0.0

	for _, bet := range table.CurrentBets {
		totalWagered += bet.Amount
		
		if e.isBetWinner(bet, number) {
			payout := e.calculatePayout(bet)
			result.Winnings[bet.UserID] += payout
			totalPaidOut += payout
		}
	}

	result.TotalWagered = totalWagered
	result.TotalPaidOut = totalPaidOut

	// Add to history
	table.SpinHistory = append(table.SpinHistory, result)
	
	// Clear current bets and reset status
	table.CurrentBets = []RouletteBet{}
	table.Status = "betting"

	return result, nil
}

// GetHistory returns spin history for a table
func (e *RouletteEngine) GetHistory(ctx context.Context, tableID uuid.UUID, limit int) ([]*RouletteResult, error) {
	e.mu.RLock()
	table, exists := e.tables[tableID]
	e.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("table not found")
	}

	table.mu.RLock()
	defer table.mu.RUnlock()

	history := table.SpinHistory
	if len(history) > limit {
		history = history[len(history)-limit:]
	}

	return history, nil
}

// Helper methods

func (e *RouletteEngine) validateBet(betType RouletteBetType, numbers []int, amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("bet amount must be positive")
	}

	switch betType {
	case BetStraight:
		if len(numbers) != 1 || numbers[0] < 0 || numbers[0] > 36 {
			return fmt.Errorf("straight bet requires exactly one number (0-36)")
		}
	case BetSplit:
		if len(numbers) != 2 {
			return fmt.Errorf("split bet requires exactly two numbers")
		}
	case BetStreet:
		if len(numbers) != 3 {
			return fmt.Errorf("street bet requires exactly three numbers")
		}
	case BetCorner:
		if len(numbers) != 4 {
			return fmt.Errorf("corner bet requires exactly four numbers")
		}
	}

	return nil
}

func (e *RouletteEngine) generateSeed() string {
	// Generate cryptographically secure random seed
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func (e *RouletteEngine) generateNumber(seed string) int {
	// Use seed to generate deterministic number (0-36)
	hash := sha256.Sum256([]byte(seed + time.Now().String()))
	// Convert hash to number 0-36
	num := int(hash[0]) % 37
	return num
}

func (e *RouletteEngine) generateHash(seed string, number int) string {
	data := fmt.Sprintf("%s:%d", seed, number)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func (e *RouletteEngine) getNumberColor(number int) string {
	if number == 0 {
		return "green"
	}
	
	redNumbers := []int{1, 3, 5, 7, 9, 12, 14, 16, 18, 19, 21, 23, 25, 27, 30, 32, 34, 36}
	for _, red := range redNumbers {
		if number == red {
			return "red"
		}
	}
	
	return "black"
}

func (e *RouletteEngine) getDozen(number int) int {
	if number == 0 {
		return 0
	}
	if number <= 12 {
		return 1
	}
	if number <= 24 {
		return 2
	}
	return 3
}

func (e *RouletteEngine) getColumn(number int) int {
	if number == 0 {
		return 0
	}
	return ((number - 1) % 3) + 1
}

func (e *RouletteEngine) isBetWinner(bet RouletteBet, number int) bool {
	switch bet.Type {
	case BetStraight:
		return bet.Numbers[0] == number
		
	case BetSplit:
		for _, n := range bet.Numbers {
			if n == number {
				return true
			}
		}
		return false
		
	case BetRed:
		return e.getNumberColor(number) == "red"
		
	case BetBlack:
		return e.getNumberColor(number) == "black"
		
	case BetEven:
		return number > 0 && number%2 == 0
		
	case BetOdd:
		return number > 0 && number%2 == 1
		
	case BetLow:
		return number >= 1 && number <= 18
		
	case BetHigh:
		return number >= 19 && number <= 36
		
	case BetDozen1:
		return number >= 1 && number <= 12
		
	case BetDozen2:
		return number >= 13 && number <= 24
		
	case BetDozen3:
		return number >= 25 && number <= 36
		
	case BetColumn1:
		return e.getColumn(number) == 1
		
	case BetColumn2:
		return e.getColumn(number) == 2
		
	case BetColumn3:
		return e.getColumn(number) == 3
	}
	
	return false
}

func (e *RouletteEngine) calculatePayout(bet RouletteBet) float64 {
	payouts := map[RouletteBetType]float64{
		BetStraight: 35.0,  // 35:1
		BetSplit:    17.0,  // 17:1
		BetStreet:   11.0,  // 11:1
		BetCorner:   8.0,   // 8:1
		BetSixLine:  5.0,   // 5:1
		BetRed:      1.0,   // 1:1
		BetBlack:    1.0,   // 1:1
		BetEven:     1.0,   // 1:1
		BetOdd:      1.0,   // 1:1
		BetLow:      1.0,   // 1:1
		BetHigh:     1.0,   // 1:1
		BetDozen1:   2.0,   // 2:1
		BetDozen2:   2.0,   // 2:1
		BetDozen3:   2.0,   // 2:1
		BetColumn1:  2.0,   // 2:1
		BetColumn2:  2.0,   // 2:1
		BetColumn3:  2.0,   // 2:1
	}
	
	multiplier, exists := payouts[bet.Type]
	if !exists {
		return 0
	}
	
	// Return original bet + winnings
	return bet.Amount * (multiplier + 1)
}

// GetTableStatus returns current table status
func (e *RouletteEngine) GetTableStatus(ctx context.Context, tableID uuid.UUID) (string, error) {
	e.mu.RLock()
	table, exists := e.tables[tableID]
	e.mu.RUnlock()

	if !exists {
		return "", fmt.Errorf("table not found")
	}

	table.mu.RLock()
	defer table.mu.RUnlock()

	return table.Status, nil
}
