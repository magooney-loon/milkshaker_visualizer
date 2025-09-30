package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawWave creates flowing sine wave patterns across the terminal
func DrawWave(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	basePhase := GetBasePhase()

	// Wave parameters that respond to audio
	numWaves := 3 + int(peak*4) // 3-7 waves based on audio
	if numWaves > 8 {
		numWaves = 8
	}

	// Wave characters for different intensities
	waveChars := []rune{'·', '˙', '∘', '○', '●', '◉', '~', '≈', '∽', '〰'}

	// Draw multiple wave layers
	for waveIndex := 0; waveIndex < numWaves; waveIndex++ {
		// Each wave has different properties
		waveSpeed := 0.5 + float64(waveIndex)*0.3 + peak*0.5
		amplitude := 3 + float64(waveIndex)*2 + peak*8
		frequency := 0.1 + float64(waveIndex)*0.05 + peak*0.1
		verticalOffset := height/2 + int(float64(waveIndex-numWaves/2)*3)

		// Phase offset for this wave
		phaseOffset := float64(waveIndex) * math.Pi / 3

		// Draw wave across screen width
		for x := 0; x < width; x++ {
			// Calculate wave position
			t := basePhase*waveSpeed + phaseOffset
			waveX := float64(x) * frequency

			// Primary sine wave
			primaryY := amplitude * math.Sin(waveX+t)

			// Add harmonics for complexity
			harmonic1 := amplitude * 0.3 * math.Sin(waveX*2+t*1.5)
			harmonic2 := amplitude * 0.2 * math.Sin(waveX*3+t*0.7)

			totalY := primaryY + harmonic1 + harmonic2
			finalY := verticalOffset + int(totalY)

			// Ensure y is in bounds
			if finalY >= 0 && finalY < height {
				// Calculate wave intensity
				distanceFromCenter := math.Abs(float64(finalY - height/2))
				maxDistance := float64(height) / 2
				centerIntensity := 1.0 - distanceFromCenter/maxDistance

				// Wave amplitude intensity
				amplitudeRatio := math.Abs(totalY) / (amplitude * 1.5)
				amplitudeIntensity := 1.0 - amplitudeRatio

				// Combine intensities
				totalIntensity := (centerIntensity*0.6 + amplitudeIntensity*0.4) * (0.4 + peak*0.6)

				// Add some sparkle at wave peaks
				if amplitudeRatio > 0.8 && rng.Float64() < 0.3*peak {
					totalIntensity += 0.4
				}

				if totalIntensity > 0.15 {
					// Character selection based on intensity
					var waveChar rune
					if totalIntensity < 0.2 {
						waveChar = waveChars[0] // dot
					} else if totalIntensity < 0.35 {
						waveChar = waveChars[1] // small dot
					} else if totalIntensity < 0.5 {
						waveChar = waveChars[2] // circle outline
					} else if totalIntensity < 0.65 {
						waveChar = waveChars[3] // circle
					} else if totalIntensity < 0.8 {
						waveChar = waveChars[4] // filled circle
					} else if totalIntensity < 0.9 {
						waveChar = waveChars[6] // wave tilde
					} else {
						waveChar = waveChars[7] // double wave
					}

					// Color based on wave position and time
					hueBase := math.Mod(float64(waveIndex)*0.15+basePhase*0.1, 1)
					hueShift := math.Sin(waveX*0.5+t*0.3) * 0.1
					hue := math.Mod(hueBase+hueShift, 1)

					saturation := 0.5 + peak*0.3 + totalIntensity*0.2
					brightness := 0.3 + totalIntensity*0.5 + peak*0.2

					waveColor := HSVToRGB(hue, saturation, brightness)
					screen.SetContent(x, finalY, waveChar, nil, tcell.StyleDefault.Foreground(waveColor))
				}
			}

			// Add vertical wave lines at high peaks
			if peak > 0.6 && waveIndex < 2 && x%8 == 0 {
				lineHeight := int(amplitude * 0.8)
				startY := finalY - lineHeight/2
				endY := finalY + lineHeight/2

				if startY < 0 {
					startY = 0
				}
				if endY >= height {
					endY = height - 1
				}

				// Recalculate intensity for vertical lines
				lineBaseIntensity := (0.4 + peak*0.6)
				lineIntensity := (peak - 0.6) * 2.5 * lineBaseIntensity
				if lineIntensity > 0.3 {
					// Calculate line colors
					lineHueBase := math.Mod(float64(waveIndex)*0.15+basePhase*0.1, 1)
					lineHueShift := math.Sin(waveX*0.5+t*0.3) * 0.1
					lineHue := math.Mod(lineHueBase+lineHueShift+0.1, 1)
					lineSaturation := (0.5 + peak*0.3 + lineBaseIntensity*0.2) * 0.8

					for lineY := startY; lineY <= endY; lineY++ {
						if lineY >= 0 && lineY < height {
							distFromWave := math.Abs(float64(lineY - finalY))
							lineCharIntensity := lineIntensity * (1.0 - distFromWave/float64(lineHeight/2+1))

							if lineCharIntensity > 0.2 {
								var lineChar rune
								if lineCharIntensity < 0.4 {
									lineChar = '│'
								} else if lineCharIntensity < 0.7 {
									lineChar = '┃'
								} else {
									lineChar = '║'
								}

								lineColor := HSVToRGB(lineHue, lineSaturation, lineCharIntensity*0.8)
								screen.SetContent(x, lineY, lineChar, nil, tcell.StyleDefault.Foreground(lineColor))
							}
						}
					}
				}
			}
		}
	}

	// Add wave interference patterns at very high peaks
	if peak > 0.8 {
		drawWaveInterference(screen, width, height, basePhase, peak)
	}

	// Add flowing particles that follow the waves
	if peak > 0.4 {
		drawWaveParticles(screen, width, height, basePhase, peak, rng, numWaves)
	}
}

// drawWaveInterference creates interference patterns between waves
func drawWaveInterference(screen tcell.Screen, width, height int, phase, peak float64) {
	interferenceChars := []rune{'·', '∘', '○', '◇', '◆', '✦', '✧', '✶'}

	// Create interference grid
	gridSize := 6
	for gx := 0; gx < width/gridSize; gx++ {
		for gy := 0; gy < height/gridSize; gy++ {
			x := gx * gridSize
			y := gy * gridSize

			if x >= width || y >= height {
				continue
			}

			// Calculate interference from multiple wave sources
			interference := 0.0
			for sourceX := width / 4; sourceX < width; sourceX += width / 3 {
				sourceY := height / 2

				dx := float64(x - sourceX)
				dy := float64(y - sourceY)
				distance := math.Sqrt(dx*dx + dy*dy)

				waveValue := math.Sin(distance*0.1-phase*2) / (distance*0.05 + 1)
				interference += waveValue
			}

			interferenceIntensity := math.Abs(interference) * (peak - 0.8) * 5

			if interferenceIntensity > 0.3 {
				charIndex := int(interferenceIntensity * float64(len(interferenceChars)-1))
				if charIndex >= len(interferenceChars) {
					charIndex = len(interferenceChars) - 1
				}

				interferenceChar := interferenceChars[charIndex]

				hue := math.Mod(interference*0.5+phase*0.05, 1)
				if hue < 0 {
					hue = -hue
				}

				color := HSVToRGB(hue, 0.8, interferenceIntensity*0.8)
				screen.SetContent(x, y, interferenceChar, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

// drawWaveParticles adds flowing particles that follow wave patterns
func drawWaveParticles(screen tcell.Screen, width, height int, phase, peak float64, rng *rand.Rand, numWaves int) {
	particleChars := []rune{'·', '∘', '○', '●', '◉', '✦', '*', '⋆'}

	// Number of particles based on audio intensity
	numParticles := int((peak - 0.4) * 50)
	if numParticles > 30 {
		numParticles = 30
	}

	for p := 0; p < numParticles; p++ {
		// Particle follows a wave path
		particlePhase := phase*2 + float64(p)*0.3
		particleX := int(float64(width) * math.Mod(particlePhase*0.1, 1))

		// Calculate Y position based on wave equation
		waveY := float64(height / 2)
		for w := 0; w < Min(numWaves, 3); w++ {
			amplitude := 5 + peak*10
			frequency := 0.15 + float64(w)*0.05
			waveComponent := amplitude * math.Sin(float64(particleX)*frequency+particlePhase+float64(w)*math.Pi/2)
			waveY += waveComponent * (1.0 - float64(w)*0.3)
		}

		particleY := int(waveY)

		if particleX >= 0 && particleX < width && particleY >= 0 && particleY < height {
			// Particle intensity based on wave strength and position
			particleIntensity := (peak - 0.4) * 1.7

			// Add some randomness
			if rng.Float64() < 0.3 {
				particleIntensity *= 0.5
			}

			if particleIntensity > 0.2 {
				charIndex := int(particleIntensity * float64(len(particleChars)-1))
				if charIndex >= len(particleChars) {
					charIndex = len(particleChars) - 1
				}

				particleChar := particleChars[charIndex]

				// Color flows with the particle movement
				hue := math.Mod(particlePhase*0.1+float64(p)*0.1, 1)
				saturation := 0.7 + particleIntensity*0.2
				brightness := 0.4 + particleIntensity*0.5

				particleColor := HSVToRGB(hue, saturation, brightness)
				screen.SetContent(particleX, particleY, particleChar, nil, tcell.StyleDefault.Foreground(particleColor))
			}
		}
	}
}
