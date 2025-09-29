package audio

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

// Player handles audio capture and processing
type Player struct {
	paStream         *portaudio.Stream
	peakLevel        float64
	mutex            sync.RWMutex
	lastAudioTime    time.Time
	running          bool
	sensitivity      float64
	devices          []*portaudio.DeviceInfo
	currentDeviceIdx int
	updateInfoFunc   func()
}

// NewPlayer creates a new audio player
func NewPlayer() *Player {
	return &Player{
		sensitivity:   1.0,
		lastAudioTime: time.Now(),
	}
}

// Initialize sets up the audio system
func (p *Player) Initialize() error {
	// Automatically detect and set the active audio monitor
	p.setupCurrentAudioMonitor()

	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %v", err)
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return fmt.Errorf("failed to get audio devices: %v", err)
	}

	p.devices = make([]*portaudio.DeviceInfo, 0)
	for _, device := range devices {
		if device.MaxInputChannels > 0 {
			p.devices = append(p.devices, device)
			fmt.Printf("Found input device: %s (%d channels)\n", device.Name, device.MaxInputChannels)
		}
	}

	if len(p.devices) == 0 {
		return fmt.Errorf("no input devices found")
	}

	// Find the best device (prefer one with many channels, likely a monitor)
	var selectedDevice *portaudio.DeviceInfo
	for _, device := range p.devices {
		if device.MaxInputChannels >= 32 {
			selectedDevice = device
			fmt.Printf("Auto-selected: %s (%d channels)\n", device.Name, device.MaxInputChannels)
			break
		}
	}

	if selectedDevice == nil {
		defaultInput, _ := portaudio.DefaultInputDevice()
		selectedDevice = defaultInput
		fmt.Printf("Using default: %s\n", selectedDevice.Name)
	}

	// Set current device index
	for i, device := range p.devices {
		if device == selectedDevice {
			p.currentDeviceIdx = i
			break
		}
	}

	return p.openStream(selectedDevice)
}

// openStream opens an audio stream with the given device
func (p *Player) openStream(device *portaudio.DeviceInfo) error {
	if p.paStream != nil {
		p.paStream.Close()
	}

	streamParams := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   device,
			Channels: 2,
			Latency:  device.DefaultLowInputLatency,
		},
		SampleRate:      44100,
		FramesPerBuffer: 1024,
	}

	var err error
	p.paStream, err = portaudio.OpenStream(streamParams, p.audioCallback)
	if err != nil {
		return fmt.Errorf("failed to open audio stream: %v", err)
	}

	return nil
}

// audioCallback processes incoming audio data
func (p *Player) audioCallback(inputBuffer [][]float32) {
	if len(inputBuffer) == 0 {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

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

	// Apply sensitivity
	peak *= p.sensitivity

	// Clamp to reasonable range
	if peak > 1.0 {
		peak = 1.0
	}

	p.peakLevel = peak

	if peak > 0.0001 {
		p.lastAudioTime = time.Now()
	}
}

// Start begins audio capture
func (p *Player) Start() error {
	if p.paStream == nil {
		return fmt.Errorf("audio stream not initialized")
	}

	err := p.paStream.Start()
	if err != nil {
		return fmt.Errorf("failed to start audio stream: %v", err)
	}

	p.running = true
	return nil
}

// Stop stops audio capture
func (p *Player) Stop() {
	if p.paStream != nil && p.running {
		p.paStream.Stop()
		p.running = false
	}
}

// Restart stops and starts audio capture
func (p *Player) Restart() {
	p.Stop()
	time.Sleep(100 * time.Millisecond)
	p.Start()
}

// Cleanup cleans up audio resources
func (p *Player) Cleanup() {
	p.Stop()
	if p.paStream != nil {
		p.paStream.Close()
	}
	portaudio.Terminate()
}

// IsCapturing returns true if currently capturing audio
func (p *Player) IsCapturing() bool {
	return p.running
}

// GetPeakLevel returns the current audio peak level
func (p *Player) GetPeakLevel() float64 {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.peakLevel
}

// GetVolumePercentage returns peak level as percentage
func (p *Player) GetVolumePercentage() float64 {
	return p.GetPeakLevel() * 100
}

// GetSensitivity returns current sensitivity setting
func (p *Player) GetSensitivity() float64 {
	return p.sensitivity
}

// IncreaseSensitivity increases audio sensitivity
func (p *Player) IncreaseSensitivity() {
	p.sensitivity = math.Min(p.sensitivity+0.1, 5.0)
	if p.updateInfoFunc != nil {
		p.updateInfoFunc()
	}
}

// DecreaseSensitivity decreases audio sensitivity
func (p *Player) DecreaseSensitivity() {
	p.sensitivity = math.Max(p.sensitivity-0.1, 0.1)
	if p.updateInfoFunc != nil {
		p.updateInfoFunc()
	}
}

// GetCurrentDeviceName returns name of current audio device
func (p *Player) GetCurrentDeviceName() string {
	if p.currentDeviceIdx >= 0 && p.currentDeviceIdx < len(p.devices) {
		return p.devices[p.currentDeviceIdx].Name
	}
	return "Unknown"
}

// CycleDevice switches to next available input device
func (p *Player) CycleDevice() {
	if len(p.devices) <= 1 {
		return
	}

	p.currentDeviceIdx = (p.currentDeviceIdx + 1) % len(p.devices)
	nextDevice := p.devices[p.currentDeviceIdx]

	fmt.Printf("Switching to device: %s\n", nextDevice.Name)

	// Stop current stream
	p.Stop()

	// Open stream with new device
	if err := p.openStream(nextDevice); err != nil {
		fmt.Printf("Failed to switch device: %v\n", err)
		return
	}

	// Restart
	p.Start()

	if p.updateInfoFunc != nil {
		p.updateInfoFunc()
	}
}

// GetCurrentTrack returns a placeholder track info
func (p *Player) GetCurrentTrack() string {
	peak := p.GetPeakLevel()
	if peak > 0.1 {
		return "🎵 Audio Detected"
	} else if peak > 0.01 {
		return "🔉 Low Audio"
	} else {
		return "🔇 No Audio"
	}
}

// SetUpdateInfoFunc sets callback for UI updates
func (p *Player) SetUpdateInfoFunc(fn func()) {
	p.updateInfoFunc = fn
}

// setupCurrentAudioMonitor automatically configures PulseAudio monitor
func (p *Player) setupCurrentAudioMonitor() {
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
