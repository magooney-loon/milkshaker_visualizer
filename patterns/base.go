package patterns

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Common utilities and types for all patterns

// HSVToRGB converts HSV color values to RGB tcell.Color
func HSVToRGB(h, s, v float64) tcell.Color {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return tcell.NewRGBColor(int32(r*255), int32(g*255), int32(b*255))
}

// Abs returns absolute value of integer
func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Min returns minimum of two integers
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// GetBasePhase returns current time-based phase for animations
func GetBasePhase() float64 {
	return float64(time.Now().UnixNano()) / 1e9
}

// RandomRune returns a random character from a predefined set
func RandomRune(rng *rand.Rand) rune {
	runes := []rune{'*', '+', 'x', 'o', '~', '@', '#', '$', '%', '&'}
	return runes[rng.Intn(len(runes))]
}

// PatternFunc defines the signature for pattern drawing functions
type PatternFunc func(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64)
