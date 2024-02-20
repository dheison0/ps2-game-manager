package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/rivo/tview"
)

var app *tview.Application
var pages *tview.Pages

func Init() *tview.Application {
	app = tview.NewApplication()
	pages = tview.NewPages()
	app.EnableMouse(true)
	app.SetRoot(pages, true)
	Menu()
	return app
}

func Menu() {
	games := manager.GetAll()
	list := tview.NewList()
	list.SetTitle(" PS2 Game Manager ")
	list.SetBorder(true)
	gamePage := func() {
		actions := GameActions(list.GetCurrentItem())
		pages.AddAndSwitchToPage("Game", actions, true)
	}
	for i, game := range games {
		description := "Image: " + game.GetImage()
		if !game.IsCoverInstalled() {
			description += "  !!! Cover is missing !!!"
		}
		list.AddItem(game.GetName(), description, rune(i+'a'), gamePage)
	}
	list.AddItem("Get covers", "Download missing game covers", '"', DownloadCovers)
	list.AddItem("Quit", "Exit program", '!', app.Stop)
	pages.AddAndSwitchToPage("Menu", list, true)
}

func DownloadCovers() {
	games := manager.GetAll()
	textBox := tview.NewTextView()
	textBox.SetChangedFunc(func() {
		app.Draw()
	})
	var toDownload []manager.Game
	for _, g := range games {
		if !g.IsCoverInstalled() {
			toDownload = append(toDownload, g)
		}
	}
	textBox.SetTitle("Downloading covers...")
	pages.AddAndSwitchToPage("DownloadProgress", textBox, true)
	completed := 0
	errors := 0
	for _, g := range toDownload {
		textBox.SetText(
			fmt.Sprintf(
				"Download covers for %d games\nSuccess: %d  Errors: %d",
				len(toDownload), completed, errors,
			),
		)
		err := g.DownloadCover()
		if err == nil {
			completed += 1
		} else {
			errors += 1
		}
	}
	dialog := tview.NewModal()
	dialog.SetText(fmt.Sprintf("Download process finished with %d errors and %d success.", errors, completed))
	dialog.AddButtons([]string{"Done"})
	dialog.SetDoneFunc(func(_ int, _ string) { Menu() })
	pages.AddAndSwitchToPage("Download", dialog, true)
}
