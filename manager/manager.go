package manager

import (
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"ps2manager/utils"
	"strings"
)

const (
	SYSTEM_CONFIG_NAME = "/SYSTEM.CNF"
)

var configFile string
var workingDir string

var games []Game

func InitManager(dir string) error {
	configFile = path.Join(dir, "ul.cfg")
	workingDir = dir
	if err := os.MkdirAll(path.Join(workingDir, "ART"), os.ModePerm); err != nil {
		return err
	}
	_, err := os.Stat(configFile)
	if os.IsNotExist(err) {
		file, err := os.Create(configFile)
		if err != nil {
			return err
		}
		file.Close()
	} else if err != nil {
		return err
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
		game := NewGameFromBytes(data[offset:offset+GameConfigSize], workingDir)
		if len(game.Parts.Files) < int(game.Config.Parts) {
			fmt.Printf("Game '%s' is missing files\n", game.GetName())
		}
		games = append(games, game)
	}
	return nil
}

func GetAll() []Game {
	return games
}

func Get(index int) Game {
	return games[index]
}

func Add(game Game) error {
	games = append(games, game)
	return WriteChanges()
}

func Install(isoPath, name string, progress chan int) error {
	systemCnf, err := utils.ReadFileFromISO(isoPath, SYSTEM_CONFIG_NAME)
	if err != nil {
		return err
	}
	isoFile, _ := os.Stat(isoPath)
	image := strings.Split(strings.Split(string(systemCnf), ":\\")[1], ";")[0]
	size := isoFile.Size()
	game := NewGame(name, image, size, workingDir)
	isoAsReader, err := os.Open(isoPath)
	if err != nil {
		return err
	}
	if err = writeGameParts(isoAsReader, game, size, progress); err != nil {
		return err
	}
	return Add(game)
}

func writeGameParts(data io.Reader, game Game, size int64, progress chan int) error {
	totalRead, percent := 0, 0
	chunk := make([]byte, DEFAULT_CHUNK_SIZE)
	for _, part := range game.Parts.Files {
		partFile, err := os.Create(part)
		if err != nil {
			return err
		}
		toRead := MAX_GAME_PART_SIZE
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
			return errors.New("a game with the same name already exists")
		}
	}
	if err := games[index].Rename(newName); err != nil {
		return err
	}
	return WriteChanges()
}

func Delete(index int) error {
	game := games[index]
	if err := game.Parts.RemoveAll(); err != nil {
		return err
	}
	games = append(games[:index], games[index+1:]...)
	return WriteChanges()
}

func WriteChanges() error {
	data := make([]byte, len(games)*GameConfigSize)
	for i, game := range games {
		copy(data[i*GameConfigSize:], game.Config.AsBytes())
	}
	return os.WriteFile(configFile, data, 0644)
}
