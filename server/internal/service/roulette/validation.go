package roulette

import (
	"strconv"
	"strings"

	"github.com/jokeoa/goigaming/internal/core/domain"
)

func ValidateBetType(betType string) bool {
	return domain.IsValidBetType(betType)
}

func ValidateBetValue(betType, betValue string) bool {
	switch betType {
	case "straight":
		n, err := strconv.Atoi(betValue)
		return err == nil && n >= 0 && n <= 36
	case "split":
		return validateNumberList(betValue, 2)
	case "street":
		return validateNumberList(betValue, 3)
	case "corner":
		return validateNumberList(betValue, 4)
	case "line":
		return validateNumberList(betValue, 6)
	case "dozen":
		return betValue == "1" || betValue == "2" || betValue == "3"
	case "column":
		return betValue == "1" || betValue == "2" || betValue == "3"
	case "red":
		return betValue == "red"
	case "black":
		return betValue == "black"
	case "odd":
		return betValue == "odd"
	case "even":
		return betValue == "even"
	case "high":
		return betValue == "high"
	case "low":
		return betValue == "low"
	default:
		return false
	}
}

func validateNumberList(value string, expectedCount int) bool {
	parts := strings.Split(value, ",")
	if len(parts) != expectedCount {
		return false
	}
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil || n < 0 || n > 36 {
			return false
		}
	}
	return true
}
