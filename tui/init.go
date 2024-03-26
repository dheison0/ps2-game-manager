package tui

import "github.com/rivo/tview"

var (
	// Base
	app   *tview.Application
	pages *tview.Pages

	// Screens
	menu            *MenuScreen
	fileSelector    *FileSelectorScreen
	errorDialog     *ErrorDialogScreen
	install         *InstallScreen
	installProgress *InstallProgressScreen
)

func Init() *tview.Application {
	app = tview.NewApplication()
	pages = tview.NewPages()
	menu = NewMenuScreen()
	fileSelector = NewFileSelectorScreen()
	errorDialog = NewErrorDialogScreen()
	install = NewInstallScreen()
	installProgress = NewInstallProgressScreen()

	app.SetRoot(pages, true)
	app.EnableMouse(true)

	menu.Show()
	return app
}