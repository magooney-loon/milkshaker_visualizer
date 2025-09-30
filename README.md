# MILKSHAKER VISUALIZER

A real-time audio visualizer that responds to system audio or microphone input.

DEMO VIDEO: https://youtu.be/qguHrVe7_T4

## Controls
- `+/-`: Increase/Decrease sensitivity
- `D`: Cycle audio I/O
- `P`: Cycle visualizors
- `X`: Random visualizor
- `Ctrl+C`: Quit

### Audio Issues on Linux
- Check if PulseAudio/PipeWire is running: `systemctl --user status pulseaudio`
- Monitor sources may be suspended - start playing audio to activate them
- Some systems require explicit loopback setup, run the binary with --help flag

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
