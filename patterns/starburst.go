package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type StarburstParticle struct {
	x, y      float64
	vx, vy    float64
	life      float64
	maxLife   float64
	intensity float64
	hue       float64
	size      int
	char      rune
	trail     []Point
}

type Point struct {
	x, y float64
}

type Lightning struct {
	segments  []Point
	intensity float64
	life      float64
	maxLife   float64
	hue       float64
	thickness int
}

type Shockwave struct {
	radius    float64
	maxRadius float64
	intensity float64
	life      float64
	maxLife   float64
	hue       float64
	centerX   int
	centerY   int
}

type Spiral struct {
	angle     float64
	radius    float64
	speed     float64
	intensity float64
	hue       float64
	direction int
}

var (
	// Particle systems
	starburstParticles []StarburstParticle
	maxStarParticles   = 200

	// Lightning system
	lightningBolts []Lightning
	maxLightning   = 15

	// Shockwave system
	shockwaves    []Shockwave
	maxShockwaves = 8

	// Spiral system
	spirals    []Spiral
	maxSpirals = 12

	// Animation phases
	explosionPhase float64 = 0.0
	lightningPhase float64 = 0.0
	// spiralPhase         float64 = 0.0
	shockwavePhase      float64 = 0.0
	starburstLastUpdate time.Time

	// Peak tracking for better responsiveness
	starPeakHistory []float64
	maxStarHistory  = 20
)

// DrawStarburst creates an EPIC explosive starburst with lightning, particles, and shockwaves
func DrawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	now := time.Now()
	elapsed := now.Sub(starburstLastUpdate).Seconds()
	if elapsed < 1.0/60.0 { // 60 FPS limit
		return
	}
	starburstLastUpdate = now

	// Track peak history for explosive effects
	starPeakHistory = append(starPeakHistory, peak)
	if len(starPeakHistory) > maxStarHistory {
		starPeakHistory = starPeakHistory[1:]
	}

	// Calculate peak derivatives for explosion detection
	peakMomentum := 0.0
	if len(starPeakHistory) > 5 {
		recent := starPeakHistory[len(starPeakHistory)-3:]
		old := starPeakHistory[len(starPeakHistory)-6 : len(starPeakHistory)-3]
		recentAvg := 0.0
		oldAvg := 0.0
		for _, p := range recent {
			recentAvg += p
		}
		for _, p := range old {
			oldAvg += p
		}
		recentAvg /= float64(len(recent))
		oldAvg /= float64(len(old))
		peakMomentum = recentAvg - oldAvg
	}

	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := math.Min(float64(width), float64(height)) / 2.0

	// Update all animation phases with audio reactivity
	speedMultiplier := 1.0 + peak*4.0 + math.Max(0, peakMomentum*10.0)
	explosionPhase += elapsed * speedMultiplier * 3.0
	lightningPhase += elapsed * speedMultiplier * 8.0
	spiralPhase += elapsed * speedMultiplier * 2.0
	shockwavePhase += elapsed * speedMultiplier * 6.0

	// Update all effect systems
	updateStarburstParticles(elapsed, peak, peakMomentum, width, height, centerX, centerY, rng)
	updateLightning(elapsed, peak, peakMomentum, centerX, centerY, maxRadius, rng)
	updateShockwaves(elapsed, peak, peakMomentum, centerX, centerY, rng)
	updateSpirals(elapsed, peak, speedMultiplier, rng)

	// Draw base starburst rays with EPIC enhancements
	drawEpicRays(screen, width, height, centerX, centerY, maxRadius, peak, peakMomentum, basePhase, rng)

	// Draw all effect layers
	drawShockwaves(screen, width, height)
	drawSpirals(screen, width, height, centerX, centerY, peak)
	drawStarburstParticles(screen, width, height)
	drawLightning(screen, width, height)

	// Draw explosive center core
	drawExplosiveCore(screen, centerX, centerY, peak, peakMomentum, basePhase)

	// Draw energy rings
	drawEnergyRings(screen, centerX, centerY, maxRadius, peak, basePhase)
}

func drawEpicRays(screen tcell.Screen, width, height, centerX, centerY int, maxRadius, peak, peakMomentum, basePhase float64, rng *rand.Rand) {
	// Explosive ray count
	baseRays := 12
	bonusRays := int(peak*24) + int(math.Max(0, peakMomentum)*30)
	totalRays := baseRays + bonusRays
	if totalRays > 60 {
		totalRays = 60
	}

	// Epic ray characters
	rayChars := [][]rune{
		{'·', '∘', '○', '●', '◉', '⬢', '★', '✦', '✧', '✯', '⟡'}, // Intensity progression
		{'─', '━', '═', '▬', '■', '█'},                          // Horizontal
		{'│', '┃', '║', '▮', '█'},                               // Vertical
		{'╱', '╲', '⟋', '⟍', '⧸', '⧹'},                          // Diagonals
		{'◢', '◣', '◤', '◥', '▰', '▱'},                          // Triangular
		{'※', '✱', '✲', '✳', '✴', '✵', '✶', '✷', '✸', '✹'},      // Star variants
		{'⚡', '🌟', '💫', '⭐', '🌠'},                               // Special effects
	}

	for rayIndex := 0; rayIndex < totalRays; rayIndex++ {
		// Dynamic ray angles with chaos factor
		chaosAngle := math.Sin(explosionPhase*0.1+float64(rayIndex)*0.3) * (peakMomentum * 2.0)
		baseAngle := float64(rayIndex) * 2 * math.Pi / float64(totalRays)
		rayAngle := baseAngle + chaosAngle

		// Rotation effects
		rotationSpeed := 0.2 + peak*0.8
		if rayIndex%2 == 0 {
			rotationSpeed = -rotationSpeed // Counter-rotate alternate rays
		}
		finalAngle := rayAngle + basePhase*rotationSpeed

		// Dynamic ray length with explosive bursts
		rayPersonality := float64(rayIndex) * 1.618 // Golden ratio
		baseLength := maxRadius * (0.4 + peak*0.6)

		// Explosive length boost
		explosiveBurst := 0.0
		if peakMomentum > 0.1 {
			burstPhase := explosionPhase*3.0 + rayPersonality
			explosiveBurst = math.Max(0, math.Sin(burstPhase)) * peakMomentum * 2.0
		}

		lengthVariation := 0.7 + 0.6*math.Sin(spiralPhase*0.2+rayPersonality)
		rayLength := baseLength * lengthVariation * (1.0 + explosiveBurst)

		// Ray width based on intensity
		rayWidth := 1
		if peak > 0.3 {
			rayWidth = 2
		}
		if peak > 0.6 {
			rayWidth = 3
		}
		if peak > 0.8 && peakMomentum > 0.2 {
			rayWidth = 4
		}

		// Draw ray with multiple segments
		raySteps := int(rayLength / 1.2)
		if raySteps < 5 {
			raySteps = 5
		}

		for step := 1; step <= raySteps; step++ {
			stepRatio := float64(step) / float64(raySteps)
			currentRadius := stepRatio * rayLength

			if currentRadius < 3 {
				continue
			}

			// Draw ray with width
			for w := -rayWidth / 2; w <= rayWidth/2; w++ {
				for h := -rayWidth / 2; h <= rayWidth/2; h++ {
					x := centerX + int(currentRadius*math.Cos(finalAngle)) + w
					y := centerY + int(currentRadius*math.Sin(finalAngle)) + h

					if x >= 0 && x < width && y >= 0 && y < height {
						// Distance-based intensity with explosive boosts
						intensity := (1.0 - stepRatio*0.7) * (0.3 + peak*0.7)

						// Tip explosion effect
						if step == raySteps && intensity > 0.4 {
							intensity += explosiveBurst * 0.5
						}

						// Pulsing effect along rays
						pulsePhase := lightningPhase*2.0 + currentRadius*0.1
						pulseFactor := 1.0 + math.Sin(pulsePhase)*0.3*peak
						finalIntensity := intensity * pulseFactor

						if finalIntensity < 0.1 {
							continue
						}

						// Epic character selection
						var finalChar rune
						morphLevel := finalIntensity + peakMomentum*0.5

						if morphLevel < 0.15 {
							finalChar = rayChars[0][0] // ·
						} else if morphLevel < 0.25 {
							finalChar = rayChars[0][1] // ∘
						} else if morphLevel < 0.4 {
							finalChar = rayChars[0][2] // ○
						} else if morphLevel < 0.55 {
							finalChar = rayChars[0][3] // ●
						} else if morphLevel < 0.7 {
							// Direction-based characters
							normalizedAngle := math.Mod(finalAngle, 2*math.Pi)
							if normalizedAngle > -math.Pi/8 && normalizedAngle <= math.Pi/8 {
								finalChar = rayChars[1][int(morphLevel*float64(len(rayChars[1])))]
							} else if normalizedAngle > 3*math.Pi/8 && normalizedAngle <= 5*math.Pi/8 {
								finalChar = rayChars[2][int(morphLevel*float64(len(rayChars[2])))]
							} else {
								finalChar = rayChars[3][int(morphLevel*float64(len(rayChars[3])))]
							}
						} else if morphLevel < 0.85 {
							// Star variants
							starIndex := int(morphLevel * float64(len(rayChars[5])))
							if starIndex >= len(rayChars[5]) {
								starIndex = len(rayChars[5]) - 1
							}
							finalChar = rayChars[5][starIndex]
						} else {
							// EXPLOSIVE special effects
							if step == raySteps {
								finalChar = rayChars[6][rng.Intn(len(rayChars[6]))]
							} else {
								finalChar = rayChars[0][len(rayChars[0])-1] // ⟡
							}
						}

						// Epic color system
						baseHue := float64(rayIndex)*0.05 + explosionPhase*0.1
						hueVariation := math.Sin(currentRadius*0.05+lightningPhase*0.3) * 0.15
						explosiveHueShift := peakMomentum * 0.3 * math.Sin(explosionPhase*5.0)
						finalHue := math.Mod(baseHue+hueVariation+explosiveHueShift, 1.0)

						saturation := 0.6 + peak*0.35 + finalIntensity*0.2
						saturation = math.Max(0.3, math.Min(1.0, saturation))

						value := 0.3 + finalIntensity*0.6 + peak*0.2 + explosiveBurst*0.3
						value = math.Max(0.1, math.Min(1.0, value))

						rayColor := HSVToRGB(finalHue, saturation, value)
						screen.SetContent(x, y, finalChar, nil, tcell.StyleDefault.Foreground(rayColor))
					}
				}
			}
		}
	}
}

func updateStarburstParticles(elapsed, peak, peakMomentum float64, width, height, centerX, centerY int, rng *rand.Rand) {
	// Spawn particles from ray tips and explosive events
	spawnRate := peak*12.0 + peakMomentum*20.0
	if len(starburstParticles) < maxStarParticles && rng.Float64() < spawnRate*elapsed {
		// Random spawn angle
		angle := rng.Float64() * 2 * math.Pi
		spawnRadius := 20.0 + rng.Float64()*60.0

		particle := StarburstParticle{
			x:         float64(centerX) + spawnRadius*math.Cos(angle),
			y:         float64(centerY) + spawnRadius*math.Sin(angle),
			vx:        math.Cos(angle) * (40.0 + peak*80.0 + peakMomentum*100.0) * (0.5 + rng.Float64()),
			vy:        math.Sin(angle) * (40.0 + peak*80.0 + peakMomentum*100.0) * (0.5 + rng.Float64()),
			life:      1.0,
			maxLife:   0.8 + rng.Float64()*2.2,
			intensity: 0.7 + rng.Float64()*0.3 + peak*0.5,
			hue:       math.Mod(explosionPhase*0.1+rng.Float64()*0.4, 1.0),
			size:      1 + rng.Intn(3) + int(peak*2),
			char:      []rune{'·', '∘', '○', '●', '★', '✦', '✧', '⟡', '◉'}[rng.Intn(9)],
			trail:     make([]Point, 0, 8),
		}
		starburstParticles = append(starburstParticles, particle)
	}

	// Update existing particles
	for i := len(starburstParticles) - 1; i >= 0; i-- {
		p := &starburstParticles[i]

		// Add current position to trail
		p.trail = append(p.trail, Point{p.x, p.y})
		if len(p.trail) > 8 {
			p.trail = p.trail[1:]
		}

		// Physics
		p.x += p.vx * elapsed
		p.y += p.vy * elapsed
		p.life -= elapsed / p.maxLife

		// Air resistance and gravity
		p.vx *= 0.97
		p.vy *= 0.97
		p.vy += 15.0 * elapsed // Light gravity

		// Remove dead or off-screen particles
		if p.life <= 0 || p.x < -50 || p.x >= float64(width+50) || p.y < -50 || p.y >= float64(height+50) {
			starburstParticles = append(starburstParticles[:i], starburstParticles[i+1:]...)
		}
	}
}

func drawStarburstParticles(screen tcell.Screen, width, height int) {
	for _, p := range starburstParticles {
		// Draw particle trail
		for j, trailPoint := range p.trail {
			x, y := int(trailPoint.x), int(trailPoint.y)
			if x >= 0 && x < width && y >= 0 && y < height {
				trailIntensity := float64(j) / float64(len(p.trail)) * p.life * 0.5
				if trailIntensity > 0.1 {
					saturation := 0.4 + trailIntensity*0.4
					value := trailIntensity * 0.8
					color := HSVToRGB(p.hue, saturation, value)
					screen.SetContent(x, y, '·', nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}

		// Draw main particle
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			alpha := p.life * p.intensity
			if alpha > 0.1 {
				saturation := 0.7 + alpha*0.3
				value := alpha * 0.9
				color := HSVToRGB(p.hue, saturation, value)
				screen.SetContent(x, y, p.char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func updateLightning(elapsed, peak, peakMomentum float64, centerX, centerY int, maxRadius float64, rng *rand.Rand) {
	// Spawn lightning on strong beats
	if len(lightningBolts) < maxLightning && (peak > 0.4 || peakMomentum > 0.15) && rng.Float64() < (peak+peakMomentum)*2.0*elapsed {
		// Create lightning bolt from center to random point
		angle := rng.Float64() * 2 * math.Pi
		targetRadius := maxRadius * (0.6 + rng.Float64()*0.4)
		targetX := float64(centerX) + targetRadius*math.Cos(angle)
		targetY := float64(centerY) + targetRadius*math.Sin(angle)

		// Generate zigzag path
		segments := make([]Point, 0, 20)
		segments = append(segments, Point{float64(centerX), float64(centerY)})

		numSegments := 8 + rng.Intn(12)
		for i := 1; i <= numSegments; i++ {
			ratio := float64(i) / float64(numSegments)

			// Base position along line
			baseX := float64(centerX) + (targetX-float64(centerX))*ratio
			baseY := float64(centerY) + (targetY-float64(centerY))*ratio

			// Add random zigzag
			perpX := -(targetY - float64(centerY))
			perpY := targetX - float64(centerX)
			perpLen := math.Sqrt(perpX*perpX + perpY*perpY)
			if perpLen > 0 {
				perpX /= perpLen
				perpY /= perpLen
			}

			zigzag := (rng.Float64() - 0.5) * 30.0 * (1.0 - ratio*0.3)
			finalX := baseX + perpX*zigzag
			finalY := baseY + perpY*zigzag

			segments = append(segments, Point{finalX, finalY})
		}

		lightning := Lightning{
			segments:  segments,
			intensity: 0.8 + peak*0.2 + peakMomentum*0.5,
			life:      1.0,
			maxLife:   0.1 + rng.Float64()*0.2,
			hue:       math.Mod(lightningPhase*0.2+rng.Float64()*0.1, 1.0),
			thickness: 1 + int(peak*2) + int(peakMomentum*3),
		}
		lightningBolts = append(lightningBolts, lightning)
	}

	// Update existing lightning
	for i := len(lightningBolts) - 1; i >= 0; i-- {
		l := &lightningBolts[i]
		l.life -= elapsed / l.maxLife
		l.intensity *= 0.95 // Fade out

		if l.life <= 0 || l.intensity < 0.05 {
			lightningBolts = append(lightningBolts[:i], lightningBolts[i+1:]...)
		}
	}
}

func drawLightning(screen tcell.Screen, width, height int) {
	lightningChars := []rune{'│', '┃', '║', '█', '▌', '▐', '▄', '▀', '⚡'}

	for _, bolt := range lightningBolts {
		for i := 0; i < len(bolt.segments)-1; i++ {
			p1 := bolt.segments[i]
			p2 := bolt.segments[i+1]

			// Draw line between segments
			dx := p2.x - p1.x
			dy := p2.y - p1.y
			dist := math.Sqrt(dx*dx + dy*dy)
			steps := int(dist)

			for step := 0; step <= steps; step++ {
				t := float64(step) / float64(steps)
				x := int(p1.x + dx*t)
				y := int(p1.y + dy*t)

				// Draw with thickness
				for w := -bolt.thickness / 2; w <= bolt.thickness/2; w++ {
					for h := -bolt.thickness / 2; h <= bolt.thickness/2; h++ {
						finalX := x + w
						finalY := y + h

						if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {
							intensity := bolt.intensity * bolt.life * (1.0 - t*0.1)

							charIndex := int(intensity * float64(len(lightningChars)))
							if charIndex >= len(lightningChars) {
								charIndex = len(lightningChars) - 1
							}
							char := lightningChars[charIndex]

							saturation := 0.9
							value := intensity
							color := HSVToRGB(bolt.hue, saturation, value)

							screen.SetContent(finalX, finalY, char, nil, tcell.StyleDefault.Foreground(color))
						}
					}
				}
			}
		}
	}
}

func updateShockwaves(elapsed, peak, peakMomentum float64, centerX, centerY int, rng *rand.Rand) {
	// Create shockwaves on explosive beats
	if len(shockwaves) < maxShockwaves && peakMomentum > 0.2 && rng.Float64() < peakMomentum*4.0*elapsed {
		shockwave := Shockwave{
			radius:    5.0,
			maxRadius: 100.0 + peak*150.0,
			intensity: 0.8 + peakMomentum*0.2,
			life:      1.0,
			maxLife:   1.0 + rng.Float64()*1.5,
			hue:       math.Mod(shockwavePhase*0.1+rng.Float64()*0.3, 1.0),
			centerX:   centerX,
			centerY:   centerY,
		}
		shockwaves = append(shockwaves, shockwave)
	}

	// Update existing shockwaves
	for i := len(shockwaves) - 1; i >= 0; i-- {
		s := &shockwaves[i]
		s.radius += (s.maxRadius / s.maxLife) * elapsed
		s.life -= elapsed / s.maxLife

		if s.life <= 0 || s.radius > s.maxRadius {
			shockwaves = append(shockwaves[:i], shockwaves[i+1:]...)
		}
	}
}

func drawShockwaves(screen tcell.Screen, width, height int) {
	waveChars := []rune{'∘', '○', '◦', '●', '▫', '▪', '■', '█'}

	for _, wave := range shockwaves {
		points := int(wave.radius * 4) // More points for larger waves
		if points < 12 {
			points = 12
		}

		for i := 0; i < points; i++ {
			angle := float64(i) * 2 * math.Pi / float64(points)
			x := wave.centerX + int(wave.radius*math.Cos(angle))
			y := wave.centerY + int(wave.radius*math.Sin(angle))

			if x >= 0 && x < width && y >= 0 && y < height {
				intensity := wave.intensity * wave.life * (1.0 - wave.radius/wave.maxRadius)

				if intensity > 0.1 {
					charIndex := int(intensity * float64(len(waveChars)))
					if charIndex >= len(waveChars) {
						charIndex = len(waveChars) - 1
					}
					char := waveChars[charIndex]

					saturation := 0.8
					value := intensity
					color := HSVToRGB(wave.hue, saturation, value)

					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func updateSpirals(elapsed, peak, speedMultiplier float64, rng *rand.Rand) {
	// Maintain active spirals based on audio intensity
	targetSpirals := int(peak*8) + 2
	if targetSpirals > maxSpirals {
		targetSpirals = maxSpirals
	}

	// Add spirals if needed
	for len(spirals) < targetSpirals {
		spiral := Spiral{
			angle:     rng.Float64() * 2 * math.Pi,
			radius:    10.0 + rng.Float64()*20.0,
			speed:     0.5 + rng.Float64()*2.0,
			intensity: 0.6 + rng.Float64()*0.4,
			hue:       math.Mod(spiralPhase*0.05+rng.Float64()*1.0, 1.0),
			direction: []int{-1, 1}[rng.Intn(2)],
		}
		spirals = append(spirals, spiral)
	}

	// Update spirals
	for i := 0; i < len(spirals); i++ {
		s := &spirals[i]
		s.angle += float64(s.direction) * s.speed * speedMultiplier * elapsed
		s.radius += s.speed * 10.0 * elapsed
		s.intensity *= 0.995

		// Reset spiral when it gets too far or weak
		if s.radius > 150 || s.intensity < 0.2 {
			s.radius = 10.0 + rng.Float64()*20.0
			s.intensity = 0.6 + rng.Float64()*0.4 + peak*0.3
			s.hue = math.Mod(spiralPhase*0.05+rng.Float64()*1.0, 1.0)
		}
	}

	// Remove excess spirals
	if len(spirals) > targetSpirals {
		spirals = spirals[:targetSpirals]
	}
}

func drawSpirals(screen tcell.Screen, width, height, centerX, centerY int, peak float64) {
	spiralChars := []rune{'·', '∘', '○', '◦', '●', '✧', '✦', '★'}

	for _, spiral := range spirals {
		// Draw multiple arms of the spiral
		arms := 3 + int(peak*2)
		for arm := 0; arm < arms; arm++ {
			armAngle := spiral.angle + float64(arm)*2*math.Pi/float64(arms)
			x := centerX + int(spiral.radius*math.Cos(armAngle))
			y := centerY + int(spiral.radius*math.Sin(armAngle))

			if x >= 0 && x < width && y >= 0 && y < height {
				charIndex := int(spiral.intensity * float64(len(spiralChars)))
				if charIndex >= len(spiralChars) {
					charIndex = len(spiralChars) - 1
				}
				char := spiralChars[charIndex]

				saturation := 0.7 + spiral.intensity*0.3
				value := spiral.intensity * 0.8
				color := HSVToRGB(spiral.hue, saturation, value)

				screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func drawExplosiveCore(screen tcell.Screen, centerX, centerY int, peak, peakMomentum, basePhase float64) {
	// Dynamic core size based on audio
	coreSize := 2 + int(peak*8) + int(peakMomentum*10)
	if coreSize > 12 {
		coreSize = 12
	}

	coreChars := []rune{'·', '∘', '○', '◦', '●', '◉', '⬢', '⬡', '★', '✦', '✧', '✯', '⟡', '◈', '◊'}

	for radius := 0; radius <= coreSize; radius++ {
		if radius == 0 {
			// Center point - most explosive
			intensity := 0.8 + peak*0.2 + peakMomentum*0.5
			charIndex := int(intensity * float64(len(coreChars)))
			if charIndex >= len(coreChars) {
				charIndex = len(coreChars) - 1
			}
			char := coreChars[charIndex]

			hue := math.Mod(basePhase*0.1+explosionPhase*0.3, 1.0)
			saturation := 0.9
			value := 0.7 + intensity*0.3

			color := HSVToRGB(hue, saturation, value)
			screen.SetContent(centerX, centerY, char, nil, tcell.StyleDefault.Foreground(color))
		} else {
			// Core rings with explosive effects
			ringIntensity := (1.0 - float64(radius)/float64(coreSize)) * (0.6 + peak*0.4)
			if ringIntensity > 0.2 {
				// Pulsing ring effect
				pulseEffect := 1.0 + math.Sin(explosionPhase*4.0+float64(radius)*0.5)*0.3*peak
				finalIntensity := ringIntensity * pulseEffect

				points := radius * 8 // More points for larger rings
				for i := 0; i < points; i++ {
					angle := float64(i) * 2 * math.Pi / float64(points)
					// Add ring wobble
					wobble := math.Sin(lightningPhase*2.0+angle*3.0) * 0.5 * peak
					finalRadius := float64(radius) + wobble

					x := centerX + int(finalRadius*math.Cos(angle))
					y := centerY + int(finalRadius*math.Sin(angle))

					screenWidth, screenHeight := screen.Size()
					if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
						var ringChar rune
						morphLevel := finalIntensity + peakMomentum*0.3

						if morphLevel < 0.3 {
							ringChar = '∘'
						} else if morphLevel < 0.5 {
							ringChar = '○'
						} else if morphLevel < 0.7 {
							ringChar = '●'
						} else if morphLevel < 0.85 {
							ringChar = '◉'
						} else {
							ringChar = '★'
						}

						ringHue := math.Mod(basePhase*0.05+float64(radius)*0.1+explosionPhase*0.2, 1.0)
						ringColor := HSVToRGB(ringHue, 0.8, finalIntensity)
						screen.SetContent(x, y, ringChar, nil, tcell.StyleDefault.Foreground(ringColor))
					}
				}
			}
		}
	}
}
func drawEnergyRings(screen tcell.Screen, centerX, centerY int, maxRadius, peak, basePhase float64) {
	// Draw multiple energy rings at different radii
	numRings := 3 + int(peak*4)
	if numRings > 8 {
		numRings = 8
	}

	ringChars := []rune{'·', '∘', '○', '◦', '●', '◉', '⬢', '★', '✦', '✧'}

	for ring := 1; ring <= numRings; ring++ {
		ringRadius := (float64(ring) / float64(numRings)) * maxRadius * (0.7 + peak*0.3)
		ringIntensity := (1.0 - float64(ring)/float64(numRings)) * (0.4 + peak*0.6)

		// Skip weak rings
		if ringIntensity < 0.15 {
			continue
		}

		// Ring animation phase
		ringPhase := basePhase*0.8 + float64(ring)*0.3
		energyPulse := 1.0 + math.Sin(ringPhase*3.0)*0.4*peak

		finalIntensity := ringIntensity * energyPulse

		// Number of points on ring
		points := int(ringRadius * 1.5)
		if points < 16 {
			points = 16
		}
		if points > 80 {
			points = 80
		}

		for i := 0; i < points; i++ {
			angle := float64(i) * 2 * math.Pi / float64(points)

			// Add energy fluctuations
			fluctuation := math.Sin(angle*5.0+ringPhase*2.0) * 2.0 * peak
			actualRadius := ringRadius + fluctuation

			x := centerX + int(actualRadius*math.Cos(angle))
			y := centerY + int(actualRadius*math.Sin(angle))

			// Bounds check
			screenWidth, screenHeight := screen.Size()
			if x >= 0 && x < screenWidth && y >= 0 && y < screenHeight {
				// Sparkling effect - not all points are always visible
				sparkleChance := 0.6 + peak*0.3 + math.Sin(ringPhase*4.0+angle*2.0)*0.2
				if sparkleChance > 0.5 {
					// Character selection based on intensity
					charIndex := int(finalIntensity * float64(len(ringChars)))
					if charIndex >= len(ringChars) {
						charIndex = len(ringChars) - 1
					}
					char := ringChars[charIndex]

					// Dynamic coloring
					baseHue := float64(ring)*0.15 + explosionPhase*0.1
					hueShift := math.Sin(angle*2.0+ringPhase*1.5) * 0.1
					finalHue := math.Mod(baseHue+hueShift, 1.0)

					saturation := 0.6 + finalIntensity*0.3 + peak*0.1
					value := finalIntensity*0.8 + peak*0.2

					color := HSVToRGB(finalHue, saturation, value)
					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
