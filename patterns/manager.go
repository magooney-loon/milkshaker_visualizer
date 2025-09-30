package patterns

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// Manager handles pattern selection and drawing
type Manager struct {
	patterns []PatternFunc
}

// NewManager creates a new pattern manager with all available patterns
func NewManager() *Manager {
	return &Manager{
		patterns: []PatternFunc{
			DrawField,
			DrawSpiral,
			DrawStarburst,
			DrawGeometry,
			DrawLogo,
		},
	}
}

// DrawRandomPattern selects and draws a random pattern based on amplitude
func (m *Manager) DrawRandomPattern(screen tcell.Screen, rng *rand.Rand, color tcell.Color, amplitude float64) {
	width, height := screen.Size()
	char := RandomRune(rng)

	patternIndex := int(amplitude * float64(len(m.patterns)))
	if patternIndex >= len(m.patterns) {
		patternIndex = len(m.patterns) - 1
	}

	pattern := m.patterns[patternIndex]
	pattern(screen, width, height, color, char, rng, amplitude)
}

// GetPatternCount returns the number of available patterns
func (m *Manager) GetPatternCount() int {
	return len(m.patterns)
}

// DrawPattern draws a specific pattern by index
func (m *Manager) DrawPattern(screen tcell.Screen, patternIndex int, color tcell.Color, rng *rand.Rand, amplitude float64) {
	if patternIndex < 0 || patternIndex >= len(m.patterns) {
		return
	}

	width, height := screen.Size()
	char := RandomRune(rng)
	pattern := m.patterns[patternIndex]
	pattern(screen, width, height, color, char, rng, amplitude)
}
