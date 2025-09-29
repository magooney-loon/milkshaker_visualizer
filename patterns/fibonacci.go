package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FibonacciVisualizer represents the main Fibonacci spiral visualizer
type FibonacciVisualizer struct {
	*tview.Box
	points     []float64
	fibonacci  []int
	angle      float64
	scale      float64
	depth      int
	colorCache map[int]tcell.Color
	sinCache   []float64
	cosCache   []float64
	lastUpdate time.Time
	patternMgr *Manager
}

// NewFibonacciVisualizer creates a new Fibonacci visualizer
func NewFibonacciVisualizer() *FibonacciVisualizer {
	v := &FibonacciVisualizer{
		Box:        tview.NewBox(),
		points:     make([]float64, 24),
		fibonacci:  generateFibonacci(25),
		angle:      0,
		scale:      1,
		depth:      4,
		colorCache: make(map[int]tcell.Color),
		sinCache:   make([]float64, 360),
		cosCache:   make([]float64, 360),
		lastUpdate: time.Now(),
		patternMgr: NewManager(),
	}

	// Initialize points array with default values
	for i := range v.points {
		v.points[i] = 0.1 // Small default amplitude
	}

	// Pre-calculate sin/cos values for performance
	for i := 0; i < 360; i++ {
		angle := float64(i) * math.Pi / 180
		v.sinCache[i] = math.Sin(angle)
		v.cosCache[i] = math.Cos(angle)
	}

	return v
}

// Draw renders the Fibonacci visualizer with organic, dynamic effects
func (v *FibonacciVisualizer) Draw(screen tcell.Screen) {
	now := time.Now()
	elapsed := now.Sub(v.lastUpdate).Seconds()
	v.lastUpdate = now

	x, y, width, height := v.GetInnerRect()
	centerX, centerY := x+width/2, y+height/2
	baseScale := math.Min(float64(width), float64(height)) / 300
	basePhase := GetBasePhase()

	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Organic character set with natural feel
	chars := []rune{'⋅', '∘', '◦', '○', '●', '◉', '⬡', '⬢', '◇', '◆', '✧', '✦', '·', '˙', '∙', '°', '⁘', '⁛', '⁝', '⚬', '⚭', '⚮'}

	// Multi-layered organic fibonacci spiral
	for layer := 0; layer < v.depth; layer++ {
		layerPhase := basePhase * (0.3 + float64(layer)*0.2)
		layerScale := baseScale * (1.2 - float64(layer)*0.15)

		// Generate organic Fibonacci spiral points
		for i := 0; i < len(v.fibonacci)-2 && i < len(v.points); i++ {
			amplitude := v.points[i%len(v.points)]

			// Safety check for division by zero
			var fibRatio float64 = 1.618 // Default to golden ratio
			if v.fibonacci[i] != 0 {
				fibRatio = float64(v.fibonacci[i+1]) / float64(v.fibonacci[i])
			}

			// Complex wave functions for organic movement
			wave1 := amplitude * math.Sin(layerPhase*1.7+float64(i)*0.3)
			wave2 := amplitude * 0.6 * math.Cos(layerPhase*1.2+float64(i)*0.7)
			wave3 := amplitude * 0.4 * math.Sin(layerPhase*2.3+float64(i)*0.2)
			organicOffset := wave1 + wave2 + wave3

			// Safety check for NaN/Inf values
			if math.IsNaN(organicOffset) || math.IsInf(organicOffset, 0) {
				organicOffset = 0
			}

			// Pulsing radius effect based on golden ratio
			baseRadius := float64(v.fibonacci[i]) * layerScale * v.scale
			pulseRadius := baseRadius * (1 + 0.15*math.Sin(layerPhase*1.8+float64(i)*goldenAngle))
			finalRadius := pulseRadius + organicOffset*5

			// Organic angle calculation with subtle variations
			rotationSpeed := 0.4 * (1 - float64(layer)*0.1)
			angleVariation := math.Sin(layerPhase*0.8+float64(i)*0.4) * 0.1
			spiralAngle := v.angle*rotationSpeed + float64(i)*goldenAngle + float64(layer)*0.3 + angleVariation

			// Calculate position with organic breathing effect
			breathe := 1 + 0.08*math.Sin(layerPhase*2.1+float64(layer)*1.2)
			// spiralX := float64(centerX) + finalRadius*math.Cos(spiralAngle)*breathe
			// spiralY := float64(centerY) + finalRadius*math.Sin(spiralAngle)*breathe

			// Draw organic spiral segments with varying density
			segmentDensity := 8 + int(amplitude*12)
			// Safety check to prevent excessive loops
			if segmentDensity > 30 {
				segmentDensity = 30
			}
			for segment := 0; segment < segmentDensity; segment++ {
				segmentRatio := float64(segment) / float64(segmentDensity)

				// Organic position interpolation with micro-variations
				microWave := 0.3 * math.Sin(layerPhase*4+float64(i*segment)*0.1)
				interpolatedRadius := finalRadius * (0.7 + segmentRatio*0.6 + microWave)
				interpolatedAngle := spiralAngle + segmentRatio*goldenAngle*0.3

				segmentX := float64(centerX) + interpolatedRadius*math.Cos(interpolatedAngle)*breathe
				segmentY := float64(centerY) + interpolatedRadius*math.Sin(interpolatedAngle)*breathe

				// Dynamic character selection based on organic factors
				charPhase := float64(i*segment+layer) + organicOffset*2
				if math.IsNaN(charPhase) || math.IsInf(charPhase, 0) {
					charPhase = 0
				}
				charIndex := int(math.Abs(charPhase)) % len(chars)
				displayChar := chars[charIndex]

				// Subtle, organic color generation
				colorKey := i*100 + layer*10 + segment
				color, exists := v.colorCache[colorKey]
				if !exists {
					color = v.getOrganicColor(i, layer, segment, amplitude, fibRatio, organicOffset, layerPhase)
					v.colorCache[colorKey] = color
				}

				// Add organic transparency effect
				if segmentRatio > 0.8 || amplitude < 0.1 {
					// Fade out at edges and low amplitude
					intensity := math.Max(0.3, 1-segmentRatio*1.5) * math.Max(0.3, amplitude*2)
					if intensity < 0.5 {
						displayChar = '·'
					}
				}

				// Safety check for screen coordinates
				screenX, screenY := int(segmentX), int(segmentY)
				if !math.IsNaN(segmentX) && !math.IsInf(segmentX, 0) && !math.IsNaN(segmentY) && !math.IsInf(segmentY, 0) {
					if screenX >= 0 && screenX < width && screenY >= 0 && screenY < height {
						screen.SetContent(screenX, screenY, displayChar, nil, tcell.StyleDefault.Foreground(color))
					}
				}

				// Add organic branching effects at golden ratio points
				if math.Mod(float64(i), goldenRatio*2) < 0.1 && amplitude > 0.4 && screenX >= 0 && screenX < width && screenY >= 0 && screenY < height {
					v.drawOrganicBranch(screen, screenX, screenY, interpolatedAngle, color, amplitude, layer, width, height)
				}
			}

			// Add spiral and starburst patterns on top of fibonacci spiral
			if amplitude > 0.3 {
				// Generate pattern color based on layer and fibonacci position
				patternColor := v.getOrganicColor(i, layer, 0, amplitude, fibRatio, organicOffset, layerPhase)
				rng := rand.New(rand.NewSource(time.Now().UnixNano() + int64(i*layer)))
				v.patternMgr.DrawRandomPattern(screen, rng, patternColor, amplitude)
			}
		}
	}

	// Smooth, organic angle updates
	v.angle += (0.15 + elapsed*0.1) * elapsed
	v.angle = math.Mod(v.angle, 2*math.Pi)

	// Periodic color cache refresh for organic color evolution
	if math.Mod(v.angle, math.Pi/4) < 0.01 {
		v.colorCache = make(map[int]tcell.Color)
	}
}

// drawOrganicBranch creates small organic branches at fibonacci points
func (v *FibonacciVisualizer) drawOrganicBranch(screen tcell.Screen, x, y int, baseAngle float64, color tcell.Color, amplitude float64, layer int, width, height int) {
	branchChars := []rune{'·', '˙', '∘', '◦'}
	branchLength := 2 + int(amplitude*4)

	// Create 2-3 organic branches
	numBranches := 2 + layer%2
	for branch := 0; branch < numBranches; branch++ {
		branchAngle := baseAngle + (float64(branch)*2-1)*0.6 + math.Sin(GetBasePhase()*2)*0.2

		for step := 1; step <= branchLength; step++ {
			branchX := x + int(float64(step)*math.Cos(branchAngle))
			branchY := y + int(float64(step)*math.Sin(branchAngle))

			// Safety check for branch coordinates
			if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
				charIndex := (step + branch) % len(branchChars)
				branchChar := branchChars[charIndex]

				// Fade branch color
				intensity := 1.0 - float64(step)/float64(branchLength*2)
				if intensity > 0.3 {
					screen.SetContent(branchX, branchY, branchChar, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

// UpdateWithPeak updates the visualizer with audio peak data
func (v *FibonacciVisualizer) UpdateWithPeak(peak float64) {
	// Organic point distribution based on fibonacci ratios
	goldenRatio := (1 + math.Sqrt(5)) / 2
	for i := range v.points {
		phaseOffset := float64(i) * math.Pi / (goldenRatio * 8)
		v.points[i] = peak * math.Sin(phaseOffset) * (1 + 0.3*math.Cos(phaseOffset*goldenRatio))
	}

	// Subtle scale and depth changes
	v.scale = 0.8 + peak*0.4
	v.depth = 3 + int(peak*4)
	if v.depth > 7 {
		v.depth = 7 // Cap depth for performance
	}
}

// getOrganicColor generates subtle, organic colors
func (v *FibonacciVisualizer) getOrganicColor(i, layer, segment int, amplitude, fibRatio, organicOffset, phase float64) tcell.Color {
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Base hue from fibonacci position and golden ratio
	baseHue := float64(i) / float64(len(v.fibonacci)) * goldenRatio

	// Organic hue variations
	hueShift := organicOffset*0.05 + math.Sin(phase*0.7+float64(layer)*0.4)*0.1
	hue := math.Mod(baseHue+hueShift+v.angle*0.08, 1)

	// Subtle saturation based on amplitude and layer depth
	saturation := 0.5 + amplitude*0.3 - float64(layer)*0.08
	saturation = math.Max(0.2, math.Min(0.9, saturation))

	// Organic value variations
	baseValue := 0.6 + amplitude*0.3
	valueVariation := math.Sin(phase*1.3+float64(segment)*0.2) * 0.1
	value := baseValue + valueVariation - float64(layer)*0.08
	value = math.Max(0.3, math.Min(0.95, value))

	return HSVToRGB(hue, saturation, value)
}

// generateFibonacci generates Fibonacci sequence up to n terms
func generateFibonacci(n int) []int {
	fib := make([]int, n)
	fib[0], fib[1] = 1, 1
	for i := 2; i < n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
		// Cap large values to prevent overflow
		if fib[i] > 10000 {
			fib[i] = 10000
		}
	}
	return fib
}
