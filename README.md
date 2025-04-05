
# 🚀 Go Typer

**The sleek, blazing-fast terminal typing test inspired by [MonkeyType](https://monkeytype.com/)!**

Go Typer brings the popular web-based typing experience of MonkeyType to your terminal with a beautiful, customizable interface. Master your typing skills right in your command line without the distractions of a browser.

![250405_12h57m34s_screenshot](https://github.com/user-attachments/assets/309aac27-5c1e-4468-947d-123e0b33efe0)
![250405_12h57m45s_screenshot](https://github.com/user-attachments/assets/ebebc4d5-b895-48fb-81dd-3b920b5825bd)
![250405_12h58m44s_screenshot](https://github.com/user-attachments/assets/16c52833-0e4e-432e-b7e2-5dc04a60bec2)

## ✨ Features

- **⚡ MonkeyType-Style Gameplay**: Space bar to advance between words, just like the web favorite!
- **📊 Real-time WPM & Accuracy Tracking**: Watch your stats update live as you type
- **🎮 Multiple Game Modes**: Choose between normal mode (with punctuation) or simple mode for beginners
- **🎨 Gorgeous Themes**: Customize your experience with beautiful color schemes
- **📏 Flexible Text Lengths**: Practice with short, medium, long, or very long passages
- **⚙️ Performance Tuning**: Adjust refresh rates from 5-60 FPS for any hardware
- **🌈 Stunning Gradient Animations**: Eye-catching color flows throughout the interface
- **📝 Cursor Options**: Choose your preferred cursor style (block or underline)
- **💻 100% Terminal-Based**: No browser needed - perfect for developers and terminal enthusiasts

## 🖥️ Terminal Requirements

Go Typer works best in terminals that support "TrueColor" (24-bit color). It's been tested extensively on Linux but runs great on macOS and Windows too!

Verify your terminal supports TrueColor by running:

```bash
printf "\x1b[38;2;255;0;0mTRUECOLOR\x1b[0m\n"
```

If you see "TRUECOLOR" in red, you're good to go! If not, check out [this compatibility guide](https://gist.github.com/weimeng23/60b51b30eb758bd7a2a648436da1e562).

## 🚀 Installation

### Download Binaries (Quickest Start)

Grab pre-built binaries from the [Releases](https://github.com/prime-run/go-typer/releases) page for instant typing pleasure!

### Go Install (For Go Users)

```bash
go install github.com/prime-run/go-typer@latest
```

### Clone and Build (From Source)

```bash
git clone https://github.com/prime-run/go-typer.git
cd go-typer
go build -o bin/go-typer
./bin/go-typer
```

### Make (Unix/Linux)

```bash
git clone https://github.com/prime-run/go-typer.git
cd go-typer
make build
./bin/go-typer
```

### Docker (Container)

```bash
# Pull from Docker Hub
docker pull primerun/go-typer:latest

# Run in container
docker run -it --rm primerun/go-typer:latest

# Build locally
git clone https://github.com/prime-run/go-typer.git
cd go-typer
docker build -t go-typer .
docker run -it --rm go-typer
```

## 🛠️ Built With

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) [![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/spf13/cobra) [![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea) [![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## 🎮 How to Play

1. Launch Go Typer:
   ```bash
   go-typer
   ```

2. **Navigate** through the menu to start typing
3. **Type** the text as displayed - it's colored to show you what's correct and incorrect
4. Press **spacebar** to advance to the next word (just like on MonkeyType!)
5. Your **WPM**, **accuracy**, and timing are tracked in real-time
6. Complete the passage to see your final stats

### 🎯 Keyboard Controls

- **↑/↓ or j/k**: Navigate through menu items
- **Enter**: Select menu item
- **Esc**: Go back to previous screen
- **Space**: Advance to the next word while typing
- **Tab**: Restart current typing exercise
- **q or Ctrl+C**: Quit the application

## ⚙️ Configuration

Go-Typer automatically saves your preferences in your user config directory:
- Linux/BSD: `~/.config/go-typer/settings.json`
- macOS: `~/Library/Application Support/go-typer/settings.json`
- Windows: `%AppData%\go-typer\settings.json`

### 🔧 Customize Your Experience

- **Theme**: Pick from eye-catching color schemes
- **Cursor Style**: Block or underline, based on your preference
- **Game Mode**: Normal mode (with punctuation) or Simple mode for beginners
- **Numbers**: Toggle inclusion of numbers in typing tests
- **Text Length**: Choose from Short (1 quote), Medium (2 quotes), Long (3 quotes), or Very Long (5 quotes)
- **Refresh Rate**: Fine-tune animation smoothness from 5 FPS (battery-saving) to 60 FPS (ultra-smooth)

## 🎨 Themes

Go Typer includes beautiful themes inspired by popular coding and typing interfaces.

Built-in themes:

- `default`: Clean light theme with green/blue highlights
- `dark`: Sleek dark theme with purple/blue accents
- `monochrome`: Minimalist black and white theme for distraction-free typing

### 🖌️ Create Your Own Theme

Customize your typing experience by creating a theme file:

```yaml
# UI Elements
help_text: "#626262" # Help text at the bottom
timer: "#FFDB58" # Timer display
border: "#7F9ABE" # Border color for containers

# Text Display
text_dim: "#555555" # Untyped text
text_preview: "#7F9ABE" # Preview text color
text_correct: "#00FF00" # Correctly typed text
text_error: "#FF0000" # Incorrectly typed characters
text_partial_error: "#FF8C00" # Correct characters in error words

# Cursor
cursor_fg: "#FFFFFF" # Cursor foreground color
cursor_bg: "#00AAFF" # Cursor background color
cursor_underline: "#00AAFF" # Underline cursor color

# Miscellaneous
padding: "#888888" # Padding elements color
```

## 🔄 Related Projects

**togo**: A terminal-based todo manager built with the same technology stack!  
[Check out togo on GitHub](https://github.com/prime-run/togo)

## 🤝 Contributing

Love Go Typer? Contributions are always welcome!

- Check out our [todos](todos.md) for upcoming features


## 📜 License

[MIT](LICENSE)
