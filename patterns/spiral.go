package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawSpiral creates organic, fibonacci-harmonious multi-armed spirals with 3D depth and perspective
func DrawSpiral(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2.5

	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// 3D depth layers - create multiple Z planes
	numDepthLayers := 3 + int(peak*2)
	if numDepthLayers > 5 {
		numDepthLayers = 5
	}

	// Process each depth layer from back to front for proper 3D layering
	for depthLayer := numDepthLayers - 1; depthLayer >= 0; depthLayer-- {
		depthRatio := float64(depthLayer) / float64(numDepthLayers-1)

		// 3D perspective effects
		perspectiveScale := 0.5 + depthRatio*0.7 // Back layers smaller
		depthMaxRadius := maxRadius * perspectiveScale

		// 3D Z-axis rotation and movement - gentler than starburst
		zRotation := basePhase * (0.1 + depthRatio*0.2) * (1 - 2*float64(depthLayer%2))
		zBobbing := math.Sin(basePhase*1.2+float64(depthLayer)*goldenAngle) * 0.08

		// Depth-based phase offset for parallax
		depthPhase := basePhase * (0.6 + depthRatio*0.4)

		// 3D gentle number of spiral arms per depth layer
		baseArms := 2 + depthLayer%2
		dynamicArms := int(peak * (1.5 + float64(depthLayer)*0.3))
		numArms := baseArms + dynamicArms
		if numArms > 3+depthLayer {
			numArms = 3 + depthLayer
		}

		// 3D character sets based on depth - gentler progression
		var chars []rune
		if depthLayer < 2 { // Foreground layers - slightly more visible
			chars = []rune{'◦', '○', '●', '⋅', '∘', '·', '˙', '∙', '°'}
		} else if depthLayer < 4 { // Mid layers - very subtle
			chars = []rune{'⋅', '∘', '◦', '·', '˙', '∙', '°'}
		} else { // Background layers - barely visible
			chars = []rune{'⋅', '·', '˙', '∘'}
		}

		for arm := 0; arm < numArms; arm++ {
			// 3D arm phase with depth influence
			armPhase := float64(arm)*2*math.Pi/float64(numArms) +
				float64(depthLayer)*goldenAngle*0.2

			// 3D very slow, gentle rotation with depth variation
			rotationSpeed := (0.1 * (1 + float64(arm%2)*0.05)) * (0.8 + depthRatio*0.4)
			armRotation := depthPhase * rotationSpeed * goldenRatio

			// 3D very gentle arm characteristics scaled by depth
			armAmplitude := peak * (0.3 + float64(arm)*0.05) * perspectiveScale
			armFrequency := (0.2 + float64(arm)*0.02 + peak*0.05) * (1 + depthRatio*0.1)

			// 3D minimal layers for smooth, laid-back feel
			numLayers := 1 + int(peak*1.2*(0.7+depthRatio*0.5))
			if numLayers > 2 {
				numLayers = 2
			}

			for layer := 0; layer < numLayers; layer++ {
				layerPhase := depthPhase * (0.5 + float64(layer)*0.15)
				layerScale := (1.0 - float64(layer)*0.2) * perspectiveScale

				// 3D start from center with organic growth and depth influence
				radius := (2.0 + float64(layer)*3) * perspectiveScale
				angle := armPhase + armRotation + float64(layer)*0.4 + zRotation*0.5

				// 3D very gentle, sparse step calculation with depth scaling
				baseStep := (1.2 + peak*0.4) * perspectiveScale

				for radius < depthMaxRadius*layerScale {
					// 3D very gentle wave functions with depth influence
					wave1 := armAmplitude * 0.5 * math.Sin(armFrequency*angle+layerPhase*0.8+zBobbing)
					wave2 := armAmplitude * 0.3 * math.Cos(armFrequency*1.2*angle+layerPhase*0.6+zRotation*0.3)
					wave3 := armAmplitude * 0.2 * math.Sin(armFrequency*0.8*angle+layerPhase*1.1+zBobbing*2)
					organicOffset := wave1 + wave2 + wave3

					// Safety check
					if math.IsNaN(organicOffset) || math.IsInf(organicOffset, 0) {
						organicOffset = 0
					}

					// 3D very gentle breathing effect with depth influence
					breathe := 1 + 0.03*math.Sin(layerPhase*1.0+float64(arm)*goldenAngle+
						float64(layer)*0.4+zBobbing*3)*perspectiveScale

					// 3D very subtle radius variations with depth scaling
					microWave := 0.1 * math.Sin(layerPhase*2+radius*0.05+zRotation*0.5) * perspectiveScale
					finalRadius := (radius + organicOffset*1.5 + microWave) * breathe

					// 3D very gentle angle variations with depth influence
					angleVariation := math.Sin(layerPhase*0.4+radius*0.02+zBobbing) * 0.03 * perspectiveScale
					depthAngleShift := depthRatio * math.Cos(depthPhase*0.3+float64(arm)*goldenAngle) * 0.05
					finalAngle := angle + angleVariation + depthAngleShift

					// 3D micro-positioning with depth parallax
					microX := 0.2 * math.Sin(layerPhase*2.5+radius*0.04) * perspectiveScale
					microY := (0.2*math.Cos(layerPhase*2.8+radius*0.045) + zBobbing) * perspectiveScale

					x := centerX + int(finalRadius*math.Cos(finalAngle)+microX)
					y := centerY + int(finalRadius*math.Sin(finalAngle)+microY)

					// Bounds check
					if x >= 0 && x < width && y >= 0 && y < height {
						// 3D organic character selection with depth influence
						charPhase := float64(arm)*goldenRatio + radius*0.1 + float64(layer)*1.3 +
							float64(depthLayer)*2.0 + organicOffset*0.5
						charIndex := int(charPhase) % len(chars)
						displayChar := chars[charIndex]

						// 3D organic color generation with depth-based hues
						hue := math.Mod(float64(arm)/float64(numArms)*goldenRatio+depthPhase*0.08+
							organicOffset*0.03+float64(depthLayer)*0.1, 1)

						// 3D saturation with depth attenuation - more muted
						saturation := (0.3 + peak*0.2 - float64(layer)*0.05) * (0.5 + depthRatio*0.5)
						saturation = math.Max(0.05, math.Min(0.6, saturation))

						// 3D brightness with depth-based dimming - softer
						baseValue := (0.4 + peak*0.2 + organicOffset*0.02 - float64(layer)*0.05) *
							(0.3 + depthRatio*0.7)
						depthDimming := 1.0 - (1.0-depthRatio)*0.3 // Back layers dimmer but gentler
						value := baseValue * depthDimming
						value = math.Max(0.15, math.Min(0.7, value))

						spiralColor := HSVToRGB(hue, saturation, value)

						// 3D very gentle transparency effect with depth
						distanceRatio := radius / (depthMaxRadius * layerScale)
						depthTransparency := 0.4 + depthRatio*0.6 // Back layers more transparent

						if distanceRatio > 0.6 || peak < 0.3 {
							intensity := math.Max(0.15, 1-distanceRatio*0.8) *
								math.Max(0.15, peak*1.5) * depthTransparency
							if intensity < 0.8 {
								displayChar = '·'
							}
						}

						screen.SetContent(x, y, displayChar, nil, tcell.StyleDefault.Foreground(spiralColor))

						// 3D very rare, gentle branching with depth considerations
						if math.Mod(radius, goldenRatio*15) < 0.3 && peak > 0.6 && depthLayer < 3 {
							draw3DSpiralBranch(screen, x, y, finalAngle, spiralColor, peak, layer,
								depthLayer, perspectiveScale, width, height)
						}
					}

					// 3D gentle, sparse radius progression with depth scaling
					radius += baseStep + organicOffset*0.05*perspectiveScale
					angle += (0.06 + peak*0.02) * goldenAngle * 0.3 * (1 + depthRatio*0.2)
				}
			}
		}
	}
}

// draw3DSpiralBranch creates very subtle 3D organic branches for spiral arms with depth effects
func draw3DSpiralBranch(screen tcell.Screen, x, y int, baseAngle float64, color tcell.Color,
	amplitude float64, layer, depthLayer int, perspectiveScale float64, width, height int) {

	branchChars := []rune{'·', '˙', '∘'}
	goldenRatio := (1 + math.Sqrt(5)) / 2
	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	basePhase := GetBasePhase()

	// 3D very minimal branching scaled by depth
	branchLength := int(float64(1+int(amplitude*1.5)) * perspectiveScale)
	if branchLength > 3 {
		branchLength = 3
	}
	if branchLength < 1 {
		branchLength = 1
	}

	// 3D create very minimal branches with depth variation
	numBranches := 1 // Keep it super minimal for laid-back feel

	for branch := 0; branch < numBranches; branch++ {
		branchPersonality := float64(layer)*goldenRatio + float64(branch)*1.5 + float64(depthLayer)*0.8

		// 3D very gentle branch angle with depth influence
		organicVariation := math.Sin(basePhase*0.8+branchPersonality) * 0.08 * perspectiveScale
		depthAngleVariation := float64(depthLayer) * 0.05 * math.Cos(basePhase*0.6+branchPersonality)
		branchAngle := baseAngle + (float64(branch)*2-1)*0.2*goldenRatio*perspectiveScale +
			organicVariation + depthAngleVariation

		for step := 1; step <= branchLength; step++ {
			stepRatio := float64(step) / float64(branchLength)

			// 3D organic branch curve with depth influence
			branchCurve := math.Sin(stepRatio*math.Pi*1.2+branchPersonality) * 0.3 * perspectiveScale
			depthCurve := math.Cos(basePhase*1.5+float64(depthLayer)*goldenAngle+stepRatio) * 0.1

			// 3D branch angle curves organically with depth
			currentAngle := branchAngle + branchCurve + depthCurve

			// 3D organic branch radius with flowing variations and perspective
			baseRadius := float64(step) * (0.4 + math.Sin(stepRatio*math.Pi)*0.2) * perspectiveScale
			radiusFlow := 1 + 0.08*math.Sin(basePhase*2.0+branchPersonality+stepRatio*2)*perspectiveScale
			branchRadius := baseRadius * radiusFlow

			branchX := x + int(branchRadius*math.Cos(currentAngle))
			branchY := y + int(branchRadius*math.Sin(currentAngle))

			if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
				charIndex := (step + branch + layer + depthLayer) % len(branchChars)
				branchChar := branchChars[charIndex]

				// 3D very gentle fade with depth attenuation
				flowIntensity := 1.0 + math.Abs(branchCurve)*0.3
				depthAttenuation := 0.5 + float64(depthLayer)/8.0*0.5 // Very gentle dimming
				intensity := (1.0 - stepRatio*0.8) * flowIntensity * goldenRatio * 0.4 *
					perspectiveScale * depthAttenuation

				if intensity > 0.5 { // Higher threshold for gentler appearance
					screen.SetContent(branchX, branchY, branchChar, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
