package manager

import (
	"bytes"
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

func (g *GameConfig) AsBytes() []byte {
	data := bytes.Join([][]byte{
		g.Name[:],
		g.Image[:],
		{byte(g.Parts), byte(g.Media)},
		g.Padding[:],
	}, []byte{})
	return data
}

func (g *GameConfig) FromBytes(data []byte) {
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
}
