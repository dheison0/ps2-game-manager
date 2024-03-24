package tui

import "github.com/rivo/tview"

var (
	app             *tview.Application
	pages           *tview.Pages
	menu            *tview.List
	fileSelector    *tview.List
	installForm     *tview.Form
	installProgress *tview.TextView
)

func Init() *tview.Application {
	SetupMenu()
	SetupInstall()
	SetupInstallProgress()
	SetupFileSelector()

	pages = tview.NewPages()
	pages.AddPage("menu", menu, true, true)
	pages.AddPage("fileSelector", fileSelector, true, false)
	pages.AddPage("installForm", installForm, true, false)
	pages.AddPage("installProgress", installProgress, true, false)

	app = tview.NewApplication()
	app.SetRoot(pages, true)
	app.EnableMouse(true)
	return app
}
