package tui

import "github.com/rivo/tview"

var (
	app             *tview.Application
	pages           *tview.Pages
	menu            *tview.List
	fileSelector    *tview.List
	installForm     *tview.Form
	installProgress *tview.Modal
)

func Init() *tview.Application {
	SetupMenu()
	SetupInstall()
	SetupFileSelector()

	pages = tview.NewPages()
	pages.AddPage("menu", menu, true, true)
	pages.AddPage("fileSelector", fileSelector, true, false)
	pages.AddPage("installForm", installForm, true, false)
	pages.AddPage("installProgress", installProgress, true, false)

	app = tview.NewApplication()
	app.SetRoot(pages, true)
	return app
}
