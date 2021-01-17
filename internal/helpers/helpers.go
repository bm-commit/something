package helpers

import (
	"math"
)

// Round value to specific unit
func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
