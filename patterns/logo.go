package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
)

var (
	logoGradientPhase    float64 = 0.0
	logoGradientStrength float64 = 0.0
	logoLastUpdate       time.Time
)

// DrawLogo creates a smooth gradient logo with fixed size
func DrawLogo(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	now := time.Now()
	elapsed := now.Sub(logoLastUpdate).Seconds()
	if elapsed < 1.0/60.0 { // 60 FPS limit
		return
	}
	logoLastUpdate = now

	logoFrames := []string{
		" __    __     __     __         __  __     ______     __  __     ______     __  __     ______     ______    ",
		"/\\ \"-./  \\   /\\ \\   /\\ \\       /\\ \\/ /    /\\  ___\\   /\\ \\_\\ \\   /\\  __ \\   /\\ \\/ /    /\\  ___\\   /\\  == \\   ",
		"\\ \\ \\-./\\ \\  \\ \\ \\  \\ \\ \\____  \\ \\  _\"-.  \\ \\___  \\  \\ \\  __ \\  \\ \\  __ \\  \\ \\  _\"-.  \\ \\  __\\   \\ \\  __<   ",
		" \\ \\_\\ \\ \\_\\  \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\  \\/\\_____\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\ ",
		"  \\/_/  \\/_/   \\/_/   \\/_____/   \\/_/\\/_/   \\/_____/   \\/_/\\/_/   \\/_/\\/_/   \\/_/\\/_/   \\/_____/   \\/_/ /_/ ",
	}

	// Smooth logo strength progression
	targetStrength := 0.0
	if peak > 0.15 {
		// Smooth sine wave for visibility with long period
		cycleSpeed := 0.3 // Very slow cycle
		logoGradientPhase += elapsed * cycleSpeed
		targetStrength = (math.Sin(logoGradientPhase) + 1.0) / 2.0 * math.Min(1.0, (peak-0.15)*2.5)
	}

	// Smooth interpolation to target strength
	smoothing := 0.95 // Higher = smoother transition
	logoGradientStrength = logoGradientStrength*smoothing + targetStrength*(1.0-smoothing)

	// Don't render if too weak
	if logoGradientStrength < 0.05 {
		return
	}

	// Fixed positioning - center of screen
	logoHeight := len(logoFrames)
	logoWidth := len(logoFrames[0])
	startY := (height - logoHeight) / 2
	startX := (width - logoWidth) / 2

	// Smooth breathing effect
	basePhase := GetBasePhase()
	breathe := 1.0 + math.Sin(basePhase*1.2)*0.03 // Very subtle breathing

	for i, line := range logoFrames {
		for j, char := range line {
			if char == ' ' {
				continue
			}

			finalX := startX + j
			finalY := startY + i

			// Bounds check
			if finalX < 0 || finalX >= width || finalY < 0 || finalY >= height {
				continue
			}

			// Smooth radial gradient from center
			centerY := float64(logoHeight) / 2.0
			centerX := float64(logoWidth) / 2.0
			dy := float64(i) - centerY
			dx := float64(j) - centerX
			distanceFromCenter := math.Sqrt(dx*dx + dy*dy)
			maxDistance := math.Sqrt(centerX*centerX + centerY*centerY)

			// Smooth radial falloff
			radialGradient := 1.0 - math.Min(1.0, distanceFromCenter/maxDistance*1.2)
			radialGradient = math.Max(0.0, radialGradient)

			// Smooth wave gradient across logo
			waveGradient := (math.Sin(float64(j)*0.1+basePhase*0.8) + 1.0) / 2.0

			// Vertical gradient
			verticalGradient := (math.Sin(float64(i)*0.3+basePhase*0.5) + 1.0) / 2.0

			// Combine gradients
			combinedGradient := (radialGradient*0.6 + waveGradient*0.2 + verticalGradient*0.2) * breathe

			// Final intensity
			finalIntensity := logoGradientStrength * combinedGradient

			// Skip very weak pixels
			if finalIntensity < 0.1 {
				continue
			}

			// Smooth color based on intensity and position
			baseHue := 0.25 + peak*0.1 // Golden to green range
			hueVariation := math.Sin(float64(i*j)*0.01+basePhase*0.1) * 0.05
			finalHue := baseHue + hueVariation

			saturation := 0.4 + peak*0.3 + finalIntensity*0.2
			saturation = math.Max(0.2, math.Min(0.9, saturation))

			value := 0.3 + finalIntensity*0.6 + peak*0.1
			value = math.Max(0.1, math.Min(1.0, value))

			logoColor := HSVToRGB(finalHue, saturation, value)

			// Intelligent character fading based on intensity - 5 levels for consistency
			var displayChar rune
			if finalIntensity < 0.1 {
				displayChar = '·' // Barely visible dot
			} else if finalIntensity < 0.2 {
				displayChar = '˙' // Small dot
			} else if finalIntensity < 0.35 {
				displayChar = '∘' // Circle outline
			} else if finalIntensity < 0.5 {
				displayChar = '◦' // Larger circle
			} else {
				displayChar = rune(char) // Original logo character
			}

			screen.SetContent(finalX, finalY, displayChar, nil, tcell.StyleDefault.Foreground(logoColor))
		}
	}
}

// DrawLogoLayer draws the logo as an integrated pattern layer with smooth gradients
func DrawLogoLayer(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64, depthLayer int) {
	// Only show in middle layers
	if depthLayer < 1 || depthLayer > 3 {
		return
	}

	// Reduce intensity for depth layers
	originalStrength := logoGradientStrength
	depthScale := 1.0 - float64(depthLayer)*0.2
	logoGradientStrength *= depthScale * 0.7

	DrawLogo(screen, width, height, color, char, rng, peak*depthScale)

	// Restore original strength
	logoGradientStrength = originalStrength
}
