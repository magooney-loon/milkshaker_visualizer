package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawSpiral creates organic procedural flow patterns with counter-rotating streams
func DrawSpiral(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := math.Sqrt(float64(width*width+height*height)) / 1.8 // Expand to use more screen space

	// 3D depth layers for organic flows
	numDepthLayers := 3 + int(peak*2)
	if numDepthLayers > 5 {
		numDepthLayers = 5
	}

	// Process each depth layer from back to front
	for depthLayer := numDepthLayers - 1; depthLayer >= 0; depthLayer-- {
		depthRatio := float64(depthLayer) / float64(numDepthLayers-1)
		perspectiveScale := 0.5 + depthRatio*0.7
		depthMaxRadius := maxRadius * perspectiveScale

		// Organic phase with depth influence
		depthPhase := basePhase * (0.4 + depthRatio*0.5)

		// 3D character sets based on depth - very subtle
		var chars []rune
		if depthLayer < 2 {
			chars = []rune{'◦', '○', '⋅', '∘', '·', '˙', '∙'}
		} else if depthLayer < 4 {
			chars = []rune{'⋅', '∘', '◦', '·', '˙'}
		} else {
			chars = []rune{'⋅', '·', '˙'}
		}

		// PROCEDURAL ORGANIC FLOWS instead of geometric spirals
		drawOrganicFlows(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)

		// COUNTER-ROTATING STREAMS
		drawCounterStreams(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)

		// ORGANIC GROWTH TENDRILS
		drawGrowthTendrils(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)
	}
}

// drawOrganicFlows creates flowing organic patterns using procedural generation
func drawOrganicFlows(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2
	goldenAngle := math.Pi * (3 - math.Sqrt(5))

	// Number of organic flow streams - slightly more prominent
	numFlows := 3 + int(peak*2)
	if numFlows > 5 {
		numFlows = 5
	}

	for flowIndex := 0; flowIndex < numFlows; flowIndex++ {
		flowPersonality := float64(flowIndex)*goldenRatio + float64(depthLayer)*1.3

		// Organic flow direction - some clockwise, some counter-clockwise
		flowDirection := float64(1 - 2*(flowIndex%2))

		// Procedural starting angle
		baseAngle := flowPersonality * goldenAngle * 1.7

		// Organic rotation speed
		rotationSpeed := (0.15 + float64(flowIndex)*0.05) * flowDirection
		currentAngle := baseAngle + phase*rotationSpeed

		// Flow characteristics
		flowAmplitude := peak * (0.2 + float64(flowIndex)*0.03) * scale
		flowFrequency := 0.15 + float64(flowIndex)*0.02

		// Generate organic flow path using procedural methods
		stepCount := int(maxRadius / (3.0 + float64(depthLayer)*0.5))

		for step := 0; step < stepCount; step++ {
			stepRatio := float64(step) / float64(stepCount)
			// Start from minimum radius to avoid center density
			minRadius := maxRadius * 0.15
			radius := minRadius + stepRatio*(maxRadius-minRadius)*(0.9+peak*0.4)

			// PROCEDURAL ORGANIC CURVATURE - not geometric
			// Use multiple noise-like functions for natural flow
			flowNoise1 := flowAmplitude * math.Sin(radius*flowFrequency*0.8+phase*1.2+flowPersonality)
			flowNoise2 := flowAmplitude * 0.6 * math.Cos(radius*flowFrequency*1.3+phase*0.9+flowPersonality*0.7)
			flowNoise3 := flowAmplitude * 0.4 * math.Sin(radius*flowFrequency*2.1+phase*1.6+flowPersonality*0.4)

			// Combine for organic curvature
			organicCurvature := flowNoise1 + flowNoise2 + flowNoise3

			// Organic breathing and pulsing
			breathe := 1 + 0.04*math.Sin(phase*1.5+flowPersonality+radius*0.03)*scale

			// Natural angle deviation - not spiral-like
			angleDeviation := organicCurvature * 0.12
			finalAngle := currentAngle + angleDeviation

			// Organic radius with natural variations
			radiusVariation := 0.1 * math.Sin(phase*2.2+radius*0.06+flowPersonality)
			finalRadius := (radius + radiusVariation) * breathe

			// Micro-positioning for organic feel
			microX := 0.2 * math.Sin(phase*2.8+radius*0.04) * scale
			microY := 0.2 * math.Cos(phase*3.1+radius*0.045) * scale

			x := centerX + int(finalRadius*math.Cos(finalAngle)+microX)
			y := centerY + int(finalRadius*math.Sin(finalAngle)+microY)

			if x >= 0 && x < width && y >= 0 && y < height {
				// Calculate flow intensity for intelligent character fading
				flowIntensity := math.Abs(organicCurvature) * (1.0 - stepRatio*0.4) * scale * peak

				// Organic character selection
				charPhase := flowPersonality + radius*0.08 + organicCurvature*0.4
				charIndex := int(math.Abs(charPhase)) % len(chars)
				baseDisplayChar := chars[charIndex]

				// Intelligent character fading based on intensity - more visible
				var displayChar rune
				if flowIntensity < 0.08 {
					displayChar = '·' // Barely visible dot
				} else if flowIntensity < 0.18 {
					displayChar = '˙' // Small dot
				} else if flowIntensity < 0.3 {
					displayChar = '∘' // Circle outline
				} else if flowIntensity < 0.45 {
					displayChar = '◦' // Larger circle
				} else {
					displayChar = baseDisplayChar // Full character set
				}

				// Organic color - enhanced
				hue := math.Mod(flowPersonality*0.4+phase*0.06+organicCurvature*0.02, 1)
				saturation := (0.3 + peak*0.18) * (0.45 + stepRatio*0.4) * scale
				saturation = math.Max(0.08, math.Min(0.6, saturation))

				value := (0.35 + peak*0.18) * (0.35 + stepRatio*0.5) * scale
				value = math.Max(0.15, math.Min(0.7, value))

				flowColor := HSVToRGB(hue, saturation, value)

				// Additional fading for very weak areas - less aggressive
				if stepRatio > 0.85 || flowIntensity < 0.06 {
					displayChar = '·'
				}

				screen.SetContent(x, y, displayChar, nil, tcell.StyleDefault.Foreground(flowColor))
			}

			// Organic angle progression - not geometric
			angleStep := (0.04 + peak*0.02) * goldenAngle * 0.4 * flowDirection
			currentAngle += angleStep + organicCurvature*0.05
		}
	}
}

// drawCounterStreams creates counter-rotating organic streams
func drawCounterStreams(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Counter-streams - slightly more prominent
	numStreams := 2 + int(peak*1.5)
	if numStreams > 4 {
		numStreams = 4
	}

	for streamIndex := 0; streamIndex < numStreams; streamIndex++ {
		streamPersonality := float64(streamIndex)*goldenRatio*2.1 + float64(depthLayer)*0.9

		// Alternating directions for organic counter-flow
		streamDirection := float64(1 - 2*(streamIndex%2))

		// Organic starting position
		startAngle := streamPersonality * 1.3
		rotationSpeed := 0.12 * streamDirection * (0.8 + float64(streamIndex)*0.1)

		// Stream characteristics
		streamAmplitude := peak * (0.15 + float64(streamIndex)*0.02) * scale

		// Create organic stream path - expand and avoid center
		streamLength := int(maxRadius * 1.2)
		stepSize := 2.5 + float64(depthLayer)*0.3
		minStartPos := maxRadius * 0.12

		for pos := minStartPos; pos < float64(streamLength); pos += stepSize {
			posRatio := pos / float64(streamLength)

			// Organic stream curvature - not spiral
			streamCurve := streamAmplitude * math.Sin(pos*0.06+phase*1.4+streamPersonality*0.8)
			organicWiggle := streamAmplitude * 0.5 * math.Cos(pos*0.09+phase*1.1+streamPersonality*0.5)

			// Current angle with organic deviation
			currentAngle := startAngle + phase*rotationSpeed + streamCurve*0.08 + organicWiggle*0.05

			// Organic radius with natural pulsing - expanded range
			baseRadius := pos * (1.1 + 0.3*math.Sin(phase*1.6+streamPersonality))
			radiusPulse := 1 + 0.04*math.Sin(phase*2.0+pos*0.04)*scale
			finalRadius := baseRadius * radiusPulse

			x := centerX + int(finalRadius*math.Cos(currentAngle))
			y := centerY + int(finalRadius*math.Sin(currentAngle))

			if x >= 0 && x < width && y >= 0 && y < height {
				// Calculate stream intensity for intelligent character fading
				streamIntensity := math.Abs(streamCurve) * (1.0 - posRatio*0.3) * scale * peak

				// Character selection
				charPhase := streamPersonality + pos*0.07 + streamCurve*0.3
				charIndex := int(math.Abs(charPhase)) % len(chars)
				baseStreamChar := chars[charIndex]

				// Intelligent character fading based on intensity - more visible
				var streamChar rune
				if streamIntensity < 0.06 {
					streamChar = '·' // Barely visible dot
				} else if streamIntensity < 0.15 {
					streamChar = '˙' // Small dot
				} else if streamIntensity < 0.28 {
					streamChar = '∘' // Circle outline
				} else if streamIntensity < 0.42 {
					streamChar = '◦' // Larger circle
				} else {
					streamChar = baseStreamChar // Full character set
				}

				// Organic color - enhanced
				hue := math.Mod(streamPersonality*0.3+phase*0.05+streamCurve*0.02, 1)
				saturation := (0.25 + peak*0.12) * (0.65 + posRatio*0.3) * scale
				saturation = math.Max(0.06, math.Min(0.5, saturation))

				value := (0.3 + peak*0.12) * (0.55 + posRatio*0.4) * scale
				value = math.Max(0.12, math.Min(0.6, value))

				streamColor := HSVToRGB(hue, saturation, value)

				// Subtle transparency
				if posRatio > 0.8 || math.Abs(streamCurve) < streamAmplitude*0.4 {
					streamChar = '·'
				}

				screen.SetContent(x, y, streamChar, nil, tcell.StyleDefault.Foreground(streamColor))
			}
		}
	}
}

// drawGrowthTendrils creates organic growth patterns like plant tendrils
func drawGrowthTendrils(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2
	goldenAngle := math.Pi * (3 - math.Sqrt(5))

	// Organic tendrils based on peak - slightly more prominent
	numTendrils := 2 + int(peak*2)
	if numTendrils > 5 {
		numTendrils = 5
	}

	for tendrilIndex := 0; tendrilIndex < numTendrils; tendrilIndex++ {
		tendrilPersonality := float64(tendrilIndex)*goldenRatio*1.6 + float64(depthLayer)*1.1

		// Organic growth direction
		growthAngle := tendrilPersonality * goldenAngle * 0.7

		// Growth characteristics
		tendrilAmplitude := peak * (0.1 + float64(tendrilIndex)*0.02) * scale
		growthSpeed := 0.08 + float64(tendrilIndex)*0.02

		// Organic growth path - expand and avoid center
		maxGrowthSteps := int(maxRadius * 0.9)
		currentAngle := growthAngle
		minGrowthStart := maxRadius * 0.1

		for growth := minGrowthStart; growth < float64(maxGrowthSteps); growth += 2.8 + float64(depthLayer)*0.2 {
			growthRatio := growth / float64(maxGrowthSteps)

			// Organic tendril curvature - like plant growth
			growthCurve := tendrilAmplitude * math.Sin(growth*0.08+phase*1.3+tendrilPersonality)
			organicTwist := tendrilAmplitude * 0.4 * math.Cos(growth*0.12+phase*0.9+tendrilPersonality*0.6)

			// Natural growth angle changes
			angleChange := (growthCurve + organicTwist) * 0.06
			currentAngle += angleChange + growthSpeed*goldenAngle*0.2

			// Organic growth radius - expanded range
			growthRadius := growth * (1.2 + 0.15*math.Sin(phase*1.8+tendrilPersonality))

			// Organic breathing
			breathe := 1 + 0.03*math.Sin(phase*2.3+growth*0.05)*scale
			finalRadius := growthRadius * breathe

			x := centerX + int(finalRadius*math.Cos(currentAngle))
			y := centerY + int(finalRadius*math.Sin(currentAngle))

			if x >= 0 && x < width && y >= 0 && y < height {
				// Calculate tendril intensity for intelligent character fading
				tendrilIntensity := math.Abs(growthCurve) * (1.0 - growthRatio*0.5) * scale * peak

				// Character selection
				charPhase := tendrilPersonality + growth*0.06 + growthCurve*0.2
				charIndex := int(math.Abs(charPhase)) % len(chars)
				baseTendrilChar := chars[charIndex]

				// Intelligent character fading based on intensity - more visible
				var tendrilChar rune
				if tendrilIntensity < 0.05 {
					tendrilChar = '·' // Barely visible dot
				} else if tendrilIntensity < 0.12 {
					tendrilChar = '˙' // Small dot
				} else if tendrilIntensity < 0.22 {
					tendrilChar = '∘' // Circle outline
				} else if tendrilIntensity < 0.38 {
					tendrilChar = '◦' // Larger circle
				} else {
					tendrilChar = baseTendrilChar // Full character set
				}

				// Subtle organic color - enhanced
				hue := math.Mod(tendrilPersonality*0.2+phase*0.04, 1)
				saturation := (0.18 + peak*0.1) * (0.35 + growthRatio*0.5) * scale
				saturation = math.Max(0.04, math.Min(0.4, saturation))

				value := (0.25 + peak*0.1) * (0.45 + growthRatio*0.4) * scale
				value = math.Max(0.1, math.Min(0.5, value))

				tendrilColor := HSVToRGB(hue, saturation, value)

				// Additional fading for very weak areas - less aggressive
				if growthRatio > 0.92 || tendrilIntensity < 0.04 {
					tendrilChar = '·'
				}

				screen.SetContent(x, y, tendrilChar, nil, tcell.StyleDefault.Foreground(tendrilColor))

				// Rare organic branching
				if int(growth)%15 == 0 && peak > 0.5 && depthLayer < 2 {
					drawTendrilBranch(screen, x, y, currentAngle, tendrilColor, peak, scale, width, height)
				}
			}
		}
	}
}

// drawTendrilBranch creates small organic branches from growth tendrils
func drawTendrilBranch(screen tcell.Screen, x, y int, baseAngle float64, color tcell.Color,
	amplitude, scale float64, width, height int) {

	branchChars := []rune{'·', '˙', '∘'}
	basePhase := GetBasePhase()
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Very small organic branches
	branchLength := int(float64(1+int(amplitude*2)) * scale)
	if branchLength > 3 {
		branchLength = 3
	}
	if branchLength < 1 {
		branchLength = 1
	}

	// Single subtle branch
	branchPersonality := baseAngle + float64(x+y)*0.01

	for step := 1; step <= branchLength; step++ {
		stepRatio := float64(step) / float64(branchLength)

		// Organic branch curve
		branchCurve := math.Sin(stepRatio*math.Pi*1.1+branchPersonality) * 0.2 * scale
		branchAngle := baseAngle + branchCurve + math.Sin(basePhase*1.5+branchPersonality)*0.1

		// Organic branch length
		branchRadius := float64(step) * (0.3 + math.Sin(stepRatio*math.Pi)*0.1) * scale

		branchX := x + int(branchRadius*math.Cos(branchAngle))
		branchY := y + int(branchRadius*math.Sin(branchAngle))

		if branchX >= 0 && branchX < width && branchY >= 0 && branchY < height {
			charIndex := step % len(branchChars)
			baseBranchChar := branchChars[charIndex]

			// Calculate branch intensity for intelligent character fading
			branchIntensity := (1.0 - stepRatio*0.6) * amplitude * scale * goldenRatio * 0.4

			// Intelligent character fading for branches
			var branchChar rune
			if branchIntensity < 0.1 {
				branchChar = '·'
			} else if branchIntensity < 0.2 {
				branchChar = '˙'
			} else if branchIntensity < 0.3 {
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
