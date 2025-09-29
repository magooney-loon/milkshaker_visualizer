package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"milkshaker/audio"
	"milkshaker/patterns"

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
		case "test-monitor":
			testMonitorSource()
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
	tester := audio.NewTester()
	tester.Run()
}

func testMonitorSource() {
	fmt.Println("MONITOR SOURCE TEST")
	fmt.Println("===================")
	fmt.Println("This will test if the monitor source is properly configured.")
	fmt.Println()

	player := audio.NewPlayer()

	// Initialize but don't start the full visualizer
	if err := player.Initialize(); err != nil {
		fmt.Printf("❌ Failed to initialize: %v\n", err)
		return
	}
	defer player.Cleanup()

	fmt.Println("\nTesting audio capture from configured source...")
	fmt.Println("Play some music and you should see audio levels below:")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Start capture
	if err := player.Start(); err != nil {
		fmt.Printf("❌ Failed to start capture: %v\n", err)
		return
	}

	// Monitor for a bit
	for i := 0; i < 100; i++ {
		time.Sleep(100 * time.Millisecond)
		peak := player.GetPeakLevel()

		if peak > 0.01 {
			fmt.Printf("🎵 STRONG audio detected: %.3f\n", peak)
		} else if peak > 0.001 {
			fmt.Printf("🔉 Medium audio detected: %.3f\n", peak)
		} else if peak > 0.0001 {
			fmt.Printf("🔈 Low audio detected: %.6f\n", peak)
		} else {
			fmt.Printf("🔇 No audio: %.8f\n", peak)
		}
	}

	player.Stop()
	fmt.Println("\nTest completed!")
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
	fmt.Println("  go run . test-monitor    # Test monitor source configuration")
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
	player := audio.NewPlayer()

	if err := player.Initialize(); err != nil {
		log.Fatalf("Failed to initialize audio player: %v", err)
	}
	defer player.Cleanup()

	app := tview.NewApplication()
	visualizer := patterns.NewFibonacciVisualizer()

	app.SetAfterDrawFunc(func(screen tcell.Screen) {
		width, height := screen.Size()
		visualizer.SetRect(0, 0, width, height)
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

	player.SetUpdateInfoFunc(updateInfo)
	fullScreenVisualizer := tview.NewBox().SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		visualizer.SetRect(x, y, width, height)

		// Update visualizer with current audio peak
		peak := player.GetPeakLevel()
		visualizer.UpdateWithPeak(peak)

		visualizer.Draw(screen)
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
