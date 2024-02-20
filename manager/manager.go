package manager

import (
	"errors"
	"fmt"
	"os"
	"path"
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
		game, err := NewGameFromBytes(data[offset:offset+GameConfigSize], workingDir)
		if err != nil {
			return err
		}
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
