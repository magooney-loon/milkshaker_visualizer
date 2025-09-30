package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

type WaveParticle struct {
	x, y      float64
	vx, vy    float64
	life      float64
	maxLife   float64
	intensity float64
	hue       float64
	size      float64
	char      rune
}

type Ripple struct {
	x, y      float64
	radius    float64
	maxRadius float64
	intensity float64
	life      float64
	maxLife   float64
	hue       float64
	frequency float64
}

type FlowField struct {
	x, y      float64
	angle     float64
	magnitude float64
	life      float64
}

var (
	// Minimalist particle system
	waveParticles    []WaveParticle
	maxWaveParticles = 10 // Much fewer particles for clean wireframe

	// Gentle ripple system
	ripples    []Ripple
	maxRipples = 4

	// Flow field for organic movement
	flowField []FlowField

	// Animation phases
	wavePhase      float64 = 0.0
	liquidPhase    float64 = 0.0
	ripplePhase    float64 = 0.0
	waveLastUpdate time.Time

	// Peak tracking
	wavePeakHistory []float64
	maxWaveHistory  = 9
)

// DrawWave creates a minimalistic yet epic flowing liquid wave experience
func DrawWave(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	now := time.Now()
	elapsed := now.Sub(waveLastUpdate).Seconds()
	if elapsed < 1.0/520.0 { // 520 FPS limit
		return
	}
	waveLastUpdate = now

	// Track peak history for smooth responsiveness
	wavePeakHistory = append(wavePeakHistory, peak)
	if len(wavePeakHistory) > maxWaveHistory {
		wavePeakHistory = wavePeakHistory[1:]
	}

	// Calculate smooth peak average
	avgPeak := 0.0
	for _, p := range wavePeakHistory {
		avgPeak += p
	}
	avgPeak /= float64(len(wavePeakHistory))

	// Update phases with slow, meditative audio reactivity
	speedMultiplier := 0.3 + avgPeak*0.8 + peak*0.4
	wavePhase += elapsed * speedMultiplier * 0.6
	liquidPhase += elapsed * speedMultiplier * 0.3
	ripplePhase += elapsed * speedMultiplier * 0.9

	// Update systems
	updateWaveParticles(elapsed, peak, avgPeak, width, height, rng)
	updateRipples(elapsed, peak, avgPeak, width, height, rng)
	updateFlowField(elapsed, peak, width, height)

	// Draw main liquid waves
	drawLiquidWaves(screen, width, height, peak, avgPeak, rng)

	// Draw flowing particles
	drawWaveParticles(screen, width, height)

	// Draw gentle ripples
	drawRipples(screen, width, height)

	// Draw subtle flow field effects
	drawFlowEffects(screen, width, height, peak)
}

func drawLiquidWaves(screen tcell.Screen, width, height int, peak, avgPeak float64, rng *rand.Rand) {
	basePhase := GetBasePhase()

	// Clean wireframe character set for clear wave lines
	waveChars := []rune{'·', '-', '─', '━', '═', '~', '≈', '∿'}

	// Fewer, smoother waves (2-4 based on audio)
	numWaves := 2 + int(avgPeak*2)
	if numWaves > 4 {
		numWaves = 4
	}

	for waveIndex := 0; waveIndex < numWaves; waveIndex++ {
		// Depth-based wave properties for parallax effect
		depthLayer := float64(waveIndex) / float64(numWaves)

		waveSpeed := (0.15 + float64(waveIndex)*0.1 + avgPeak*0.2) * (1.0 - depthLayer*0.4)
		amplitude := 2.0 + float64(waveIndex)*1.0 + peak*3.0
		frequency := 0.05 + float64(waveIndex)*0.02 // Slower frequency changes
		verticalOffset := height/2 + int(float64(waveIndex-numWaves/2)*1)

		// Phase offset for smooth layered depth
		phaseOffset := float64(waveIndex) * math.Pi / 3.5

		// Draw smooth depth-layered wave across screen
		for x := 0; x < width; x++ {
			waveX := float64(x) * frequency
			t := basePhase*waveSpeed + phaseOffset

			// Primary liquid wave with depth-based smoothing
			primaryY := amplitude * math.Sin(waveX+t) * (0.7 + depthLayer*0.3)

			// Subtle harmonics with depth variation
			harmonic1 := amplitude * (0.15 - depthLayer*0.05) * math.Sin(waveX*1.618+t*0.6)
			harmonic2 := amplitude * (0.08 - depthLayer*0.03) * math.Sin(waveX*0.618+t*0.9)

			// Gentle liquid distortion that moves slower in deeper layers
			liquidDistort := amplitude * 0.06 * math.Sin(waveX*0.2+liquidPhase*(0.8-depthLayer*0.3))

			totalY := primaryY + harmonic1 + harmonic2 + liquidDistort
			finalY := verticalOffset + int(totalY)

			// Draw multiple points for smoother wave thickness
			thickness := 1 + int(peak*1)
			for dy := -thickness; dy <= thickness; dy++ {
				drawY := finalY + dy
				if drawY >= 0 && drawY < height {
					// Smooth intensity calculation
					distanceFromCore := math.Abs(float64(dy))
					coreIntensity := 1.0 - (distanceFromCore / float64(thickness+1))

					// Wave amplitude intensity
					amplitudeRatio := math.Abs(totalY) / amplitude
					waveIntensity := (1.0 - amplitudeRatio*0.3) * (0.3 + avgPeak*0.7)

					// Combine for final intensity
					totalIntensity := coreIntensity * waveIntensity

					// Minimal sparkles to keep clean wireframe look
					if amplitudeRatio > 0.8 && rng.Float64() < 0.05*peak {
						totalIntensity += 0.2
					}

					if totalIntensity > 0.1 {
						// Clean wireframe character progression
						var waveChar rune
						morphLevel := totalIntensity + peak*0.2

						if morphLevel < 0.2 {
							waveChar = waveChars[0] // ·
						} else if morphLevel < 0.35 {
							waveChar = waveChars[1] // -
						} else if morphLevel < 0.5 {
							waveChar = waveChars[2] // ─
						} else if morphLevel < 0.65 {
							waveChar = waveChars[3] // ━
						} else if morphLevel < 0.8 {
							waveChar = waveChars[4] // ═
						} else if morphLevel < 0.9 {
							waveChar = waveChars[5] // ~
						} else {
							waveChar = waveChars[6] // ≈
						}

						// Depth-based liquid color flow
						baseHue := 0.5 + float64(waveIndex)*0.08 + liquidPhase*0.03 // Slower color changes
						hueFlow := math.Sin(waveX*0.2+t*0.3) * 0.06                 // Gentler flow
						depthHue := depthLayer * 0.05                               // Deeper layers slightly different hue
						finalHue := math.Mod(baseHue+hueFlow+depthHue+peak*0.08, 1.0)

						saturation := (0.3 + avgPeak*0.25 + totalIntensity*0.15) * (0.8 + depthLayer*0.2)
						saturation = math.Max(0.15, math.Min(0.7, saturation))

						value := (0.25 + totalIntensity*0.4 + peak*0.15) * (1.0 - depthLayer*0.2)
						value = math.Max(0.08, math.Min(0.8, value))

						waveColor := HSVToRGB(finalHue, saturation, value)
						screen.SetContent(x, drawY, waveChar, nil, tcell.StyleDefault.Foreground(waveColor))
					}
				}
			}

			// Subtle vertical flow lines at moderate peaks
			if peak > 0.6 && waveIndex == 0 && x%16 == 0 {
				drawVerticalFlow(screen, x, finalY, height, amplitude*0.3, peak, waveX, t)
			}
		}
	}
}

func drawVerticalFlow(screen tcell.Screen, x, centerY, height int, flowHeight, peak, waveX, t float64) {
	flowChars := []rune{'│', '┆', '┊', '︙'}

	startY := centerY - int(flowHeight/2)
	endY := centerY + int(flowHeight/2)

	if startY < 0 {
		startY = 0
	}
	if endY >= height {
		endY = height - 1
	}

	flowIntensity := (peak - 0.4) * 2.0
	if flowIntensity > 0.3 {
		// Flowing animation
		flowOffset := math.Sin(t*2.0+waveX*0.5) * 2.0

		for y := startY; y <= endY; y++ {
			adjustedY := y + int(flowOffset)
			if adjustedY >= 0 && adjustedY < height {
				distFromCenter := math.Abs(float64(y - centerY))
				lineIntensity := flowIntensity * (1.0 - distFromCenter/flowHeight)

				if lineIntensity > 0.2 {
					charIndex := int(lineIntensity * float64(len(flowChars)))
					if charIndex >= len(flowChars) {
						charIndex = len(flowChars) - 1
					}
					char := flowChars[charIndex]

					hue := math.Mod(0.55+liquidPhase*0.03, 1.0)
					saturation := 0.3 + lineIntensity*0.3
					value := lineIntensity * 0.7

					color := HSVToRGB(hue, saturation, value)
					screen.SetContent(x, adjustedY, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func updateWaveParticles(elapsed, peak, avgPeak float64, width, height int, rng *rand.Rand) {
	// Minimal particles to reduce visual noise
	spawnRate := avgPeak * 0.5
	if len(waveParticles) < maxWaveParticles && rng.Float64() < spawnRate*elapsed {
		// Spawn from wave areas with depth variation
		spawnX := rng.Float64() * float64(width)
		spawnY := float64(height/2) + (rng.Float64()-0.5)*float64(height/8)
		depthFactor := rng.Float64() // Random depth for parallax

		particle := WaveParticle{
			x:         spawnX,
			y:         spawnY,
			vx:        (rng.Float64() - 0.5) * 4.0 * (1.0 + peak*0.5) * (1.0 - depthFactor*0.3),
			vy:        (rng.Float64() - 0.5) * 2.0 * (1.0 + avgPeak*0.5) * (1.0 - depthFactor*0.3),
			life:      1.0,
			maxLife:   3.0 + rng.Float64()*4.0, // Longer life for slower movement
			intensity: (0.4 + rng.Float64()*0.2 + avgPeak*0.15) * (0.6 + depthFactor*0.4),
			hue:       math.Mod(0.5+liquidPhase*0.03+rng.Float64()*0.15, 1.0),
			size:      (1.0 + rng.Float64()*1.5) * (1.0 - depthFactor*0.3),
			char:      []rune{'·', '∘', '○', '●', '◉'}[rng.Intn(5)],
		}
		waveParticles = append(waveParticles, particle)
	}

	// Update particles with liquid physics
	for i := len(waveParticles) - 1; i >= 0; i-- {
		p := &waveParticles[i]

		// Liquid flow physics
		p.x += p.vx * elapsed
		p.y += p.vy * elapsed
		p.life -= elapsed / p.maxLife

		// Very gentle wave-following behavior for meditative flow
		waveInfluence := math.Sin(p.x*0.08+wavePhase*0.7) * 1.0 * elapsed
		p.vy += waveInfluence

		// Higher liquid viscosity for slower, smoother movement
		p.vx *= 0.96
		p.vy *= 0.93

		// Remove dead particles
		if p.life <= 0 || p.x < -10 || p.x >= float64(width+10) || p.y < -10 || p.y >= float64(height+10) {
			waveParticles = append(waveParticles[:i], waveParticles[i+1:]...)
		}
	}
}

func drawWaveParticles(screen tcell.Screen, width, height int) {
	for _, p := range waveParticles {
		x, y := int(p.x), int(p.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			alpha := p.life * p.intensity
			if alpha > 0.1 {
				saturation := 0.4 + alpha*0.4
				value := alpha * 0.8
				color := HSVToRGB(p.hue, saturation, value)

				screen.SetContent(x, y, p.char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

func updateRipples(elapsed, peak, avgPeak float64, width, height int, rng *rand.Rand) {
	// Create minimal ripples to keep focus on wave lines
	if len(ripples) < maxRipples && rng.Float64() < peak*0.3*elapsed {
		ripple := Ripple{
			x:         rng.Float64() * float64(width),
			y:         float64(height/2) + (rng.Float64()-0.5)*float64(height/8),
			radius:    1.0,
			maxRadius: 12.0 + peak*15.0,
			intensity: 0.4 + peak*0.3,
			life:      1.0,
			maxLife:   2.5 + rng.Float64()*3.5, // Longer lived ripples
			hue:       math.Mod(0.52+ripplePhase*0.05+rng.Float64()*0.12, 1.0),
			frequency: 0.8 + rng.Float64()*1.5, // Slower frequency
		}
		ripples = append(ripples, ripple)
	}

	// Update ripples with much slower expansion
	for i := len(ripples) - 1; i >= 0; i-- {
		r := &ripples[i]
		r.radius += (r.maxRadius / r.maxLife) * elapsed * 0.4 // Much slower expansion
		r.life -= elapsed / r.maxLife

		if r.life <= 0 || r.radius > r.maxRadius {
			ripples = append(ripples[:i], ripples[i+1:]...)
		}
	}
}

func drawRipples(screen tcell.Screen, width, height int) {
	rippleChars := []rune{'∘', '○', '◦', '●'}

	for _, ripple := range ripples {
		points := int(ripple.radius * 3)
		if points < 8 {
			points = 8
		}
		if points > 24 {
			points = 24
		}

		for i := 0; i < points; i++ {
			angle := float64(i) * 2 * math.Pi / float64(points)

			// Gentle ripple distortion for smooth meditative effect
			distortion := math.Sin(angle*ripple.frequency+ripplePhase*1.2) * 1.0
			actualRadius := ripple.radius + distortion

			x := int(ripple.x + actualRadius*math.Cos(angle))
			y := int(ripple.y + actualRadius*math.Sin(angle))

			if x >= 0 && x < width && y >= 0 && y < height {
				intensity := ripple.intensity * ripple.life * (1.0 - ripple.radius/ripple.maxRadius)

				if intensity > 0.15 {
					charIndex := int(intensity * float64(len(rippleChars)))
					if charIndex >= len(rippleChars) {
						charIndex = len(rippleChars) - 1
					}
					char := rippleChars[charIndex]

					saturation := 0.3 + intensity*0.4
					value := intensity * 0.6
					color := HSVToRGB(ripple.hue, saturation, value)

					screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}

func updateFlowField(elapsed, peak float64, width, height int) {
	targetFields := int(peak*20) + 5
	if targetFields > 30 {
		targetFields = 30
	}

	// Maintain flow field
	for len(flowField) < targetFields {
		field := FlowField{
			x:         math.Mod(wavePhase*10.0, float64(width)),
			y:         float64(height/2) + math.Sin(liquidPhase)*float64(height/4),
			angle:     liquidPhase + math.Pi/4,
			magnitude: 0.5 + peak*0.5,
			life:      1.0,
		}
		flowField = append(flowField, field)
	}

	// Update flow field with slower, more meditative movement
	for i := 0; i < len(flowField); i++ {
		f := &flowField[i]
		f.x += math.Cos(f.angle) * f.magnitude * elapsed * 4.0
		f.y += math.Sin(f.angle) * f.magnitude * elapsed * 2.0
		f.angle += elapsed * 0.2 // Much slower rotation
		f.life -= elapsed * 0.15 // Longer lived

		// Wrap around
		if f.x < 0 {
			f.x = float64(width)
		}
		if f.x >= float64(width) {
			f.x = 0
		}
	}

	// Remove excess fields
	if len(flowField) > targetFields {
		flowField = flowField[:targetFields]
	}
}

func drawFlowEffects(screen tcell.Screen, width, height int, peak float64) {
	if peak < 0.6 {
		return
	}

	flowChars := []rune{'·', '˙'}

	for _, field := range flowField {
		x, y := int(field.x), int(field.y)
		if x >= 0 && x < width && y >= 0 && y < height {
			intensity := field.magnitude * field.life * (peak - 0.6) * 1.0

			if intensity > 0.4 {
				charIndex := int(intensity * float64(len(flowChars)))
				if charIndex >= len(flowChars) {
					charIndex = len(flowChars) - 1
				}
				char := flowChars[charIndex]

				hue := math.Mod(0.48+field.angle*0.1, 1.0)
				saturation := 0.1 + intensity*0.2
				value := intensity * 0.3

				color := HSVToRGB(hue, saturation, value)
				screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}
