package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawStarburst creates organic, fibonacci-harmonious rotating rays with natural branching
func DrawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2.8 // More contained like fibonacci

	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Organic number of rays based on golden ratio and peak
	baseRays := 8
	dynamicRays := int(peak * 12)
	numRays := baseRays + dynamicRays
	if numRays > 16 { // Keep it subtle and harmonious
		numRays = 16
	}

	rayAngleStep := 2 * math.Pi / float64(numRays)

	// Organic character set matching fibonacci aesthetic
	chars := []rune{'∙', '•', '●', '◉', '⋅', '∘', '◦', '○', '⬢', '⬡', '◇', '◆', '✧', '✦', '·', '˙', '°', '⁘', '⁛', '⁝'}

	for rayIndex := 0; rayIndex < numRays; rayIndex++ {
		baseAngle := float64(rayIndex) * rayAngleStep

		// Organic rotation with golden ratio influence
		rayPhase := basePhase * (0.4 + float64(rayIndex%3)*0.15) * goldenRatio
		organicAngle := baseAngle + rayPhase + float64(rayIndex)*goldenAngle*0.1

		// Organic ray characteristics
		rayAmplitude := peak * (0.7 + float64(rayIndex%4)*0.1)
		rayFrequency := 0.3 + float64(rayIndex)*0.03 + peak*0.08

		// Multi-layered organic beams per ray
		numBeams := 1 + int(peak*2)
		if numBeams > 3 {
			numBeams = 3 // Keep subtle
		}

		for beam := 0; beam < numBeams; beam++ {
			beamPhase := basePhase * (0.5 + float64(beam)*0.2)
			beamOffset := (float64(beam) - float64(numBeams)/2) * 0.04 // Subtle beam spread
			finalBeamAngle := organicAngle + beamOffset

			// Organic ray length with golden ratio pulsing
			basePulse := 0.7 + 0.3*math.Sin(beamPhase*2.1+float64(rayIndex)*goldenAngle)
			rayLength := maxRadius * basePulse * (0.9 + peak*0.4)

			// Organic step size for natural density
			baseStep := 0.7 + peak*0.5

			for radius := 1.5; radius < rayLength; radius += baseStep {
				// Complex organic wave functions like fibonacci
				distancePhase := radius * 0.08
				wave1 := rayAmplitude * math.Sin(distancePhase*rayFrequency+beamPhase*1.7)
				wave2 := rayAmplitude * 0.6 * math.Cos(distancePhase*rayFrequency*1.5+beamPhase*1.2)
				wave3 := rayAmplitude * 0.4 * math.Sin(distancePhase*rayFrequency*0.7+beamPhase*2.3)
				organicOffset := wave1 + wave2 + wave3

				// Safety check for NaN/Inf values
				if math.IsNaN(organicOffset) || math.IsInf(organicOffset, 0) {
					organicOffset = 0
				}

				// Golden ratio breathing effect
				breathe := 1 + 0.06*math.Sin(beamPhase*2.0+float64(rayIndex)*goldenAngle+radius*0.02)

				// Organic radius with micro-variations
				microWave := 0.25 * math.Sin(beamPhase*3.5+radius*0.12+float64(beam)*0.8)
				finalRadius := (radius + organicOffset*2.5 + microWave) * breathe

				// Organic angle with subtle variations
				angleVariation := math.Sin(beamPhase*0.9+radius*0.04) * 0.06
				rayAngle := finalBeamAngle + angleVariation

				x := centerX + int(finalRadius*math.Cos(rayAngle))
				y := centerY + int(finalRadius*math.Sin(rayAngle))

				// Bounds check
				if x >= 0 && x < width && y >= 0 && y < height {
					// Organic character selection
					charPhase := float64(rayIndex)*goldenRatio + radius*0.15 + float64(beam)*2.1 + organicOffset*0.5
					if math.IsNaN(charPhase) || math.IsInf(charPhase, 0) {
						charPhase = 0
					}
					charIndex := int(math.Abs(charPhase)) % len(chars)
					rayChar := chars[charIndex]

					// Organic color generation harmonious with fibonacci
					colorPhase := float64(rayIndex)/float64(numRays)*goldenRatio + radius*0.008 + beamPhase*0.12
					hue := math.Mod(colorPhase+organicOffset*0.02, 1)

					saturation := 0.6 + peak*0.25 - radius/(maxRadius*2)
					saturation = math.Max(0.2, math.Min(0.85, saturation))

					baseValue := 0.8 - radius/maxRadius*0.3
					valueVariation := math.Sin(beamPhase*1.4+radius*0.08) * 0.08
					value := baseValue + peak*0.15 + valueVariation + organicOffset*0.03
					value = math.Max(0.25, math.Min(0.95, value))

					rayColor := HSVToRGB(hue, saturation, value)

					// Organic transparency effect
					distanceRatio := radius / rayLength
					if distanceRatio > 0.75 || peak < 0.15 {
						intensity := math.Max(0.3, 1-distanceRatio*1.3) * math.Max(0.3, peak*3)
						if intensity < 0.5 {
							rayChar = '·'
						}
					}

					screen.SetContent(x, y, rayChar, nil, tcell.StyleDefault.Foreground(rayColor))

					// Add organic branching at golden ratio intervals
					if math.Mod(radius, goldenRatio*10) < 0.8 && peak > 0.35 && beam == 0 {
						drawStarburstBranch(screen, x, y, rayAngle, rayColor, peak, rayIndex, width, height)
					}
				}

				// Organic step progression with subtle variations
				stepVariation := 0.1 * math.Sin(beamPhase*2.5+radius*0.1)
				baseStep = 0.7 + peak*0.5 + stepVariation
				if baseStep < 0.3 {
					baseStep = 0.3
				}
			}
		}
	}
}

// drawStarburstBranch creates organic sub-rays that branch off main rays
func drawStarburstBranch(screen tcell.Screen, x, y int, baseAngle float64, color tcell.Color, amplitude float64, rayIndex int, width, height int) {
	branchChars := []rune{'·', '˙', '∘', '◦', '⋅'}
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Organic branch length
	branchLength := 2 + int(amplitude*4)
	if branchLength > 6 {
		branchLength = 6 // Keep subtle
	}

	// Create 1-3 organic branches per ray
	numBranches := 1 + int(amplitude*2)
	if numBranches > 3 {
		numBranches = 3
	}

	for branch := 0; branch < numBranches; branch++ {
		// Organic branch angles using golden ratio
		branchAngleOffset := (float64(branch) - float64(numBranches)/2 + 0.5) * 0.5 * goldenRatio
		organicVariation := math.Sin(GetBasePhase()*1.8+float64(rayIndex)*0.3) * 0.2
		branchAngle := baseAngle + branchAngleOffset + organicVariation

		for step := 1; step <= branchLength; step++ {
			// Organic branch progression
			stepRatio := float64(step) / float64(branchLength)
			organicCurve := math.Sin(stepRatio*math.Pi) * 0.3 // Natural curve

			branchRadius := float64(step) * (0.8 + organicCurve)
			branchX := x + int(branchRadius*math.Cos(branchAngle))
			branchY := y + int(branchRadius*math.Sin(branchAngle))

			if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
				charIndex := (step + branch + rayIndex) % len(branchChars)
				branchChar := branchChars[charIndex]

				// Organic fade with golden ratio
				intensity := (1.0 - stepRatio) * goldenRatio * 0.6
				if intensity > 0.3 {
					screen.SetContent(branchX, branchY, branchChar, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
