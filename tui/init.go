package tui

import "github.com/rivo/tview"

var (
	// Base
	app   *tview.Application
	pages *tview.Pages

	// Screens
	menu               *MenuScreen
	fileSelector       *FileSelectorScreen
	errorDialog        *ErrorDialogScreen
	install            *InstallScreen
	installProgress    *InstallProgressScreen
	covers             *CoverDownloadScreen
	actionsMenu        *ActionsMenuScreen
	actionsDelete      *ActionsDeleteScreen
	actionsRename      *ActionsRenameScreen
	gameExport         *GameExporterScreen
	gameExportProgress *GameExportProgressScreen
	cleanerScreen      *CleanerScreen
)

func Init() *tview.Application {
	app = tview.NewApplication()
	pages = tview.NewPages()
	menu = NewMenuScreen()
	fileSelector = NewFileSelectorScreen()
	errorDialog = NewErrorDialogScreen()
	install = NewInstallScreen()
	installProgress = NewInstallProgressScreen()
	covers = NewCoverDownloadScreen()
	actionsMenu = NewActionsMenuScreen()
	actionsDelete = NewActionsDeleteScreen()
	actionsRename = NewActionsRenameScreen()
	gameExport = NewGameExporterScreen()
	gameExportProgress = NewGameExportProgressScreen()
	cleanerScreen = NewCleanerScreen()

	app.SetRoot(pages, true)
	app.EnableMouse(true)

	menu.Show()
	return app
}
