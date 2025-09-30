package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type FibonacciParticle struct {
	x, y      float64
	vx, vy    float64
	life      float64
	maxLife   float64
	intensity float64
	hue       float64
	size      float64
	char      rune
	fibIndex  int
}

type GoldenRatio struct {
	x, y      float64
	radius    float64
	angle     float64
	intensity float64
	life      float64
	maxLife   float64
}

type SacredGeometry struct {
	centerX   int
	centerY   int
	radius    float64
	angles    []float64
	intensity float64
	life      float64
	pattern   int
}

type NumberSequence struct {
	x, y      int
	number    int
	intensity float64
	life      float64
	maxLife   float64
	hue       float64
}

var (
	// Mathematical particle system
	fibParticles    []FibonacciParticle
	maxFibParticles = 120

	// Golden ratio effects
	goldenRatios    []GoldenRatio
	maxGoldenRatios = 20

	// Sacred geometry patterns
	sacredGeometry []SacredGeometry
	maxSacredGeo   = 8

	// Number sequences
	numberSequences []NumberSequence
	maxNumbers      = 15

	// Animation phases
	goldenPhase   float64 = 0.0
	spiralPhase   float64 = 0.0
	mathPhase     float64 = 0.0
	fibLastUpdate time.Time

	// Peak tracking for mathematical beauty
	fibPeakHistory []float64
	maxFibHistory  = 25

	// Golden ratio constant
	goldenRatio = (1 + math.Sqrt(5)) / 2
	goldenAngle = 2 * math.Pi / (goldenRatio * goldenRatio)
)

// DrawFibonacci creates an epic mathematical fibonacci visualization with sacred geometry
func DrawFibonacci(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	now := time.Now()
	elapsed := now.Sub(fibLastUpdate).Seconds()
	if elapsed < 1.0/180.0 { // 180 FPS limit
		return
	}
	fibLastUpdate = now

	// Track peak history for mathematical progression
	fibPeakHistory = append(fibPeakHistory, peak)
	if len(fibPeakHistory) > maxFibHistory {
		fibPeakHistory = fibPeakHistory[1:]
	}

	// Calculate mathematical progression
	mathProgression := 0.0
	if len(fibPeakHistory) > 10 {
		recent := fibPeakHistory[len(fibPeakHistory)-5:]
		avgRecent := 0.0
		for _, p := range recent {
			avgRecent += p
		}
		avgRecent /= float64(len(recent))
		mathProgression = avgRecent
	}

	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()

	// Update mathematical phases with golden ratio timing
	speedMultiplier := 1.0 + peak*2.0 + mathProgression*1.5
	goldenPhase += elapsed * speedMultiplier * goldenRatio * 0.5
	spiralPhase += elapsed * speedMultiplier * 1.618
	mathPhase += elapsed * speedMultiplier * 2.618

	// Update all mathematical systems
	updateFibonacciParticles(elapsed, peak, mathProgression, width, height, centerX, centerY, rng)
	updateGoldenRatios(elapsed, peak, mathProgression, centerX, centerY, rng)
	updateSacredGeometry(elapsed, peak, mathProgression, centerX, centerY, width, height, rng)
	updateNumberSequences(elapsed, peak, mathProgression, width, height, rng)

	// Draw main fibonacci spiral with enhancements
	drawEpicFibonacciSpiral(screen, width, height, centerX, centerY, peak, mathProgression, basePhase, rng)

	// Draw mathematical effects
	drawGoldenRatios(screen, width, height)
	drawSacredGeometry(screen, width, height)
	drawFibonacciParticles(screen, width, height)
	drawNumberSequences(screen, width, height)

	// Draw mathematical core
	drawMathematicalCore(screen, centerX, centerY, peak, mathProgression, basePhase)
}

func drawEpicFibonacciSpiral(screen tcell.Screen, width, height, centerX, centerY int, peak, mathProgression, basePhase float64, rng *rand.Rand) {
	// Dynamic spiral parameters
	maxRadius := math.Min(float64(width), float64(height)) / 2.5
	peakScale := 0.5 + peak*0.8 + mathProgression*0.3

	// Number of fibonacci terms with mathematical progression
	maxTerms := 18 + int(peak*12) + int(mathProgression*8)
	if maxTerms > 40 {
		maxTerms = 40
	}

	// Generate fibonacci sequence
	fib := make([]int, maxTerms)
	if maxTerms >= 1 {
		fib[0] = 1
	}
	if maxTerms >= 2 {
		fib[1] = 1
	}
	for i := 2; i < maxTerms; i++ {
		fib[i] = fib[i-1] + fib[i-2]
		if fib[i] > 2000 {
			fib[i] = 2000
		}
	}

	// Mathematical character progression
	mathChars := [][]rune{
		{'·', '∘', '○', '●', '◉', '⬢', '◆', '★', '✦', '✧'},  // Basic progression
		{'0', '1', '2', '3', '5', '8', '#', '@', '%', '&'},  // Fibonacci numbers (ASCII safe)
		{'◢', '◣', '◤', '◥', '▰', '▱', '◐', '◑', '◒', '◓'},  // Geometric shapes
		{'+', 'x', '*', '~', '=', '-', '|', '/', '\\', '^'}, // Sacred geometry (ASCII safe)
	}

	// Multiple spiral arms with golden ratio spacing
	numArms := 3 + int(peak*3) + int(mathProgression*2)
	if numArms > 8 {
		numArms = 8
	}

	for arm := 0; arm < numArms; arm++ {
		armOffset := float64(arm) * goldenAngle
		armPhase := spiralPhase*0.1 + armOffset

		// Draw enhanced fibonacci spiral
		for i := 4; i < len(fib); i++ {
			if fib[i] == 0 {
				continue
			}

			// Calculate spiral position using golden ratio mathematics
			fibRadius := math.Pow(float64(fib[i]), 0.618) * maxRadius * peakScale / 15.0
			spiralAngle := float64(i)*goldenAngle + armPhase

			// Golden ratio variations
			goldenVariation := 1 + 0.2*math.Sin(armPhase*goldenRatio+float64(i)*0.618)
			finalRadius := fibRadius * goldenVariation

			// Mathematical density based on fibonacci properties
			density := int(math.Log(float64(fib[i]))) + 3
			if density > 20 {
				density = 20
			}

			// Draw mathematical arc
			for point := 0; point < density; point++ {
				pointRatio := float64(point) / float64(density)

				// Golden ratio interpolation
				var prevRadius float64
				if i > 4 {
					prevRadius = math.Pow(float64(fib[i-1]), 0.618) * maxRadius * peakScale / 15.0
				}

				interpRadius := prevRadius + (finalRadius-prevRadius)*math.Pow(pointRatio, 1/goldenRatio)
				interpAngle := spiralAngle - goldenAngle*(1-pointRatio)

				x := centerX + int(interpRadius*math.Cos(interpAngle))
				y := centerY + int(interpRadius*math.Sin(interpAngle))

				if x >= 0 && x < width && y >= 0 && y < height {
					// Mathematical intensity calculation
					fibIntensity := (1.0 - pointRatio*0.4) * peakScale
					termIntensity := 1.0 - math.Pow(float64(i)/float64(len(fib)), 0.618)
					goldenIntensity := math.Sin(armPhase*goldenRatio+float64(point)*0.618)*0.3 + 0.7
					totalIntensity := fibIntensity * termIntensity * goldenIntensity * (0.4 + peak*0.6)

					// Character selection based on mathematical properties
					var finalChar rune
					morphLevel := totalIntensity + mathProgression*0.3

					if morphLevel < 0.15 {
						finalChar = mathChars[0][0] // ·
					} else if morphLevel < 0.3 {
						finalChar = mathChars[0][1] // ∘
					} else if morphLevel < 0.45 {
						finalChar = mathChars[0][2] // ○
					} else if morphLevel < 0.6 {
						finalChar = mathChars[0][3] // ●
					} else if morphLevel < 0.75 {
						// Use fibonacci numbers or geometric shapes
						if i < len(mathChars[1]) && peak > 0.4 {
							finalChar = mathChars[1][i%len(mathChars[1])]
						} else {
							finalChar = mathChars[2][point%len(mathChars[2])]
						}
					} else if morphLevel < 0.9 {
						finalChar = mathChars[3][int(math.Mod(float64(i*point), float64(len(mathChars[3]))))]
					} else {
						// Epic mathematical symbols
						epicChars := []rune{'φ', '∞', '∑', '∏', '∫', '∂', '√', '∆', '∇', '⊕'}
						finalChar = epicChars[int(math.Mod(goldenPhase*3.7+float64(i), float64(len(epicChars))))]
					}

					// Golden ratio color system
					hueBase := float64(i)/float64(len(fib))*goldenRatio + armOffset/(2*math.Pi)
					hueShift := math.Sin(spiralPhase*0.3+float64(i)*0.618) * 0.1
					mathHue := math.Sin(mathPhase*0.1+interpAngle) * 0.08
					finalHue := math.Mod(hueBase+hueShift+mathHue, 1.0)

					saturation := 0.6 + peak*0.3 + totalIntensity*0.2
					saturation = math.Max(0.3, math.Min(0.9, saturation))

					value := 0.3 + totalIntensity*0.6 + mathProgression*0.2
					value = math.Max(0.1, math.Min(1.0, value))

					spiralColor := HSVToRGB(finalHue, saturation, value)

					if totalIntensity > 0.12 {
						screen.SetContent(x, y, finalChar, nil, tcell.StyleDefault.Foreground(spiralColor))
					}
				}
			}

			// Draw golden ratio connecting lines
			if peak > 0.5 && i > 4 && arm == 0 && i%3 == 0 {
				drawGoldenConnections(screen, centerX, centerY, finalRadius, spiralAngle, fib, i, maxRadius, peakScale, armPhase, width, height, peak)
			}
		}
	}
}

func updateFibonacciParticles(elapsed, peak, mathProgression float64, width, height, centerX, centerY int, rng *rand.Rand) {
	// Spawn mathematical particles
	spawnRate := peak*8.0 + mathProgression*6.0
	if len(fibParticles) < maxFibParticles && rng.Float64() < spawnRate*elapsed {
		// Spawn from fibonacci positions
		fibIndex := 3 + rng.Intn(15)
		fibValue := 1
		for i := 0; i < fibIndex; i++ {
			if i < 2 {
				fibValue = 1
			} else {
				prevFib := fibValue
				fibValue = fibValue + prevFib
			}
		}

		angle := float64(fibIndex) * goldenAngle
		radius := math.Sqrt(float64(fibValue)) * 3.0

		particle := FibonacciParticle{
			x:         float64(centerX) + radius*math.Cos(angle),
			y:         float64(centerY) + radius*math.Sin(angle),
			vx:        math.Cos(angle+math.Pi/2) * (15.0 + peak*25.0) * (0.5 + rng.Float64()),
			vy:        math.Sin(angle+math.Pi/2) * (15.0 + peak*25.0) * (0.5 + rng.Float64()),
			life:      1.0,
			maxLife:   1.5 + rng.Float64()*2.5,
			intensity: 0.6 + rng.Float64()*0.4 + mathProgression*0.3,
			hue:       math.Mod(goldenPhase*0.1+float64(fibIndex)*0.618, 1.0),
			size:      1.0 + rng.Float64()*2.0 + peak,
			char:      []rune{'·', '∘', '○', '●', '◉', '⬢', '★', '✦'}[rng.Intn(8)],
			fibIndex:  fibIndex,
		}
		fibParticles = append(fibParticles, particle)
	}

	// Update particles with golden ratio physics
	for i := len(fibParticles) - 1; i >= 0; i-- {
		p := &fibParticles[i]

		// Golden ratio spiral motion
		p.x += p.vx * elapsed
		p.y += p.vy * elapsed
		p.life -= elapsed / p.maxLife

		// Fibonacci spiral attraction
		centerDx := float64(centerX) - p.x
		centerDy := float64(centerY) - p.y
		centerDist := math.Sqrt(centerDx*centerDx + centerDy*centerDy)

		if centerDist > 0 {
			spiralForce := 10.0 * elapsed / centerDist
			p.vx += centerDx * spiralForce
			p.vy += centerDy * spiralForce
		}

		// Golden ratio rotation
		rotationSpeed := goldenAngle * elapsed * 0.5
		newVx := p.vx*math.Cos(rotationSpeed) - p.vy*math.Sin(rotationSpeed)
		newVy := p.vx*math.Sin(rotationSpeed) + p.vy*math.Cos(rotationSpeed)
		p.vx = newVx * 0.98
		p.vy = newVy * 0.98

		// Remove dead particles
		if p.life <= 0 || centerDist > 200 {
			fibParticles = append(fibParticles[:i], fibParticles[i+1:]...)
		}
	}
}

func drawFibonacciParticles(screen tcell.Screen, width, height int) {
	for _, p := range fibParticles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			alpha := p.life * p.intensity
			if alpha > 0.1 {
				// Use fibonacci-based characters
				fibChars := []rune{'·', '∘', '○', '●', '◉', '⬢', '★', '✦', 'φ', '∞'}
				charIndex := p.fibIndex % len(fibChars)
				char := fibChars[charIndex]

				saturation := 0.7 + alpha*0.3
				value := alpha * 0.9
				color := HSVToRGB(p.hue, saturation, value)

				screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func updateGoldenRatios(elapsed, peak, mathProgression float64, centerX, centerY int, rng *rand.Rand) {
	// Create golden ratio patterns
	if len(goldenRatios) < maxGoldenRatios && rng.Float64() < peak*2.0*elapsed {
		golden := GoldenRatio{
			x:         float64(centerX) + (rng.Float64()-0.5)*100,
			y:         float64(centerY) + (rng.Float64()-0.5)*100,
			radius:    5.0 + rng.Float64()*20.0,
			angle:     rng.Float64() * 2 * math.Pi,
			intensity: 0.7 + peak*0.3,
			life:      1.0,
			maxLife:   2.0 + rng.Float64()*3.0,
		}
		goldenRatios = append(goldenRatios, golden)
	}

	// Update golden ratios
	for i := len(goldenRatios) - 1; i >= 0; i-- {
		g := &goldenRatios[i]
		g.radius += goldenRatio * 5.0 * elapsed
		g.angle += goldenAngle * elapsed
		g.life -= elapsed / g.maxLife

		if g.life <= 0 || g.radius > 100 {
			goldenRatios = append(goldenRatios[:i], goldenRatios[i+1:]...)
		}
	}
}

func drawGoldenRatios(screen tcell.Screen, width, height int) {
	goldenChars := []rune{'φ', '∞', '◯', '⊙', '⊚', '⊛', '⊜', '⊝'}

	for _, golden := range goldenRatios {
		points := int(golden.radius * goldenRatio)
		if points < 6 {
			points = 6
		}
		if points > 24 {
			points = 24
		}

		for i := 0; i < points; i++ {
			angle := float64(i)*goldenAngle + golden.angle
			x := int(golden.x + golden.radius*math.Cos(angle))
			y := int(golden.y + golden.radius*math.Sin(angle))

			if x >= 0 && x < width && y >= 0 && y < height {
				intensity := golden.intensity * golden.life * (1.0 - golden.radius/100.0)

				if intensity > 0.2 {
					charIndex := int(intensity * float64(len(goldenChars)))
					if charIndex >= len(goldenChars) {
						charIndex = len(goldenChars) - 1
					}
					char := goldenChars[charIndex]

					hue := math.Mod(goldenPhase*0.05+angle/(2*math.Pi), 1.0)
					saturation := 0.8
					value := intensity
					color := HSVToRGB(hue, saturation, value)

					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func updateSacredGeometry(elapsed, peak, mathProgression float64, centerX, centerY, width, height int, rng *rand.Rand) {
	targetPatterns := int(mathProgression*4) + 2
	if targetPatterns > maxSacredGeo {
		targetPatterns = maxSacredGeo
	}

	// Add sacred geometry patterns
	for len(sacredGeometry) < targetPatterns {
		pattern := SacredGeometry{
			centerX:   centerX + rng.Intn(width/4) - width/8,
			centerY:   centerY + rng.Intn(height/4) - height/8,
			radius:    20.0 + rng.Float64()*40.0,
			angles:    make([]float64, 5+rng.Intn(8)),
			intensity: 0.6 + mathProgression*0.4,
			life:      1.0,
			pattern:   rng.Intn(4),
		}

		// Initialize angles based on pattern
		for i := range pattern.angles {
			switch pattern.pattern {
			case 0: // Pentagon (golden ratio)
				pattern.angles[i] = float64(i) * 2 * math.Pi / 5
			case 1: // Fibonacci spiral points
				pattern.angles[i] = float64(i) * goldenAngle
			case 2: // Golden rectangle
				pattern.angles[i] = float64(i) * math.Pi / 2
			case 3: // Sacred geometry
				pattern.angles[i] = float64(i) * 2 * math.Pi / goldenRatio
			}
		}

		sacredGeometry = append(sacredGeometry, pattern)
	}

	// Update patterns
	for i := 0; i < len(sacredGeometry); i++ {
		s := &sacredGeometry[i]
		s.radius += elapsed * 5.0
		s.life -= elapsed * 0.2
		for j := range s.angles {
			s.angles[j] += elapsed * goldenAngle * 0.1
		}
	}

	// Remove excess patterns
	if len(sacredGeometry) > targetPatterns {
		sacredGeometry = sacredGeometry[:targetPatterns]
	}
}

func drawSacredGeometry(screen tcell.Screen, width, height int) {
	sacredChars := []rune{'◯', '△', '▽', '◊', '⬟', '⬠', '⬡', '⟐', '⟑', '⟒'}

	for _, geo := range sacredGeometry {
		if geo.life <= 0 {
			continue
		}

		for _, angle := range geo.angles {
			x := geo.centerX + int(geo.radius*math.Cos(angle))
			y := geo.centerY + int(geo.radius*math.Sin(angle))

			if x >= 0 && x < width && y >= 0 && y < height {
				intensity := geo.intensity * geo.life * (1.0 - geo.radius/100.0)

				if intensity > 0.2 {
					charIndex := geo.pattern % len(sacredChars)
					char := sacredChars[charIndex]

					hue := math.Mod(mathPhase*0.03+angle/(2*math.Pi), 1.0)
					saturation := 0.6 + intensity*0.3
					value := intensity * 0.8
					color := HSVToRGB(hue, saturation, value)

					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func updateNumberSequences(elapsed, peak, mathProgression float64, width, height int, rng *rand.Rand) {
	// Spawn fibonacci numbers
	if len(numberSequences) < maxNumbers && rng.Float64() < mathProgression*2.0*elapsed {
		// Generate fibonacci number
		fibIndex := 1 + rng.Intn(12)
		fibNumber := 1
		if fibIndex >= 2 {
			a, b := 1, 1
			for i := 2; i < fibIndex; i++ {
				a, b = b, a+b
			}
			fibNumber = b
		}

		number := NumberSequence{
			x:         rng.Intn(width),
			y:         rng.Intn(height),
			number:    fibNumber,
			intensity: 0.7 + mathProgression*0.3,
			life:      1.0,
			maxLife:   3.0 + rng.Float64()*2.0,
			hue:       math.Mod(goldenPhase*0.08+float64(fibIndex)*0.618, 1.0),
		}
		numberSequences = append(numberSequences, number)
	}

	// Update numbers
	for i := len(numberSequences) - 1; i >= 0; i-- {
		n := &numberSequences[i]
		n.life -= elapsed / n.maxLife

		if n.life <= 0 {
			numberSequences = append(numberSequences[:i], numberSequences[i+1:]...)
		}
	}
}

func drawNumberSequences(screen tcell.Screen, width, height int) {
	for _, num := range numberSequences {
		if num.x >= 0 && num.x < width && num.y >= 0 && num.y < height {
			intensity := num.intensity * num.life

			if intensity > 0.3 {
				// Use fibonacci number as character if small enough
				var char rune
				if num.number < 10 {
					char = rune('0' + num.number)
				} else {
					mathSymbols := []rune{'φ', '∑', '∏', '∞', '√', '∆', '∇', '∂'}
					char = mathSymbols[num.number%len(mathSymbols)]
				}

				saturation := 0.8
				value := intensity * 0.9
				color := HSVToRGB(num.hue, saturation, value)

				screen.SetContent(num.x, num.y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func drawGoldenConnections(screen tcell.Screen, centerX, centerY int, radius, angle float64, fib []int, index int, maxRadius, peakScale, armPhase float64, width, height int, peak float64) {
	if index < 5 {
		return
	}

	prevRadius := math.Pow(float64(fib[index-1]), 0.618) * maxRadius * peakScale / 15.0
	prevAngle := float64(index-1)*goldenAngle + armPhase

	startX := centerX + int(prevRadius*math.Cos(prevAngle))
	startY := centerY + int(prevRadius*math.Sin(prevAngle))
	endX := centerX + int(radius*math.Cos(angle))
	endY := centerY + int(radius*math.Sin(angle))

	// Draw golden ratio connecting line
	dx := endX - startX
	dy := endY - startY
	steps := int(math.Sqrt(float64(dx*dx + dy*dy)))

	connectChars := []rune{'·', '∘', '─', '━', '═'}

	for step := 0; step <= steps; step++ {
		if steps == 0 {
			continue
		}
		t := float64(step) / float64(steps)
		x := startX + int(float64(dx)*t)
		y := startY + int(float64(dy)*t)

		if x >= 0 && x < width && y >= 0 && y < height {
			intensity := (1.0 - t*0.3) * (peak - 0.5) * 2.0

			if intensity > 0.3 {
				charIndex := int(intensity * float64(len(connectChars)))
				if charIndex >= len(connectChars) {
					charIndex = len(connectChars) - 1
				}
				char := connectChars[charIndex]

				hue := math.Mod(t*goldenRatio, 1.0)
				saturation := 0.6 + intensity*0.3
				value := intensity * 0.8
				color := HSVToRGB(hue, saturation, value)

				screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func drawMathematicalCore(screen tcell.Screen, centerX, centerY int, peak, mathProgression, basePhase float64) {
	coreRadius := 2 + int(peak*6) + int(mathProgression*4)
	if coreRadius > 10 {
		coreRadius = 10
	}

	coreChars := []rune{'∘', '○', '●', '◉', '⬢', '⬡', '★', '✦', 'φ', '∞'}

	for radius := 1; radius <= coreRadius; radius++ {
		coreIntensity := (1.0 - float64(radius-1)/float64(coreRadius)) * (0.5 + peak*0.5)

		if coreIntensity > 0.2 {
			// Golden ratio point distribution
			fibPoints := int(float64(radius) * goldenRatio * 3)
			for point := 0; point < fibPoints; point++ {
				angle := float64(point) * goldenAngle
				angle += mathPhase * 0.05 // Slow mathematical rotation

				x := centerX + int(float64(radius)*math.Cos(angle))
				y := centerY + int(float64(radius)*math.Sin(angle))

				screenWidth, screenHeight := screen.Size()
				if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
					charIndex := int(coreIntensity * float64(len(coreChars)))
					if charIndex >= len(coreChars) {
						charIndex = len(coreChars) - 1
					}
					char := coreChars[charIndex]

					hue := math.Mod(goldenPhase*0.02+float64(radius)*0.618+angle/(2*math.Pi), 1.0)
					saturation := 0.8 + mathProgression*0.2
					value := coreIntensity + peak*0.3
					color := HSVToRGB(hue, saturation, value)

					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
