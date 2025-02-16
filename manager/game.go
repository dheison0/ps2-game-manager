package manager

import (
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path"
	"ps2manager/config"
	"ps2manager/utils"
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
	GameConfigSize    = 64 // Every game config takes 64 bytes on ul.cfg file
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

// NewGameConfig creates a new configuration for a game that isn't instaled yet
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
	g.Update()

	return g
}

// NewGameConfigFromBytes loads a game configuration from data given from a
// previous game installation
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

	g.Update()
	return g
}

// Update updates the game name hash, generate file names based on new hash and
// add padding byte
func (g *GameConfig) Update() {
	g.refreshNameHash()
	g.generateAllFileNames()
	g.Padding[4] = PaddingByte
}

// refreshNameHash creates a new CRC32 hash of the game name and update that field
func (g *GameConfig) refreshNameHash() {
	g.NameHash = utils.Crc32(g.Name[:])
}

// generateAllFileNames uses the path, name hash, image and part number
// informations to generate the file names of all game parts on disk
func (g *GameConfig) generateAllFileNames() {
	image := g.GetImage()
	for partNumber := int8(0); partNumber < g.Parts; partNumber++ {
		g.Files = append(g.Files, g.generatePartName(partNumber))
	}
	g.CoverPath = path.Join(g.GamePath, "ART", image+"_COV.jpg")
}

// generatePartName generates a game part name with full path
func (g *GameConfig) generatePartName(partNumber int8) string {
	partName := fmt.Sprintf("ul.%s.%s.%02d", g.NameHash, g.GetImage(), partNumber)
	return path.Join(g.GamePath, partName)
}

// AsBytes returns a game config as an slice of bytes to be saved or
// transfered over network
func (g *GameConfig) AsBytes() []byte {
	return slices.Concat(
		g.Name[:],
		g.Image[:],
		[]byte{byte(g.Parts), byte(g.Media)},
		g.Padding[:],
	)
}

// Rename renames a game and it's files on disk
func (g *GameConfig) Rename(name string) error {
	// This's necessary because if we just copy the name it won't delete the end when the new name is smaller than old one
	newName := make([]byte, len(g.Name))
	copy(newName, []byte(name))
	copy(g.Name[:], newName)

	g.refreshNameHash()
	for idx, oldFile := range g.Files {
		newFile := g.generatePartName(int8(idx))
		if err := os.Rename(oldFile, newFile); err != nil {
			return err
		}
		g.Files[idx] = newFile
	}
	return nil
}

// DeleteFiles deletes all files of the game, including it's cover
func (g *GameConfig) DeleteFiles() error {
	for _, file := range g.Files {
		if err := os.Remove(file); err != nil {
			return err
		}
	}
	os.Remove(g.CoverPath) // It must not exists, this's why I won't check errors
	return nil
}

// IsCoverInstalled checks if a cover is installed for the game
func (g *GameConfig) IsCoverInstalled() bool {
	return utils.FileExists(g.CoverPath)
}

// GetName returns game name as string
func (g *GameConfig) GetName() string {
	return utils.BytesToString(g.Name[:])
}

// GetImage returns game image as string
func (g *GameConfig) GetImage() string {
	return utils.BytesToString(g.Image[3:])
}

// DownloadCover tries to download game cover for remote github repository
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
	resizedCover, err := utils.ResizeJPGKeepingAspectRatio(originalCover, CoverMaxWidth, CoverMaxHeight)
	if err != nil {
		return err
	}
	return os.WriteFile(g.CoverPath, resizedCover, 0644)
}

// ExportAsISO exports the game as an ISO file
func (g *GameConfig) ExportAsISO(outputFile string, progress chan int, errChan chan error) {
	file, err := os.Create(outputFile)
	if err != nil {
		errChan <- err
		return
	}
	defer file.Close()
	actualPercentage := 0
	writtenSize := int64(0)
	totalSize, err := utils.GetFilesSizeSum(g.Files)
	if err != nil {
		errChan <- errors.New("failed to get sum of all game parts size: " + err.Error())
		return
	}
	buffer := make([]byte, config.BUFFER_SIZE)
	for _, f := range g.Files { // this will read all parts
		part, err := os.Open(f)
		if err != nil {
			errChan <- errors.New("error opening file '" + f + "': " + err.Error())
			return
		}
		for { // read part data in chunks
			readSize, err := part.Read(buffer)
			if err == io.EOF {
				break
			} else if err != nil {
				errChan <- errors.New("failed to read from file '" + f + "': " + err.Error())
				part.Close()
				return
			} else if _, err = file.Write(buffer[:readSize]); err != nil {
				errChan <- errors.New("fail writing chunk to iso file: " + err.Error())
				return
			}
			writtenSize += int64(readSize)
			newPercent := int(math.Floor(float64(writtenSize) / float64(totalSize) * 100))
			if newPercent > actualPercentage {
				if err = file.Sync(); err != nil {
					errChan <- errors.New("fail to flush data on ISO file: " + err.Error())
					part.Close()
					return
				}
				actualPercentage = newPercent
				progress <- actualPercentage
			}
		}
		part.Close()
	}
	if err = file.Sync(); err != nil {
		errChan <- errors.New("final file sync failed: " + err.Error())
	}
}
