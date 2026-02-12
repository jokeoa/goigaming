package templates

import (
    "fmt"
    "html/template"
    
    "github.com/shopspring/decimal"
)

// GetFuncMap возвращает кастомные template functions
func GetFuncMap() template.FuncMap {
    return template.FuncMap{
        "formatDecimal": formatDecimal,
        "isRed":        isRed,
        "iterate":      iterate,
    }
}

func formatDecimal(d decimal.Decimal) string {
    return d.StringFixed(2)
}

func isRed(number int) bool {
    redNumbers := map[int]bool{
        1: true, 3: true, 5: true, 7: true, 9: true,
        12: true, 14: true, 16: true, 18: true,
        19: true, 21: true, 23: true, 25: true, 27: true,
        30: true, 32: true, 34: true, 36: true,
    }
    return redNumbers[number]
}

func iterate(start, end int) []int {
    result := make([]int, end-start)
    for i := range result {
        result[i] = start + i
    }
    return result
}
