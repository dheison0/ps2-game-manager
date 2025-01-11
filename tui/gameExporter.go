package tui

import (
	"fmt"
	"path"
	"ps2manager/manager"
	"ps2manager/utils"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type GameExporterScreen struct {
	ID               string
	root             *tview.Flex
	info             *tview.TextView
	pathSelectorView *tview.TextView
	nameForm         *tview.Form

	fileName     string
	outputFolder string

	game *manager.GameConfig
}

func NewGameExporterScreen() *GameExporterScreen {
	screen := GameExporterScreen{ID: "exporter", outputFolder: "."}
	screen.root = tview.NewFlex()
	screen.root.
		SetDirection(tview.FlexRow).
		SetBorderPadding(1, 1, 2, 2).
		SetBorder(true).
		SetTitle(" Export game as ISO ").
		SetTitleAlign(tview.AlignCenter)

	screen.info = tview.NewTextView().SetDynamicColors(true)

	screen.pathSelectorView = tview.NewTextView()
	pathSelector := tview.NewFlex().
		SetDirection(tview.FlexRowCSS).
		AddItem(screen.pathSelectorView, 0, 1, false).
		AddItem(
			tview.NewButton("Change").SetSelectedFunc(func() {
				fileSelector.SetSelectFileConfig(screen.onOutputFolderChanged, true, "")
				fileSelector.Show()
			}),
			8, 0, false,
		)

	screen.nameForm = tview.NewForm().
		AddButton("Export", screen.onExportPress).
		AddButton("Cancel", actionsMenu.Show)

	screen.root.
		AddItem(screen.info, 4, 1, false).
		AddItem(pathSelector, 1, 0, false).
		AddItem(screen.nameForm, 0, 1, true)

	pages.AddPage(screen.ID, screen.root, true, false)
	return &screen
}

func (s *GameExporterScreen) onExportPress() {
	outputFile := path.Join(s.outputFolder, s.fileName)
	progress := make(chan int)
	errChan := make(chan error)
	go s.game.ExportAsISO(outputFile, progress, errChan)
	gameExportProgress.Show()
	gameExportProgress.ShowProgress(outputFile, progress, errChan)
}

func (s *GameExporterScreen) onOutputFolderChanged(newFolder string) {
	s.outputFolder = newFolder
	s.updateScreen()
	s.Show()
}

func (s *GameExporterScreen) Show() {
	pages.SwitchToPage(s.ID)
}

func (s *GameExporterScreen) updateScreen() {
	size, err := utils.GetFilesSizeSum(s.game.Files)
	if err != nil {
		errorDialog.SetMessage("Error loading game info, error:\n" + err.Error())
		errorDialog.Show()
		return
	}

	s.info.SetText(fmt.Sprintf(
		"[white]Name: [purple]%s\n[white]Size: [purple]%s",
		s.game.GetName(),
		utils.FileSizeToHumanReadable(size),
	))

	s.pathSelectorView.SetText("Output folder: " + s.outputFolder)

	s.fileName = s.game.GetName() + ".iso"
	s.nameForm.
		Clear(false).
		AddInputField(
			"File name:",
			s.fileName, 0,
			func(_ string, r rune) bool { return r != '/' && r != '\\' },
			func(t string) { s.fileName = t },
		)
}

func (s *GameExporterScreen) SetGame(game *manager.GameConfig) {
	s.game = game
	s.updateScreen()
}

type GameExportProgressScreen struct {
	ID   string
	root *tview.TextView
}

func NewGameExportProgressScreen() *GameExportProgressScreen {
	screen := &GameExportProgressScreen{ID: "exportProgress"}
	screen.root = tview.NewTextView()
	screen.root.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true).
		SetBorder(true).
		SetBorderPadding(3, 3, 3, 3).
		SetTitle(" Exporting game... ")
	pages.AddPage(screen.ID, screen.root, true, false)
	return screen
}

func (s *GameExportProgressScreen) ShowProgress(fileName string, progress chan int, errChan chan error) {
	s.root.SetText(fmt.Sprintf("[white]Exportation of [purple]%s[white] is starting...", fileName)).SetDoneFunc(nil)
	for {
		app.ForceDraw()
		select {
		case p := <-progress:
			if p == 100 {
				goto exportationComplete
			}
			s.root.SetText(fmt.Sprintf("[white]Exportation of '[purple]%s[white]' is [purple]%d%%[white] done...", fileName, p))
		case err := <-errChan:
			errorDialog.SetMessage(fmt.Sprintf("Failed to export game:\n%s", err.Error()))
			errorDialog.Show()
		}
	}
exportationComplete:
	s.root.
		SetText(fmt.Sprintf(
			"[white]'[purple]%s[white]' was exported with success!\n\n[green]Press any key to go back...",
			fileName,
		)).
		SetDoneFunc(func(_ tcell.Key) { actionsMenu.Show() })
}

func (s *GameExportProgressScreen) Show() {
	pages.SwitchToPage(s.ID)
}
