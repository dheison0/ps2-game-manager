# PS2 Game Manager

![ps2manager](https://github.com/dheison0/ps2-game-manager/assets/98592255/3dbb1fdf-70a0-4091-a037-1e262fb1e291)

A simple tool for managing your PlayStation 2 game collection in a clean and easy way.

## What it is

PS2 Game Manager helps you keep your PS2 games organized in a folder prepared for [Open PS2 Loader][OPL]. It is designed for people who want a quick and straightforward way to manage their library from a terminal.

## Main features

- Install games from ISO files
- Keep your collection organized
- Rename games when needed
- Remove games you no longer want
- Download cover art for the library
- Export games back to ISO format
- Clean up unused files from the game folder

## Supported platforms

This project is intended for Unix-like systems, such as:

- Linux
- macOS

Windows is not officially supported at the moment, and the application is designed around Unix-style paths and filesystem behavior.

## Local build

Requirements:

- [Go 1.22+](https://go.dev)

Clone the project and build it locally:

```bash
git clone --depth=1 https://github.com/dheison0/ps2-game-manager
cd ps2-game-manager
go build -o ps2manager .
```

Run it with the folder where your games are stored:

```bash
./ps2manager /path/to/your/games
```

[OPL]: <https://github.com/ps2homebrew/Open-PS2-Loader>
