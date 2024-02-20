package manager

import (
	"errors"
	"os"
	"path"
)

var configFile string
var workingDir string

type Game struct {
	Config GameConfig
	Parts  GameFiles
}

func NewGame(data []byte) (Game, error) {
	game := Game{}
	game.Config.FromBytes(data)
	files, _ := os.ReadDir(workingDir)
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		game.Parts.Files = append(game.Parts.Files, path.Join(workingDir, f.Name()))
	}
	return game, nil
}

var games []Game

func InitManager(dir string) error {
	configFile = path.Join(dir, "ul.cfg")
	workingDir = dir
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
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
		game, err := NewGame(data[offset : offset+GameConfigSize])
		if err != nil {
			return err
		}
		games = append(games, game)
	}
	return nil
}

func GetAllGames() []Game {
	return games
}

func GetGame(index int) Game {
	return games[index]
}

func UpdateGameConfig(index int, gc GameConfig) error {
	for _, g := range games {
		if g.Config.Name == gc.Name {
			return errors.New("game with the same name already exists")
		}
	}
	games[index].Config = gc
	games[index].Parts.UpdateHash(string(gc.Name[:]))
	return WriteChanges()
}

func RemoveGame(index int) error {
	game := games[index]
	game.Parts.RemoveAll()
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
