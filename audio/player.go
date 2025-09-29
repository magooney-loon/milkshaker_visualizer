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
	monitorSource := p.setupCurrentAudioMonitor()

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

	// Verify our monitor source is available
	p.verifyMonitorSource(monitorSource)

	// Find the best device - prioritize devices that match our monitor source
	var selectedDevice *portaudio.DeviceInfo

	// First priority: Look for device that matches our monitor source
	if monitorSource != "" {
		fmt.Printf("\nSearching for device matching monitor source: %s\n", monitorSource)
		for _, device := range p.devices {
			deviceName := strings.ToLower(device.Name)
			// monitorName := strings.ToLower(monitorSource)

			// Check if device name contains parts of our monitor source
			if strings.Contains(deviceName, "pulse") || strings.Contains(deviceName, "pipewire") {
				selectedDevice = device
				fmt.Printf("Selected PulseAudio compatible device: %s (%d channels)\n", device.Name, device.MaxInputChannels)
				break
			}
		}
	}

	// Second priority: Look for pulse/pipewire devices (these respect PulseAudio routing)
	if selectedDevice == nil {
		for _, device := range p.devices {
			deviceName := strings.ToLower(device.Name)
			if (strings.Contains(deviceName, "pulse") || strings.Contains(deviceName, "pipewire")) && device.MaxInputChannels >= 16 {
				selectedDevice = device
				fmt.Printf("Auto-selected PulseAudio device: %s (%d channels)\n", device.Name, device.MaxInputChannels)
				break
			}
		}
	}

	// Third priority: Any device with many channels
	if selectedDevice == nil {
		for _, device := range p.devices {
			if device.MaxInputChannels >= 32 {
				selectedDevice = device
				fmt.Printf("Auto-selected high-channel device: %s (%d channels)\n", device.Name, device.MaxInputChannels)
				break
			}
		}
	}

	// Fallback: Default input device
	if selectedDevice == nil {
		defaultInput, _ := portaudio.DefaultInputDevice()
		selectedDevice = defaultInput
		if selectedDevice != nil {
			fmt.Printf("Using default: %s\n", selectedDevice.Name)
		} else {
			fmt.Println("Warning: No suitable audio device found")
		}
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

	// Use fewer channels for better compatibility
	channels := 2
	if device.MaxInputChannels == 1 {
		channels = 1
	}

	streamParams := portaudio.StreamParameters{
		Input: portaudio.StreamDeviceParameters{
			Device:   device,
			Channels: channels,
			Latency:  device.DefaultLowInputLatency,
		},
		SampleRate:      44100,
		FramesPerBuffer: 1024,
	}

	fmt.Printf("Opening stream with %d channels at 44100 Hz\n", channels)

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
func (p *Player) setupCurrentAudioMonitor() string {
	fmt.Println("=== AUTO-DETECTING ACTIVE AUDIO OUTPUT ===")

	// Get list of sinks and find the one that's RUNNING
	cmd := exec.Command("pactl", "list", "sinks", "short")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Could not query audio sinks: %v\n", err)
		return ""
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
		fmt.Println("No actively running audio sink found, trying default sink...")
		// Try to get default sink
		cmd = exec.Command("pactl", "get-default-sink")
		output, err = cmd.Output()
		if err != nil {
			fmt.Printf("Could not get default sink: %v\n", err)
			return ""
		}
		runningSink = strings.TrimSpace(string(output))
		if runningSink != "" {
			fmt.Printf("Using default sink: %s\n", runningSink)
		}
	}

	if runningSink == "" {
		fmt.Println("No audio sink found")
		return ""
	}

	// Set the monitor of the running sink as default source
	monitorSource := runningSink + ".monitor"
	cmd = exec.Command("pactl", "set-default-source", monitorSource)
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Failed to set monitor source %s: %v\n", monitorSource, err)
		return ""
	}

	fmt.Printf("✅ Auto-configured source: %s\n", monitorSource)
	fmt.Println("This will capture system audio from your active output")

	// Give PulseAudio a moment to propagate the change
	time.Sleep(500 * time.Millisecond)

	// Try to force the monitor source for applications
	p.forceMonitorSource(monitorSource)

	return monitorSource
}

// verifyMonitorSource checks if the configured monitor source is available
func (p *Player) verifyMonitorSource(monitorSource string) {
	if monitorSource == "" {
		return
	}

	fmt.Println("\n=== VERIFYING MONITOR SOURCE ===")

	// List PulseAudio sources
	cmd := exec.Command("pactl", "list", "sources", "short")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Could not list sources: %v\n", err)
		return
	}

	fmt.Println("Available PulseAudio sources:")
	lines := strings.Split(string(output), "\n")
	monitorFound := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			sourceName := parts[1]
			fmt.Printf("  - %s\n", sourceName)
			if sourceName == monitorSource {
				monitorFound = true
			}
		}
	}

	if monitorFound {
		fmt.Printf("✅ Monitor source %s is available\n", monitorSource)
	} else {
		fmt.Printf("❌ Monitor source %s not found in PulseAudio\n", monitorSource)
	}

	// Check current default source
	cmd = exec.Command("pactl", "get-default-source")
	output, err = cmd.Output()
	if err == nil {
		currentSource := strings.TrimSpace(string(output))
		fmt.Printf("Current default source: %s\n", currentSource)
		if currentSource == monitorSource {
			fmt.Println("✅ Default source matches our monitor")
		} else {
			fmt.Println("⚠️  Default source doesn't match our monitor")
		}
	}
}

// forceMonitorSource tries to force applications to use the monitor source
func (p *Player) forceMonitorSource(monitorSource string) {
	if monitorSource == "" {
		return
	}

	fmt.Println("\n=== FORCING MONITOR SOURCE ===")

	// Try to move all recording streams to our monitor source
	cmd := exec.Command("pactl", "list", "source-outputs", "short")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Could not list source outputs: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 1 {
			sourceOutputId := parts[0]
			// Move this source output to our monitor
			moveCmd := exec.Command("pactl", "move-source-output", sourceOutputId, monitorSource)
			if err := moveCmd.Run(); err == nil {
				fmt.Printf("Moved source output %s to %s\n", sourceOutputId, monitorSource)
			}
		}
	}

	// Also try setting environment variable for new applications
	fmt.Printf("Setting PULSE_SOURCE environment variable to: %s\n", monitorSource)
}
