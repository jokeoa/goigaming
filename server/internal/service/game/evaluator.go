package game

import (
	"github.com/chehsunliu/poker"
	"github.com/jokeoa/goigaming/internal/core/domain"
)

type HandRank struct {
	Score int32
	Name  string
}

func toLibCard(c domain.Card) poker.Card {
	rankStr := string(c.Rank)
	suitStr := string(c.Suit)
	return poker.NewCard(rankStr + suitStr)
}

func toLibCards(cards []domain.Card) []poker.Card {
	result := make([]poker.Card, len(cards))
	for i, c := range cards {
		result[i] = toLibCard(c)
	}
	return result
}

func EvaluateHand(cards []domain.Card) HandRank {
	libCards := toLibCards(cards)
	score := poker.Evaluate(libCards)
	return HandRank{
		Score: score,
		Name:  poker.RankString(score),
	}
}

func BestHand(holeCards, community []domain.Card) HandRank {
	all := make([]domain.Card, 0, len(holeCards)+len(community))
	all = append(all, holeCards...)
	all = append(all, community...)
	return EvaluateHand(all)
}

func CompareHands(hands []HandRank) []int {
	if len(hands) == 0 {
		return nil
	}

	bestScore := hands[0].Score
	for _, h := range hands[1:] {
		if h.Score < bestScore {
			bestScore = h.Score
		}
	}

	var winners []int
	for i, h := range hands {
		if h.Score == bestScore {
			winners = append(winners, i)
		}
	}

	return winners
}
