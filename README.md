# PS2 Game Manager

![Sample gif](https://imgur.com/a/SmdsL5s)

A simple [Open PS2 Loader](https://github.com/ps2homebrew/Open-PS2-Loader) game manager create for use on command line,
more specififc my server that's where I keep my ROMs

This can:

  - List games;
  - Rename and;
  - Delete from disk.

This will can:

  - Install new games and;
  - Insert game covers.

# Build and run

Install [Go](https://go.dev) 1.21 or newer and run:
```bash
go build .
```

To run just set `ul.cfg` file:
```bash
./ps2manager path/to/ul.cfg
```