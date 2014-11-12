package packer

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
)

type Atlas struct {
	DimX, DimY int
	Space      int
	sprites    []sprite
	img        *image.RGBA
	data       []byte
}

// readDir takes a directory and processes the sprites in it\
// if it encounters a non image file it crashes
func (a *Atlas) ReadDir(dir string) error {
	var area int

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}

	a.sprites = make([]sprite, len(files))

	for i := range a.sprites {
		s, err := readSprite(dir, files[i].Name(), a.Space)
		if err != nil {
			return err
		}
		area += s.area
		a.sprites[i] = s
	}

	if area > a.DimX*a.DimY {
		return fmt.Errorf("atlas to small")
	}

	return nil
}

// readSprite takes a file and returns the sprite
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

// pack finds a way to layout all the atlas' sprites so that
// they fit inside the atlas nicely
func (a *Atlas) Pack() error {

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(a.sprites)))

	a.img = image.NewRGBA(image.Rect(0, 0, a.DimX, a.DimY))

	n := node{rect: image.Rect(0, 0, a.DimX, a.DimY)}

	a.data = []byte("[")

	for i := range a.sprites {
		s := &a.sprites[i]
		node := n.insert(s)

		// need to remove all space added around the sprite
		s.Size.X -= a.Space
		s.Size.Y -= a.Space

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

// write creates name.png and name.json
// for the atlas texture and data
func (a *Atlas) Write(imageName, dataName string) error {

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
