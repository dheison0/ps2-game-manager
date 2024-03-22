# PS2 Game Manager

![Sample gif](https://imgur.com/a/SmdsL5s)

A simple [Open PS2 Loader][OPL] game manager create for use on command line,
more specific my server that's where I keep my ROMs

This program can:

  - View/Rename/Delete installed games;
  - Install new games from ISO;
  - Download and insert game covers(thanks to [ps2-covers]).

# Build and run

Install [Go] 1.21 or newer then run:

```bash
git clone --depth=1 https://github.com/dheison0/ps2-game-manager
cd ps2-game-manager
go build
```

To run just set the game files path:
```bash
./ps2manager games/path
```

[OPL]: <https://github.com/ps2homebrew/Open-PS2-Loader>
[ps2-covers]: <https://github.com/xlenore/ps2-covers>
[Go]: <https://go.dev>
