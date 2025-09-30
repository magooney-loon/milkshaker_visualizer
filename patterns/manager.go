package patterns

import (
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

// Visualizator represents a group of patterns that work together
type Visualizator struct {
	Name     string
	Patterns []PatternFunc
	Enabled  []bool // Which patterns in the group are currently enabled
}

// Manager handles visualizator selection and pattern drawing
type Manager struct {
	visualizators   []Visualizator
	currentIndex    int
	shuffleEnabled  bool
	rng             *rand.Rand
	lastShuffleTime time.Time
	shuffleDuration time.Duration
}

// NewManager creates a new pattern manager with predefined visualizators
func NewManager() *Manager {
	visualizators := []Visualizator{
		{
			Name:     "Milkshaker",
			Patterns: []PatternFunc{DrawLogo},
			Enabled:  []bool{true},
		},
		{
			Name:     "Starburst",
			Patterns: []PatternFunc{DrawStarburst},
			Enabled:  []bool{true},
		},
		{
			Name:     "Fibonacci",
			Patterns: []PatternFunc{DrawFibonacci},
			Enabled:  []bool{true},
		},
		{
			Name:     "Wave",
			Patterns: []PatternFunc{DrawWave},
			Enabled:  []bool{true},
		},
		{
			Name:     "MixMax",
			Patterns: []PatternFunc{DrawStarburst, DrawFibonacci, DrawWave, DrawLogo},
			Enabled:  []bool{true, true, true, true},
		},
	}

	return &Manager{
		visualizators:   visualizators,
		currentIndex:    0,
		shuffleEnabled:  false,
		rng:             rand.New(rand.NewSource(42)),
		lastShuffleTime: time.Now(),
		shuffleDuration: 27 * time.Second,
	}
}

// GetCurrentVisualizatorName returns the name of the current visualizator
func (m *Manager) GetCurrentVisualizatorName() string {
	if m.currentIndex >= 0 && m.currentIndex < len(m.visualizators) {
		return m.visualizators[m.currentIndex].Name
	}
	return "Unknown"
}

// CycleVisualizator switches to the next visualizator group
func (m *Manager) CycleVisualizator() {
	if len(m.visualizators) > 1 {
		m.currentIndex = (m.currentIndex + 1) % len(m.visualizators)
	}
}

// ToggleShuffle toggles shuffle mode on/off
func (m *Manager) ToggleShuffle() {
	m.shuffleEnabled = !m.shuffleEnabled
	if m.shuffleEnabled {
		m.lastShuffleTime = time.Now() // Reset timer when enabling shuffle
	}
}

// IsShuffleEnabled returns whether shuffle is currently enabled
func (m *Manager) IsShuffleEnabled() bool {
	return m.shuffleEnabled
}

// ShuffleCurrentVisualizator randomly enables/disables patterns in current visualizator
func (m *Manager) ShuffleCurrentVisualizator() {
	if m.currentIndex < 0 || m.currentIndex >= len(m.visualizators) {
		return
	}

	current := &m.visualizators[m.currentIndex]

	// Ensure at least one pattern stays enabled
	enabledCount := 0
	for _, enabled := range current.Enabled {
		if enabled {
			enabledCount++
		}
	}

	// If only one pattern enabled, enable some more randomly
	if enabledCount <= 1 {
		for i := range current.Enabled {
			current.Enabled[i] = m.rng.Float64() < 0.7 // 70% chance to enable
		}
	} else {
		// Randomly toggle patterns
		for i := range current.Enabled {
			if m.rng.Float64() < 0.4 { // 40% chance to toggle
				current.Enabled[i] = !current.Enabled[i]
			}
		}
	}

	// Ensure at least one pattern is still enabled
	hasEnabled := false
	for _, enabled := range current.Enabled {
		if enabled {
			hasEnabled = true
			break
		}
	}

	if !hasEnabled && len(current.Enabled) > 0 {
		// Enable first pattern as fallback
		current.Enabled[0] = true
	}
}

// DrawCurrentVisualizator draws all enabled patterns in the current visualizator
func (m *Manager) DrawCurrentVisualizator(screen tcell.Screen, color tcell.Color, rng *rand.Rand, amplitude float64) {
	if m.currentIndex < 0 || m.currentIndex >= len(m.visualizators) {
		return
	}

	width, height := screen.Size()
	char := RandomRune(rng)
	current := m.visualizators[m.currentIndex]

	// Auto-shuffle: cycle visualizators every 27 seconds when shuffle is enabled
	if m.shuffleEnabled {
		if time.Since(m.lastShuffleTime) >= m.shuffleDuration {
			m.CycleVisualizator()
			m.lastShuffleTime = time.Now()
		}
	}

	// Draw all enabled patterns
	for i, pattern := range current.Patterns {
		if i < len(current.Enabled) && current.Enabled[i] {
			pattern(screen, width, height, color, char, rng, amplitude)
		}
	}
}

// GetVisualizatorCount returns the number of available visualizators
func (m *Manager) GetVisualizatorCount() int {
	return len(m.visualizators)
}

// GetCurrentVisualizatorIndex returns the current visualizator index
func (m *Manager) GetCurrentVisualizatorIndex() int {
	return m.currentIndex
}

// SetVisualizator sets the current visualizator by index
func (m *Manager) SetVisualizator(index int) {
	if index >= 0 && index < len(m.visualizators) {
		m.currentIndex = index
	}
}

// TogglePatternInCurrent toggles a specific pattern in current visualizator
func (m *Manager) TogglePatternInCurrent(patternIndex int) {
	if m.currentIndex < 0 || m.currentIndex >= len(m.visualizators) {
		return
	}

	current := &m.visualizators[m.currentIndex]
	if patternIndex >= 0 && patternIndex < len(current.Enabled) {
		current.Enabled[patternIndex] = !current.Enabled[patternIndex]

		// Ensure at least one pattern stays enabled
		hasEnabled := false
		for _, enabled := range current.Enabled {
			if enabled {
				hasEnabled = true
				break
			}
		}

		if !hasEnabled {
			// Re-enable the toggled pattern
			current.Enabled[patternIndex] = true
		}
	}
}

// GetCurrentPatternStates returns the enabled states of patterns in current visualizator
func (m *Manager) GetCurrentPatternStates() []bool {
	if m.currentIndex < 0 || m.currentIndex >= len(m.visualizators) {
		return nil
	}

	current := m.visualizators[m.currentIndex]
	states := make([]bool, len(current.Enabled))
	copy(states, current.Enabled)
	return states
}

// GetCurrentPatternNames returns the names of patterns in current visualizator
func (m *Manager) GetCurrentPatternNames() []string {
	if m.currentIndex < 0 || m.currentIndex >= len(m.visualizators) {
		return nil
	}

	current := m.visualizators[m.currentIndex]
	names := make([]string, len(current.Patterns))

	for i, pattern := range current.Patterns {
		// Simple name mapping based on function comparison
		// This is hacky but functional for our purposes
		switch {
		case isSameFunction(pattern, DrawStarburst):
			names[i] = "Starburst"
		case isSameFunction(pattern, DrawFibonacci):
			names[i] = "Fibonacci"
		case isSameFunction(pattern, DrawLogo):
			names[i] = "Logo"
		case isSameFunction(pattern, DrawWave):
			names[i] = "Wave"
		default:
			names[i] = "Unknown"
		}
	}

	// Fallback to generic names if function comparison fails
	for i, name := range names {
		if name == "Unknown" {
			names[i] = "Pattern " + string(rune('A'+i))
		}
	}

	return names
}

// Helper function to compare functions (basic implementation)
func isSameFunction(f1, f2 PatternFunc) bool {
	// This is a simple comparison - in a real implementation you might want
	// a more robust way to identify functions
	return &f1 == &f2
}
