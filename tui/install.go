package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type InstallScreen struct {
	root *tview.Form
}

func NewInstallScreen() *InstallScreen {
	screen := &InstallScreen{}
	screen.root = tview.NewForm()
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 1, 1).
		SetTitle(" Install new game ")
	pages.AddPage("install", screen.root, true, false)
	return screen
}

func (i *InstallScreen) Show() {
	pages.SwitchToPage("install")
}

func (i *InstallScreen) NewForm(isoFile string) {
	gameName := ""
	i.root.
		Clear(true).
		AddTextView("ISO file:", isoFile, 0, 1, false, false).
		AddInputField(
			"Game Name:",
			gameName,
			manager.MaxNameSize,
			manager.CheckIfAcceptName,
			func(t string) { gameName = t },
		).
		AddButton("Install", func() {
			progress := make(chan int, 100)
			iError := make(chan error)
			go func() { iError <- manager.Install(isoFile, gameName, progress) }()
			installProgress.Show()
			installProgress.SetProgressSource(gameName, progress, iError)
		}).
		AddButton("Cancel", menu.Show)
}

type InstallProgressScreen struct {
	root *tview.TextView
}

func NewInstallProgressScreen() *InstallProgressScreen {
	screen := &InstallProgressScreen{}
	screen.root = tview.NewTextView()
	screen.root.
		SetDynamicColors(true).
		SetTextAlign(tview.AlignCenter).
		SetTitle(" Installation progress ").
		SetBorder(true).
		SetBorderPadding(3, 3, 3, 3)
	pages.AddPage("installProgress", screen.root, true, false)
	return screen
}

func (s *InstallProgressScreen) Show() {
	pages.SwitchToPage("installProgress")
}

func (s *InstallProgressScreen) SetProgressSource(gameName string, progress chan int, err chan error) {
	s.root.
		SetText("[white]Installation of [purple]" + gameName + "[white] is starting...").
		SetDoneFunc(nil)
	for {
		app.ForceDraw()
		select {
		case iError := <-err:
			if iError != nil {
				errorDialog.Show()
				errorDialog.SetMessage(fmt.Sprintf("[white]Failed to install [purple]%s:\n[red]%s", gameName, iError.Error()))
				return
			}
			goto installationComplete
		case percent := <-progress:
			s.root.SetText(fmt.Sprintf(
				"[white]Installation of [purple]%s[white] is [blue]%d%%[white] complete...",
				gameName, percent,
			))
		}
	}
installationComplete:
	s.root.
		SetText("[purple]" + gameName + "[white] was installed with success!\n\n[green]Press any key to go back...").
		SetDoneFunc(func(_ tcell.Key) {
			menu.UpdateItemList()
			menu.Show()
		})
}
