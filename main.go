package main

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	MusicPlayerMain()
}

type MusicPlayer struct {
	streamer         beep.StreamSeekCloser
	ctrl             *beep.Ctrl
	volume           *effects.Volume
	format           beep.Format
	done             chan bool
	isPlaying        bool
	currentTrack     string
	tracks           []string
	visualizer       *FibonacciVisualizer
	stopVisualizer   chan bool
	visualizerTicker *time.Ticker
	playbackLock     sync.Mutex
	updateInfoFunc   func()
}

func NewMusicPlayer() *MusicPlayer {
	tracks := loadTracks()
	return &MusicPlayer{
		done:           make(chan bool),
		tracks:         tracks,
		visualizer:     NewFibonacciVisualizer(),
		stopVisualizer: make(chan bool),
	}
}

func (mp *MusicPlayer) SetUpdateInfoFunc(updateFunc func()) {
	mp.updateInfoFunc = updateFunc
}

func loadTracks() []string {
	var tracks []string
	files, err := os.ReadDir("media")
	if err != nil {
		log.Printf("Error reading media directory: %v", err)
		return tracks
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".mp3") {
			tracks = append(tracks, filepath.Join("media", file.Name()))
		}
	}

	if len(tracks) == 0 {
		log.Printf("No tracks found in the media directory.")
	}

	return tracks
}

func (mp *MusicPlayer) Play() error {
	mp.playbackLock.Lock()
	defer mp.playbackLock.Unlock()

	if mp.isPlaying {
		if err := mp.Stop(); err != nil {
			log.Printf("Error stopping current track: %v", err)
		}
	}

	if len(mp.tracks) == 0 {
		return fmt.Errorf("no tracks available to play")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	audioFile := mp.tracks[rng.Intn(len(mp.tracks))]
	mp.currentTrack = audioFile

	f, err := os.Open(audioFile)
	if err != nil {
		return fmt.Errorf("error opening audio file: %v", err)
	}

	if mp.updateInfoFunc != nil {
		mp.updateInfoFunc()
	}

	mp.streamer, mp.format, err = mp3.Decode(f)
	if err != nil {
		return fmt.Errorf("error decoding audio file: %v", err)
	}

	mp.ctrl = &beep.Ctrl{Streamer: mp.streamer, Paused: false}
	mp.volume = &effects.Volume{
		Streamer: mp.ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}

	analyzer := NewPeakAnalyzer(mp.volume)

	speaker.Init(mp.format.SampleRate, mp.format.SampleRate.N(time.Second/10))

	speaker.Play(beep.Seq(analyzer, beep.Callback(func() {
		mp.done <- true
	})))

	mp.SetVolume(0)

	mp.isPlaying = true
	mp.visualizerTicker = time.NewTicker(time.Second / 60)

	select {
	case mp.stopVisualizer <- true:
	default:
	}

	go func() {
		for {
			select {
			case <-mp.stopVisualizer:
				mp.visualizerTicker.Stop()
				return
			case <-mp.visualizerTicker.C:
				peak := analyzer.Peak
				mp.visualizer.UpdateWithPeak(peak)
			}
		}
	}()

	go func() {
		<-mp.done
		mp.Play()
	}()

	return nil
}

func (mp *MusicPlayer) Stop() error {
	if !mp.isPlaying {
		return fmt.Errorf("not playing")
	}
	mp.stopVisualizer <- true
	if mp.visualizerTicker != nil {
		mp.visualizerTicker.Stop()
	}
	speaker.Clear()
	err := mp.streamer.Close()
	if err != nil {
		return fmt.Errorf("error stopping playback: %v", err)
	}

	mp.isPlaying = false
	mp.currentTrack = ""
	return nil
}

func (mp *MusicPlayer) Shuffle() error {
	return mp.Play()
}

func (mp *MusicPlayer) SetVolume(percentage float64) {
	if mp.volume != nil {
		percentage = math.Max(-100, math.Min(100, percentage))
		volume := math.Log2(float64(percentage)/100 + 1)
		speaker.Lock()
		mp.volume.Volume = volume
		speaker.Unlock()
	}
}

func (mp *MusicPlayer) GetVolumePercentage() float64 {
	if mp.volume == nil {
		return 0
	}
	return (math.Pow(2, mp.volume.Volume) - 1) * 100
}

func MusicPlayerMain() {
	player := NewMusicPlayer()
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
		infoTextNowPlaying.SetText(filepath.Base(player.currentTrack))
		infoTextVolume.SetText(fmt.Sprintf("Volume: %.0f%%", player.GetVolumePercentage()))
	}

	visualizer := player.visualizer
	player.SetUpdateInfoFunc(updateInfo)
	fullScreenVisualizer := tview.NewBox().SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		visualizer.SetRect(x, y, width, height)
		visualizer.Draw(screen)
		animateLogo(screen, x, y, width, height)
		tview.Print(screen, infoTextNowPlaying.GetText(true), x, y, width, tview.AlignCenter, tcell.ColorWhite)

		tview.Print(screen, infoTextVolume.GetText(true), x, y+1, width, tview.AlignCenter, tcell.ColorWhite)

		tview.Print(screen, "MILKSHAKER PLAYER", x, height-2, width, tview.AlignCenter, tcell.ColorGreen)
		tview.Print(screen, "S (Shuffle), U/D (Volume), Q (Quit)", x, height-1, width, tview.AlignCenter, tcell.ColorGreenYellow)

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
			player.Play()
		case 'u', 'U':
			player.SetVolume(player.GetVolumePercentage() + 10)
		case 'd', 'D':
			player.SetVolume(player.GetVolumePercentage() - 10)
		case 'q', 'Q':
			if err := player.Stop(); err != nil {
				log.Printf("Error stopping playback: %v", err)
			}
			app.Stop()
		}
		updateInfo()
		return event
	})

	if err := app.SetRoot(fullScreenVisualizer, true).SetFocus(fullScreenVisualizer).Run(); err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}

type PeakAnalyzer struct {
	Streamer beep.Streamer
	Peak     float64
	decay    float64
}

func NewPeakAnalyzer(streamer beep.Streamer) *PeakAnalyzer {
	return &PeakAnalyzer{Streamer: streamer, decay: 0.99}
}

func (a *PeakAnalyzer) Stream(samples [][2]float64) (n int, ok bool) {
	n, ok = a.Streamer.Stream(samples)
	a.Peak *= a.decay
	for i := 0; i < n; i++ {
		a.Peak = math.Max(a.Peak, math.Abs(samples[i][0]))
		a.Peak = math.Max(a.Peak, math.Abs(samples[i][1]))
	}
	return n, ok
}

func (a *PeakAnalyzer) Err() error {
	return nil
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
