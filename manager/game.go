package manager

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"ps2manager/utils"
	"ps2manager/utils/oplCRC32"
	"slices"
	"strings"
)

const ( // Config data
	MediaCD     = 0x12
	MediaDVD    = 0x14
	PaddingByte = 0x08
)

const ( // Sizes
	MaxCDSize         = 700000000  // 667.57MiB is the maximum physical CD capacity
	MaxPartSize       = 1073741824 // 1GiB
	MaxNameSize       = 32
	MaxImageSize      = 15
	PaddingConfigSize = 15
	ConfigSize        = 64 // Every game config takes 64 bytes on ul.cfg file
)

const ( // Covers
	// Covers are being extracted from https://github.com/xlenore/ps2-covers repository
	CoverDownloadUnformattedUrl = "https://raw.githubusercontent.com/xlenore/ps2-covers/main/covers/default/%s.jpg"
	CoverMaxWidth               = 360
	CoverMaxHeight              = 640
)

type GameConfig struct {
	// Basic
	Name    [MaxNameSize]byte
	Image   [MaxImageSize]byte
	Parts   int8
	Media   int8
	Padding [PaddingConfigSize]byte

	// Extra
	NameHash  string
	Files     []string
	CoverPath string
	GamePath  string
}

const GameConfigSize = 64

func NewGameConfig(name, image, path string, size int64) *GameConfig {
	g := &GameConfig{GamePath: path}

	copy(g.Name[:], []byte(name))
	copy(g.Image[:], []byte("ul."+image))
	if size <= MaxCDSize {
		g.Media = MediaCD
		g.Parts = 1
	} else {
		g.Media = MediaDVD
		g.Parts = int8(math.Ceil(float64(size) / float64(MaxPartSize)))
	}
	g.setup()

	return g
}

func NewGameConfigFromBytes(data []byte, path string) *GameConfig {
	g := &GameConfig{GamePath: path}

	offset := 0
	copy(g.Name[:], data[offset:])
	offset += len(g.Name)
	copy(g.Image[:], data[offset:])
	offset += len(g.Image)
	g.Parts = int8(data[offset])
	offset++
	g.Media = int8(data[offset])
	offset++
	copy(g.Padding[:], data[offset:])

	g.setup()
	return g
}

func (g *GameConfig) setup() {
	g.updateHash()
	g.generateFileNames()
	g.Padding[4] = PaddingByte
}

func (g *GameConfig) updateHash() {
	g.NameHash = oplCRC32.Crc32(utils.BytesToString(g.Name[:]))
}

func (g *GameConfig) generateFileNames() {
	image := g.GetImage()
	for partNumber := int8(0); partNumber < g.Parts; partNumber++ {
		partName := fmt.Sprintf("ul.%s.%s.%02d", g.NameHash, image, partNumber)
		g.Files = append(g.Files, path.Join(g.GamePath, partName))
	}
	g.CoverPath = path.Join(g.GamePath, "ART", image+"_COV.jpg")
}

func (g *GameConfig) AsBytes() []byte {
	return slices.Concat(
		g.Name[:],
		g.Image[:],
		[]byte{byte(g.Parts), byte(g.Media)},
		g.Padding[:],
	)
}

func (g *GameConfig) Rename(name string) error {
	// This's necessary because if we just copy the name it won't delete the end when the new name is smaller than old one
	newName := make([]byte, MaxNameSize)
	copy(newName, []byte(name))
	copy(g.Name[:], newName)

	oldHash := g.NameHash // save to rename files
	g.updateHash()
	newHash := g.NameHash
	for index, fileName := range g.Files {
		newFileName := strings.Replace(fileName, oldHash, newHash, 1)
		if err := os.Rename(fileName, newFileName); err != nil {
			return err
		}
		g.Files[index] = newFileName
	}
	return nil
}

func (g *GameConfig) DeleteFiles() error {
	for _, file := range g.Files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	os.Remove(g.CoverPath) // It must not exists, this's why I won't check errors
	return nil
}

func (g *GameConfig) IsCoverInstalled() bool {
	return utils.FileExists(g.CoverPath)
}

func (g *GameConfig) GetName() string {
	return utils.BytesToString(g.Name[:])
}

func (g *GameConfig) GetImage() string {
	return utils.BytesToString(g.Image[3:])
}

func (g *GameConfig) DownloadCover() error {
	// Make image name looks like on the website
	gameImage := strings.Replace(g.GetImage(), "_", "-", 1)
	gameImage = strings.Replace(gameImage, ".", "", 1)

	response, err := http.Get(fmt.Sprintf(CoverDownloadUnformattedUrl, gameImage))
	if err != nil {
		return ErrCoverRequestFailed
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return ErrCoverNotFound
	}
	originalCover, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	resizedCover, err := utils.ResizeJPGToMax(originalCover, CoverMaxWidth, CoverMaxHeight)
	if err != nil {
		return err
	}
	return os.WriteFile(g.CoverPath, resizedCover, 0644)
}
