package engine

import (
	"math/rand"
)

func RandomRangeInt(a int, b int) int {
	return a + rand.Intn(b-a+1)
}

func RandomRangeF64(a float64, b float64) float64 {
	return a + rand.Float64()*(b-a)
}
