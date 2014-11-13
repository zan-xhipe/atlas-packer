package packer

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

// Atlas hold dimension and spacing for final image
type Atlas struct {
	DimX, DimY int
	Space      int
}

// Pack the sprite slice so that it fits in the dimensions provided by atlas
func Pack(sprites []Sprite, atlas Atlas) (*image.RGBA, []byte, error) {

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(byArea(sprites)))

	img := image.NewRGBA(image.Rect(0, 0, atlas.DimX, atlas.DimY))

	n := node{rect: image.Rect(0, 0, atlas.DimX, atlas.DimY)}

	data := []byte("[")

	for i := range sprites {
		s := &sprites[i]

		// add separator to each sprite
		s.Size.X += atlas.Space
		s.Size.Y += atlas.Space

		node := n.insert(s)

		// need to remove all space added around the sprite
		s.Size.X -= atlas.Space
		s.Size.Y -= atlas.Space

		j, err := json.MarshalIndent(s, "", " ")
		if err != nil {
			return img, data, err
		}
		data = append(data, '\n')
		data = append(data, j...)

		if node != nil {
			draw.Draw(img, node.rect, s.img, image.ZP, draw.Src)
		} else {
			return img, data, fmt.Errorf("could not place %s\n", s.Name)
		}

	}

	data = append(data, "]"...)

	return img, data, nil

}

// ReadSprites returns all spritees found in the provided directory
func ReadSprites(dir string) ([]Sprite, error) {

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	sprites := make([]Sprite, len(files))

	for i := range sprites {
		s, err := readSprite(dir, files[i].Name())
		if err != nil {
			return sprites, err
		}
		sprites[i] = s
	}

	return sprites, nil
}

func readSprite(dir, name string) (Sprite, error) {
	var s Sprite

	path := path.Join(dir, name)
	reader, err := os.Open(path)
	if err != nil {
		return s, err
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return s, err
	}

	s.Name = strings.TrimSuffix(name, filepath.Ext(name))
	s.img = img
	rect := img.Bounds()
	s.Size = size{rect.Dx(), rect.Dy()}
	s.area = s.Size.X * s.Size.Y
	return s, nil
}
