package manager

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"ps2manager/utils"
	"ps2manager/utils/oplCRC32"
	"strings"
)

const (
	MediaCD     = 0x12
	MediaDVD    = 0x14
	PaddingByte = 0x08
)

type GameConfig struct {
	Name     [32]byte
	Image    [15]byte
	Parts    int8
	Media    int8
	Padding  [15]byte
	NameHash string
}

const GameConfigSize = 64

func (g *GameConfig) AsBytes() []byte {
	g.Padding[4] = PaddingByte
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

func (g *GameConfig) RegenerateHash() {
	g.NameHash = oplCRC32.Crc32(utils.BytesToString(g.Name[:]))
}

type GameFiles struct {
	Files []string
}

func (f *GameFiles) GenerateFileNames(hash, image, dataDir string, parts int8) {
	for i := int8(0); i < parts; i++ {
		f.Files = append(f.Files, path.Join(dataDir, fmt.Sprintf("ul.%s.%s.%02d", hash, image, i)))
	}
}

func (f *GameFiles) UpdateHash(hash string) error {
	if len(f.Files) == 0 {
		return nil
	}
	oldHash := strings.Split(f.Files[0], ".")[1]
	for i, p := range f.Files {
		newName := strings.Replace(p, oldHash, hash, 1)
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
