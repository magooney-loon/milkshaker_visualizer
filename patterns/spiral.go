package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawSpiral creates organic, fibonacci-harmonious multi-armed rotating spirals
func DrawSpiral(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2.5 // More contained like fibonacci

	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Gentle number of spiral arms - very laid back
	baseArms := 2
	dynamicArms := int(peak * 2)
	numArms := baseArms + dynamicArms
	if numArms > 4 {
		numArms = 4 // Keep it very subtle and smooth
	}

	// Very gentle character set for laid-back feel
	chars := []rune{'⋅', '∘', '◦', '·', '˙', '∙', '°'}

	for arm := 0; arm < numArms; arm++ {
		armPhase := float64(arm) * 2 * math.Pi / float64(numArms)

		// Very slow, gentle rotation for laid-back feel
		rotationSpeed := 0.1 * (1 + float64(arm%2)*0.05) // Very subtle speed variation
		armRotation := basePhase * rotationSpeed * goldenRatio

		// Very gentle arm characteristics
		armAmplitude := peak * (0.3 + float64(arm)*0.05)    // Much more subtle
		armFrequency := 0.2 + float64(arm)*0.02 + peak*0.05 // Very gentle frequency

		// Minimal layers for smooth, laid-back feel
		numLayers := 1 + int(peak*1.5)
		if numLayers > 2 {
			numLayers = 2 // Keep very minimal
		}

		for layer := 0; layer < numLayers; layer++ {
			layerPhase := basePhase * (0.5 + float64(layer)*0.15)
			layerScale := 1.0 - float64(layer)*0.2

			// Start from center with organic growth
			radius := 2.0 + float64(layer)*3
			angle := armPhase + armRotation + float64(layer)*0.4

			// Very gentle, sparse step calculation
			baseStep := 1.2 + peak*0.4 // Larger steps for more spacing

			for radius < maxRadius*layerScale {
				// Very gentle wave functions for smooth movement
				wave1 := armAmplitude * 0.5 * math.Sin(armFrequency*angle+layerPhase*0.8)
				wave2 := armAmplitude * 0.3 * math.Cos(armFrequency*1.2*angle+layerPhase*0.6)
				wave3 := armAmplitude * 0.2 * math.Sin(armFrequency*0.8*angle+layerPhase*1.1)
				organicOffset := wave1 + wave2 + wave3

				// Very gentle breathing effect
				breathe := 1 + 0.03*math.Sin(layerPhase*1.0+float64(arm)*goldenAngle+float64(layer)*0.4)

				// Very subtle radius variations
				microWave := 0.1 * math.Sin(layerPhase*2+radius*0.05)
				finalRadius := (radius + organicOffset*1.5 + microWave) * breathe

				// Very gentle angle variations
				angleVariation := math.Sin(layerPhase*0.4+radius*0.02) * 0.03
				finalAngle := angle + angleVariation

				x := centerX + int(finalRadius*math.Cos(finalAngle))
				y := centerY + int(finalRadius*math.Sin(finalAngle))

				// Bounds check
				if x >= 0 && x < width && y >= 0 && y < height {
					// Organic character selection
					charPhase := float64(arm)*goldenRatio + radius*0.1 + float64(layer)*1.3
					charIndex := int(charPhase) % len(chars)
					displayChar := chars[charIndex]

					// Organic color generation harmonious with fibonacci
					hue := math.Mod(float64(arm)/float64(numArms)*goldenRatio+basePhase*0.08+organicOffset*0.03, 1)
					saturation := 0.3 + peak*0.2 - float64(layer)*0.05 // More muted
					saturation = math.Max(0.1, math.Min(0.6, saturation))

					value := 0.4 + peak*0.2 + organicOffset*0.02 - float64(layer)*0.05 // Softer
					value = math.Max(0.2, math.Min(0.7, value))

					spiralColor := HSVToRGB(hue, saturation, value)

					// Very gentle transparency effect
					distanceRatio := radius / (maxRadius * layerScale)
					if distanceRatio > 0.6 || peak < 0.3 {
						intensity := math.Max(0.2, 1-distanceRatio*0.8) * math.Max(0.2, peak*1.5)
						if intensity < 0.8 {
							displayChar = '·'
						}
					}

					screen.SetContent(x, y, displayChar, nil, tcell.StyleDefault.Foreground(spiralColor))

					// Very rare, gentle branching
					if math.Mod(radius, goldenRatio*15) < 0.3 && peak > 0.6 {
						drawSpiralBranch(screen, x, y, finalAngle, spiralColor, peak, layer, width, height)
					}
				}

				// Gentle, sparse radius progression
				radius += baseStep + organicOffset*0.05
				angle += (0.06 + peak*0.02) * goldenAngle * 0.3 // Slower angle progression
			}
		}
	}
}

// drawSpiralBranch creates very subtle organic branches for spiral arms
func drawSpiralBranch(screen tcell.Screen, x, y int, baseAngle float64, color tcell.Color, amplitude float64, layer int, width, height int) {
	branchChars := []rune{'·', '˙', '∘'}
	branchLength := 1 + int(amplitude*1.5) // Very minimal branching
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Create very minimal branches
	numBranches := 1
	for branch := 0; branch < numBranches; branch++ {
		// Very gentle branch angle
		branchAngle := baseAngle + (float64(branch)*2-1)*0.2*goldenRatio + math.Sin(GetBasePhase()*0.8)*0.08

		for step := 1; step <= branchLength; step++ {
			branchX := x + int(float64(step)*0.4*math.Cos(branchAngle)) // Very subtle offset
			branchY := y + int(float64(step)*0.4*math.Sin(branchAngle))

			if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
				charIndex := (step + branch) % len(branchChars)
				branchChar := branchChars[charIndex]

				// Very gentle fade
				intensity := 1.0 - float64(step)/(float64(branchLength)*1.2)
				if intensity > 0.6 {
					screen.SetContent(branchX, branchY, branchChar, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
