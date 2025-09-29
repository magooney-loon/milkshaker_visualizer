package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawField creates an organic 3D field that connects and fills depth between patterns
func DrawField(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Very subtle field characters for organic background feel
	fieldChars := []rune{'⋅', '·', '˙', '∘', '◦', '⁚', '⁛', '⁝'}

	// Create multiple depth layers for 3D feel
	numDepthLayers := 3 + int(peak*2)
	if numDepthLayers > 5 {
		numDepthLayers = 5
	}

	for depthLayer := 0; depthLayer < numDepthLayers; depthLayer++ {
		layerDepth := float64(depthLayer) / float64(numDepthLayers)
		layerPhase := basePhase * (0.2 + layerDepth*0.3) // Slower for background layers
		layerScale := 1.0 - layerDepth*0.3               // Smaller for distant layers

		// Organic field grid with golden ratio spacing
		gridSpacing := int(12 + layerDepth*8) // Sparser for background layers

		for x := gridSpacing; x < width-gridSpacing; x += gridSpacing {
			for y := gridSpacing; y < height-gridSpacing; y += gridSpacing {

				// Calculate distance from center for organic field influence
				dx := float64(x - centerX)
				dy := float64(y - centerY)
				distanceFromCenter := math.Sqrt(dx*dx + dy*dy)
				maxDistance := math.Sqrt(float64(width*width+height*height)) / 2
				centerInfluence := 1.0 - (distanceFromCenter / maxDistance)

				// Organic field calculations using perlin-noise-like functions
				fieldX := float64(x) * 0.02
				fieldY := float64(y) * 0.02

				// Multiple octaves of organic noise for depth
				noise1 := math.Sin(fieldX+layerPhase*0.8) * math.Cos(fieldY+layerPhase*0.6)
				noise2 := 0.5 * math.Sin(fieldX*2+layerPhase*1.2) * math.Cos(fieldY*2+layerPhase*0.9)
				noise3 := 0.25 * math.Sin(fieldX*4+layerPhase*1.5) * math.Cos(fieldY*4+layerPhase*1.1)
				organicNoise := noise1 + noise2 + noise3

				// Golden ratio influence for natural distribution
				goldenPhase := (fieldX+fieldY)*goldenRatio + layerPhase
				goldenInfluence := math.Sin(goldenPhase) * 0.3

				// Combine all organic influences
				fieldStrength := (organicNoise + goldenInfluence) * centerInfluence * layerScale * peak

				// Organic breathing effect that varies across the field
				breathePhase := layerPhase*1.3 + fieldX*0.1 + fieldY*0.1
				breathe := 1 + 0.1*math.Sin(breathePhase)*centerInfluence

				// Final field strength with breathing
				finalStrength := fieldStrength * breathe

				// Only draw if field strength is above threshold (creates organic gaps)
				strengthThreshold := 0.1 + layerDepth*0.1
				if math.Abs(finalStrength) > strengthThreshold {

					// Organic position offset for non-grid-like feel
					offsetX := int(finalStrength * 3)
					offsetY := int(goldenInfluence * 2)
					finalX := x + offsetX
					finalY := y + offsetY

					// Bounds check
					if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {

						// Organic character selection based on field properties
						charPhase := finalStrength*5 + goldenPhase + float64(depthLayer)*2.1
						charIndex := int(math.Abs(charPhase)*goldenRatio) % len(fieldChars)
						fieldChar := fieldChars[charIndex]

						// Very subtle colors that create depth
						hue := math.Mod(goldenPhase*0.1+layerPhase*0.05, 1)

						// Background layers are more muted
						saturation := (0.1 + peak*0.2) * (1 - layerDepth*0.3)
						saturation = math.Max(0.05, math.Min(0.4, saturation))

						// Depth-based brightness
						value := (0.2 + peak*0.2 + math.Abs(finalStrength)*0.1) * (1 - layerDepth*0.4)
						value = math.Max(0.1, math.Min(0.5, value))

						fieldColor := HSVToRGB(hue, saturation, value)

						// Extra subtle for background layers
						if layerDepth > 0.5 && math.Abs(finalStrength) < 0.3 {
							fieldChar = '·'
						}

						screen.SetContent(finalX, finalY, fieldChar, nil, tcell.StyleDefault.Foreground(fieldColor))

						// Add organic connecting tendrils between field points
						if math.Abs(finalStrength) > 0.3 && peak > 0.4 {
							drawFieldTendril(screen, finalX, finalY, finalStrength, goldenPhase, fieldColor, width, height)
						}
					}
				}
			}
		}

		// Add organic flowing lines that connect across the field
		if peak > 0.3 {
			drawFieldFlows(screen, centerX, centerY, width, height, layerPhase, layerDepth, peak)
		}
	}
}

// drawFieldTendril creates subtle organic connections between field points
func drawFieldTendril(screen tcell.Screen, x, y int, strength, phase float64, color tcell.Color, width, height int) {
	tendrilChars := []rune{'⋅', '·', '˙'}
	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Very short, subtle tendrils
	tendrilLength := 2 + int(math.Abs(strength)*2)
	if tendrilLength > 4 {
		tendrilLength = 4
	}

	// Organic tendril direction
	tendrilAngle := phase*goldenRatio + math.Sin(GetBasePhase()*0.5)*0.3

	for step := 1; step <= tendrilLength; step++ {
		// Organic curve in tendril
		stepPhase := float64(step) / float64(tendrilLength)
		curve := math.Sin(stepPhase*math.Pi) * 0.5

		tendrilX := x + int(float64(step)*math.Cos(tendrilAngle)+curve*math.Sin(tendrilAngle))
		tendrilY := y + int(float64(step)*math.Sin(tendrilAngle)-curve*math.Cos(tendrilAngle))

		if tendrilX >= 0 && tendrilX < width && tendrilY >= 0 && tendrilY < height {
			charIndex := step % len(tendrilChars)
			tendrilChar := tendrilChars[charIndex]

			// Fade along tendril
			intensity := 1.0 - stepPhase
			if intensity > 0.5 {
				screen.SetContent(tendrilX, tendrilY, tendrilChar, nil, tcell.StyleDefault.Foreground(color))
			}
		}
	}
}

// drawFieldFlows creates flowing organic lines across the field for connectivity
func drawFieldFlows(screen tcell.Screen, centerX, centerY, width, height int, phase, depth, peak float64) {
	goldenRatio := (1 + math.Sqrt(5)) / 2

	flowChars := []rune{'⋅', '·', '˙', '∘'}

	// Number of flow lines based on peak
	numFlows := 2 + int(peak*3)
	if numFlows > 5 {
		numFlows = 5
	}

	for flow := 0; flow < numFlows; flow++ {
		flowAngle := float64(flow) * 2 * math.Pi / float64(numFlows)
		flowPhase := phase * (0.3 + float64(flow)*0.1)

		// Start from center area
		startRadius := 20.0 + float64(flow)*10

		// Flow outward with organic curves
		maxFlowRadius := float64(Min(width, height)) / 3

		for radius := startRadius; radius < maxFlowRadius; radius += 3 + peak*2 {
			// Organic flow curve
			flowCurve := math.Sin(radius*0.05+flowPhase*2) * 0.2
			organicAngle := flowAngle + flowPhase*0.3 + flowCurve

			flowX := centerX + int(radius*math.Cos(organicAngle))
			flowY := centerY + int(radius*math.Sin(organicAngle))

			if flowX >= 0 && flowX < width && flowY >= 0 && flowY < height {
				charIndex := (int(radius) + flow) % len(flowChars)
				flowChar := flowChars[charIndex]

				// Very subtle flow colors
				hue := math.Mod(float64(flow)/float64(numFlows)*goldenRatio+phase*0.1, 1)
				saturation := 0.1 + peak*0.1
				value := 0.15 + peak*0.1 - depth*0.05
				value = math.Max(0.05, math.Min(0.3, value))

				flowColor := HSVToRGB(hue, saturation, value)

				// Distance fade
				distanceRatio := (radius - startRadius) / (maxFlowRadius - startRadius)
				if distanceRatio < 0.8 {
					screen.SetContent(flowX, flowY, flowChar, nil, tcell.StyleDefault.Foreground(flowColor))
				}
			}
		}
	}
}
