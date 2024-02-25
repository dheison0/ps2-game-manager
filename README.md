# PS2 Game Manager

![Sample gif](https://imgur.com/a/SmdsL5s)

A simple [Open PS2 Loader](https://github.com/ps2homebrew/Open-PS2-Loader) game manager create for use on command line,
more specific my server that's where I keep my ROMs

This can:

  - List games;
  - Rename;
  - Delete from disk;
  - Insert game covers and;
  - Install new games;

# Build and run

Install [Go](https://go.dev) 1.21 or newer and run:
```bash
go build .
```

To run just set `ul.cfg` file:
```bash
./ps2manager path/to/ul.cfg
```
