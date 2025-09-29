package main

import (
	"fmt"
	"sync"
	"time"
)

type AudioPlayer struct {
	visualizer       *FibonacciVisualizer
	stopVisualizer   chan bool
	visualizerTicker *time.Ticker
	captureLock      sync.Mutex
	updateInfoFunc   func()
	peakAnalyzer     *SystemPeakAnalyzer
	peakSensitivity  float64
	audioManager     *AudioManager
}

type SystemPeakAnalyzer struct {
	Peak  float64
	decay float64
	mutex sync.RWMutex
}

func NewSystemPeakAnalyzer() *SystemPeakAnalyzer {
	return &SystemPeakAnalyzer{
		decay: 0.95,
	}
}

func (a *SystemPeakAnalyzer) UpdatePeak(peak float64) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.Peak *= a.decay
	if peak > a.Peak {
		a.Peak = peak
	}
}

func (a *SystemPeakAnalyzer) GetPeak() float64 {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.Peak
}

func NewAudioPlayer() *AudioPlayer {
	return &AudioPlayer{
		visualizer:      NewFibonacciVisualizer(),
		stopVisualizer:  make(chan bool),
		peakAnalyzer:    NewSystemPeakAnalyzer(),
		peakSensitivity: 1.0,
		audioManager:    NewAudioManager(),
	}
}

func (ap *AudioPlayer) SetUpdateInfoFunc(updateFunc func()) {
	ap.updateInfoFunc = updateFunc
}

func (ap *AudioPlayer) Initialize() error {
	return ap.audioManager.Initialize()
}

func (ap *AudioPlayer) Start() error {
	ap.captureLock.Lock()
	defer ap.captureLock.Unlock()

	if ap.updateInfoFunc != nil {
		ap.updateInfoFunc()
	}

	// Start visualizer ticker
	ap.visualizerTicker = time.NewTicker(time.Second / 60) // 60 FPS

	// Clear any previous stop signal
	select {
	case <-ap.stopVisualizer:
	default:
	}

	// Start visualizer goroutine
	go func() {
		for {
			select {
			case <-ap.stopVisualizer:
				ap.visualizerTicker.Stop()
				return
			case <-ap.visualizerTicker.C:
				// Get peak from audio manager
				rawPeak := ap.audioManager.GetPeakLevel()
				peak := rawPeak * ap.peakSensitivity / 100.0

				// Update peak analyzer
				ap.peakAnalyzer.UpdatePeak(peak)

				// Update visualizer
				visualPeak := ap.peakAnalyzer.GetPeak()
				ap.visualizer.UpdateWithPeak(visualPeak)
			}
		}
	}()

	// Start audio capture
	return ap.audioManager.StartCapture()
}

func (ap *AudioPlayer) Stop() error {
	ap.captureLock.Lock()
	defer ap.captureLock.Unlock()

	// Stop visualizer
	select {
	case ap.stopVisualizer <- true:
	default:
	}

	if ap.visualizerTicker != nil {
		ap.visualizerTicker.Stop()
	}

	if ap.updateInfoFunc != nil {
		ap.updateInfoFunc()
	}

	return ap.audioManager.StopCapture()
}

func (ap *AudioPlayer) Cleanup() {
	ap.Stop()
	ap.audioManager.Cleanup()
}

func (ap *AudioPlayer) IsCapturing() bool {
	return ap.audioManager.IsCapturing()
}

func (ap *AudioPlayer) GetCurrentTrack() string {
	if ap.audioManager.IsCapturing() {
		timeSinceAudio := ap.audioManager.GetTimeSinceLastAudio()
		if timeSinceAudio > 5*time.Second {
			return fmt.Sprintf("Live - No Audio (%.0fs)", timeSinceAudio.Seconds())
		}
		return "Live - System Audio"
	}
	return "Stopped"
}

func (ap *AudioPlayer) GetVolumePercentage() float64 {
	return ap.audioManager.GetPeakLevel()
}

func (ap *AudioPlayer) IncreaseSensitivity() {
	if ap.peakSensitivity < 5.0 {
		ap.peakSensitivity += 0.2
	}
}

func (ap *AudioPlayer) DecreaseSensitivity() {
	if ap.peakSensitivity > 0.2 {
		ap.peakSensitivity -= 0.2
	}
}

func (ap *AudioPlayer) GetSensitivity() float64 {
	return ap.peakSensitivity
}

func (ap *AudioPlayer) Restart() error {
	if ap.IsCapturing() {
		if err := ap.Stop(); err != nil {
			return err
		}
		time.Sleep(100 * time.Millisecond)
	}
	return ap.Start()
}

func (ap *AudioPlayer) CycleDevice() {
	ap.audioManager.CycleDevice()
	if ap.updateInfoFunc != nil {
		ap.updateInfoFunc()
	}
}

func (ap *AudioPlayer) GetCurrentDeviceName() string {
	return ap.audioManager.GetCurrentDeviceName()
}
