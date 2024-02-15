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
	games := manager.ReadFromFile(os.Args[1])
	fmt.Printf("%v games found!\n", len(games))
	root := tui.Menu(games)
	app := tview.NewApplication().SetRoot(root, true)
	if err := app.Run(); err != nil {
		panic(err)
	}
}
