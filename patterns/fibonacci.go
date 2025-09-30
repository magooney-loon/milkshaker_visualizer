package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawFibonacci creates a clean mathematical fibonacci spiral with terminal aesthetic
func DrawFibonacci(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()

	// Golden ratio and angle for perfect fibonacci spiral
	goldenRatio := (1 + math.Sqrt(5)) / 2
	goldenAngle := 2 * math.Pi / (goldenRatio * goldenRatio) // ~2.4 radians

	// Scale based on screen size and peak
	maxRadius := math.Min(float64(width), float64(height)) / 3
	peakScale := 0.4 + peak*0.8 // Scales from 40% to 120%

	// Number of fibonacci terms to draw
	maxTerms := 15 + int(peak*10)
	if maxTerms > 30 {
		maxTerms = 30
	}

	// Generate fibonacci sequence
	fib := make([]int, maxTerms)
	if maxTerms >= 1 {
		fib[0] = 1
	}
	if maxTerms >= 2 {
		fib[1] = 1
	}
	for i := 2; i < maxTerms; i++ {
		fib[i] = fib[i-1] + fib[i-2]
		// Cap to prevent overflow
		if fib[i] > 1000 {
			fib[i] = 1000
		}
	}

	// Terminal-style characters for spiral
	spiralChars := []rune{'·', '∘', '○', '●', '◉', '⬢', '◆', '★'}

	// Draw multiple spiral arms for fullness
	numArms := 2 + int(peak*2)
	if numArms > 5 {
		numArms = 5
	}

	for arm := 0; arm < numArms; arm++ {
		armOffset := float64(arm) * 2 * math.Pi / float64(numArms)
		armPhase := basePhase*0.1 + armOffset

		// Draw fibonacci spiral points
		for i := 3; i < len(fib); i++ { // Start from 3 to avoid center clutter
			if fib[i] == 0 {
				continue
			}

			// Calculate spiral position using fibonacci numbers
			fibRadius := math.Sqrt(float64(fib[i])) * maxRadius * peakScale / 10
			spiralAngle := float64(i)*goldenAngle + armPhase

			// Individual point animation
			pointPhase := armPhase + float64(i)*0.3
			radiusVariation := 1 + 0.15*math.Sin(pointPhase)
			finalRadius := fibRadius * radiusVariation

			// Calculate number of points to draw along this fibonacci arc
			arcPoints := int(finalRadius / 3)
			if arcPoints < 2 {
				arcPoints = 2
			}
			if arcPoints > 15 {
				arcPoints = 15
			}

			// Draw arc of points for this fibonacci term
			for point := 0; point < arcPoints; point++ {
				pointRatio := float64(point) / float64(arcPoints)

				// Interpolate radius from previous fibonacci number
				var prevRadius float64
				if i > 3 {
					prevRadius = math.Sqrt(float64(fib[i-1])) * maxRadius * peakScale / 10
				} else {
					prevRadius = 0
				}

				interpRadius := prevRadius + (finalRadius-prevRadius)*pointRatio
				interpAngle := spiralAngle - goldenAngle*(1-pointRatio)

				x := centerX + int(interpRadius*math.Cos(interpAngle))
				y := centerY + int(interpRadius*math.Sin(interpAngle))

				if x >= 0 && x < width && y >= 0 && y < height {
					// Calculate intensity based on fibonacci position and peak
					fibIntensity := (1.0 - pointRatio*0.5) * peakScale
					termIntensity := 1.0 - float64(i)/float64(len(fib))*0.6
					totalIntensity := fibIntensity * termIntensity * (0.5 + peak*0.5)

					// Character selection based on intensity and position
					charPhase := float64(i)*goldenRatio + float64(point)*0.5 + armOffset
					charIndex := int(math.Abs(charPhase)) % len(spiralChars)
					baseChar := spiralChars[charIndex]

					// Intelligent character progression
					var finalChar rune
					if totalIntensity < 0.15 {
						finalChar = '·'
					} else if totalIntensity < 0.3 {
						finalChar = '∘'
					} else if totalIntensity < 0.5 {
						finalChar = '○'
					} else if totalIntensity < 0.7 {
						finalChar = '●'
					} else {
						finalChar = baseChar
					}

					// Mathematical color based on golden ratio
					hueBase := float64(i)/float64(len(fib))*goldenRatio + armOffset/(2*math.Pi)
					hueShift := math.Sin(armPhase*0.5+float64(i)*0.2) * 0.05
					hue := math.Mod(hueBase+hueShift, 1)

					// Clean saturation and brightness
					saturation := 0.6 + peak*0.3*totalIntensity
					brightness := 0.4 + totalIntensity*0.5 + peak*0.2

					spiralColor := HSVToRGB(hue, saturation, brightness)

					if totalIntensity > 0.1 {
						screen.SetContent(x, y, finalChar, nil, tcell.StyleDefault.Foreground(spiralColor))
					}
				}
			}

			// Draw connecting lines at high peaks for mathematical beauty
			if peak > 0.6 && i > 3 && arm == 0 {
				prevRadius := math.Sqrt(float64(fib[i-1])) * maxRadius * peakScale / 10
				prevAngle := float64(i-1)*goldenAngle + armPhase

				// Draw line from previous fibonacci point to current
				startX := centerX + int(prevRadius*math.Cos(prevAngle))
				startY := centerY + int(prevRadius*math.Sin(prevAngle))
				endX := centerX + int(finalRadius*math.Cos(spiralAngle))
				endY := centerY + int(finalRadius*math.Sin(spiralAngle))

				drawFibonacciLine(screen, startX, startY, endX, endY, width, height, peak)
			}
		}
	}

	// Draw golden ratio rectangles at very high peaks
	if peak > 0.8 {
		drawGoldenRectangles(screen, centerX, centerY, maxRadius*peakScale, basePhase, peak, width, height)
	}

	// Central fibonacci core
	coreRadius := 1 + int(peak*3)
	if coreRadius > 5 {
		coreRadius = 5
	}

	coreChars := []rune{'∘', '○', '●', '◉', '⬢'}

	for radius := 1; radius <= coreRadius; radius++ {
		coreIntensity := (1.0 - float64(radius-1)/float64(coreRadius)) * peak
		if coreIntensity > 0.2 {
			// Draw core ring based on fibonacci numbers
			fibPoints := int(float64(radius) * goldenRatio * 2)
			for point := 0; point < fibPoints; point++ {
				angle := float64(point) * 2 * math.Pi / float64(fibPoints)
				angle += basePhase * 0.05 // Slow rotation

				x := centerX + int(float64(radius)*math.Cos(angle))
				y := centerY + int(float64(radius)*math.Sin(angle))

				if x >= 0 && x < width && y >= 0 && y < height {
					charIndex := int(coreIntensity * float64(len(coreChars)-1))
					if charIndex >= len(coreChars) {
						charIndex = len(coreChars) - 1
					}
					coreChar := coreChars[charIndex]

					coreHue := math.Mod(basePhase*0.02+float64(radius)*0.1, 1)
					coreColor := HSVToRGB(coreHue, 0.8, coreIntensity)

					screen.SetContent(x, y, coreChar, nil, tcell.StyleDefault.Foreground(coreColor))
				}
			}
		}
	}
}

// drawFibonacciLine draws a connecting line between fibonacci points
func drawFibonacciLine(screen tcell.Screen, x1, y1, x2, y2, width, height int, peak float64) {
	lineChars := []rune{'·', '˙', '∘', '─', '│', '╱', '╲'}

	dx := x2 - x1
	dy := y2 - y1
	steps := int(math.Sqrt(float64(dx*dx + dy*dy)))

	if steps < 2 {
		return
	}

	for step := 1; step < steps; step++ {
		t := float64(step) / float64(steps)
		x := x1 + int(float64(dx)*t)
		y := y1 + int(float64(dy)*t)

		if x >= 0 && x < width && y >= 0 && y < height {
			intensity := (1.0 - t*0.5) * (peak - 0.6) * 2.5

			if intensity > 0.3 {
				// Choose line character based on direction
				var lineChar rune
				if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
					lineChar = lineChars[3] // horizontal
				} else {
					lineChar = lineChars[4] // vertical
				}

				if intensity < 0.5 {
					lineChar = lineChars[0] // dot
				}

				lineHue := math.Mod(t*0.5, 1)
				lineColor := HSVToRGB(lineHue, 0.6, intensity*0.7)

				screen.SetContent(x, y, lineChar, nil, tcell.StyleDefault.Foreground(lineColor))
			}
		}
	}
}

// drawGoldenRectangles draws golden ratio rectangles at high peaks
func drawGoldenRectangles(screen tcell.Screen, centerX, centerY int, maxSize float64, phase, peak float64, width, height int) {
	goldenRatio := (1 + math.Sqrt(5)) / 2
	rectChars := []rune{'┌', '┐', '└', '┘', '├', '┤', '┬', '┴', '│', '─'}

	numRects := 2 + int((peak-0.8)*10)
	if numRects > 4 {
		numRects = 4
	}

	for rect := 0; rect < numRects; rect++ {
		rectPhase := phase*0.1 + float64(rect)*0.5
		size := maxSize * (0.3 + float64(rect)*0.2) * (peak - 0.6) * 2

		// Golden ratio rectangle dimensions
		rectWidth := int(size)
		rectHeight := int(size / goldenRatio)

		// Rotation for organic feel
		rotation := math.Sin(rectPhase) * 0.2

		// Draw rectangle outline
		for side := 0; side < 4; side++ {
			var startX, startY, endX, endY int

			switch side {
			case 0: // Top
				startX, startY = centerX-rectWidth/2, centerY-rectHeight/2
				endX, endY = centerX+rectWidth/2, centerY-rectHeight/2
			case 1: // Right
				startX, startY = centerX+rectWidth/2, centerY-rectHeight/2
				endX, endY = centerX+rectWidth/2, centerY+rectHeight/2
			case 2: // Bottom
				startX, startY = centerX+rectWidth/2, centerY+rectHeight/2
				endX, endY = centerX-rectWidth/2, centerY+rectHeight/2
			case 3: // Left
				startX, startY = centerX-rectWidth/2, centerY+rectHeight/2
				endX, endY = centerX-rectWidth/2, centerY-rectHeight/2
			}

			// Apply rotation
			cos, sin := math.Cos(rotation), math.Sin(rotation)
			rotateX := func(x, y int) int {
				fx, fy := float64(x-centerX), float64(y-centerY)
				return centerX + int(fx*cos-fy*sin)
			}
			rotateY := func(x, y int) int {
				fx, fy := float64(x-centerX), float64(y-centerY)
				return centerY + int(fx*sin+fy*cos)
			}

			rotStartX, rotStartY := rotateX(startX, startY), rotateY(startX, startY)
			rotEndX, rotEndY := rotateX(endX, endY), rotateY(endX, endY)

			// Draw line
			dx := rotEndX - rotStartX
			dy := rotEndY - rotStartY
			steps := int(math.Sqrt(float64(dx*dx + dy*dy)))

			for step := 0; step <= steps; step++ {
				if steps == 0 {
					continue
				}
				t := float64(step) / float64(steps)
				x := rotStartX + int(float64(dx)*t)
				y := rotStartY + int(float64(dy)*t)

				if x >= 0 && x < width && y >= 0 && y < height {
					rectChar := rectChars[8] // vertical line
					if math.Abs(float64(dx)) > math.Abs(float64(dy)) {
						rectChar = rectChars[9] // horizontal line
					}

					rectIntensity := (peak - 0.8) * 5 * (1.0 - float64(rect)*0.2)
					rectHue := math.Mod(phase*0.03+float64(rect)*0.25, 1)
					rectColor := HSVToRGB(rectHue, 0.7, rectIntensity)

					if rectIntensity > 0.3 {
						screen.SetContent(x, y, rectChar, nil, tcell.StyleDefault.Foreground(rectColor))
					}
				}
			}
		}
	}
}
