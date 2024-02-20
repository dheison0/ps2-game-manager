package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/gdamore/tcell/v2"
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
	textView := tview.NewTextView().
		SetRegions(true).
		SetTextAlign(tview.AlignCenter).
		SetChangedFunc(func() { app.ForceDraw() }).
		SetDynamicColors(true)
	textView.SetTitle("Downloading covers...").
		SetBorder(true).
		SetBorderPadding(2, 2, 2, 2)
	pages.AddAndSwitchToPage("DownloadProgress", textView, true)

	games := manager.GetAll()
	var toDownload []manager.Game
	for _, g := range games {
		if !g.IsCoverInstalled() {
			toDownload = append(toDownload, g)
		}
	}

	completed := 0
	errors := 0

	for _, g := range toDownload {
		textView.SetText(
			fmt.Sprintf(
				"Download covers for [cyan]%d[white] games\nSuccess: [green]%d[white]  Errors: [red]%d",
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

	textView.SetDoneFunc(func(_ tcell.Key) { Menu() })
	textView.SetText(
		fmt.Sprintf(
			`Download covers for [cyan]%d[white] games
Success: [green]%d[white]  Errors: [red]%d

[blue]Press enter to go back...`,
			len(toDownload), completed, errors,
		),
	)
}
