package tui

import (
	"ps2manager/manager"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type CoverDownloadScreen struct {
	root *tview.TextView
}

func NewCoverDownloadScreen() *CoverDownloadScreen {
	screen := &CoverDownloadScreen{}
	screen.root = tview.NewTextView()
	screen.root.
		SetBorder(true).
		SetBorderPadding(1, 1, 2, 2).
		SetTitle(" Downloading Covers ")
	screen.root.
		SetTextAlign(tview.AlignCenter).
		SetDynamicColors(true)
	pages.AddPage("covers", screen.root, true, false)
	return screen
}

func (c *CoverDownloadScreen) Show() {
	pages.SwitchToPage("covers")
}

func (c *CoverDownloadScreen) AppendText(text string) {
	oldText := c.root.GetText(false)
	c.root.SetText(oldText + text)
	app.ForceDraw()
}

func (c *CoverDownloadScreen) DownloadMissingCovers() {
	games := manager.GetAll()
	var missingCovers []*manager.GameConfig
	for _, g := range games {
		if !g.IsCoverInstalled() {
			missingCovers = append(missingCovers, g)
		}
	}
	c.root.SetDoneFunc(nil)
	c.root.SetText("")
	defer c.AppendText("\n\n[blue]Press enter to go back...[default]")
	defer c.root.SetDoneFunc(func(_ tcell.Key) { menu.Show() })
	if len(missingCovers) == 0 {
		c.AppendText("[purple]There's no missing covers :)[default]")
		return
	}
	for _, g := range missingCovers {
		c.AppendText("\nDownload cover for [purple]" + g.GetName() + "[default]...")
		if err := g.DownloadCover(); err != nil {
			c.AppendText(" [red]Error: " + err.Error() + "[default]")
		}
	}
	c.AppendText("\nComplete!")
}
