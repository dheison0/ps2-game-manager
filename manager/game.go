package manager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"ps2manager/utils"
	"strings"
)

const COVER_UNFORMATED_URL = "https://raw.githubusercontent.com/xlenore/ps2-covers/main/covers/default/%s.jpg"

type Game struct {
	Config         GameConfig
	Parts          GameFiles
	CoverPath      string
	FilesDirectory string
}

func NewGameFromBytes(data []byte, workDir string) (Game, error) {
	game := Game{FilesDirectory: workDir}
	game.Config.FromBytes(data)
	files, _ := os.ReadDir(workDir)
	for _, f := range files {
		if !f.IsDir() && strings.Contains(f.Name(), game.GetImage()) {
			game.Parts.Files = append(game.Parts.Files, path.Join(workDir, f.Name()))
		}
	}
	game.CoverPath = path.Join(workDir, "ART", game.GetImage()+"_COV.jpg")
	return game, nil
}

func (g *Game) IsCoverInstalled() bool {
	_, err := os.Stat(g.CoverPath)
	return err == nil
}

func (g *Game) DownloadCover() error {
	image := strings.Replace(g.GetImage(), "_", "-", 1)
	image = strings.Replace(image, ".", "", 1)
	response, err := http.Get(fmt.Sprintf(COVER_UNFORMATED_URL, image))
	if err != nil {
		return ErrCoverRequestFailed
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ErrCoverNotFound
	}
	cover, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	return os.WriteFile(g.CoverPath, cover, 0644)
}

func (g *Game) Rename(name string) error {
	nameBytes := [len(g.Config.Name)]byte{}
	copy(nameBytes[:], []byte(name))
	g.Config.Name = nameBytes
	return g.Parts.UpdateHash(name)
}

func (g *Game) GetName() string {
	return utils.BytesToString(g.Config.Name[:])
}

func (g *Game) GetImage() string {
	return utils.BytesToString(g.Config.Image[3:])
}
