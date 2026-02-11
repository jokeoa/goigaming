package services

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jokeoa/igaming/models"
	"github.com/jokeoa/igaming/repository"
)

type RouletteService struct {
	tableRepo *repository.RouletteTableRepository
	roundRepo *repository.RouletteRoundRepository
	betRepo   *repository.RouletteBetRepository
}

func NewRouletteService(
	tableRepo *repository.RouletteTableRepository,
	roundRepo *repository.RouletteRoundRepository,
	betRepo *repository.RouletteBetRepository,
) *RouletteService {
	return &RouletteService{
		tableRepo: tableRepo,
		roundRepo: roundRepo,
		betRepo:   betRepo,
	}
}

func (s *RouletteService) PlaceBet(ctx context.Context, userID uuid.UUID, tableID uuid.UUID, betType, betValue string, amount float64) (models.RouletteBet, error) {
	currentRound, err := s.roundRepo.FindCurrent(ctx, tableID)
	if err != nil {
		return models.RouletteBet{}, fmt.Errorf("no active round: %w", err)
	}

	if currentRound.BettingEndsAt != nil && time.Now().After(*currentRound.BettingEndsAt) {
		return models.RouletteBet{}, fmt.Errorf("betting time ended")
	}

	validBetTypes := map[string]bool{
		"straight": true, "split": true, "street": true, "corner": true,
		"line": true, "dozen": true, "column": true, "red": true,
		"black": true, "odd": true, "even": true, "high": true, "low": true,
	}
	if !validBetTypes[betType] {
		return models.RouletteBet{}, fmt.Errorf("invalid bet type: %s", betType)
	}

	bet := models.RouletteBet{
		RoundID:  currentRound.ID,
		UserID:   userID,
		BetType:  betType,
		BetValue: betValue,
		Amount:   amount,
		Status:   "pending",
	}

	return s.betRepo.Create(ctx, bet)
}

func (s *RouletteService) Spin(ctx context.Context, tableID uuid.UUID) (models.RouletteRound, error) {
	currentRound, err := s.roundRepo.FindCurrent(ctx, tableID)
	if err != nil {
		return models.RouletteRound{}, fmt.Errorf("no active round: %w", err)
	}

	result, seed := generateRouletteResult()
	color := getNumberColor(result)

	now := time.Now()
	currentRound.Result = &result
	currentRound.ResultColor = &color
	currentRound.SeedRevealed = &seed
	currentRound.SettledAt = &now

	updatedRound, err := s.roundRepo.Update(ctx, currentRound)
	if err != nil {
		return models.RouletteRound{}, err
	}

	bets, err := s.betRepo.FindByRoundId(ctx, currentRound.ID)
	if err != nil {
		return models.RouletteRound{}, err
	}

	for _, bet := range bets {
		payout := calculatePayout(bet, result, color)
		bet.Payout = payout
		if payout > 0 {
			bet.Status = "won"
		} else {
			bet.Status = "lost"
		}
		_, _ = s.betRepo.Update(ctx, bet)
	}

	return updatedRound, nil
}

func (s *RouletteService) GetHistory(ctx context.Context, tableID uuid.UUID, limit int) ([]models.RouletteRound, error) {
	round, err := s.roundRepo.FindCurrent(ctx, tableID)
	if err != nil {
		return []models.RouletteRound{}, nil
	}
	return []models.RouletteRound{round}, nil
}

func (s *RouletteService) StartNewRound(ctx context.Context, tableID uuid.UUID, roundNumber int, bettingDuration time.Duration) (models.RouletteRound, error) {
	seed := generateSeed()
	seedHash := hashSeed(seed)
	bettingEndsAt := time.Now().Add(bettingDuration)

	round := models.RouletteRound{
		TableID:       tableID,
		RoundNumber:   roundNumber,
		SeedHash:      &seedHash,
		BettingEndsAt: &bettingEndsAt,
	}

	return s.roundRepo.Create(ctx, round)
}

func generateRouletteResult() (int, string) {
	seed := generateSeed()
	hash := sha256.Sum256([]byte(seed))
	num := new(big.Int).SetBytes(hash[:])
	result := new(big.Int).Mod(num, big.NewInt(37)).Int64()
	return int(result), seed
}

func generateSeed() string {
	bytes := make([]byte, 32)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func hashSeed(seed string) string {
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])
}

func getNumberColor(num int) string {
	if num == 0 {
		return "green"
	}
	redNumbers := map[int]bool{
		1: true, 3: true, 5: true, 7: true, 9: true, 12: true,
		14: true, 16: true, 18: true, 19: true, 21: true, 23: true,
		25: true, 27: true, 30: true, 32: true, 34: true, 36: true,
	}
	if redNumbers[num] {
		return "red"
	}
	return "black"
}

func calculatePayout(bet models.RouletteBet, result int, color string) float64 {
	won := false

	switch bet.BetType {
	case "straight":
		if fmt.Sprintf("%d", result) == bet.BetValue {
			return bet.Amount * 35
		}
	case "red":
		if color == "red" {
			won = true
		}
	case "black":
		if color == "black" {
			won = true
		}
	case "odd":
		if result%2 == 1 && result != 0 {
			won = true
		}
	case "even":
		if result%2 == 0 && result != 0 {
			won = true
		}
	case "low":
		if result >= 1 && result <= 18 {
			won = true
		}
	case "high":
		if result >= 19 && result <= 36 {
			won = true
		}
	case "dozen":
		dozen := (result - 1) / 12
		if result > 0 && fmt.Sprintf("%d", dozen+1) == bet.BetValue {
			return bet.Amount * 2
		}
	case "column":
		if result > 0 {
			column := ((result - 1) % 3) + 1
			if fmt.Sprintf("%d", column) == bet.BetValue {
				return bet.Amount * 2
			}
		}
	}

	if won {
		return bet.Amount
	}

	return 0
}
