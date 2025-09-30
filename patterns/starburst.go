package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawStarburst creates a retro terminal CLI-style starburst with ASCII line art
func DrawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := math.Min(float64(width), float64(height)) / 2.2

	// Clean peak scaling for terminal feel
	peakScale := 0.3 + peak*0.7 // Always 30% visible, scales to 100%

	// Number of primary rays - classic 8-way plus extras
	baseRays := 8
	extraRays := int(peak * 8) // Add more rays at high peaks
	totalRays := baseRays + extraRays
	if totalRays > 20 {
		totalRays = 20
	}

	// Terminal-style ray characters for different directions
	rayChars := map[string]rune{
		"horizontal": '─',
		"vertical":   '│',
		"diag1":      '╱',
		"diag2":      '╲',
		"thick_h":    '━',
		"thick_v":    '┃',
		"dot":        '·',
		"bullet":     '•',
		"star":       '★',
		"plus":       '+',
		"x":          '×',
	}

	// Draw primary rays
	for rayIndex := 0; rayIndex < totalRays; rayIndex++ {
		rayAngle := float64(rayIndex) * 2 * math.Pi / float64(totalRays)

		// Add some rotation for movement
		rotation := basePhase * 0.1 * (1 - 2*float64(rayIndex%2)) // Alternate directions
		finalAngle := rayAngle + rotation

		// Ray length based on peak with some individual variation
		rayPersonality := float64(rayIndex) * 1.618 // Golden ratio for variation
		lengthVariation := 0.8 + 0.4*math.Sin(basePhase*0.3+rayPersonality)
		rayLength := maxRadius * peakScale * lengthVariation

		// Determine ray character based on angle
		var rayChar rune
		normalizedAngle := math.Mod(finalAngle, 2*math.Pi)

		if (normalizedAngle > -math.Pi/8 && normalizedAngle <= math.Pi/8) ||
			(normalizedAngle > 7*math.Pi/8 || normalizedAngle <= -7*math.Pi/8) {
			rayChar = rayChars["horizontal"]
		} else if normalizedAngle > 3*math.Pi/8 && normalizedAngle <= 5*math.Pi/8 {
			rayChar = rayChars["vertical"]
		} else if normalizedAngle > math.Pi/8 && normalizedAngle <= 3*math.Pi/8 {
			rayChar = rayChars["diag1"]
		} else if normalizedAngle > 5*math.Pi/8 && normalizedAngle <= 7*math.Pi/8 {
			rayChar = rayChars["diag2"]
		} else if normalizedAngle > -3*math.Pi/8 && normalizedAngle <= -math.Pi/8 {
			rayChar = rayChars["diag2"]
		} else {
			rayChar = rayChars["diag1"]
		}

		// Use thicker characters for high peaks
		if peak > 0.7 {
			if rayChar == rayChars["horizontal"] {
				rayChar = rayChars["thick_h"]
			} else if rayChar == rayChars["vertical"] {
				rayChar = rayChars["thick_v"]
			}
		}

		// Draw ray segments
		raySteps := int(rayLength / 1.5)
		if raySteps < 3 {
			raySteps = 3
		}

		for step := 1; step <= raySteps; step++ {
			stepRatio := float64(step) / float64(raySteps)
			currentRadius := stepRatio * rayLength

			if currentRadius < 2 {
				continue
			}

			x := centerX + int(currentRadius*math.Cos(finalAngle))
			y := centerY + int(currentRadius*math.Sin(finalAngle))

			if x >= 0 && x < width && y >= 0 && y < height {
				// Distance-based intensity
				intensity := (1.0 - stepRatio*0.8) * peakScale

				// Character progression based on intensity and distance
				var finalChar rune
				if intensity < 0.2 {
					finalChar = rayChars["dot"]
				} else if intensity < 0.4 {
					finalChar = rayChars["bullet"]
				} else if step == raySteps && intensity > 0.6 {
					finalChar = rayChars["star"] // Star at ray tips
				} else {
					finalChar = rayChar
				}

				// Terminal-style colors
				hue := math.Mod(float64(rayIndex)*0.1+basePhase*0.05, 1)
				saturation := 0.6 + peak*0.3
				brightness := 0.4 + intensity*0.5

				rayColor := HSVToRGB(hue, saturation, brightness)
				screen.SetContent(x, y, finalChar, nil, tcell.StyleDefault.Foreground(rayColor))
			}
		}
	}

	// Draw center core
	coreRadius := 1 + int(peak*3)
	if coreRadius > 4 {
		coreRadius = 4
	}

	coreChars := []rune{'○', '●', '◉', '⬢', '★'}

	for radius := 0; radius <= coreRadius; radius++ {
		if radius == 0 {
			// Center point
			coreIntensity := 0.5 + peak*0.5
			charIndex := int(coreIntensity * float64(len(coreChars)-1))
			if charIndex >= len(coreChars) {
				charIndex = len(coreChars) - 1
			}
			coreChar := coreChars[charIndex]

			coreHue := math.Mod(basePhase*0.03, 1)
			coreColor := HSVToRGB(coreHue, 0.8, 0.7+peak*0.3)
			screen.SetContent(centerX, centerY, coreChar, nil, tcell.StyleDefault.Foreground(coreColor))
		} else {
			// Core rings
			intensity := (1.0 - float64(radius)/float64(coreRadius)) * peak * 0.8
			if intensity > 0.2 {
				points := radius * 6 // 6 points per radius
				for i := 0; i < points; i++ {
					angle := float64(i) * 2 * math.Pi / float64(points)
					x := centerX + int(float64(radius)*math.Cos(angle))
					y := centerY + int(float64(radius)*math.Sin(angle))

					if x >= 0 && x < width && y >= 0 && y < height {
						var ringChar rune
						if intensity > 0.6 {
							ringChar = '●'
						} else if intensity > 0.4 {
							ringChar = '◦'
						} else {
							ringChar = '∘'
						}

						ringHue := math.Mod(basePhase*0.02, 1)
						ringColor := HSVToRGB(ringHue, 0.7, intensity)
						screen.SetContent(x, y, ringChar, nil, tcell.StyleDefault.Foreground(ringColor))
					}
				}
			}
		}
	}

	// Add secondary rays between main rays at high peaks
	if peak > 0.5 && totalRays >= 8 {
		secondaryRays := totalRays / 2
		for rayIndex := 0; rayIndex < secondaryRays; rayIndex++ {
			rayAngle := (float64(rayIndex) + 0.5) * 2 * math.Pi / float64(totalRays)
			rotation := basePhase * 0.05
			finalAngle := rayAngle + rotation

			secondaryLength := maxRadius * peakScale * 0.6 // Shorter than main rays
			raySteps := int(secondaryLength / 2.0)

			for step := 1; step <= raySteps; step++ {
				stepRatio := float64(step) / float64(raySteps)
				currentRadius := stepRatio * secondaryLength

				if currentRadius < 3 {
					continue
				}

				x := centerX + int(currentRadius*math.Cos(finalAngle))
				y := centerY + int(currentRadius*math.Sin(finalAngle))

				if x >= 0 && x < width && y >= 0 && y < height {
					intensity := (1.0 - stepRatio*0.9) * (peak - 0.5) * 2 // Only visible above 50% peak

					if intensity > 0.3 {
						secondaryChar := rayChars["dot"]
						if intensity > 0.6 {
							secondaryChar = rayChars["bullet"]
						}

						secondaryHue := math.Mod(float64(rayIndex)*0.15+basePhase*0.03, 1)
						secondaryColor := HSVToRGB(secondaryHue, 0.5, intensity*0.7)
						screen.SetContent(x, y, secondaryChar, nil, tcell.StyleDefault.Foreground(secondaryColor))
					}
				}
			}
		}
	}

	// Terminal-style border effect at very high peaks
	if peak > 0.8 {
		borderRadius := int(maxRadius * 0.9)
		borderPoints := borderRadius * 4 // Sparse border

		for i := 0; i < borderPoints; i++ {
			angle := float64(i) * 2 * math.Pi / float64(borderPoints)
			x := centerX + int(float64(borderRadius)*math.Cos(angle))
			y := centerY + int(float64(borderRadius)*math.Sin(angle))

			if x >= 0 && x < width && y >= 0 && y < height {
				borderIntensity := (peak - 0.8) * 5 // Scale from 0 to 1 when peak is 0.8 to 1.0

				if borderIntensity > 0.5 && rng.Float64() < 0.3 { // Sparse, random border
					borderChar := rayChars["plus"]
					if borderIntensity > 0.8 {
						borderChar = rayChars["x"]
					}

					borderHue := math.Mod(basePhase*0.08, 1)
					borderColor := HSVToRGB(borderHue, 0.9, borderIntensity)
					screen.SetContent(x, y, borderChar, nil, tcell.StyleDefault.Foreground(borderColor))
				}
			}
		}
	}
}
