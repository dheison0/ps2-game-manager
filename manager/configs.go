package manager

import (
	"bytes"
	"os"
	"ps2manager/utils/oplCRC32"
	"strings"
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

type GameFiles struct {
	Files []string
}

func (f *GameFiles) UpdateHash(name string) error {
	if len(f.Files) == 0 {
		return nil
	}
	oldHash := strings.Split(f.Files[0], ".")[1]
	newHash := oplCRC32.Crc32(name)
	for i, p := range f.Files {
		newName := strings.Replace(p, oldHash, newHash, 1)
		if err := os.Rename(p, newName); err != nil {
			return err
		}
		f.Files[i] = newName
	}
	return nil
}

func (f *GameFiles) RemoveAll() error {
	for _, f := range f.Files {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}
