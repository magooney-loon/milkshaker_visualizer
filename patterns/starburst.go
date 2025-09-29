package patterns

import (
	"math"
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// DrawStarburst creates dynamic rotating rays with branching effects
func DrawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := GetBasePhase()
	maxRadius := float64(Min(width, height)) / 2

	// Dynamic number of rays based on peak intensity
	numRays := 12 + int(peak*16)
	rayAngleStep := 2 * math.Pi / float64(numRays)

	for rayIndex := 0; rayIndex < numRays; rayIndex++ {
		baseAngle := float64(rayIndex) * rayAngleStep

		// Each ray rotates at different speeds
		rayRotation := basePhase * (0.5 + float64(rayIndex%3)*0.3)
		finalAngle := baseAngle + rayRotation

		// Multiple beams per ray for thickness effect
		beamCount := 1 + int(peak*3)
		for beam := 0; beam < beamCount; beam++ {
			beamAngle := finalAngle + (float64(beam)-float64(beamCount)/2)*0.05

			// Variable ray length with pulsing
			rayLength := maxRadius * (0.6 + 0.4*math.Sin(basePhase*3+float64(rayIndex)*0.2))

			// Dynamic step size for ray density
			stepSize := 0.8 + peak*0.7

			for radius := 2.0; radius < rayLength; radius += stepSize {
				// Complex wave patterns along each ray
				distancePhase := radius * 0.1
				wave1 := 8.0 * peak * math.Sin(distancePhase+basePhase*2)
				wave2 := 4.0 * peak * math.Cos(distancePhase*1.5+basePhase*1.3)
				wave3 := 2.0 * peak * math.Sin(distancePhase*3+basePhase*0.8)

				// Branching effect - rays can split
				branchOffset := wave1 + wave2 + wave3

				// Main ray
				mainX := centerX + int((radius+branchOffset)*math.Cos(beamAngle))
				mainY := centerY + int((radius+branchOffset)*math.Sin(beamAngle))

				// Distance-based character selection
				chars := []rune{'∙', '•', '●', '◉', '⬢', '⬡', '◆', '◇', '★', '☆', '✦', '✧', '✩', '✪', '✫', '✬', '✭', '✮', '✯', '✰'}
				charIndex := (rayIndex + int(radius*2) + beam) % len(chars)
				rayChar := chars[charIndex]

				// Dynamic color based on distance and ray index
				colorPhase := float64(rayIndex)/float64(numRays) + radius*0.01 + basePhase*0.2
				hue := math.Mod(colorPhase, 1)
				saturation := 0.8 + peak*0.2
				value := 0.9 - radius/maxRadius*0.4 + peak*0.1
				rayColor := HSVToRGB(hue, saturation, math.Max(0.1, value))

				screen.SetContent(mainX, mainY, rayChar, nil, tcell.StyleDefault.Foreground(rayColor))

				// Add branching sub-rays at certain intervals
				if int(radius)%15 == 0 && peak > 0.3 {
					for branch := 0; branch < 2; branch++ {
						branchAngle := beamAngle + (float64(branch)*2-1)*0.3
						branchRadius := radius * 0.3
						branchX := centerX + int(branchRadius*math.Cos(branchAngle))
						branchY := centerY + int(branchRadius*math.Sin(branchAngle))
						screen.SetContent(branchX, branchY, '·', nil, tcell.StyleDefault.Foreground(rayColor))
					}
				}
			}
		}
	}
}
