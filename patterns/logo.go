package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Particle struct {
	x, y      float64
	vx, vy    float64
	life      float64
	maxLife   float64
	intensity float64
	hue       float64
	char      rune
}

type GlitchBlock struct {
	x, y        int
	width       int
	height      int
	offsetX     int
	offsetY     int
	intensity   float64
	duration    float64
	maxDuration float64
}

type Sparkle struct {
	x, y      int
	intensity float64
	life      float64
	maxLife   float64
	hue       float64
	phase     float64
}

var (
	logoGradientPhase    float64 = 0.0
	logoGradientStrength float64 = 0.0
	logoLastUpdate       time.Time

	// Particle system
	particles    []Particle
	maxParticles = 150

	// Glitch system
	glitchBlocks []GlitchBlock
	glitchTimer  float64 = 0.0

	// Sparkle system
	sparkles    []Sparkle
	maxSparkles = 50

	// Rainbow wave effects
	rainbowPhase float64 = 0.0
	pulsePhase   float64 = 0.0

	// Intensity tracking
	peakHistory []float64
	maxHistory  = 30
)

// DrawLogo creates an epic dynamic logo with particles, glitches, and rainbow effects
func DrawLogo(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	now := time.Now()
	elapsed := now.Sub(logoLastUpdate).Seconds()
	if elapsed < 1.0/60.0 { // 60 FPS limit
		return
	}
	logoLastUpdate = now

	// Track peak history for more responsive effects
	peakHistory = append(peakHistory, peak)
	if len(peakHistory) > maxHistory {
		peakHistory = peakHistory[1:]
	}

	avgPeak := 0.0
	for _, p := range peakHistory {
		avgPeak += p
	}
	avgPeak /= float64(len(peakHistory))

	logoFrames := []string{
		" __    __     __     __         __  __     ______     __  __     ______     __  __     ______     ______    ",
		"/\\ \"-./  \\   /\\ \\   /\\ \\       /\\ \\/ /    /\\  ___\\   /\\ \\_\\ \\   /\\  __ \\   /\\ \\/ /    /\\  ___\\   /\\  == \\   ",
		"\\ \\ \\-./\\ \\  \\ \\ \\  \\ \\ \\____  \\ \\  _\"-.  \\ \\___  \\  \\ \\  __ \\  \\ \\  __ \\  \\ \\  _\"-.  \\ \\  __\\   \\ \\  __<   ",
		" \\ \\_\\ \\ \\_\\  \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\  \\/\\_____\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\ ",
		"  \\/_/  \\/_/   \\/_/   \\/_____/   \\/_/\\/_/   \\/_____/   \\/_/\\/_/   \\/_/\\/_/   \\/_/\\/_/   \\/_____/   \\/_/ /_/ ",
	}

	// Update phases with audio reactivity
	speedMultiplier := 1.0 + peak*3.0
	logoGradientPhase += elapsed * speedMultiplier * 2.0
	rainbowPhase += elapsed * speedMultiplier * 1.5
	pulsePhase += elapsed * speedMultiplier * 4.0

	// Dynamic logo strength with explosive peaks
	targetStrength := 0.0
	if peak > 0.1 {
		baseStrength := math.Min(1.0, (peak-0.1)*1.8)
		explosiveBoost := 0.0
		if peak > 0.3 {
			explosiveBoost = math.Pow((peak-0.3)*2.0, 2.0) * 0.5
		}
		waveEffect := (math.Sin(logoGradientPhase*0.5) + 1.0) / 2.0
		targetStrength = (baseStrength + explosiveBoost) * (0.7 + waveEffect*0.3)
	}

	// Smooth but responsive interpolation
	smoothing := 0.88 - peak*0.3 // More responsive during peaks
	logoGradientStrength = logoGradientStrength*smoothing + targetStrength*(1.0-smoothing)

	// Update particle system
	updateParticles(elapsed, peak, width, height, rng)

	// Update glitch system
	updateGlitchSystem(elapsed, peak, rng)

	// Update sparkle system
	updateSparkles(elapsed, peak, width, height, rng)

	if logoGradientStrength < 0.02 {
		// Still show particles and sparkles even when logo is dim
		drawParticles(screen, width, height)
		drawSparkles(screen, width, height)
		return
	}

	// Fixed positioning - center of screen
	logoHeight := len(logoFrames)
	logoWidth := len(logoFrames[0])
	startY := (height - logoHeight) / 2
	startX := (width - logoWidth) / 2

	// Dynamic breathing and pulsing effects
	basePhase := GetBasePhase()
	breathe := 1.0 + math.Sin(basePhase*1.2+pulsePhase*0.3)*0.08*logoGradientStrength

	// Explosive pulse effect on beats
	beatPulse := 1.0
	if peak > 0.25 {
		beatPulse += math.Pow(peak-0.25, 2.0) * 0.4 * math.Sin(pulsePhase*8.0)
	}

	// Draw logo with enhanced effects
	for i, line := range logoFrames {
		for j, logoChar := range line {
			if logoChar == ' ' {
				continue
			}

			// Apply glitch offset if in glitch block
			finalX := startX + j
			finalY := startY + i

			for _, glitch := range glitchBlocks {
				if finalX >= glitch.x && finalX < glitch.x+glitch.width &&
					finalY >= glitch.y && finalY < glitch.y+glitch.height {
					finalX += int(float64(glitch.offsetX) * glitch.intensity)
					finalY += int(float64(glitch.offsetY) * glitch.intensity)
					break
				}
			}

			// Bounds check
			if finalX < 0 || finalX >= width || finalY < 0 || finalY >= height {
				continue
			}

			// Enhanced gradient calculations
			centerY := float64(logoHeight) / 2.0
			centerX := float64(logoWidth) / 2.0
			dy := float64(i) - centerY
			dx := float64(j) - centerX
			distanceFromCenter := math.Sqrt(dx*dx + dy*dy)
			maxDistance := math.Sqrt(centerX*centerX + centerY*centerY)

			// Multi-layered radial gradient
			radialGradient := 1.0 - math.Min(1.0, distanceFromCenter/maxDistance*1.1)
			radialGradient = math.Max(0.0, radialGradient)

			// Complex wave patterns
			waveX := math.Sin(float64(j)*0.08+rainbowPhase*1.2) * 0.4
			waveY := math.Sin(float64(i)*0.12+rainbowPhase*0.8) * 0.3
			waveCircular := math.Sin(distanceFromCenter*0.3+rainbowPhase*2.0) * 0.3
			combinedWave := (waveX + waveY + waveCircular + 3.0) / 6.0

			// Vertical energy cascade
			energyFlow := (math.Sin(float64(i)*0.25+rainbowPhase*1.5-peak*10.0) + 1.0) / 2.0

			// Audio-reactive horizontal sweep
			horizontalSweep := (math.Sin(float64(j)*0.06+rainbowPhase*2.2+peak*15.0) + 1.0) / 2.0

			// Combine all gradients with breathing
			combinedGradient := (radialGradient*0.4 + combinedWave*0.3 + energyFlow*0.2 + horizontalSweep*0.1) * breathe * beatPulse

			// Final intensity with explosive peaks
			finalIntensity := logoGradientStrength * combinedGradient

			// Skip very weak pixels
			if finalIntensity < 0.08 {
				continue
			}

			// Rainbow color cycling with audio reactivity
			baseHue := (rainbowPhase * 0.1) + (float64(j) * 0.01) + (peak * 0.3)
			baseHue = math.Mod(baseHue, 1.0)

			// Add color variation across the logo
			hueVariation := math.Sin(float64(i*j)*0.02+rainbowPhase*0.15) * 0.12
			hueVariation += math.Sin(distanceFromCenter*0.1+rainbowPhase*0.3) * 0.08
			finalHue := math.Mod(baseHue+hueVariation, 1.0)

			// Dynamic saturation and brightness
			saturation := 0.5 + peak*0.4 + finalIntensity*0.3 + math.Sin(rainbowPhase*0.7)*0.1
			saturation = math.Max(0.3, math.Min(1.0, saturation))

			value := 0.4 + finalIntensity*0.5 + peak*0.2 + math.Sin(pulsePhase*0.5)*0.1
			value = math.Max(0.2, math.Min(1.0, value))

			logoColor := HSVToRGB(finalHue, saturation, value)

			// Enhanced character morphing with more stages
			var displayChar rune
			morphLevel := finalIntensity + peak*0.2

			if morphLevel < 0.1 {
				displayChar = '·'
			} else if morphLevel < 0.18 {
				displayChar = '˙'
			} else if morphLevel < 0.28 {
				displayChar = '∘'
			} else if morphLevel < 0.38 {
				displayChar = '○'
			} else if morphLevel < 0.5 {
				displayChar = '●'
			} else if morphLevel < 0.65 {
				// Glitch characters for mid-intensity
				glitchChars := []rune{'▓', '▒', '░', '█', '▄', '▀', '■', '□'}
				displayChar = glitchChars[int(float64(len(glitchChars))*math.Mod(rainbowPhase*7.7+float64(i*j), 1.0))]
			} else if morphLevel < 0.8 {
				displayChar = rune(logoChar)
			} else {
				// Explosive characters for high intensity
				explodeChars := []rune{'★', '✦', '✧', '✯', '✪', '✫', '✬', '⟡', '◉', '◎'}
				displayChar = explodeChars[int(float64(len(explodeChars))*math.Mod(pulsePhase*5.3+float64(i*j), 1.0))]
			}

			screen.SetContent(finalX, finalY, displayChar, nil, tcell.StyleDefault.Foreground(logoColor))
		}
	}

	// Draw particle effects
	drawParticles(screen, width, height)

	// Draw sparkle effects
	drawSparkles(screen, width, height)

	// Draw glitch overlay effects
	drawGlitchOverlay(screen, width, height)
}

func updateParticles(elapsed, peak float64, width, height int, rng *rand.Rand) {
	// Spawn new particles based on audio intensity
	spawnRate := peak * 8.0 // More particles during peaks
	if len(particles) < maxParticles && rng.Float64() < spawnRate*elapsed {
		// Spawn from logo area
		logoHeight := 5
		logoWidth := 110
		startY := (height - logoHeight) / 2
		startX := (width - logoWidth) / 2

		particle := Particle{
			x:         float64(startX + rng.Intn(logoWidth)),
			y:         float64(startY + rng.Intn(logoHeight)),
			vx:        (rng.Float64() - 0.5) * 60.0 * (1.0 + peak),
			vy:        (rng.Float64() - 0.5) * 40.0 * (1.0 + peak),
			life:      1.0,
			maxLife:   1.0 + rng.Float64()*2.0,
			intensity: 0.8 + rng.Float64()*0.2,
			hue:       math.Mod(rainbowPhase*0.1+rng.Float64()*0.3, 1.0),
			char:      []rune{'*', '·', '○', '●', '✦', '✧', '▓', '░'}[rng.Intn(8)],
		}
		particles = append(particles, particle)
	}

	// Update existing particles
	for i := len(particles) - 1; i >= 0; i-- {
		p := &particles[i]
		p.x += p.vx * elapsed
		p.y += p.vy * elapsed
		p.life -= elapsed / p.maxLife

		// Gravity and air resistance
		p.vy += 20.0 * elapsed // Light gravity
		p.vx *= 0.98           // Air resistance
		p.vy *= 0.98

		// Remove dead particles
		if p.life <= 0 || p.x < 0 || p.x >= float64(width) || p.y < 0 || p.y >= float64(height) {
			particles = append(particles[:i], particles[i+1:]...)
		}
	}
}

func drawParticles(screen tcell.Screen, width, height int) {
	for _, p := range particles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			alpha := p.life * p.intensity
			if alpha > 0.05 {
				saturation := 0.7 + alpha*0.3
				value := alpha * 0.9
				color := HSVToRGB(p.hue, saturation, value)
				screen.SetContent(x, y, p.char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func updateGlitchSystem(elapsed, peak float64, rng *rand.Rand) {
	glitchTimer += elapsed

	// Trigger glitches on strong beats
	glitchThreshold := 0.4 - float64(len(glitchBlocks))*0.05
	if peak > glitchThreshold && glitchTimer > 0.1 && rng.Float64() < peak*0.7 {
		if len(glitchBlocks) < 8 {
			glitch := GlitchBlock{
				x:           rng.Intn(110),
				y:           rng.Intn(5),
				width:       3 + rng.Intn(8),
				height:      1 + rng.Intn(3),
				offsetX:     rng.Intn(7) - 3,
				offsetY:     rng.Intn(3) - 1,
				intensity:   0.3 + peak*0.7,
				duration:    0.0,
				maxDuration: 0.05 + rng.Float64()*0.15,
			}
			glitchBlocks = append(glitchBlocks, glitch)
		}
		glitchTimer = 0.0
	}

	// Update existing glitch blocks
	for i := len(glitchBlocks) - 1; i >= 0; i-- {
		g := &glitchBlocks[i]
		g.duration += elapsed
		g.intensity *= 0.95 // Fade out

		if g.duration >= g.maxDuration || g.intensity < 0.05 {
			glitchBlocks = append(glitchBlocks[:i], glitchBlocks[i+1:]...)
		}
	}
}

func updateSparkles(elapsed, peak float64, width, height int, rng *rand.Rand) {
	// Spawn sparkles around the logo area
	if len(sparkles) < maxSparkles && rng.Float64() < peak*2.0*elapsed {
		logoHeight := 5
		logoWidth := 110
		centerY := height / 2
		centerX := width / 2

		// Spawn in expanded area around logo
		margin := 20
		sparkle := Sparkle{
			x:         centerX - logoWidth/2 - margin + rng.Intn(logoWidth+margin*2),
			y:         centerY - logoHeight/2 - margin + rng.Intn(logoHeight+margin*2),
			intensity: 0.7 + rng.Float64()*0.3,
			life:      1.0,
			maxLife:   0.5 + rng.Float64()*1.5,
			hue:       math.Mod(rainbowPhase*0.1+rng.Float64()*1.0, 1.0),
			phase:     rng.Float64() * math.Pi * 2,
		}
		sparkles = append(sparkles, sparkle)
	}

	// Update existing sparkles
	for i := len(sparkles) - 1; i >= 0; i-- {
		s := &sparkles[i]
		s.life -= elapsed / s.maxLife
		s.phase += elapsed * 8.0

		if s.life <= 0 {
			sparkles = append(sparkles[:i], sparkles[i+1:]...)
		}
	}
}

func drawSparkles(screen tcell.Screen, width, height int) {
	sparkleChars := []rune{'✦', '✧', '★', '✪', '✫', '✬', '⋆', '∗', '◦', '·'}

	for _, s := range sparkles {
		if s.x >= 0 && s.x < width && s.y >= 0 && s.y < height {
			twinkle := (math.Sin(s.phase) + 1.0) / 2.0
			alpha := s.life * s.intensity * twinkle

			if alpha > 0.1 {
				saturation := 0.8 + alpha*0.2
				value := alpha
				color := HSVToRGB(s.hue, saturation, value)

				charIndex := int(float64(len(sparkleChars)) * twinkle)
				if charIndex >= len(sparkleChars) {
					charIndex = len(sparkleChars) - 1
				}

				screen.SetContent(s.x, s.y, sparkleChars[charIndex], nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func drawGlitchOverlay(screen tcell.Screen, width, height int) {
	// Additional glitch effects like random noise pixels
	for _, glitch := range glitchBlocks {
		if glitch.intensity > 0.3 {
			// Add some random noise in glitch areas
			noiseChars := []rune{'▓', '▒', '░', '█', '▄', '▀', '■', '□', '▤', '▥', '▦', '▧', '▨', '▩'}

			for dy := 0; dy < glitch.height; dy++ {
				for dx := 0; dx < glitch.width; dx++ {
					if rand.Float64() < 0.3 {
						x := glitch.x + dx + glitch.offsetX
						y := glitch.y + dy + glitch.offsetY

						if x >= 0 && x < width && y >= 0 && y < height {
							char := noiseChars[rand.Intn(len(noiseChars))]
							hue := math.Mod(rainbowPhase*0.15+rand.Float64()*0.1, 1.0)
							saturation := 0.4 + glitch.intensity*0.4
							value := glitch.intensity * 0.7
							color := HSVToRGB(hue, saturation, value)

							screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
						}
					}
				}
			}
		}
	}
}

// DrawLogoLayer draws the logo as an integrated pattern layer with all the epic effects
func DrawLogoLayer(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64, depthLayer int) {
	// Show in multiple layers with different intensities
	if depthLayer < 0 || depthLayer > 4 {
		return
	}

	// Reduce intensity for depth layers but keep effects
	originalStrength := logoGradientStrength
	depthScale := 1.0 - float64(depthLayer)*0.15
	logoGradientStrength *= depthScale * (0.5 + peak*0.3)

	DrawLogo(screen, width, height, color, char, rng, peak*depthScale)

	// Restore original strength
	logoGradientStrength = originalStrength
}
