package tui

import (
	"fmt"
	"ps2manager/cleaner"
	"ps2manager/utils"

	"github.com/rivo/tview"
)

type CleanerScreen struct {
	root *tview.Modal
	id   string
}

func NewCleanerScreen() *CleanerScreen {
	screen := &CleanerScreen{id: "cleanerScreen"}
	screen.root = tview.NewModal()
	screen.root.
		SetTitle("Clear unused files").
		SetTitleAlign(tview.AlignCenter).
		SetBorder(true)
	pages.AddPage(screen.id, screen.root, true, false)
	return screen
}

func (s *CleanerScreen) Update() {
  s.root.ClearButtons()
	files, err := cleaner.GetUnusedGameFiles()
	if err != nil {
		errorDialog.SetMessage(fmt.Sprintf("Failed to load unused files!\n%s", err.Error()))
		return
	}
	text := "The following files are unused:\n\n"
	var totalSize int64 = 0
	for _, f := range files {
		totalSize += f.Size
		text += fmt.Sprintf("  - [white]%s [cyan]%s\n", f.File, utils.FileSizeToHumanReadable(f.Size))
	}
	if totalSize == 0 {
		s.root.SetText("There's no unused files :)")
		s.root.AddButtons([]string{"Ok"})
		s.root.SetDoneFunc(func(_ int, _ string) { menu.Show() })
		return
	}
	s.root.SetText(text)
  s.root.AddButtons([]string{"Delete", "Cancel"})
  s.root.SetDoneFunc(func(buttonIndex int, buttonLabel string) {
    if buttonIndex == 1 {
      menu.Show()
      return
    }
    if err := cleaner.DeleteFiles(files); err != nil {
      errorDialog.SetMessage("Failed to delete a file!\n"+err.Error())
      errorDialog.Show()
    } else {
      menu.Show()
    }
  })
}

func (s *CleanerScreen) Show() {
	pages.SwitchToPage(s.id)
}
