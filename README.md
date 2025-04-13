# 🚀 Go Typer

**The sleek, fast terminal typing game inspired by [MonkeyType](https://monkeytype.com/)!**

Go Typer brings the popular web-based typing experience of MonkeyType to your terminal with a beautiful, customizable interface. Master your typing skills right in your terminal (where it actually matters 😉) without a browser.
(online multiplayer type racer, coming soon)

## 🛠️ Built With

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) [![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/spf13/cobra) [![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea) [![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

<h2><sub> 📷 </sub> Screenshots</h2>
<table align="center">
  <tr>
    <td colspan="3"><img src="https://github.com/user-attachments/assets/309aac27-5c1e-4468-947d-123e0b33efe0" alt="Go Typer Main Screen"></td>
  </tr>
  <tr>
    <td><img src="https://github.com/user-attachments/assets/16c52833-0e4e-432e-b7e2-5dc04a60bec2" alt="Go Typer Theme Selection"></td>
    <td><img src="https://github.com/user-attachments/assets/ebebc4d5-b895-48fb-81dd-3b920b5825bd" alt="Go Typer Typing Session"></td>
    <td align="center"><img src="https://github.com/user-attachments/assets/16c52833-0e4e-432e-b7e2-5dc04a60bec2" alt="Go Typer Settings"></td>
  </tr>
</table>

## ✨ Features

- **⚡ Standard-Style Gameplay**: Space bar to advance between words, just like the web favorite!
- **📊 WPM & Accuracy Tracking**: Watch your stats update when you done typing
- **🎮 Multiple Game Modes**: Choose between normal mode (with punctuation) or simple mode for beginners
- **🎨 Gorgeous Themes**: Customize your experience with beautiful color schemes
- **📏 Flexible Text Lengths**: Practice with short, medium, long, or very long passages
- **⚙️ Performance Tuning**: Adjust refresh rates from 1-60 FPS for any terminals (or modify it in code for any value)
- **📝 Cursor Options**: Choose your preferred cursor style (block or underline)
- **💻 100% Terminal-Based**: No browser needed - perfect for developers and terminal enthusiasts.

### Demo video

[![📹 DEMO video](https://github.com/user-attachments/assets/644a3feb-5758-4d3e-bd0d-878abde63787)](https://github.com/user-attachments/assets/644a3feb-5758-4d3e-bd0d-878abde63787)

## 🖥️ Terminal Requirements

Go Typer works best in terminals that support "TrueColor" (24-bit color). It's been tested extensively on Linux but runs great on macOS and Windows too!
(if you see inconsistency with colors, animations or functionality with different terminal emulators please open an issue in this repo)

Verify your terminal supports TrueColor by running:

```bash
printf "\x1b[38;2;255;0;0mTRUECOLOR\x1b[0m\n"
```

If you see "TRUECOLOR" in red, you're good to go\! If not, check out [this compatibility guide](https://gist.github.com/weimeng23/60b51b30eb758bd7a2a648436da1e562).

> [\!CAUTION]
> Please avoid launching the game through `tmux` as it might cause unexpected behavior.

> [\!TIP]
> I recommend using terminal emulators like [`alacritty`](https://github.com/alacritty/alacritty) or [`kitty`](https://github.com/kovidgoyal/kitty), as they are GPU accelerated and generally offer better performance.

## 🚀 Installation

Choose the installation method that suits you best:

<details>
<summary><b>⬇️ Download Binaries (Quickest Start)</b></summary>

Download the latest pre-built binaries for your operating system from the [Releases](https://github.com/prime-run/go-typer/releases) page. Here's a simplified way to download and install (rootless):

**Linux (x86_64):**

```bash
wget https://github.com/prime-run/go-typer/releases/download/v1.0.2/go-typer_1.0.2_linux_x86_64.tar.gz
mkdir -p ~/.local/bin
tar -xzf go-typer_*.tar.gz -C ~/.local/bin go-typer
```

**macOS (Intel x86_64):**

```bash
wget https://github.com/prime-run/go-typer/releases/download/v1.0.2/go-typer_1.0.2_macOS_intel.tar.gz
mkdir -p ~/.local/bin
tar -xzf go-typer_*.tar.gz -C ~/.local/bin go-typer
```

**macOS (Apple Silicon arm64):**

```bash
wget https://github.com/prime-run/go-typer/releases/download/v1.0.2/go-typer_1.0.2_macOS_apple-silicon.tar.gz
mkdir -p ~/.local/bin
tar -xzf go-typer_*.tar.gz -C ~/.local/bin go-typer
```

After downloading and extracting, ensure that `~/.local/bin` is in your system's `PATH` environment variable. You can usually do this by adding the following line to your shell's configuration file (e.g., `.bashrc`, `.zshrc`):

```bash
export PATH="$HOME/.local/bin:$PATH"
```

Then, reload your shell configuration:

```bash
source ~/.bashrc  # For Bash
# or
source ~/.zshrc  # For Zsh
```

Now you should be able to run Go Typer by simply typing `go-typer` in your terminal.

</details>

<details>
<summary><b>⚙️ Go Install (For Go Users)</b></summary>

> [!NOTE]  
> [go](https://go.dev/doc/install) version > v1.24 is required

```bash
go install github.com/prime-run/go-typer@latest
```

Make sure you have Go installed and your `GOPATH/bin` or `GOBIN` is in your system's `PATH`.

</details>

<details>
<summary><b>🛠️ Clone and Build (From Source)</b></summary>

```bash
git clone https://github.com/prime-run/go-typer.git
cd go-typer
go build -o bin/go-typer
./bin/go-typer
```

</details>

<details>
<summary><b>🔨 Make (Unix/Linux)</b></summary>

```bash
git clone https://github.com/prime-run/go-typer.git
cd go-typer
make
./bin/go-typer
```

</details>

<details>
<summary><b>🐳 Docker (Container)</b></summary>

```bash
git clone https://github.com/prime-run/go-typer.git
cd go-typer
docker build -t go-typer .

# Run in container
docker run -it --rm go-typer

```

</details>

## 🎮 How to Play

1.  Launch Go Typer:

    ```bash
    go-typer
    ```

2.  **Navigate** through the menu using the arrow keys or `j`/`k`.

3.  Press **Enter** to select a menu item and start typing.

4.  **Type** the text as displayed. Correctly typed characters will be highlighted.

5.  Press **spacebar** to advance to the next word, just like on MonkeyType\!

6.  Your **WPM**, **accuracy**, and time are tracked in real-time at the bottom of the screen.

7.  Complete the passage to see your final statistics.

### 🎯 Keyboard Controls

- **↑/↓ or j/k**: Navigate through menu items
- **Enter**: Select menu item
- **Esc**: Go back to the previous screen
- **Space**: Advance to the next word while typing
- **Tab**: Restart the current typing exercise
- **q or Ctrl+C**: Quit the application

## ⚙️ Configuration

Go Typer automatically saves your preferences in your user config directory:

- Linux/BSD: `~/.config/go-typer/settings.json`
- macOS: `~/Library/Application Support/go-typer/settings.json`
- Windows: `%AppData%\go-typer\settings.json`

![image](https://github.com/user-attachments/assets/fec6e04c-57d7-4d63-ae24-fc9dff73d923)

You can directly edit the `settings.json` file to customize the following options:

- **theme**: Pick from eye-catching color schemes (`default`, `dark`, `monochrome`) or create your own.
- **cursor_style**: Choose between `block` or `underline`.
- **game_mode**: Select `normal` (with punctuation) or `simple` for beginners.
- **include_numbers**: Set to `true` to include numbers in typing tests.
- **text_length**: Choose from `short`, `medium`, `long`, or `very_long`.
- **refresh_rate**: Fine-tune animation smoothness from `5` (battery-saving) to `60` (ultra-smooth) FPS.

## 🎨 Themes

Go Typer includes beautiful themes inspired by popular coding and typing interfaces.

Built-in themes:

- `default`: Clean light theme with green/blue highlights
- `dark`: Sleek dark theme with purple/blue accents
- `monochrome`: Minimalist black and white theme for distraction-free typing

### 🖌️ Create Your Own Theme

Further customize your typing experience by creating a custom theme file with a `toml` file under `colorschemes` directory and select it within the game settings. Here's the structure of a theme file:

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

**togo**: A terminal-based todo manager built with the same technology stack\!
[Check out togo on GitHub](https://github.com/prime-run/togo)

## 🤝 Contributing

Love Go Typer? Contributions are always welcome\!

- Check out our [todos](https://www.google.com/search?q=todos.md) for upcoming features and areas where you can help.
- Feel free to submit pull requests for bug fixes or new features.
- If you have any suggestions or find any issues, please open an issue on GitHub.

## 📜 License

[MIT](https://www.google.com/search?q=LICENSE)
