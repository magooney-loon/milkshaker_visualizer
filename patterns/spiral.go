package patterns

import (
	"math"

	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawSpiral creates dynamic multi-armed rotating spirals
func DrawSpiral(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2

	// Dynamic number of spiral arms based on peak
	numArms := 3 + int(peak*4)

	for arm := 0; arm < numArms; arm++ {
		armOffset := float64(arm) * 2 * math.Pi / float64(numArms)
		armRotation := basePhase * 0.5 * float64(1+arm%2*2-1) // Alternate rotation directions

		// Each arm has its own characteristics
		armAmplitude := (2.0 + float64(arm)*0.5) * peak
		armFrequency := 0.3 + 0.1*float64(arm) + 0.2*peak

		radius := 1.0 + float64(arm)*2
		angle := armOffset + armRotation
		angleStep := 0.08 + peak*0.05

		// Multi-layered spiral per arm
		for layer := 0; layer < 2+int(peak*2); layer++ {
			layerRadius := radius
			layerAngle := angle + float64(layer)*0.3

			for layerRadius < maxRadius {
				// Complex wave function combining multiple frequencies
				wave1 := armAmplitude * math.Sin(armFrequency*layerAngle+basePhase)
				wave2 := armAmplitude * 0.5 * math.Cos(armFrequency*2*layerAngle+basePhase*1.5)
				wave3 := armAmplitude * 0.3 * math.Sin(armFrequency*0.5*layerAngle+basePhase*0.7)
				waveOffset := wave1 + wave2 + wave3

				// Pulsing radius effect
				pulseRadius := layerRadius * (1 + 0.2*math.Sin(basePhase*2+float64(arm)*0.5))

				finalRadius := pulseRadius + waveOffset

				x := centerX + int(finalRadius*math.Cos(layerAngle))
				y := centerY + int(finalRadius*math.Sin(layerAngle))

				// Dynamic character selection based on layer and distance
				chars := []rune{'•', '◦', '○', '◎', '◉', '⚬', '⚭', '⚮', '*', '✦', '✧', '✩', '✪', '✫', '✬', '✭', '✮', '✯', '✰', '✱'}
				charIndex := (layer*arm + int(layerRadius)) % len(chars)
				displayChar := chars[charIndex]

				// Color variation based on arm and layer
				hue := float64(arm)/float64(numArms) + basePhase*0.1
				saturation := 0.7 + peak*0.3
				value := 0.6 + peak*0.4 - float64(layer)*0.1
				armColor := HSVToRGB(math.Mod(hue, 1), saturation, value)

				screen.SetContent(x, y, displayChar, nil, tcell.StyleDefault.Foreground(armColor))

				layerRadius += 0.8 + peak*0.5
				layerAngle += angleStep
			}
		}
	}
}
