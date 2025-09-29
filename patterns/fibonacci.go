package patterns

import (
	"math"
	"math/rand"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// FibonacciVisualizer represents the main Fibonacci spiral visualizer
type FibonacciVisualizer struct {
	*tview.Box
	points     []float64
	fibonacci  []int
	angle      float64
	scale      float64
	depth      int
	colorCache map[int]tcell.Color
	sinCache   []float64
	cosCache   []float64
	lastUpdate time.Time
	patternMgr *Manager
}

// NewFibonacciVisualizer creates a new Fibonacci visualizer
func NewFibonacciVisualizer() *FibonacciVisualizer {
	v := &FibonacciVisualizer{
		Box:        tview.NewBox(),
		points:     make([]float64, 18),
		fibonacci:  generateFibonacci(20),
		angle:      0,
		scale:      1,
		depth:      3,
		colorCache: make(map[int]tcell.Color),
		sinCache:   make([]float64, 360),
		cosCache:   make([]float64, 360),
		lastUpdate: time.Now(),
		patternMgr: NewManager(),
	}

	// Pre-calculate sin/cos values for performance
	for i := 0; i < 360; i++ {
		angle := float64(i) * math.Pi / 180
		v.sinCache[i] = math.Sin(angle)
		v.cosCache[i] = math.Cos(angle)
	}

	return v
}

// Draw renders the Fibonacci visualizer
func (v *FibonacciVisualizer) Draw(screen tcell.Screen) {
	now := time.Now()
	elapsed := now.Sub(v.lastUpdate).Seconds()
	v.lastUpdate = now

	x, y, width, height := v.GetInnerRect()
	centerX, centerY := x+width/2, y+height/2
	baseScale := math.Min(float64(width), float64(height)) / 200

	goldenAngle := math.Pi * (3 - math.Sqrt(5))

	// Character set for Fibonacci spiral
	chars := []rune{'•', '◦', '○', '◎', '◉', '⚬', '⚭', '⚮', '.', '·', '˙', '⋅', '∙', '⁘', '⁛', '⁝', '·', '˙', '∙', '°', '⋅', '∘', '⁖'}

	for d := 0; d < v.depth; d++ {
		for i := 0; i < len(v.fibonacci)-1; i++ {
			amplitude := v.points[i%len(v.points)]
			radius := float64(v.fibonacci[i]) * baseScale * v.scale * (1 - float64(d)*0.2) * (1 + amplitude*0.5)

			rotationDirection := float64(1 - 2*(d%2))
			angleVariation := v.sinCache[i%360] * 0.2
			angle := math.Mod(v.angle*rotationDirection+float64(i)*goldenAngle+float64(d)*0.2+angleVariation, 2*math.Pi)

			angleIndex := int(angle*180/math.Pi) % 360
			if angleIndex < 0 {
				angleIndex += 360
			}
			startX := float64(centerX) + radius*v.cosCache[angleIndex]
			startY := float64(centerY) + radius*v.sinCache[angleIndex]

			curvature := v.sinCache[(i*2)%360] * 10
			endAngle := math.Mod(angle+goldenAngle, 2*math.Pi)
			endAngleIndex := int(endAngle*180/math.Pi) % 360
			if endAngleIndex < 0 {
				endAngleIndex += 360
			}
			endX := float64(centerX) + float64(v.fibonacci[i+1])*baseScale*v.scale*(1-float64(d)*0.2)*v.cosCache[endAngleIndex]
			endY := float64(centerY) + float64(v.fibonacci[i+1])*baseScale*v.scale*(1-float64(d)*0.2)*v.sinCache[endAngleIndex]

			colorKey := i*1000 + d
			color, exists := v.colorCache[colorKey]
			if !exists {
				color = v.getColor(i, amplitude, float64(d), curvature, angleVariation)
				v.colorCache[colorKey] = color
			}

			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			charIndex := (d + i + int(amplitude*10)) % len(chars)

			// Draw simple line between Fibonacci points
			drawSimpleLine(screen, int(startX), int(startY), int(endX), int(endY), color, chars[charIndex])

			// Add random patterns
			v.patternMgr.DrawRandomPattern(screen, rng, color, amplitude)
		}
	}

	// Update angle for continuous animation
	v.angle += 0.2 * elapsed
	v.angle = math.Mod(v.angle, 2*math.Pi)

	// Clear color cache periodically for fresh colors
	if v.angle < 0.01 {
		v.colorCache = make(map[int]tcell.Color)
	}
}

// UpdateWithPeak updates the visualizer with audio peak data
func (v *FibonacciVisualizer) UpdateWithPeak(peak float64) {
	for i := range v.points {
		v.points[i] = peak * math.Sin(float64(i)*math.Pi/50)
	}
	v.scale = 1 + peak*0.2
	v.depth = 3 + int(peak*3)
}

// getColor generates colors for the Fibonacci spiral based on various parameters
func (v *FibonacciVisualizer) getColor(i int, amplitude, depth, curvature, angleVariation float64) tcell.Color {
	hue := math.Mod((float64(i)/float64(len(v.fibonacci)) + v.angle/(2*math.Pi) + curvature*0.01 + angleVariation*0.1), 1)
	saturation := 0.8 + amplitude*0.2
	value := 0.7 + amplitude*0.3 - depth*0.1
	return HSVToRGB(hue, saturation, value)
}

// generateFibonacci generates Fibonacci sequence up to n terms
func generateFibonacci(n int) []int {
	fib := make([]int, n)
	fib[0], fib[1] = 1, 1
	for i := 2; i < n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	return fib
}

// drawSimpleLine draws a basic line between two points
func drawSimpleLine(screen tcell.Screen, x1, y1, x2, y2 int, color tcell.Color, char rune) {
	dx := Abs(x2 - x1)
	dy := Abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	for {
		screen.SetContent(x1, y1, char, nil, tcell.StyleDefault.Foreground(color))
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}
