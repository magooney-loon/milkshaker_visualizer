package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/gordonklaus/portaudio"
	"github.com/rivo/tview"
)

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "devices":
			listAudioDevices()
			return
		case "setup-audio":
			setupSystemAudio()
			return
		case "test-audio":
			testAudioCapture()
			return
		case "help":
			showHelp()
			return
		}
	}
	AudioPlayerMain()
}

func listAudioDevices() {
	err := portaudio.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize PortAudio: %v", err)
	}
	defer portaudio.Terminate()

	devices, err := portaudio.Devices()
	if err != nil {
		log.Fatalf("Failed to get audio devices: %v", err)
	}

	fmt.Println("Available Audio Devices:")
	fmt.Println("========================")

	for i, device := range devices {
		fmt.Printf("[%d] %s\n", i, device.Name)
		fmt.Printf("    Max Input Channels: %d\n", device.MaxInputChannels)
		fmt.Printf("    Max Output Channels: %d\n", device.MaxOutputChannels)
		fmt.Printf("    Default Sample Rate: %.0f Hz\n", device.DefaultSampleRate)
		if device.MaxInputChannels > 0 {
			fmt.Printf("    Input Latency: %.3f ms\n", device.DefaultLowInputLatency.Seconds()*1000)
		}
		if device.MaxOutputChannels > 0 {
			fmt.Printf("    Output Latency: %.3f ms\n", device.DefaultLowOutputLatency.Seconds()*1000)
		}
		fmt.Println()
	}

	defaultInput, err := portaudio.DefaultInputDevice()
	if err == nil {
		fmt.Printf("Default Input Device: %s\n", defaultInput.Name)
	}

	defaultOutput, err := portaudio.DefaultOutputDevice()
	if err == nil {
		fmt.Printf("Default Output Device: %s\n", defaultOutput.Name)
	}
}

func setupSystemAudio() {
	fmt.Println("Setting up system audio capture for Linux...")
	fmt.Println("============================================")

	fmt.Println("Method 1: Load PulseAudio loopback module")
	fmt.Println("Run: pactl load-module module-loopback")
	fmt.Println("This creates a loopback from output to input")
	fmt.Println()

	fmt.Println("Method 2: Check available monitor sources")
	fmt.Println("Run: pactl list sources short")
	fmt.Println("Look for sources ending in '.monitor'")
	fmt.Println()

	fmt.Println("Method 3: Use pavucontrol")
	fmt.Println("1. Install: sudo apt install pavucontrol")
	fmt.Println("2. Run: pavucontrol")
	fmt.Println("3. Go to Recording tab while visualizer is running")
	fmt.Println("4. Set it to record from 'Monitor of [your output device]'")
	fmt.Println()

	fmt.Println("After setup, run the visualizer and play some music to test!")
}

func testAudioCapture() {
	fmt.Println("AUDIO CAPTURE TEST")
	fmt.Println("==================")
	fmt.Println("This will test if audio capture is working without the full visualizer.")
	fmt.Println("Press Ctrl+C to stop the test.")
	fmt.Println()

	tester := NewSimpleAudioTester()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize and start
	err := tester.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize: %v", err)
	}
	defer tester.Cleanup()

	err = tester.Start()
	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	// Wait for interrupt
	<-sigChan
	fmt.Println("\n\nStopping audio test...")
	tester.Stop()
}

type SimpleAudioTester struct {
	paStream      *portaudio.Stream
	peakLevel     float64
	mutex         sync.RWMutex
	lastAudioTime time.Time
	running       bool
}

func NewSimpleAudioTester() *SimpleAudioTester {
	return &SimpleAudioTester{
		lastAudioTime: time.Now(),
	}
}

func (sat *SimpleAudioTester) audioCallback(inputBuffer [][]float32) {
	if len(inputBuffer) == 0 {
		return
	}

	sat.mutex.Lock()
	defer sat.mutex.Unlock()

	peak := float64(0)
	sampleCount := 0

	for _, channel := range inputBuffer {
		for _, sample := range channel {
			absSample := math.Abs(float64(sample))
			if absSample > peak {
				peak = absSample
			}
			sampleCount++
		}
	}

	sat.peakLevel = peak

	if peak > 0.0001 {
		sat.lastAudioTime = time.Now()
	}

	// Print real-time audio levels
	now := time.Now()
	if peak > 0.01 {
		fmt.Printf("\r🎵 STRONG: Peak=%.4f | %s", peak, now.Format("15:04:05"))
	} else if peak > 0.001 {
		fmt.Printf("\r🔉 Medium: Peak=%.4f | %s", peak, now.Format("15:04:05"))
	} else if peak > 0.0001 {
		fmt.Printf("\r🔈 Low: Peak=%.6f | %s", peak, now.Format("15:04:05"))
	} else {
		fmt.Printf("\r🔇 Silent: Peak=%.8f | %s", peak, now.Format("15:04:05"))
	}
}

func (sat *SimpleAudioTester) Initialize() error {
	// First, automatically detect and set the active audio monitor
	sat.setupCurrentAudioMonitor()

	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %v", err)
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("failed to get audio devices: %v", err)
	}

	var selectedDevice *portaudio.DeviceInfo
	for _, device := range devices {
		if device.MaxInputChannels >= 32 {
			selectedDevice = device
			fmt.Printf("Selected: %s (%d channels)\n", device.Name, device.MaxInputChannels)
			break
		}
	}

	if selectedDevice == nil {
		defaultInput, _ := portaudio.DefaultInputDevice()
		selectedDevice = defaultInput
		fmt.Printf("Using default: %s\n", selectedDevice.Name)
	}

	streamParams := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   selectedDevice,
			Channels: 2,
			Latency:  selectedDevice.DefaultLowInputLatency,
		},
		SampleRate:      44100,
		FramesPerBuffer: 1024,
	}

	sat.paStream, err = portaudio.OpenStream(streamParams, sat.audioCallback)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}

	return nil
}

func (sat *SimpleAudioTester) Start() error {
	err := sat.paStream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}
	sat.running = true
	fmt.Println("🎤 Listening for audio... (play some music to test)")
	return nil
}

func (sat *SimpleAudioTester) Stop() {
	if sat.paStream != nil && sat.running {
		sat.paStream.Stop()
		sat.running = false
	}
}

func (sat *SimpleAudioTester) Cleanup() {
	sat.Stop()
	if sat.paStream != nil {
		sat.paStream.Close()
	}
	portaudio.Terminate()
}

func (sat *SimpleAudioTester) setupCurrentAudioMonitor() {
	fmt.Println("=== AUTO-DETECTING ACTIVE AUDIO OUTPUT ===")

	// Get list of sinks and find the one that's RUNNING
	cmd := exec.Command("pactl", "list", "sinks", "short")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Could not query audio sinks: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	var runningSink string

	for _, line := range lines {
		if strings.Contains(line, "RUNNING") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				runningSink = parts[1] // Get sink name
				fmt.Printf("Found active audio sink: %s\n", runningSink)
				break
			}
		}
	}

	if runningSink == "" {
		fmt.Println("No actively running audio sink found")
		return
	}

	// Set the monitor of the running sink as default source
	monitorSource := runningSink + ".monitor"
	cmd = exec.Command("pactl", "set-default-source", monitorSource)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to set monitor source %s: %v\n", monitorSource, err)
		return
	}

	fmt.Printf("✅ Auto-configured source: %s\n", monitorSource)
	fmt.Println("This will capture system audio from your active output")
}

func showHelp() {
	fmt.Println("MILKSHAKER VISUALIZER")
	fmt.Println("====================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run .                 # Start the visualizer")
	fmt.Println("  go run . devices         # List available audio devices")
	fmt.Println("  go run . setup-audio     # Show audio setup instructions")
	fmt.Println("  go run . test-audio      # Test audio capture without UI")
	fmt.Println("  go run . help            # Show this help")
	fmt.Println()
	fmt.Println("Controls (when running):")
	fmt.Println("  S         Start/Stop audio capture")
	fmt.Println("  R         Restart audio capture")
	fmt.Println("  +/-       Adjust sensitivity")
	fmt.Println("  D         Show available devices")
	fmt.Println("  Ctrl+C    Quit")
	fmt.Println()
	fmt.Println("For system audio capture on Linux:")
	fmt.Println("  Run: go run . setup-audio")
}

func AudioPlayerMain() {
	fmt.Println("🎵 MILKSHAKER VISUALIZER")
	fmt.Println("=======================")
	fmt.Println("Initializing audio system...")
	fmt.Println()

	player := NewAudioPlayer()

	// Initialize audio player with all logging upfront
	if err := player.Initialize(); err != nil {
		log.Fatalf("Failed to initialize audio player: %v", err)
	}
	defer player.Cleanup()

	fmt.Println()
	fmt.Println("✅ Audio system initialized successfully!")
	fmt.Println("🎤 Starting visualizer...")
	fmt.Println("Press S to start/stop | +/- for sensitivity | Ctrl+C to quit")
	fmt.Println()
	time.Sleep(2 * time.Second) // Give user time to read

	app := tview.NewApplication()
	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		width, height := screen.Size()
		player.visualizer.SetRect(0, 0, width, height)
	})

	infoTextNowPlaying := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	infoTextVolume := tview.NewTextView().
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter)

	updateInfo := func() {
		infoTextNowPlaying.SetText(player.GetCurrentTrack())
		infoTextVolume.SetText(fmt.Sprintf("Peak: %.0f%% | Sensitivity: %.1fx | Device: %s", player.GetVolumePercentage(), player.GetSensitivity(), player.GetCurrentDeviceName()))
	}

	visualizer := player.visualizer
	player.SetUpdateInfoFunc(updateInfo)
	fullScreenVisualizer := tview.NewBox().SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		visualizer.SetRect(x, y, width, height)
		visualizer.Draw(screen)
		animateLogo(screen, x, y, width, height)
		tview.Print(screen, infoTextNowPlaying.GetText(true), x, y, width, tview.AlignCenter, tcell.ColorWhite)

		tview.Print(screen, infoTextVolume.GetText(true), x, y+1, width, tview.AlignCenter, tcell.ColorWhite)

		tview.Print(screen, "MILKSHAKER VISUALIZER", x, height-2, width, tview.AlignCenter, tcell.ColorGreen)

		var statusText string
		if player.IsCapturing() {
			statusText = "R (Restart), S (Stop), +/- (Sensitivity), D (Cycle Device), Ctrl+C (Quit)"
		} else {
			statusText = "S (Start), +/- (Sensitivity), D (Cycle Device), Ctrl+C (Quit)"
		}
		tview.Print(screen, statusText, x, height-1, width, tview.AlignCenter, tcell.ColorGreenYellow)

		return x, y, width, height
	})

	updateInfo()

	go func() {
		for {
			time.Sleep(time.Second / 60) // 60 FPS
			app.Draw()
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Rune() {
		case 's', 'S':
			if player.IsCapturing() {
				player.Stop()
			} else {
				player.Start()
			}
		case 'r', 'R':
			player.Restart()
		case '+', '=':
			player.IncreaseSensitivity()
		case '-', '_':
			player.DecreaseSensitivity()

		case 'd', 'D':
			// Cycle to next audio input device
			player.CycleDevice()
		}

		// Handle Ctrl+C for quit
		if event.Key() == tcell.KeyCtrlC {
			player.Stop()
			app.Stop()
		}
		updateInfo()
		return event
	})

	if err := app.SetRoot(fullScreenVisualizer, true).SetFocus(fullScreenVisualizer).Run(); err != nil {
		fmt.Printf("\nVisualizer stopped: %v\n", err)
	}
}

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
}

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
	}
	for i := 0; i < 360; i++ {
		angle := float64(i) * math.Pi / 180
		v.sinCache[i] = math.Sin(angle)
		v.cosCache[i] = math.Cos(angle)
	}
	return v
}

func (v *FibonacciVisualizer) Draw(screen tcell.Screen) {
	now := time.Now()
	elapsed := now.Sub(v.lastUpdate).Seconds()
	v.lastUpdate = now
	x, y, width, height := v.GetInnerRect()
	centerX, centerY := x+width/2, y+height/2
	baseScale := math.Min(float64(width), float64(height)) / 200

	goldenAngle := math.Pi * (3 - math.Sqrt(5))

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
			drawFunkyLine(screen, int(startX), int(startY), int(endX), int(endY), color, chars[charIndex], amplitude)
			drawRandomPattern(screen, rng, color, amplitude)

		}
	}

	v.angle += 0.2 * elapsed
	v.angle = math.Mod(v.angle, 2*math.Pi)

	if v.angle < 0.01 {
		v.colorCache = make(map[int]tcell.Color)
	}

}

func (v *FibonacciVisualizer) UpdateWithPeak(peak float64) {
	for i := range v.points {
		v.points[i] = peak * math.Sin(float64(i)*math.Pi/50)
	}
	v.scale = 1 + peak*0.2
	v.depth = 3 + int(peak*3)
}

func (v *FibonacciVisualizer) getColor(i int, amplitude, depth, curvature, angleVariation float64) tcell.Color {
	hue := math.Mod((float64(i)/float64(len(v.fibonacci)) + v.angle/(2*math.Pi) + curvature*0.01 + angleVariation*0.1), 1)
	saturation := 0.8 + amplitude*0.2
	value := 0.7 + amplitude*0.3 - depth*0.1
	return hsvToRGB(hue, saturation, value)
}

func drawFunkyLine(screen tcell.Screen, x1, y1, x2, y2 int, color tcell.Color, char rune, amplitude float64) {

	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 >= x2 {
		sx = -1
	}
	if y1 >= y2 {
		sy = -1
	}
	err := dx - dy

	frequency := 0.2
	basePhase := float64(time.Now().UnixNano()) / 1e9

	for {
		t := float64(x1+x2+y1+y2) * frequency
		waveOffset := amplitude * math.Sin(t+basePhase)

		wx := float64(x1)
		wy := float64(y1)
		if dx > dy {
			wy += waveOffset
		} else {
			wx += waveOffset
		}

		screen.SetContent(int(wx), int(wy), char, nil, tcell.StyleDefault.Foreground(color))
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

func drawRandomPattern(screen tcell.Screen, rng *rand.Rand, color tcell.Color, amplitude float64) {
	width, height := screen.Size()
	char := randomRune(rng)

	patterns := []func(tcell.Screen, int, int, tcell.Color, rune, *rand.Rand, float64){
		drawZigZag,
		drawSpiral,
		drawStarburst,
		drawRandomWalk,
	}

	patternIndex := int(amplitude * float64(len(patterns)))
	if patternIndex >= len(patterns) {
		patternIndex = len(patterns) - 1
	}

	pattern := patterns[patternIndex]
	pattern(screen, width, height, color, char, rng, amplitude)
}

func drawZigZag(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	step := 1
	for x, y := 0, 0; x < width; x++ {
		screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
		if y >= height-1 || y <= 0 {
			step = -step
		}
		y += step
	}
}

func drawSpiral(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	radius := 1.0 * peak
	angle := 0.0
	angleStep := 0.1

	amplitude := 5.0 * peak
	frequency := 0.2 + 0.1*peak
	basePhase := float64(time.Now().UnixNano()) / 1e9

	for radius < float64(min(width, height))/2 {
		waveOffset := amplitude * math.Sin(frequency*angle+basePhase)
		x := centerX + int((radius+waveOffset)*math.Cos(angle))
		y := centerY + int((radius+waveOffset)*math.Sin(angle))
		screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
		radius += 0.1
		angle += angleStep
	}
}

func drawStarburst(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	centerX, centerY := width/2, height/2
	basePhase := float64(time.Now().UnixNano()) / 1e9
	amplitude := 5.0 * peak

	for angle := 0.0; angle < 2*math.Pi; angle += math.Pi / 8 {
		for radius := 0.0; radius < float64(min(width, height))/2; radius += peak {
			waveOffset := amplitude * math.Sin(angle+basePhase)
			x := centerX + int((radius+waveOffset)*math.Cos(angle))
			y := centerY + int((radius+waveOffset)*math.Sin(angle))
			screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
		}
	}
}

func drawRandomWalk(screen tcell.Screen, width, height int, color tcell.Color, char rune, rng *rand.Rand, peak float64) {
	x, y := width/2, height/2
	for i := 0; i < 100; i++ {
		screen.SetContent(x, y, char, nil, tcell.StyleDefault.Foreground(color))
		switch rng.Intn(4) {
		case 0:
			x++
		case 1:
			x--
		case 2:
			y++
		case 3:
			y--
		}
		if x < 0 {
			x = 0
		} else if x >= width {
			x = width - 1
		}
		if y < 0 {
			y = 0
		} else if y >= height {
			y = height - 1
		}
	}
}

func randomRune(rng *rand.Rand) rune {
	runes := []rune{'*', '+', 'x', 'o', '~', '@', '#', '$', '%', '&'}
	return runes[rng.Intn(len(runes))]
}

func generateFibonacci(n int) []int {
	fib := make([]int, n)
	fib[0], fib[1] = 1, 1
	for i := 2; i < n; i++ {
		fib[i] = fib[i-1] + fib[i-2]
	}
	return fib
}

func hsvToRGB(h, s, v float64) tcell.Color {
	i := int(h * 6)
	f := h*6 - float64(i)
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	var r, g, b float64
	switch i % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}

	return tcell.NewRGBColor(int32(r*255), int32(g*255), int32(b*255))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

const (
	logoRevealInterval  = 20 * time.Millisecond
	cycleWaitDuration   = 20 * time.Second
	stayVisibleDuration = 10 * time.Second
)

var (
	lastLogoTime  time.Time
	logoMask      [][]bool
	revealedCount int
	fadeOutCount  int
	isFadingOut   bool
	cycleEndTime  time.Time
)

func animateLogo(screen tcell.Screen, x, y, width, height int) {
	now := time.Now()
	if now.Sub(lastLogoTime) < logoRevealInterval {
		return
	}
	lastLogoTime = now

	logoFrames := []string{
		" __    __     __     __         __  __     ______     __  __     ______     __  __     ______     ______    ",
		"/\\ \"-./  \\   /\\ \\   /\\ \\       /\\ \\/ /    /\\  ___\\   /\\ \\_\\ \\   /\\  __ \\   /\\ \\/ /    /\\  ___\\   /\\  == \\   ",
		"\\ \\ \\-./\\ \\  \\ \\ \\  \\ \\ \\____  \\ \\  _\"-.  \\ \\___  \\  \\ \\  __ \\  \\ \\  __ \\  \\ \\  _\"-.  \\ \\  __\\   \\ \\  __<   ",
		" \\ \\_\\ \\ \\_\\  \\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\  \\/\\_____\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_\\ \\_\\  \\ \\_____\\  \\ \\_\\ \\_\\ ",
		"  \\/_/  \\/_/   \\/_/   \\/_____/   \\/_/\\/_/   \\/_____/   \\/_/\\/_/   \\/_/\\/_/   \\/_/\\/_/   \\/_____/   \\/_/ /_/ ",
	}

	middleY := y + (height / 2) - (len(logoFrames) / 2)
	middleX := x + (width / 2) - (len(logoFrames[0]) / 2)

	// Initialize logoMask if it's empty
	if len(logoMask) == 0 {
		logoMask = make([][]bool, len(logoFrames))
		for i := range logoMask {
			logoMask[i] = make([]bool, len(logoFrames[0]))
		}
	}

	totalNonSpaceChars := countNonSpaceChars(logoFrames)

	if cycleEndTime.IsZero() {
		cycleEndTime = now.Add(stayVisibleDuration)
	}

	if !isFadingOut {
		if revealedCount < totalNonSpaceChars {
			for {
				i := rand.Intn(len(logoMask))
				j := rand.Intn(len(logoMask[0]))
				if !logoMask[i][j] && logoFrames[i][j] != ' ' {
					logoMask[i][j] = true
					revealedCount++
					break
				}
			}
		} else if now.After(cycleEndTime) {
			isFadingOut = true
		}
	} else {
		if fadeOutCount < totalNonSpaceChars {
			for {
				i := rand.Intn(len(logoMask))
				j := rand.Intn(len(logoMask[0]))
				if logoMask[i][j] && logoFrames[i][j] != ' ' {
					logoMask[i][j] = false
					fadeOutCount++
					break
				}
			}
		} else {
			cycleEndTime = now.Add(cycleWaitDuration)
			resetCycle()
		}
	}

	for i, line := range logoFrames {
		for j, char := range line {
			if logoMask[i][j] {
				style := tcell.StyleDefault.Foreground(tcell.ColorFloralWhite)
				screen.SetContent(middleX+j, middleY+i, rune(char), nil, style)
			}
		}
	}
}

func countNonSpaceChars(logoFrames []string) int {
	count := 0
	for _, line := range logoFrames {
		for _, char := range line {
			if char != ' ' {
				count++
			}
		}
	}
	return count
}

func resetCycle() {
	for i := range logoMask {
		for j := range logoMask[i] {
			logoMask[i][j] = false
		}
	}
	revealedCount = 0
	fadeOutCount = 0
	isFadingOut = false
	lastLogoTime = time.Time{}
}
