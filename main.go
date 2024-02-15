package main

import (
	"fmt"
	"os"
	"ps2manager/manager"
	"ps2manager/tui"

	"github.com/rivo/tview"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Missing config file!")
		os.Exit(1)
	}
	gamesCount := manager.ReadFromFile(os.Args[1])
	if gamesCount == 0 {
		fmt.Println("No games found!")
		os.Exit(2)
	}
	fmt.Printf("%v games found!\n", gamesCount)
	app := tview.NewApplication()
	pages := tview.NewPages()
	pages.AddAndSwitchToPage("Menu", tui.Menu(app, pages), true)
	app.EnableMouse(true)
	app.SetRoot(pages, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
