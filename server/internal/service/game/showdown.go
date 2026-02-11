package game

import (
	"github.com/google/uuid"
	"github.com/jokeoa/goigaming/internal/core/domain"
	"github.com/shopspring/decimal"
)

type HandPlayerCards struct {
	PlayerID  uuid.UUID
	HoleCards []domain.Card
}

func DetermineWinners(handPlayers []HandPlayerCards, communityCards []domain.Card, pots []domain.Pot) domain.HandResult {
	handMap := make(map[uuid.UUID]HandRank, len(handPlayers))
	cardMap := make(map[uuid.UUID][]domain.Card, len(handPlayers))

	for _, hp := range handPlayers {
		rank := BestHand(hp.HoleCards, communityCards)
		handMap[hp.PlayerID] = rank
		cardMap[hp.PlayerID] = hp.HoleCards
	}

	var allWinners []domain.WinnerInfo
	showdownCards := make(map[uuid.UUID][]domain.Card)

	for _, pot := range pots {
		winners := resolvePot(pot, handMap)
		share := pot.Amount.Div(decimal.NewFromInt(int64(len(winners))))

		for _, winnerID := range winners {
			allWinners = append(allWinners, domain.WinnerInfo{
				PlayerID: winnerID,
				Amount:   share,
				HandRank: handMap[winnerID].Name,
			})
			showdownCards[winnerID] = cardMap[winnerID]
		}
	}

	allWinners = consolidateWinners(allWinners)

	return domain.HandResult{
		Winners:       allWinners,
		Pots:          pots,
		ShowdownCards: showdownCards,
	}
}

func resolvePot(pot domain.Pot, handMap map[uuid.UUID]HandRank) []uuid.UUID {
	if len(pot.EligibleIDs) == 0 {
		return nil
	}

	var eligibleRanks []HandRank
	for _, id := range pot.EligibleIDs {
		if rank, ok := handMap[id]; ok {
			eligibleRanks = append(eligibleRanks, rank)
		}
	}

	if len(eligibleRanks) == 0 {
		return nil
	}

	winnerIndices := CompareHands(eligibleRanks)
	winners := make([]uuid.UUID, len(winnerIndices))
	for i, idx := range winnerIndices {
		winners[i] = pot.EligibleIDs[idx]
	}

	return winners
}

func consolidateWinners(winners []domain.WinnerInfo) []domain.WinnerInfo {
	totals := make(map[uuid.UUID]domain.WinnerInfo)

	for _, w := range winners {
		if existing, ok := totals[w.PlayerID]; ok {
			totals[w.PlayerID] = domain.WinnerInfo{
				PlayerID: w.PlayerID,
				Amount:   existing.Amount.Add(w.Amount),
				HandRank: w.HandRank,
			}
		} else {
			totals[w.PlayerID] = w
		}
	}

	result := make([]domain.WinnerInfo, 0, len(totals))
	for _, w := range totals {
		result = append(result, w)
	}

	return result
}
