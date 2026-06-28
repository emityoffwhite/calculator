package ast

import "strconv"

// formatFloat форматирует число без лишних нулей после запятой,
// например 5.0 -> "5", 5.5 -> "5.5".
func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}
