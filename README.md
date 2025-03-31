# go typer :

a terminal-based typing game.

## Screenshots

![Screenshot](https://github.com/user-attachments/assets/5ac1ed81-da75-4222-8bcd-83b96497cb37)

## Requirements

- go version > 1.24.1
- modern shell, (tested on zsh and bash)
- terminal emulator, with TRUE color support
  Test TRUE color support by runnig

```shell
printf "\x1b[38;2;255;100;0mTRUECOLOR\x1b[0m\n"
```

If your terminal **emulator does NOT display the word** `TRUECOLOR` in **red**, it does not support 24-bit color, checkout [this gist](https://gist.github.com/weimeng23/60b51b30eb758bd7a2a648436da1e562).

## Installation

### Build

```bash
git clone https://github.com/prime-run/go-typer
cd go-typer
#to pick up the the ongoing changes
git checkout -b dev

go build -o go-typer
```

### docker

```bash
git clone https://github.com/prime-run/go-typer
docker build -t go-typer .
docker run --rm -it my-go-app
# above command Runs ./go-typer start inside the container
# docker run --rm -it go-typer [command] [--flag] also works

```

## Used in this project

[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev/) [![Cobra](https://img.shields.io/badge/Cobra-00ADD8?style=flat-square&logo=go&logoColor=white)](https://github.com/spf13/cobra) [![Bubble Tea](https://img.shields.io/badge/Bubble%20Tea-FF75B7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/bubbletea) [![Lip Gloss](https://img.shields.io/badge/Lip%20Gloss-FFABE7?style=flat-square&logo=go&logoColor=white)](https://github.com/charmbracelet/lipgloss)

## Usage and playing

```bash
./go-typer --help
```

#### start the game :

```bash
./go-typer start
```

and start typing

- timer starts when you type first letter
- press `Tab` to rest the same passage
- press `Esc` or `Cltr + c ` to exit

### Related

togo: a terminal-based todo manager built with same tech!  
[togo](https://github.com/prime-run/togo)

## Contributing

Contributions are always welcome!

- Please checkout `todo.md` for next steps.
- Features in `dev` branch haven't been tested for different shells and terminal emulators.

## License

[MIT](License)
please also checkout [todos](todos.md)
