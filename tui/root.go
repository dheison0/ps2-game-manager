package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

func TUI() *tview.Application {
	app = tview.NewApplication()
	pages = tview.NewPages()
	RefreshPages()
	app.EnableMouse(true)
	app.SetRoot(pages, true)
	return app
}

func RefreshPages() {
	pages.AddPage("Menu", Menu(), true, true)
}

func Menu() tview.Primitive {
	games := manager.GetAllGames()
	list := tview.NewList()
	list.SetTitle("PS2 Game Manager")
	list.SetBorder(true)
	for i := range games {
		game := games[i]
		list.AddItem(string(game.Name[:]), fmt.Sprintf("Image: %s", game.Image), rune(i+65), func() {
			index := list.GetCurrentItem()
			gamePage := GameActions(index)
			pages.AddAndSwitchToPage("Game", gamePage, true)
		})
	}
	list.AddItem("Quit", "Exit program", '"', func() {
		app.Stop()
	})
	return list
}
