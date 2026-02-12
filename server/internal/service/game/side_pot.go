package game

import (
	"sort"

	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type PotContribution struct {
	PlayerID uuid.UUID
	TotalBet decimal.Decimal
	IsAllIn  bool
	IsFolded bool
}

func CalculateSidePots(contributions []PotContribution) []domain.Pot {
	allInLevels := collectAllInLevels(contributions)

	if len(allInLevels) == 0 {
		return []domain.Pot{buildMainPot(contributions)}
	}

	var pots []domain.Pot
	prevLevel := decimal.Zero

	for _, level := range allInLevels {
		pot := buildPotAtLevel(contributions, prevLevel, level)
		if pot.Amount.IsPositive() {
			pots = append(pots, pot)
		}
		prevLevel = level
	}

	remainder := buildRemainderPot(contributions, prevLevel)
	if remainder.Amount.IsPositive() {
		pots = append(pots, remainder)
	}

	return pots
}

func collectAllInLevels(contributions []PotContribution) []decimal.Decimal {
	levelMap := make(map[string]decimal.Decimal)
	for _, c := range contributions {
		if c.IsAllIn {
			key := c.TotalBet.String()
			levelMap[key] = c.TotalBet
		}
	}

	levels := make([]decimal.Decimal, 0, len(levelMap))
	for _, v := range levelMap {
		levels = append(levels, v)
	}

	sort.Slice(levels, func(i, j int) bool {
		return levels[i].LessThan(levels[j])
	})

	return levels
}

func buildMainPot(contributions []PotContribution) domain.Pot {
	amount := decimal.Zero
	var eligible []uuid.UUID

	for _, c := range contributions {
		amount = amount.Add(c.TotalBet)
		if !c.IsFolded {
			eligible = append(eligible, c.PlayerID)
		}
	}

	return domain.Pot{
		Amount:      amount,
		EligibleIDs: eligible,
	}
}

func buildPotAtLevel(contributions []PotContribution, prevLevel, level decimal.Decimal) domain.Pot {
	amount := decimal.Zero
	var eligible []uuid.UUID

	for _, c := range contributions {
		contribution := decimal.Min(c.TotalBet, level).Sub(prevLevel)
		if contribution.IsPositive() {
			amount = amount.Add(contribution)
		}
		if !c.IsFolded && c.TotalBet.GreaterThanOrEqual(level) {
			eligible = append(eligible, c.PlayerID)
		}
	}

	return domain.Pot{
		Amount:      amount,
		EligibleIDs: eligible,
	}
}

func buildRemainderPot(contributions []PotContribution, lastLevel decimal.Decimal) domain.Pot {
	amount := decimal.Zero
	var eligible []uuid.UUID

	for _, c := range contributions {
		remainder := c.TotalBet.Sub(lastLevel)
		if remainder.IsPositive() {
			amount = amount.Add(remainder)
		}
		if !c.IsFolded && c.TotalBet.GreaterThan(lastLevel) {
			eligible = append(eligible, c.PlayerID)
		}
	}

	return domain.Pot{
		Amount:      amount,
		EligibleIDs: eligible,
	}
}
