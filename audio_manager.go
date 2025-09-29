package main

import (
	"fmt"
	"math"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/gordonklaus/portaudio"
)

type AudioManager struct {
	devices          []AudioDeviceInfo
	currentDeviceIdx int
	paStream         *portaudio.Stream
	isInitialized    bool
	isCapturing      bool
	peakLevel        float64
	mutex            sync.RWMutex
	lastAudioTime    time.Time
}

type AudioDeviceInfo struct {
	ID          string
	Name        string
	Type        string // "portaudio", "pulse_monitor"
	Channels    int
	SampleRate  float64
	PADevice    *portaudio.DeviceInfo
	PulseSource string
}

func NewAudioManager() *AudioManager {
	return &AudioManager{
		devices:       make([]AudioDeviceInfo, 0),
		lastAudioTime: time.Now(),
	}
}

func (am *AudioManager) Initialize() error {
	fmt.Println("🔧 Initializing audio system...")

	// Initialize PortAudio
	err := portaudio.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %v", err)
	}

	// Detect all available audio sources
	if err := am.detectAudioSources(); err != nil {
		return err
	}

	am.isInitialized = true
	return nil
}

func (am *AudioManager) detectAudioSources() error {
	fmt.Println("🔍 Scanning for audio sources...")

	// First, get PulseAudio/PipeWire monitor sources
	am.detectPulseMonitorSources()

	// Then get PortAudio devices
	am.detectPortAudioDevices()

	if len(am.devices) == 0 {
		return fmt.Errorf("no audio input sources found")
	}

	// Auto-select the best device
	am.selectBestDevice()

	fmt.Printf("\n📋 Found %d audio sources:\n", len(am.devices))
	for i, device := range am.devices {
		marker := ""
		if i == am.currentDeviceIdx {
			marker = " ✅"
		}
		fmt.Printf("   [%d] %s (%s)%s\n", i, device.Name, device.Type, marker)
	}

	return nil
}

func (am *AudioManager) detectPulseMonitorSources() {
	// Get running sinks first to prioritize active outputs
	runningSinks := am.getRunningAudioSinks()

	// Get all sources
	cmd := exec.Command("pactl", "list", "sources", "short")
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("⚠️  Could not query PulseAudio sources: %v\n", err)
		return
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		sourceName := parts[1]

		// Only include monitor sources
		if !strings.Contains(sourceName, ".monitor") {
			continue
		}

		// Extract base sink name
		baseSinkName := strings.TrimSuffix(sourceName, ".monitor")

		// Prioritize running sinks
		priority := 1
		for _, runningSink := range runningSinks {
			if runningSink == baseSinkName {
				priority = 0
				break
			}
		}

		device := AudioDeviceInfo{
			ID:          fmt.Sprintf("pulse_%d", len(am.devices)),
			Name:        am.formatPulseSourceName(sourceName),
			Type:        "pulse_monitor",
			Channels:    2,
			SampleRate:  48000,
			PulseSource: sourceName,
		}

		// Insert based on priority (running sinks first)
		if priority == 0 {
			// Insert at beginning (high priority)
			am.devices = append([]AudioDeviceInfo{device}, am.devices...)
		} else {
			// Append at end (lower priority)
			am.devices = append(am.devices, device)
		}
	}
}

func (am *AudioManager) getRunningAudioSinks() []string {
	cmd := exec.Command("pactl", "list", "sinks", "short")
	output, err := cmd.Output()
	if err != nil {
		return nil
	}

	var runningSinks []string
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "RUNNING") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				runningSinks = append(runningSinks, parts[1])
			}
		}
	}
	return runningSinks
}

func (am *AudioManager) formatPulseSourceName(sourceName string) string {
	// Make monitor source names more readable
	if strings.Contains(sourceName, "bluez_output") {
		return "Bluetooth Audio Monitor"
	}
	if strings.Contains(sourceName, "alsa_output") && strings.Contains(sourceName, "analog") {
		return "Built-in Audio Monitor"
	}
	if strings.Contains(sourceName, "hdmi") {
		return "HDMI Audio Monitor"
	}

	// Fallback: clean up the name
	name := strings.TrimSuffix(sourceName, ".monitor")
	if len(name) > 30 {
		return name[:27] + "..."
	}
	return name + " Monitor"
}

func (am *AudioManager) detectPortAudioDevices() {
	devices, err := portaudio.Devices()
	if err != nil {
		fmt.Printf("⚠️  Could not query PortAudio devices: %v\n", err)
		return
	}

	for _, device := range devices {
		if device.MaxInputChannels > 0 {
			am.devices = append(am.devices, AudioDeviceInfo{
				ID:         fmt.Sprintf("pa_%s", device.Name),
				Name:       am.formatPortAudioName(device.Name),
				Type:       "portaudio",
				Channels:   device.MaxInputChannels,
				SampleRate: device.DefaultSampleRate,
				PADevice:   device,
			})
		}
	}
}

func (am *AudioManager) formatPortAudioName(name string) string {
	if strings.Contains(name, "pipewire") {
		return "PipeWire"
	}
	if strings.Contains(name, "pulse") {
		return "PulseAudio"
	}
	if strings.Contains(name, "default") {
		return "Default"
	}
	if len(name) > 25 {
		return name[:22] + "..."
	}
	return name
}

func (am *AudioManager) selectBestDevice() {
	// Prioritize pulse monitor sources from running sinks
	for i, device := range am.devices {
		if device.Type == "pulse_monitor" {
			am.currentDeviceIdx = i
			return
		}
	}

	// Fallback to high-channel PortAudio device
	for i, device := range am.devices {
		if device.Type == "portaudio" && device.Channels >= 32 {
			am.currentDeviceIdx = i
			return
		}
	}

	// Last resort: first device
	am.currentDeviceIdx = 0
}

func (am *AudioManager) OpenCurrentDevice() error {
	if am.currentDeviceIdx >= len(am.devices) {
		return fmt.Errorf("invalid device index")
	}

	device := am.devices[am.currentDeviceIdx]
	fmt.Printf("🔧 Opening: %s\n", device.Name)

	// Close existing stream
	if am.paStream != nil {
		am.paStream.Close()
		am.paStream = nil
	}

	if device.Type == "pulse_monitor" {
		return am.openPulseMonitorDevice(device)
	} else {
		return am.openPortAudioDevice(device)
	}
}

func (am *AudioManager) openPulseMonitorDevice(device AudioDeviceInfo) error {
	// Set this monitor source as the default for PortAudio to use
	cmd := exec.Command("pactl", "set-default-source", device.PulseSource)
	if err := cmd.Run(); err != nil {
		fmt.Printf("⚠️  Could not set default source: %v\n", err)
	}

	// Small delay to let PulseAudio update
	time.Sleep(100 * time.Millisecond)

	// Now open with PortAudio default device
	defaultInput, err := portaudio.DefaultInputDevice()
	if err != nil {
		return fmt.Errorf("failed to get default input: %v", err)
	}

	return am.openStreamWithDevice(defaultInput, device.Name)
}

func (am *AudioManager) openPortAudioDevice(device AudioDeviceInfo) error {
	return am.openStreamWithDevice(device.PADevice, device.Name)
}

func (am *AudioManager) openStreamWithDevice(paDevice *portaudio.DeviceInfo, deviceName string) error {
	configs := []struct {
		channels   int
		sampleRate float64
		bufferSize int
		desc       string
	}{
		{2, 44100, 1024, "stereo 44.1kHz"},
		{1, 44100, 1024, "mono 44.1kHz"},
		{2, 48000, 1024, "stereo 48kHz"},
		{1, 48000, 512, "mono 48kHz"},
	}

	for _, config := range configs {
		if config.channels > paDevice.MaxInputChannels {
			continue
		}

		params := portaudio.StreamParameters{
			Input: portaudio.StreamDeviceParameters{
				Device:   paDevice,
				Channels: config.channels,
				Latency:  paDevice.DefaultLowInputLatency,
			},
			SampleRate:      config.sampleRate,
			FramesPerBuffer: config.bufferSize,
		}

		stream, err := portaudio.OpenStream(params, am.audioCallback)
		if err == nil {
			am.paStream = stream
			fmt.Printf("✅ Opened %s (%s)\n", deviceName, config.desc)
			return nil
		}
	}

	return fmt.Errorf("failed to open audio stream")
}

func (am *AudioManager) audioCallback(inputBuffer [][]float32) {
	if len(inputBuffer) == 0 {
		return
	}

	am.mutex.Lock()
	defer am.mutex.Unlock()

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

	am.peakLevel = peak
	if peak > 0.0001 {
		am.lastAudioTime = time.Now()
	}
}

func (am *AudioManager) StartCapture() error {
	if am.paStream == nil {
		if err := am.OpenCurrentDevice(); err != nil {
			return err
		}
	}

	if err := am.paStream.Start(); err != nil {
		return err
	}

	am.isCapturing = true
	return nil
}

func (am *AudioManager) StopCapture() error {
	if am.paStream != nil && am.isCapturing {
		am.paStream.Stop()
		am.isCapturing = false
	}
	return nil
}

func (am *AudioManager) CycleDevice() error {
	wasCapturing := am.isCapturing
	if wasCapturing {
		am.StopCapture()
	}

	am.currentDeviceIdx = (am.currentDeviceIdx + 1) % len(am.devices)

	if err := am.OpenCurrentDevice(); err != nil {
		return err
	}

	if wasCapturing {
		return am.StartCapture()
	}
	return nil
}

func (am *AudioManager) GetPeakLevel() float64 {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return am.peakLevel * 100.0 // Amplify for visualization
}

func (am *AudioManager) GetCurrentDeviceName() string {
	if am.currentDeviceIdx >= len(am.devices) {
		return "Unknown"
	}
	return am.devices[am.currentDeviceIdx].Name
}

func (am *AudioManager) IsCapturing() bool {
	return am.isCapturing
}

func (am *AudioManager) GetTimeSinceLastAudio() time.Duration {
	am.mutex.RLock()
	defer am.mutex.RUnlock()
	return time.Since(am.lastAudioTime)
}

func (am *AudioManager) Cleanup() {
	am.StopCapture()
	if am.paStream != nil {
		am.paStream.Close()
		am.paStream = nil
	}
	if am.isInitialized {
		portaudio.Terminate()
		am.isInitialized = false
	}
}
