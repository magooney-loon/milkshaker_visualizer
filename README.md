# MILKSHAKER VISUALIZER

A real-time audio visualizer that responds to system audio or microphone input.

## Controls
- `S`: Start/Stop audio capture
- `R`: Restart audio capture
- `M`: Toggle between Live and Demo mode
- `+/-`: Increase/Decrease sensitivity
- `D`: Cycle audio I/O
- `Ctrl+C`: Quit

### Audio Issues on Linux
- Check if PulseAudio/PipeWire is running: `systemctl --user status pulseaudio`
- Monitor sources may be suspended - start playing audio to activate them
- Some systems require explicit loopback setup (see setup instructions above)

### Build Errors
If you get PortAudio build errors:
1. Install PortAudio development libraries:
   ```bash
   # Ubuntu/Debian
   sudo apt install portaudio19-dev

   # Fedora
   sudo dnf install portaudio-devel

   # Arch
   sudo pacman -S portaudio
   ```
2. On some systems, you may need to set PKG_CONFIG_PATH

## License

MIT License

DEMO: https://www.youtube.com/watch?v=Iq3PFHrFXok

![extras](https://github.com/user-attachments/assets/dbd8940f-a651-446d-98ec-7bb7fe7a4872)
