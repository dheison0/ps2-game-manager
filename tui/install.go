package tui

import (
	"fmt"
	"ps2manager/manager"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func SetupInstall() {
	installForm = tview.NewForm()
	installForm.SetBorder(true)
	installForm.SetBorderPadding(1, 1, 1, 1)
	installForm.SetTitle(" Install new game ")
}

func UpdateInstallForm(isoFile string) {
	gameName := ""
	installForm.
		Clear(true).
		AddTextView("ISO file:", isoFile, 0, 1, false, false).
		AddInputField(
			"Game Name:",
			gameName,
			manager.MaxNameSize,
			func(t string, _ rune) bool { return len(t) <= manager.MaxNameSize },
			func(t string) { gameName = t },
		).
		AddButton("Install", func() {
			progress := make(chan int, 100)
			iError := make(chan error)
			go func() { iError <- manager.Install(isoFile, gameName, progress) }()
			pages.SwitchToPage("installProgress")
			InstallUpdateProgress(gameName, progress, iError)
		}).
		AddButton("Cancel", func() { pages.SwitchToPage("menu") })
}

func SetupInstallProgress() {
	installProgress = tview.NewTextView()
	installProgress.SetTitle(" Installation progress ")
	installProgress.SetBorder(true)
	installProgress.SetBorderPadding(3, 3, 3, 3)
	installProgress.SetTextAlign(tview.AlignCenter)
}

func InstallUpdateProgress(gameName string, progress chan int, err chan error) {
	installProgress.SetText("Installation of " + gameName + " is starting...")
	installProgress.SetDoneFunc(nil)
	for {
		select {
		case iError := <-err: // TODO: Check if there's an error on installation process
			if iError == nil {
				goto installationComplete
			}
			pages.SwitchToPage("menu")
			return
		case percent := <-progress:
			installProgress.SetText(fmt.Sprintf("Installation of %s is %d%% complete...", gameName, percent))
			app.ForceDraw()
		}
	}
installationComplete:
	installProgress.SetText(gameName + " was installed with success!\nPress any key to go back...")
	installProgress.SetDoneFunc(func(_ tcell.Key) { pages.SwitchToPage("menu") })
}
