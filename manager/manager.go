package manager

import (
	"io"
	"math"
	"os"
	"path"
	"ps2manager/utils"
	"strings"
)

const (
	SystemConfigIsoPath = "/SYSTEM.CNF"
	DefaultChunkSize    = 262144 // 256KiB
)

var dataDir, configFile string
var games []*GameConfig

func InitManager(dir string) error {
	configFile = path.Join(dir, "ul.cfg")
	dataDir = dir
	if err := os.MkdirAll(path.Join(dataDir, "ART"), os.ModePerm); err != nil {
		return err
	}
	if !utils.FileExists(configFile) {
		file, err := os.Create(configFile)
		if err != nil {
			return err
		}
		file.Close()
	}
	return ReadConfigFile()
}

func ReadConfigFile() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	for i := 0; i < len(data)/GameConfigSize; i++ {
		offset := i * GameConfigSize
		game := NewGameConfigFromBytes(data[offset:offset+GameConfigSize], dataDir)
		games = append(games, game)
	}
	return nil
}

func GetAll() []*GameConfig {
	return games
}

func Get(index int) *GameConfig {
	return games[index]
}

func Add(game *GameConfig) error {
	games = append(games, game)
	return WriteConfigChanges()
}

func Install(isoPath, name string, progress chan int) error {
	systemCnf, err := utils.ReadFileFromISO(isoPath, SystemConfigIsoPath)
	if err != nil {
		return err
	}
	isoReader, _ := os.Open(isoPath)
	isoStat, _ := isoReader.Stat()
	image := strings.Split(strings.Split(string(systemCnf), ":\\")[1], ";")[0]
	size := isoStat.Size()
	game := NewGameConfig(name, image, dataDir, size)
	for _, g := range games {
		if g.Image == game.Image {
			return ErrAlreadyInstalled
		} else if g.Name == game.Name {
			return ErrNameAlreadyExists
		}
	}
	if err = writeGameParts(isoReader, game, size, progress); err != nil {
		return err
	}
	return Add(game)
}

func writeGameParts(data io.Reader, game *GameConfig, size int64, progress chan int) error {
	totalRead, percent := 0, 0
	chunk := make([]byte, DefaultChunkSize)
	for _, partName := range game.Files {
		partFile, err := os.Create(partName)
		if err != nil {
			return err
		}
		toRead := MaxPartSize
		for toRead > 0 {
			n, err := data.Read(chunk)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			if _, err = partFile.Write(chunk[:n]); err != nil {
				return err
			}

			toRead -= n
			totalRead += n
			partFile.Sync()

			newPercent := int(math.Floor(float64(totalRead) / float64(size) * 100.0))
			if newPercent > percent {
				percent = newPercent
				progress <- percent
			}
		}
		if err = partFile.Close(); err != nil {
			return err
		}

	}
	return nil
}

func Rename(index int, newName string) error {
	for _, g := range games {
		if g.GetName() == newName {
			return ErrNameAlreadyExists
		}
	}
	if err := games[index].Rename(newName); err != nil {
		return err
	}
	return WriteConfigChanges()
}

func Delete(index int) error {
	game := games[index]
	if err := game.DeleteFiles(); err != nil {
		return err
	}
	games = append(games[:index], games[index+1:]...)
	return WriteConfigChanges()
}

func WriteConfigChanges() error {
	data := make([]byte, len(games)*GameConfigSize)
	for i, game := range games {
		copy(data[i*GameConfigSize:], game.AsBytes())
	}
	return os.WriteFile(configFile, data, 0644)
}
