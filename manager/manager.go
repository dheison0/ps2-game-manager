package manager

import (
	"io"
	"math"
	"os"
	"path"
	"ps2manager/config"
	"ps2manager/utils"
	"strings"
)

const SystemConfigIsoPath = "/SYSTEM.CNF"

var dataDir, configFile string
var games []*GameConfig

func InitManager(dir string) error {
	dataDir, configFile = dir, path.Join(dir, "ul.cfg")
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
	return readConfigFile()
}

func GetAll() []*GameConfig {
	return games
}

func Get(index int) *GameConfig {
	return games[index]
}

func Add(game *GameConfig) error {
	games = append(games, game)
	return writeConfigChanges()
}

func Install(isoPath, name string, progress chan int) error {
	systemCnf, err := utils.ReadFileFromISO(isoPath, SystemConfigIsoPath)
	if err != nil {
		return err
	}
	isoReader, _ := os.Open(isoPath)
	defer isoReader.Close()
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

func Rename(index int, newName string) error {
	for _, g := range games {
		if g.GetName() == newName {
			return ErrNameAlreadyExists
		}
	}
	if err := games[index].Rename(newName); err != nil {
		return err
	}
	return writeConfigChanges()
}

func Delete(index int) error {
	game := games[index]
	if err := game.DeleteFiles(); err != nil {
		return err
	}
	games = append(games[:index], games[index+1:]...)
	return writeConfigChanges()
}

// CheckIfAcceptName only accepts ascii characters and till the maximum name size
func CheckIfAcceptName(t string, r rune) bool {
	return len(t) <= MaxNameSize && r <= 127
}

func readConfigFile() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	for i := 0; i < len(data)/GameConfigSize; i++ {
		offset := i * GameConfigSize
		games = append(games, NewGameConfigFromBytes(data[offset:offset+GameConfigSize], dataDir))
	}
	return nil
}

func writeConfigChanges() error {
	data := make([]byte, len(games)*GameConfigSize)
	for i, game := range games {
		copy(data[i*GameConfigSize:], game.AsBytes())
	}
	return os.WriteFile(configFile, data, 0644)
}

func writeGameParts(iso io.Reader, game *GameConfig, size int64, progress chan int) error {
	written, actualProgress := 0, 0
	buffer := make([]byte, config.BUFFER_SIZE)
	progress <- actualProgress
	for _, partName := range game.Files {
		partFile, err := os.Create(partName)
		if err != nil {
			return err
		}
		toWrite := MaxPartSize
		for toWrite > 0 {
			n, err := iso.Read(buffer)
			if err == io.EOF {
				break
			} else if err != nil {
				return err
			}
			if _, err = partFile.Write(buffer[:n]); err != nil {
				return err
			}

			toWrite -= n
			written += n
			progressPercentage := int(math.Floor(float64(written) / float64(size) * 100.0))
			if progressPercentage > actualProgress {
				partFile.Sync()
				actualProgress = progressPercentage
				progress <- actualProgress
			}
		}
		if err = partFile.Sync(); err != nil {
			return err
		}
		partFile.Close()
	}
	return nil
}
