package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type atlas struct {
	dimX, dimY int
	space      int
	sprites    []sprite
	img        *image.RGBA
	data       []byte
}

func (a *atlas) readSprites(dir string, files []os.FileInfo) (int, error) {
	var area int
	a.sprites = make([]sprite, len(files))

	for i := range a.sprites {
		s, err := readSprite(dir, files[i].Name(), a.space)
		if err != nil {
			return area, err
		}
		area += s.area
		a.sprites[i] = s
	}
	return area, nil
}

func readSprite(dir, name string, space int) (s sprite, err error) {
	path := path.Join(dir, name)
	reader, err := os.Open(path)
	if err != nil {
		return
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		return
	}

	s.Name = strings.TrimSuffix(name, filepath.Ext(name))
	s.img = img
	rect := img.Bounds()
	s.Size = Size{rect.Dx() + space, rect.Dy() + space}
	s.area = s.Size.X * s.Size.Y
	return
}

func (a *atlas) pack() error {

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(a.sprites)))

	a.img = image.NewRGBA(image.Rect(0, 0, a.dimX, a.dimY))

	n := Node{rect: image.Rect(0, 0, a.dimX, a.dimY)}

	a.data = []byte("[")

	for i := range a.sprites {
		s := &a.sprites[i]
		node := n.insert(s)

		// need to remove all space added around the sprite
		s.Size.X -= a.space
		s.Size.Y -= a.space

		j, err := json.MarshalIndent(s, "", " ")
		if err != nil {
			return err
		}
		a.data = append(a.data, '\n')
		a.data = append(a.data, j...)

		if node != nil {
			draw.Draw(a.img, node.rect, s.img, image.ZP, draw.Src)
		} else {
			return fmt.Errorf("could not place %s\n", s.Name)
		}

	}

	a.data = append(a.data, []byte("]")...)

	return nil
}

func (a *atlas) write(name string) error {
	imageName := name + ".png"
	dataName := name + ".json"

	imageWriter, err := os.Create(imageName)
	if err != nil {
		return err
	}
	defer imageWriter.Close()

	dataWriter, err := os.Create(dataName)
	if err != nil {
		return err
	}
	defer dataWriter.Close()

	n, err := dataWriter.Write(a.data)
	if err != nil {
		return fmt.Errorf("position: %d error: %s", n, err)
	}

	err = png.Encode(imageWriter, a.img)
	if err != nil {
		return err
	}

	return nil
}
