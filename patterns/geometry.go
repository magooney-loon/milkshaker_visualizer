package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawGeometry creates procedural mathematical geometry loops with parametric curves
func DrawGeometry(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := math.Sqrt(float64(width*width+height*height)) / 2.2 // Use diagonal for full screen coverage

	// 3D depth layers for mathematical complexity
	numDepthLayers := 3 + int(peak*2)
	if numDepthLayers > 5 {
		numDepthLayers = 5
	}

	// Process each depth layer from back to front
	for depthLayer := numDepthLayers - 1; depthLayer >= 0; depthLayer-- {
		depthRatio := float64(depthLayer) / float64(numDepthLayers-1)
		perspectiveScale := 0.6 + depthRatio*0.5
		depthMaxRadius := maxRadius * perspectiveScale
		depthPhase := basePhase * (0.5 + depthRatio*0.4)

		// Depth-based character sets for 3D effect
		var chars []rune
		if depthLayer < 2 {
			chars = []rune{'◊', '◈', '⬟', '⬢', '⬡', '◆', '◇', '▲', '△', '□', '■', '○', '●'}
		} else if depthLayer < 4 {
			chars = []rune{'◦', '○', '⋅', '∘', '·', '˙', '▫', '▪', '▵', '▿'}
		} else {
			chars = []rune{'⋅', '·', '˙', '∘', '◦'}
		}

		// PARAMETRIC MATHEMATICAL CURVES
		drawParametricCurves(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)

		// LISSAJOUS CURVES
		drawLissajousCurves(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)

		// ROSE CURVES (Rhodonea curves)
		drawRoseCurves(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)

		// GEOMETRIC LOOPS AND CYCLES
		drawGeometricLoops(screen, centerX, centerY, depthMaxRadius, depthPhase, depthLayer,
			perspectiveScale, peak, chars, width, height)
	}
}

// drawParametricCurves creates epicycloids, hypocycloids, and other parametric mathematical curves
func drawParametricCurves(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Number of parametric curves based on peak
	numCurves := 2 + int(peak*2)
	if numCurves > 4 {
		numCurves = 4
	}

	for curveIndex := 0; curveIndex < numCurves; curveIndex++ {
		curvePersonality := float64(curveIndex)*goldenRatio*1.8 + float64(depthLayer)*2.1

		// Mathematical parameters for different curve types
		curveType := curveIndex % 3

		// Curve parameters that change with audio - scaled for full screen
		R := (0.6 + peak*0.5) * maxRadius * (0.7 + float64(curveIndex)*0.15) // Outer radius
		r := R * (0.25 + peak*0.35) * (0.5 + float64(curveIndex)*0.2)        // Inner radius
		d := r * (0.7 + peak*0.4) * (0.6 + float64(curveIndex)*0.2)          // Distance

		// Mathematical curve generation using parametric equations - outline only
		numPoints := int(20 + peak*15) // Fewer points for outline effect
		if numPoints > 35 {
			numPoints = 35
		}

		for point := 0; point < numPoints; point += 2 { // Skip points for outline gaps
			t := float64(point) / float64(numPoints) * 4 * math.Pi // Parameter t

			// Time-based parameter modulation
			tModulated := t + phase*(0.3+float64(curveIndex)*0.1)

			var x, y float64

			// Different mathematical curve types
			switch curveType {
			case 0: // Epicycloid: (R+r)*cos(t) - d*cos((R+r)/r * t)
				rRatio := (R + r) / r
				if math.Abs(r) < 0.001 { // Safety check
					r = 0.001
				}
				x = (R+r)*math.Cos(tModulated) - d*math.Cos(rRatio*tModulated)
				y = (R+r)*math.Sin(tModulated) - d*math.Sin(rRatio*tModulated)

			case 1: // Hypocycloid: (R-r)*cos(t) + d*cos((R-r)/r * t)
				rRatio := (R - r) / r
				if math.Abs(r) < 0.001 { // Safety check
					r = 0.001
				}
				x = (R-r)*math.Cos(tModulated) + d*math.Cos(rRatio*tModulated)
				y = (R-r)*math.Sin(tModulated) - d*math.Sin(rRatio*tModulated)

			case 2: // Trochoid variation
				freq := 1.0 + float64(curveIndex)*0.5 + peak*0.3
				x = R*math.Cos(tModulated*freq) + r*math.Cos(tModulated*freq*goldenRatio)
				y = R*math.Sin(tModulated*freq) + r*math.Sin(tModulated*freq*goldenRatio)
			}

			// Safety checks for mathematical calculations
			if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
				continue
			}

			// Apply perspective scaling and breathing
			breathe := 1 + 0.04*math.Sin(phase*1.8+curvePersonality)*scale
			x *= scale * breathe
			y *= scale * breathe

			// Add organic micro-variations
			microX := 0.5 * math.Sin(phase*3.2+t*0.5+curvePersonality) * scale
			microY := 0.5 * math.Cos(phase*3.7+t*0.6+curvePersonality) * scale

			finalX := centerX + int(x+microX)
			finalY := centerY + int(y+microY)

			// Bounds check
			if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {
				// Calculate curve intensity for intelligent character fading
				distanceFromCenter := math.Sqrt(x*x + y*y)
				normalizedDistance := distanceFromCenter / maxRadius
				curveIntensity := peak * (1.0 - normalizedDistance*0.6) * scale * (0.5 + math.Sin(t*2+phase)*0.3)

				// Mathematical character selection based on curve properties
				// Outline character selection - only use outline chars
				outlineChars := []rune{'·', '˙', '∘', '◦', '○'}
				charPhase := curvePersonality + t*1.5 + distanceFromCenter*0.02
				charIndex := int(math.Abs(charPhase)*goldenRatio) % len(outlineChars)

				// Simple outline character selection
				var curveChar rune
				if curveIntensity < 0.15 {
					curveChar = '·'
				} else if curveIntensity < 0.35 {
					curveChar = '∘'
				} else {
					curveChar = outlineChars[charIndex]
				}

				// Mathematical color based on curve parameters
				hue := math.Mod(curvePersonality*0.2+t*0.1+phase*0.05, 1)
				saturation := (0.4 + peak*0.3) * (0.7 + normalizedDistance*0.3) * scale
				saturation = math.Max(0.1, math.Min(0.8, saturation))

				value := (0.5 + peak*0.3) * (0.6 + (1.0-normalizedDistance)*0.4) * scale
				value = math.Max(0.2, math.Min(0.9, value))

				curveColor := HSVToRGB(hue, saturation, value)

				if curveIntensity > 0.12 { // Higher threshold for outline
					screen.SetContent(finalX, finalY, curveChar, nil, tcell.StyleDefault.Foreground(curveColor))
				}
			}
		}
	}
}

// drawLissajousCurves creates Lissajous curves using parametric equations
func drawLissajousCurves(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Number of Lissajous curves
	numCurves := 1 + int(peak*2)
	if numCurves > 3 {
		numCurves = 3
	}

	for curveIndex := 0; curveIndex < numCurves; curveIndex++ {
		curvePersonality := float64(curveIndex)*goldenRatio*2.3 + float64(depthLayer)*1.7

		// Lissajous parameters: x = A*sin(at + δ), y = B*sin(bt) - full screen scaling
		A := maxRadius * scale * (0.7 + peak*0.4)
		B := maxRadius * scale * (0.6 + peak*0.5)

		// Frequency ratios - use musical ratios for pleasing patterns
		aFreq := 1.0 + float64(curveIndex)*0.5 + peak*0.3
		bFreq := goldenRatio + float64(curveIndex)*0.3 + peak*0.2

		// Phase difference
		delta := curvePersonality + phase*0.4

		numPoints := int(25 + peak*20) // Fewer points for outline
		if numPoints > 45 {
			numPoints = 45
		}

		for point := 0; point < numPoints; point += 2 { // Skip points for outline gaps
			t := float64(point) / float64(numPoints) * 4 * math.Pi
			tModulated := t + phase*0.5

			// Lissajous equations
			x := A * math.Sin(aFreq*tModulated+delta)
			y := B * math.Sin(bFreq*tModulated)

			// Safety checks
			if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
				continue
			}

			// Organic breathing effect
			breathe := 1 + 0.03*math.Sin(phase*2.1+curvePersonality+t*0.3)*scale

			finalX := centerX + int(x*breathe)
			finalY := centerY + int(y*breathe)

			if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {
				// Calculate Lissajous intensity
				normalizedRadius := math.Sqrt(x*x+y*y) / maxRadius
				lissajousIntensity := peak * (1.0 - normalizedRadius*0.5) * scale

				// Outline character selection only
				outlineChars := []rune{'·', '∘', '◦', '○'}
				charPhase := curvePersonality + t*0.8 + normalizedRadius*2
				charIndex := int(math.Abs(charPhase)) % len(outlineChars)

				// Simple outline character selection
				var lissChar rune
				if lissajousIntensity < 0.2 {
					lissChar = '·'
				} else if lissajousIntensity < 0.4 {
					lissChar = '∘'
				} else {
					lissChar = outlineChars[charIndex]
				}

				// Mathematical color
				hue := math.Mod(curvePersonality*0.15+t*0.05+phase*0.08, 1)
				saturation := (0.3 + peak*0.25) * scale
				saturation = math.Max(0.08, math.Min(0.7, saturation))

				value := (0.4 + peak*0.25) * (1.0 - normalizedRadius*0.3) * scale
				value = math.Max(0.15, math.Min(0.8, value))

				lissColor := HSVToRGB(hue, saturation, value)

				if lissajousIntensity > 0.15 { // Higher threshold for outline
					screen.SetContent(finalX, finalY, lissChar, nil, tcell.StyleDefault.Foreground(lissColor))
				}
			}
		}
	}
}

// drawRoseCurves creates rose curves (rhodonea curves) r = cos(k*θ)
func drawRoseCurves(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Number of rose curves
	numRoses := 1 + int(peak*1.5)
	if numRoses > 3 {
		numRoses = 3
	}

	for roseIndex := 0; roseIndex < numRoses; roseIndex++ {
		rosePersonality := float64(roseIndex)*goldenRatio*1.9 + float64(depthLayer)*2.3

		// Rose curve parameter k determines number of petals
		// k = n/d where n,d are integers gives n petals if n is odd, 2n if n is even
		k := 2.0 + float64(roseIndex) + peak*2.0 + math.Sin(phase*0.3+rosePersonality)*0.5

		// Rose curve radius scaling - full screen coverage
		roseRadius := maxRadius * scale * (0.8 + peak*0.4)

		numPoints := int(30 + peak*25) // Fewer points for outline roses
		if numPoints > 55 {
			numPoints = 55
		}

		for point := 0; point < numPoints; point += 2 { // Skip points for petal outlines
			theta := float64(point) / float64(numPoints) * 4 * math.Pi
			thetaModulated := theta + phase*0.6

			// Rose curve equation: r = a * cos(k * θ)
			r := roseRadius * math.Abs(math.Cos(k*thetaModulated))

			// Convert to Cartesian coordinates
			x := r * math.Cos(thetaModulated)
			y := r * math.Sin(thetaModulated)

			// Safety checks
			if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
				continue
			}

			// Organic pulsing
			pulse := 1 + 0.05*math.Sin(phase*1.9+rosePersonality+theta*0.4)*scale

			finalX := centerX + int(x*pulse)
			finalY := centerY + int(y*pulse)

			if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {
				// Calculate rose intensity
				normalizedRadius := r / roseRadius
				roseIntensity := peak * normalizedRadius * scale * (0.6 + math.Sin(k*thetaModulated)*0.4)

				// Rose outline character selection
				roseOutlineChars := []rune{'·', '∘', '◦', '●'}
				charPhase := rosePersonality + theta*0.6 + r*0.01
				charIndex := int(math.Abs(charPhase)*goldenRatio) % len(roseOutlineChars)

				// Simple rose outline character selection
				var roseChar rune
				if roseIntensity < 0.18 {
					roseChar = '·'
				} else if roseIntensity < 0.4 {
					roseChar = '∘'
				} else {
					roseChar = roseOutlineChars[charIndex]
				}

				// Rose-specific coloring
				hue := math.Mod(rosePersonality*0.25+theta*0.08+phase*0.04, 1)
				saturation := (0.35 + peak*0.3) * normalizedRadius * scale
				saturation = math.Max(0.1, math.Min(0.75, saturation))

				value := (0.45 + peak*0.3) * normalizedRadius * scale
				value = math.Max(0.18, math.Min(0.85, value))

				roseColor := HSVToRGB(hue, saturation, value)

				if roseIntensity > 0.15 { // Higher threshold for rose outline
					screen.SetContent(finalX, finalY, roseChar, nil, tcell.StyleDefault.Foreground(roseColor))
				}
			}
		}
	}
}

// drawGeometricLoops creates mathematical loops and cycles
func drawGeometricLoops(screen tcell.Screen, centerX, centerY int, maxRadius, phase float64,
	depthLayer int, scale, peak float64, chars []rune, width, height int) {

	goldenRatio := (1 + math.Sqrt(5)) / 2

	// Number of geometric loops
	numLoops := 1 + int(peak*2)
	if numLoops > 3 {
		numLoops = 3
	}

	for loopIndex := 0; loopIndex < numLoops; loopIndex++ {
		loopPersonality := float64(loopIndex)*goldenRatio*2.7 + float64(depthLayer)*1.9

		loopType := loopIndex % 4
		loopRadius := maxRadius * scale * (0.7 + peak*0.5) * (0.7 + float64(loopIndex)*0.15)

		numPoints := int(24 + peak*20) // Fewer points for loop outlines
		if numPoints > 44 {
			numPoints = 44
		}

		for point := 0; point < numPoints; point += 2 { // Skip points for loop outlines
			t := float64(point) / float64(numPoints) * 2 * math.Pi
			tModulated := t + phase*(0.4+float64(loopIndex)*0.1)

			var x, y float64

			// Different mathematical loop types
			switch loopType {
			case 0: // Figure-8 loop (lemniscate): x = a*cos(t)/(1+sin²(t)), y = a*sin(t)*cos(t)/(1+sin²(t))
				denominator := 1 + math.Pow(math.Sin(tModulated), 2)
				if math.Abs(denominator) < 0.001 {
					denominator = 0.001
				}
				x = loopRadius * math.Cos(tModulated) / denominator
				y = loopRadius * math.Sin(tModulated) * math.Cos(tModulated) / denominator

			case 1: // Cardioid: r = a(1 + cos(θ))
				r := loopRadius * 0.5 * (1 + math.Cos(tModulated))
				x = r * math.Cos(tModulated)
				y = r * math.Sin(tModulated)

			case 2: // Limacon: r = a + b*cos(θ)
				a := loopRadius * 0.4
				b := loopRadius * 0.3 * (0.5 + peak*0.5)
				r := a + b*math.Cos(tModulated)
				x = r * math.Cos(tModulated)
				y = r * math.Sin(tModulated)

			case 3: // Folium: x³ + y³ = 3axy, parametric form
				denominator := 1 + math.Pow(math.Tan(tModulated), 3)
				if math.Abs(denominator) < 0.001 {
					denominator = 0.001
				}
				x = loopRadius * 0.8 * math.Tan(tModulated) / denominator // Reduced multiplier for screen bounds
				y = loopRadius * 0.8 * math.Pow(math.Tan(tModulated), 2) / denominator
			}

			// Safety checks
			if math.IsNaN(x) || math.IsInf(x, 0) || math.IsNaN(y) || math.IsInf(y, 0) {
				continue
			}

			// Mathematical breathing and modulation
			modulation := 1 + 0.06*math.Sin(phase*2.2+loopPersonality+t*0.8)*scale

			finalX := centerX + int(x*modulation)
			finalY := centerY + int(y*modulation)

			if finalX >= 0 && finalX < width && finalY >= 0 && finalY < height {
				// Calculate loop intensity
				distance := math.Sqrt(x*x + y*y)
				normalizedDistance := distance / loopRadius
				loopIntensity := peak * (0.8 + math.Sin(tModulated*2)*0.4) * scale * (1.0 - normalizedDistance*0.2)

				// Loop outline character selection
				loopOutlineChars := []rune{'·', '∘', '◦', '◊'}
				charPhase := loopPersonality + t*0.9 + distance*0.015
				charIndex := int(math.Abs(charPhase)) % len(loopOutlineChars)

				// Simple loop outline character selection
				var loopChar rune
				if loopIntensity < 0.2 {
					loopChar = '·'
				} else if loopIntensity < 0.45 {
					loopChar = '∘'
				} else {
					loopChar = loopOutlineChars[charIndex]
				}

				// Mathematical coloring
				hue := math.Mod(loopPersonality*0.18+t*0.06+phase*0.06, 1)
				saturation := (0.3 + peak*0.25) * (0.8 + normalizedDistance*0.2) * scale
				saturation = math.Max(0.08, math.Min(0.7, saturation))

				value := (0.4 + peak*0.25) * scale
				value = math.Max(0.15, math.Min(0.8, value))

				loopColor := HSVToRGB(hue, saturation, value)

				if loopIntensity > 0.18 { // Higher threshold for loop outline
					screen.SetContent(finalX, finalY, loopChar, nil, tcell.StyleDefault.Foreground(loopColor))
				}
			}
		}
	}
}
