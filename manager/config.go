package manager

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"unsafe"
)

// const (
// 	MediaCD  = 0x12
// 	MediaDVD = 0x14
// )

type GameConfig struct {
	Name    [32]byte
	Image   [15]byte
	Parts   int8
	Media   int8
	Padding [15]byte
}

const GameConfigSize = int(unsafe.Sizeof(GameConfig{}))

func (g GameConfig) AsBytes() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(g); err != nil {
		log.Fatalf("Failed to transform game config into slice of bytes: %v", err)
	}
	return buffer.Bytes()
}

func ReadFromFile(filename string) []GameConfig {
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}
	var games []GameConfig
	for i := 0; i < len(data)/GameConfigSize; i++ {
		var game GameConfig
		offset := i * GameConfigSize
		copy(game.Name[0:], data[offset:offset+len(game.Name)])
		offset += len(game.Name)
		copy(game.Image[0:], data[offset:offset+len(game.Image)])
		offset += len(game.Image)
		game.Parts = int8(data[offset])
		offset++
		game.Media = int8(data[offset])
		offset++
		copy(game.Padding[0:], data[offset:])
		games = append(games, game)
	}
	return games
}
