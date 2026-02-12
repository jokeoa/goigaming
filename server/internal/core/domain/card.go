package domain

import (
	"fmt"
	"strings"
)

type Suit string

const (
	SuitSpades   Suit = "s"
	SuitHearts   Suit = "h"
	SuitDiamonds Suit = "d"
	SuitClubs    Suit = "c"
)

type Rank string

const (
	RankTwo   Rank = "2"
	RankThree Rank = "3"
	RankFour  Rank = "4"
	RankFive  Rank = "5"
	RankSix   Rank = "6"
	RankSeven Rank = "7"
	RankEight Rank = "8"
	RankNine  Rank = "9"
	RankTen   Rank = "T"
	RankJack  Rank = "J"
	RankQueen Rank = "Q"
	RankKing  Rank = "K"
	RankAce   Rank = "A"
)

var allSuits = []Suit{SuitSpades, SuitHearts, SuitDiamonds, SuitClubs}

var allRanks = []Rank{
	RankTwo, RankThree, RankFour, RankFive, RankSix, RankSeven,
	RankEight, RankNine, RankTen, RankJack, RankQueen, RankKing, RankAce,
}

type Card struct {
	Rank Rank `json:"rank"`
	Suit Suit `json:"suit"`
}

func (c Card) String() string {
	return string(c.Rank) + string(c.Suit)
}

func ParseCard(s string) (Card, error) {
	s = strings.TrimSpace(s)
	if len(s) != 2 {
		return Card{}, fmt.Errorf("invalid card string: %q", s)
	}

	rank := Rank(s[0:1])
	suit := Suit(s[1:2])

	if !isValidRank(rank) {
		return Card{}, fmt.Errorf("invalid rank: %q", rank)
	}
	if !isValidSuit(suit) {
		return Card{}, fmt.Errorf("invalid suit: %q", suit)
	}

	return Card{Rank: rank, Suit: suit}, nil
}

func ParseCards(csv string) ([]Card, error) {
	csv = strings.TrimSpace(csv)
	if csv == "" {
		return nil, nil
	}

	parts := strings.Split(csv, ",")
	cards := make([]Card, 0, len(parts))
	for _, p := range parts {
		c, err := ParseCard(p)
		if err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

func CardsToString(cards []Card) string {
	strs := make([]string, len(cards))
	for i, c := range cards {
		strs[i] = c.String()
	}
	return strings.Join(strs, ",")
}

func FullDeck() []Card {
	deck := make([]Card, 0, 52)
	for _, s := range allSuits {
		for _, r := range allRanks {
			deck = append(deck, Card{Rank: r, Suit: s})
		}
	}
	return deck
}

func isValidRank(r Rank) bool {
	for _, rank := range allRanks {
		if r == rank {
			return true
		}
	}
	return false
}

func isValidSuit(s Suit) bool {
	for _, suit := range allSuits {
		if s == suit {
			return true
		}
	}
	return false
}
