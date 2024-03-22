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
	list.AddItem("Install", "Install new game", ',', InstallForm)
	list.AddItem("Get covers", "Download missing game covers", '.', DownloadCovers)
	list.AddItem("Quit", "Exit program", ';', app.Stop)
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
	var toDownload []*manager.GameConfig
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

func InstallForm() {
	form := tview.NewForm()
	form.SetBorder(true)
	form.SetTitle("Install new game")

	iso, name := "", ""
	form.AddInputField("Iso file:", iso, 0, nil, func(t string) { iso = t })
	form.AddInputField("Game name:", name, manager.MaxNameSize, nil, func(t string) { name = t })
	form.AddButton("Install", func() {
		progress := make(chan int, 100)
		errChan := make(chan error)
		textView := tview.NewTextView()
		textView.SetBorder(true)
		textView.SetBorderPadding(2, 2, 2, 2)
		textView.SetTitle("Install progress")
		textView.SetTextAlign(tview.AlignCenter)
		textView.SetText("Initializing installation...")
		textView.SetDynamicColors(true)
		pages.AddAndSwitchToPage("InstallProgress", textView, true)
		go (func() { errChan <- manager.Install(iso, name, progress) })()
		for {
			select {
			case percent := <-progress:
				textView.SetText(fmt.Sprintf("Installing '[cyan]%s[white]' is [blue]%d%%[white] done...", name, percent))
				app.ForceDraw()
			case err := <-errChan:
				if err == nil {
					goto complete
				}
				ErrorDialog(fmt.Sprintf("Failed to install '%s': %v", name, err))
				return
			}
		}
	complete:
		modal := tview.NewModal()
		modal.SetText(fmt.Sprintf("Installation of '%s' was done with success!", name))
		modal.AddButtons([]string{"Done"})
		modal.SetDoneFunc(func(_ int, _ string) { Menu() })
		pages.AddAndSwitchToPage("InstallDone", modal, true)
	})
	form.AddButton("Cancel", func() { Menu() })
	pages.AddAndSwitchToPage("InstallForm", form, true)
}
