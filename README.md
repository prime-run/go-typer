# Go Typer

A terminal-based typing game written in Go with bubbletea and lipgloss.

## Features

- Terminal UI typing game with colorful interface
- Real-time WPM calculation
- Support for different cursor styles
- Customizable themes
- Efficient and lightweight

## Screenshots

![Screenshot](https://github.com/user-attachments/assets/5ac1ed81-da75-4222-8bcd-83b96497cb37)

## Requirements

The app has been tested on linux,but should run on major platforms as long as they support "TrueColor".

You can check if your terminal emulator supports 24-bit color by running:

```bash
printf "\x1b[38;2;255;0;0mTRUECOLOR\x1b[0m\n"
```

If your terminal **emulator does NOT display the word** `TRUECOLOR` in **red**, it does not support 24-bit color, checkout [this gist](https://gist.github.com/weimeng23/60b51b30eb758bd7a2a648436da1e562).

## Installation

```
go install github.com/prime-run/go-typer@latest
```

Or clone the repository and build:

```
git clone https://github.com/prime-run/go-typer.git
cd go-typer
go build
```

## Used in this project

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) [![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/spf13/cobra) [![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea) [![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## Usage

Start the typing game:

```
go-typer start
```

### Options

- `-c, --cursor <type>`: Choose cursor type (block or underline)
- `-t, --theme <name>`: Choose a theme or path to custom theme file
- `--list-themes`: List all available themes

Examples:

```
# Start with underline cursor
go-typer start -c underline

# Use dark theme
go-typer start -t dark

# Use a custom theme file
go-typer start -t /path/to/custom/theme.yml

# List available themes
go-typer start --list-themes
```

## Themes

Go Typer supports customizable themes via YAML files in the `colorschemes` directory.

Built-in themes:

- `default`: Standard light theme with green/blue colors
- `dark`: Dark theme with purple/blue accents
- `monochrome`: Minimalist black and white theme

### Creating Custom Themes

You can create your own themes by adding a new YAML file to the `colorschemes` directory or by specifying a custom path with the `-t` flag.

Theme files use the following structure:

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

## Usage and playing

```bash
./go-typer --help
```

### Related

togo: a terminal-based todo manager built with same tech!  
[togo](https://github.com/prime-run/togo)

## Contributing

Contributions are always welcome!

- Please checkout [todos](todos.md) for next steps.
- Features in `dev` branch haven't been tested for different shells and terminal emulators.

## License

[MIT](LICENSE)
