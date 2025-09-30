package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawStarburst creates organic, flowing rays with 3D depth and perspective effects
func DrawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2.8

	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// 3D depth layers - create multiple Z planes
	numDepthLayers := 3 + int(peak*3)
	if numDepthLayers > 6 {
		numDepthLayers = 6
	}

	// Process each depth layer from back to front for proper 3D layering
	for depthLayer := numDepthLayers - 1; depthLayer >= 0; depthLayer-- {
		depthRatio := float64(depthLayer) / float64(numDepthLayers-1)

		// 3D perspective effects
		perspectiveScale := 0.4 + depthRatio*0.8 // Back layers smaller
		depthMaxRadius := maxRadius * perspectiveScale

		// 3D Z-axis rotation and movement
		zRotation := basePhase * (0.2 + depthRatio*0.3) * (1 - 2*float64(depthLayer%2))
		zBobbing := math.Sin(basePhase*1.5+float64(depthLayer)*goldenAngle) * 0.1

		// Depth-based phase offset for parallax
		depthPhase := basePhase * (0.8 + depthRatio*0.4)

		// Organic number of rays per depth layer
		baseRays := 4 + depthLayer
		dynamicRays := int(peak * (6 + float64(depthLayer)*2))
		numRays := baseRays + dynamicRays
		if numRays > 12+depthLayer*2 {
			numRays = 12 + depthLayer*2
		}

		// 3D character sets based on depth
		var chars []rune
		if depthLayer < 2 { // Foreground layers
			chars = []rune{'●', '◉', '⬢', '⬡', '◆', '✧', '✦', '∙', '•', '◦', '○'}
		} else if depthLayer < 4 { // Mid layers
			chars = []rune{'◦', '○', '⋅', '∘', '·', '˙', '°', '⁘', '⁛', '⁝'}
		} else { // Background layers
			chars = []rune{'⋅', '·', '˙', '∘', '◦'}
		}

		for rayIndex := 0; rayIndex < numRays; rayIndex++ {
			// 3D organic ray distribution
			baseAngle := float64(rayIndex) * goldenAngle * (2.1 + depthRatio*0.3)

			// 3D rotation with depth influence
			rotationDirection := float64(1-2*(rayIndex%2)) * (1 - depthRatio*0.3)
			rayPhaseSpeed := (0.3 + float64(rayIndex%4)*0.1) * (0.7 + depthRatio*0.6)

			// 3D ray personality with depth variation
			rayPersonality := float64(rayIndex)*goldenRatio + float64(depthLayer)*1.7
			rayPhase := depthPhase * rayPhaseSpeed * rotationDirection

			// 3D starting angle with Z-axis influence
			organicStartAngle := baseAngle + rayPhase + zRotation +
				math.Sin(depthPhase*0.7+rayPersonality)*0.3*perspectiveScale

			// 3D ray characteristics scaled by depth
			rayAmplitude := peak * (0.5 + float64(rayIndex%3)*0.15) * perspectiveScale
			rayFrequency := (0.25 + float64(rayIndex)*0.02 + peak*0.06) * (1 + depthRatio*0.2)

			// Multiple 3D beam layers per ray
			numBeams := 1 + int(peak*1.5*(0.8+depthRatio*0.4))
			if numBeams > 3 {
				numBeams = 3
			}

			for beam := 0; beam < numBeams; beam++ {
				beamPhase := depthPhase * (0.4 + float64(beam)*0.15)
				beamPersonality := rayPersonality + float64(beam)*1.3

				// 3D beam offset with depth perspective
				beamDepthOffset := float64(beam-1) * 0.1 * perspectiveScale

				// 3D curved ray path with depth segments
				maxSegments := int(depthMaxRadius / (2.0 + depthRatio))

				for segment := 0; segment < maxSegments; segment++ {
					segmentRatio := float64(segment) / float64(maxSegments)

					// 3D base radius with perspective
					baseRadius := segmentRatio * depthMaxRadius * (0.8 + peak*0.3)

					// 3D length variation with Z-axis influence
					lengthVariation := 0.9 + 0.4*math.Sin(beamPhase*1.8+rayPersonality+segmentRatio*3+zBobbing*2)
					currentRadius := baseRadius * lengthVariation

					if currentRadius < 1 || currentRadius > depthMaxRadius {
						continue
					}

					// 3D CURVED RAY PATH with depth influence
					curvePhase := currentRadius * (0.04 + depthRatio*0.01)

					// 3D primary curve - main flow with Z-axis influence
					primaryCurve := rayAmplitude * 0.8 * math.Sin(curvePhase*rayFrequency+beamPhase*1.5+rayPersonality+zBobbing)

					// 3D secondary curve - depth complexity
					secondaryCurve := rayAmplitude * 0.4 * math.Cos(curvePhase*rayFrequency*1.6+beamPhase*1.2+rayPersonality*0.7+zRotation*0.5)

					// 3D tertiary curve - micro Z variations
					tertiaryCurve := rayAmplitude * 0.2 * math.Sin(curvePhase*rayFrequency*2.3+beamPhase*2.0+rayPersonality*0.4+zBobbing*3)

					// Combined 3D curvature
					totalCurvature := primaryCurve + secondaryCurve + tertiaryCurve

					// Safety check
					if math.IsNaN(totalCurvature) || math.IsInf(totalCurvature, 0) {
						totalCurvature = 0
					}

					// 3D organic angle calculation with depth rotation
					baseCurvedAngle := organicStartAngle + totalCurvature*0.15*perspectiveScale

					// 3D flowing variations with depth influence
					flowVariation := math.Sin(beamPhase*0.6+currentRadius*0.03+rayPersonality+zRotation) * 0.08
					organicFlow := math.Cos(beamPhase*0.9+currentRadius*0.05+beamPersonality+zBobbing*2) * 0.05
					depthAngleShift := depthRatio * math.Sin(depthPhase*0.4+rayPersonality) * 0.1

					finalAngle := baseCurvedAngle + flowVariation + organicFlow + depthAngleShift

					// 3D breathing radius with depth perspective
					breathe := 1 + 0.05*math.Sin(beamPhase*1.7+rayPersonality+currentRadius*0.02+zBobbing*4)*perspectiveScale
					finalRadius := currentRadius * breathe

					// 3D micro-positioning with depth parallax
					microX := (0.3*math.Sin(beamPhase*3.2+currentRadius*0.08) + beamDepthOffset) * perspectiveScale
					microY := (0.3*math.Cos(beamPhase*3.7+currentRadius*0.09) + zBobbing*2) * perspectiveScale

					x := centerX + int(finalRadius*math.Cos(finalAngle)+microX)
					y := centerY + int(finalRadius*math.Sin(finalAngle)+microY)

					// Bounds check
					if x >= 0 && x < width && y >= 0 && y < height {
						// 3D character selection based on depth and flow
						// Calculate starburst intensity for intelligent character fading
						starburstIntensity := peak * (1.0 - currentRadius/depthMaxRadius*0.7) * perspectiveScale * (0.4 + math.Abs(totalCurvature)*0.6)

						// Character selection
						charPhase := rayPersonality + currentRadius*0.12 + totalCurvature*0.3 +
							float64(beam)*1.8 + float64(depthLayer)*2.5
						if math.IsNaN(charPhase) || math.IsInf(charPhase, 0) {
							charPhase = 0
						}
						charIndex := int(math.Abs(charPhase)*goldenRatio) % len(chars)
						baseRayChar := chars[charIndex]

						// Intelligent character fading based on intensity
						var rayChar rune
						if starburstIntensity < 0.1 {
							rayChar = '·' // Barely visible dot
						} else if starburstIntensity < 0.2 {
							rayChar = '˙' // Small dot
						} else if starburstIntensity < 0.35 {
							rayChar = '∘' // Circle outline
						} else if starburstIntensity < 0.5 {
							rayChar = '◦' // Larger circle
						} else if starburstIntensity < 0.7 {
							rayChar = '●' // Filled circle
						} else {
							rayChar = baseRayChar // Full character set
						}

						// 3D color with depth-based hues and perspective
						colorPhase := rayPersonality*0.3 + currentRadius*0.006 + beamPhase*0.1 +
							totalCurvature*0.02 + float64(depthLayer)*0.15
						hue := math.Mod(colorPhase, 1)

						// 3D saturation with depth attenuation
						saturation := (0.4 + peak*0.3 + math.Abs(totalCurvature)*0.1 - currentRadius/(depthMaxRadius*2.5)) *
							(0.6 + depthRatio*0.5)
						saturation = math.Max(0.1, math.Min(0.8, saturation))

						// 3D brightness with depth-based dimming
						baseValue := (0.7 - currentRadius/depthMaxRadius*0.25) * (0.4 + depthRatio*0.6)
						valueFlow := math.Sin(beamPhase*1.1+currentRadius*0.06) * 0.1 * perspectiveScale
						depthDimming := 1.0 - (1.0-depthRatio)*0.4 // Back layers dimmer
						value := (baseValue + peak*0.2 + valueFlow + math.Abs(totalCurvature)*0.05) * depthDimming
						value = math.Max(0.15, math.Min(0.9, value))

						rayColor := HSVToRGB(hue, saturation, value)

						// Additional fading for very weak areas and depth transparency
						distanceRatio := currentRadius / depthMaxRadius

						if distanceRatio > 0.8 || peak < 0.15 || starburstIntensity < 0.08 {
							rayChar = '·'
						}

						screen.SetContent(x, y, rayChar, nil, tcell.StyleDefault.Foreground(rayColor))

						// 3D organic curved branching at flow peaks
						if math.Abs(totalCurvature) > rayAmplitude*0.6 && peak > 0.4 && beam == 0 &&
							segment%(3+depthLayer) == 0 && depthLayer < 3 { // Less branching in back layers
							draw3DOrganicStarburstBranch(screen, x, y, finalAngle, totalCurvature, rayColor,
								peak, rayIndex, depthLayer, perspectiveScale, width, height)
						}
					}
				}
			}
		}
	}
}

// draw3DOrganicStarburstBranch creates 3D flowing, curved branches with depth effects
func draw3DOrganicStarburstBranch(screen tcell.Screen, x, y int, baseAngle, curvature float64, color tcell.Color,
	amplitude float64, rayIndex, depthLayer int, perspectiveScale float64, width, height int) {

	branchChars := []rune{'·', '˙', '∘', '◦', '⋅'}
	goldenRatio := (1 + math.Sqrt(5)) / 2
	goldenAngle := math.Pi * (3 - math.Sqrt(5))
	basePhase := GetBasePhase()

	// 3D branch characteristics scaled by depth
	branchLength := int(float64(2+int(amplitude*3)) * perspectiveScale)
	if branchLength > 5 {
		branchLength = 5
	}
	if branchLength < 1 {
		branchLength = 1
	}

	// 3D multiple organic branches with depth variation
	numBranches := 1 + int(amplitude*1.5*perspectiveScale)
	if numBranches > 3 {
		numBranches = 3
	}
	if numBranches < 1 {
		numBranches = 1
	}

	for branch := 0; branch < numBranches; branch++ {
		branchPersonality := float64(rayIndex)*goldenRatio + float64(branch)*2.1 + float64(depthLayer)*1.3

		// 3D branch angle influenced by main ray curvature and depth
		curvatureInfluence := curvature * 0.4 * perspectiveScale
		branchBaseAngle := baseAngle + curvatureInfluence

		// 3D organic branch angle variations with depth influence
		branchOffset := (float64(branch) - float64(numBranches)/2 + 0.5) * 0.6 * perspectiveScale
		organicVariation := math.Sin(basePhase*1.4+branchPersonality) * 0.25 * perspectiveScale
		depthAngleVariation := float64(depthLayer) * 0.1 * math.Cos(basePhase*0.8+branchPersonality)

		startAngle := branchBaseAngle + branchOffset + organicVariation + depthAngleVariation

		for step := 1; step <= branchLength; step++ {
			stepRatio := float64(step) / float64(branchLength)

			// 3D organic branch curve that flows with the main ray and depth
			branchCurve := math.Sin(stepRatio*math.Pi*1.5+branchPersonality) * 0.4 * perspectiveScale
			organicBend := math.Cos(basePhase*2.1+branchPersonality+stepRatio*2) * 0.2 * perspectiveScale
			depthCurve := math.Sin(basePhase*1.0+float64(depthLayer)*goldenAngle+stepRatio*1.5) * 0.1

			// 3D branch angle curves organically with depth influence
			currentAngle := startAngle + branchCurve + organicBend + depthCurve

			// 3D organic branch radius with flowing variations and perspective
			baseRadius := float64(step) * (0.7 + math.Sin(stepRatio*math.Pi)*0.3) * perspectiveScale
			radiusFlow := 1 + 0.1*math.Sin(basePhase*2.5+branchPersonality+stepRatio*3)*perspectiveScale
			branchRadius := baseRadius * radiusFlow

			branchX := x + int(branchRadius*math.Cos(currentAngle))
			branchY := y + int(branchRadius*math.Sin(currentAngle))

			if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
				charIndex := (step + branch + rayIndex + depthLayer) % len(branchChars)
				baseBranchChar := branchChars[charIndex]

				// Calculate branch intensity for intelligent character fading
				flowIntensity := 1.0 + math.Abs(branchCurve)*0.5
				depthAttenuation := 0.4 + float64(depthLayer)/6.0*0.6 // Back layers dimmer
				branchIntensity := (1.0 - stepRatio*0.7) * flowIntensity * goldenRatio * 0.4 *
					perspectiveScale * depthAttenuation * amplitude

				// Intelligent character fading for branches
				var branchChar rune
				if branchIntensity < 0.1 {
					branchChar = '·'
				} else if branchIntensity < 0.2 {
					branchChar = '˙'
				} else if branchIntensity < 0.35 {
					branchChar = '∘'
				} else {
					branchChar = baseBranchChar
				}

				// Only render if branch intensity is above threshold
				if branchIntensity > 0.08 {
					screen.SetContent(branchX, branchY, branchChar, nil, tcell.StyleDefault.Foreground(color))
				}
			}
		}
	}
}
