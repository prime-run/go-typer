## go typer: a terminal-based typing game

under construction ...

- ### requirements:
  `go version > go1.24.1`,
  `modern shell, (tested on zsh and bash)`,
  `terminal emulator, with TRUE color support `
  to test color support, run:

```bash
printf "\x1b[38;2;255;100;0mTRUECOLOR\x1b[0m\n"
```

If your terminal emulator does NOT display the word TRUECOLOR in red, it does not support 24-bit color. If you don't want to switch to a different terminal emulator that supports 24-bit color, checkout [this gist](https://gist.github.com/weimeng23/60b51b30eb758bd7a2a648436da1e562) .

- ## INSTALLATION:

```bash
git clone https://github.com/prime-run/go-typer
cd go-typer
#to pick up the the ongoing changes
git checkout -b dev

go build -o go-typer
```

Usage:

```bash
./go-typer --help
# to start the game demo
./go-typer start
```

please also checkout [todos](todos.md)
