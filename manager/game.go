package manager

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"ps2manager/utils"
	"strings"
)

const (
	COVER_UNFORMATTED_URL = "https://raw.githubusercontent.com/xlenore/ps2-covers/main/covers/default/%s.jpg"
	COVER_MAX_WIDTH       = 360
	COVER_MAX_HEIGHT      = 640

	MAX_CD_SIZE        = 700000000  // 700MB    | 667.57MiB
	MAX_GAME_PART_SIZE = 1073741824 // 1.07GB   | 1GiB
	DEFAULT_CHUNK_SIZE = 524288     // 524.28KB | 512KiB
)

type Game struct {
	Config    GameConfig
	Parts     GameFiles
	CoverPath string
}

func NewGame(name, image string, size int64, dataDir string) Game {
	game := Game{}
	copy(game.Config.Name[:], []byte(name))
	copy(game.Config.Image[:], []byte("ul."+image))
	game.Config.Parts = int8(math.Ceil(float64(size) / float64(MAX_GAME_PART_SIZE)))
	game.Config.RegenerateHash()
	game.Parts.GenerateFileNames(game.Config.NameHash, image, dataDir, game.Config.Parts)
	game.CoverPath = path.Join(dataDir, "ART", image+"_COV.jpg")
	if size <= MAX_CD_SIZE {
		game.Config.Media = MediaCD
	} else {
		game.Config.Media = MediaDVD
	}
	return game
}

func NewGameFromBytes(data []byte, workDir string) (Game, error) {
	game := Game{}
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
	gameImage := strings.Replace(g.GetImage(), "_", "-", 1)
	gameImage = strings.Replace(gameImage, ".", "", 1)
	response, err := http.Get(fmt.Sprintf(COVER_UNFORMATTED_URL, gameImage))
	if err != nil {
		return ErrCoverRequestFailed
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ErrCoverNotFound
	}
	coverData, _ := io.ReadAll(response.Body)
	cover := coverData
	coverOriginal, _, err := image.DecodeConfig(bytes.NewReader(coverData))
	if err != nil {
		return err
	}
	if coverOriginal.Width > COVER_MAX_WIDTH || coverOriginal.Height > COVER_MAX_HEIGHT {
		cover, err = utils.ResizeJPG(bytes.NewReader(coverData), COVER_MAX_WIDTH, COVER_MAX_HEIGHT)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(g.CoverPath, cover, 0644)
}

func (g *Game) Rename(name string) error {
	nameBytes := [len(g.Config.Name)]byte{}
	copy(nameBytes[:], []byte(name))
	g.Config.Name = nameBytes
	g.Config.RegenerateHash()
	return g.Parts.UpdateHash(g.Config.NameHash)
}

func (g *Game) GetName() string {
	return utils.BytesToString(g.Config.Name[:])
}

func (g *Game) GetImage() string {
	return utils.BytesToString(g.Config.Image[3:])
}

func (g *Game) GenerateFileNames(root string) {
	for i := int8(0); i < g.Config.Parts; i++ {
		g.Parts.Files = append(
			g.Parts.Files,
			path.Join(root, fmt.Sprintf("ul.%s.%s.%2d", g.Config.NameHash, utils.BytesToString(g.Config.Image[:]), i)),
		)
	}
}
