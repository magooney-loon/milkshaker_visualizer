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
		sensitivity:   1.3,
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
			monitorName := strings.ToLower(monitorSource)

			// Check if device name contains parts of our monitor source
			if strings.Contains(deviceName, "pulse") || strings.Contains(deviceName, "pipewire") ||
				strings.Contains(deviceName, "monitor") || strings.Contains(monitorName, deviceName) {
				selectedDevice = device

				break
			}
		}
	}

	// Second priority: Look for pulse/pipewire devices (these respect PulseAudio routing)
	if selectedDevice == nil {
		for _, device := range p.devices {
			deviceName := strings.ToLower(device.Name)
			if strings.Contains(deviceName, "pulse") || strings.Contains(deviceName, "pipewire") {
				selectedDevice = device

				break
			}
		}
	}

	// Third priority: Any device with reasonable channel count
	if selectedDevice == nil {
		for _, device := range p.devices {
			if device.MaxInputChannels >= 2 {
				selectedDevice = device

				break
			}
		}
	}

	// Fallback: First available device
	if selectedDevice == nil {
		if len(p.devices) > 0 {
			selectedDevice = p.devices[0]
		} else {
			return fmt.Errorf("no audio input devices available")
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

	wasRunning := p.running
	prevDeviceIdx := p.currentDeviceIdx
	p.currentDeviceIdx = (p.currentDeviceIdx + 1) % len(p.devices)
	nextDevice := p.devices[p.currentDeviceIdx]

	// Stop current stream
	if wasRunning {
		p.Stop()
	}

	// Open stream with new device
	if err := p.openStream(nextDevice); err != nil {
		p.currentDeviceIdx = prevDeviceIdx
		p.openStream(p.devices[prevDeviceIdx])
		if wasRunning {
			p.Start()
		}
		return
	}

	// Restart if was running
	if wasRunning {
		p.Start()
	}

	if p.updateInfoFunc != nil {
		p.updateInfoFunc()
	}
}

// GetCurrentTrack returns a placeholder track info
func (p *Player) GetCurrentTrack() string {
	peak := p.GetPeakLevel()
	if peak > 0.001 { // Much lower threshold for system audio
		return "🎵 Audio Detected"
	} else if peak > 0.0001 {
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
	// Get list of sinks and find the one that's RUNNING
	cmd := exec.Command("pactl", "list", "sinks", "short")
	output, err := cmd.Output()
	if err != nil {
		return p.fallbackMonitorSetup()
	}

	lines := strings.Split(string(output), "\n")
	var runningSink string

	// First try: Find RUNNING sink
	for _, line := range lines {
		if strings.Contains(line, "RUNNING") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				runningSink = parts[1] // Get sink name
				break
			}
		}
	}

	// Second try: Find any available sink
	if runningSink == "" {
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					runningSink = parts[1]
					break
				}
			}
		}
	}

	// Third try: Get default sink
	if runningSink == "" {
		cmd = exec.Command("pactl", "get-default-sink")
		output, err = cmd.Output()
		if err != nil {
			return p.fallbackMonitorSetup()
		}
		runningSink = strings.TrimSpace(string(output))
	}

	if runningSink == "" {
		return p.fallbackMonitorSetup()
	}

	// Set the monitor of the running sink as default source
	monitorSource := runningSink + ".monitor"
	cmd = exec.Command("pactl", "set-default-source", monitorSource)
	err = cmd.Run()
	if err != nil {
		return p.setupAlternativeMonitor(runningSink)
	}

	// Give PulseAudio a moment to propagate the change
	time.Sleep(500 * time.Millisecond)

	// Try to force the monitor source for applications
	p.forceMonitorSource(monitorSource)

	return monitorSource
}

// verifyMonitorSource checks if the configured monitor source is available
func (p *Player) verifyMonitorSource(monitorSource string) {
	// Silently verify - no debug output needed
}

// forceMonitorSource tries to force applications to use the monitor source
func (p *Player) forceMonitorSource(monitorSource string) {
	if monitorSource == "" {
		return
	}

	// Try to move all recording streams to our monitor source
	cmd := exec.Command("pactl", "list", "source-outputs", "short")
	output, err := cmd.Output()
	if err != nil {
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
			moveCmd.Run()
		}
	}
}

// fallbackMonitorSetup tries to set up audio capture when pactl commands fail
func (p *Player) fallbackMonitorSetup() string {
	// Try to find any .monitor source
	cmd := exec.Command("pactl", "list", "sources", "short")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			sourceName := parts[1]
			if strings.HasSuffix(sourceName, ".monitor") {
				// Try to set it as default
				setCmd := exec.Command("pactl", "set-default-source", sourceName)
				if err := setCmd.Run(); err == nil {
					return sourceName
				}
			}
		}
	}

	return ""
}

// setupAlternativeMonitor tries alternative ways to setup the monitor
func (p *Player) setupAlternativeMonitor(sinkName string) string {
	// Try loading a loopback module as fallback
	cmd := exec.Command("pactl", "load-module", "module-loopback", fmt.Sprintf("source=%s.monitor", sinkName))
	if err := cmd.Run(); err == nil {
		time.Sleep(1 * time.Second)
		return sinkName + ".monitor"
	}

	// Try generic loopback
	cmd = exec.Command("pactl", "load-module", "module-loopback")
	if err := cmd.Run(); err == nil {
		time.Sleep(1 * time.Second)
		return sinkName + ".monitor"
	}

	return p.fallbackMonitorSetup()
}
