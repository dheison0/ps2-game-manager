package manager

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
)

var games []GameConfig
var configFile string
var workingDir string

func InitManager(config string) error {
	configFile = config
	paths := strings.Split(configFile, "/")
	workingDir = strings.Join(paths[:len(paths)-1], "/")
	return ReadConfigFile()
}

func ReadConfigFile() error {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}
	games = []GameConfig{}
	for i := 0; i < len(data)/GameConfigSize; i++ {
		var game GameConfig
		offset := i * GameConfigSize
		game.FromBytes(data[offset : offset+GameConfigSize])
		games = append(games, game)
	}
	return nil
}

func GetAllGames() []GameConfig {
	return games
}

func GetGame(index int) GameConfig {
	return games[index]
}

func UpdateGame(index int, game GameConfig) error {
	for _, g := range games {
		if g.Name == game.Name {
			return errors.New("game with the same name already exists")
		}
	}
	games[index] = game
	return WriteChanges()
}

func RemoveGame(index int) error {
	game := games[index]
	dirContent, _ := os.ReadDir(workingDir)
	for _, file := range dirContent {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		n := bytes.IndexByte(game.Image[:], 0)
		image := strings.Split(string(game.Image[:n]), ".")[1]
		if strings.Contains(name, image) {
			if err := os.Remove(path.Join(workingDir, name)); err != nil {
				return err
			}
		}
	}
	games = append(games[:index], games[index+1:]...)
	return WriteChanges()
}

func WriteChanges() error {
	data := make([]byte, len(games)*GameConfigSize)
	for i, game := range games {
		copy(data[i*GameConfigSize:], game.AsBytes())
	}
	return os.WriteFile(configFile, data, 0644)
}
