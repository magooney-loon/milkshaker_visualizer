package audio

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Tester provides simple audio testing functionality
type Tester struct {
	player *Player
}

// NewTester creates a new audio tester
func NewTester() *Tester {
	return &Tester{
		player: NewPlayer(),
	}
}

// Run starts the audio test
func (t *Tester) Run() {
	fmt.Println("AUDIO CAPTURE TEST")
	fmt.Println("==================")
	fmt.Println("This will test if audio capture is working without the full visualizer.")
	fmt.Println("Press Ctrl+C to stop the test.")
	fmt.Println()

	// Setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Initialize and start
	err := t.player.Initialize()
	if err != nil {
		fmt.Printf("Failed to initialize: %v\n", err)
		return
	}
	defer t.player.Cleanup()

	err = t.player.Start()
	if err != nil {
		fmt.Printf("Failed to start: %v\n", err)
		return
	}

	fmt.Println("🎤 Listening for audio... (play some music to test)")

	// Monitor audio levels
	go t.monitorAudio()

	// Wait for interrupt
	<-sigChan
	fmt.Println("\n\nStopping audio test...")
	t.player.Stop()
}

// monitorAudio prints real-time audio levels
func (t *Tester) monitorAudio() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		if !t.player.IsCapturing() {
			continue
		}

		peak := t.player.GetPeakLevel()
		now := time.Now()

		if peak > 0.001 {
			fmt.Printf("\r🎵 STRONG: Peak=%.6f | %s", peak, now.Format("15:04:05"))
		} else if peak > 0.0001 {
			fmt.Printf("\r🔉 Medium: Peak=%.6f | %s", peak, now.Format("15:04:05"))
		} else if peak > 0.00001 {
			fmt.Printf("\r🔈 Low: Peak=%.8f | %s", peak, now.Format("15:04:05"))
		} else {
			fmt.Printf("\r🔇 Silent: Peak=%.10f | %s", peak, now.Format("15:04:05"))
		}
	}
}
