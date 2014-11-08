package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	var space int
	var dim string

	flag.IntVar(&space, "space", 1, "space added between images")
	flag.StringVar(&dim, "dimensions", "1024x1024", "atlas size")

	flag.Parse()

	args := flag.Args()

	inputDir := args[0]
	outputName := args[1]

	dimX, dimY, err := parseDimensions(dim)
	if err != nil {
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	sprites := make([]sprite, len(files))

	totalArea := 0

	for i := range sprites {
		s, err := readSprite(inputDir, files[i].Name(), space)
		if err != nil {
			log.Fatal(err)
		}
		totalArea += s.area
		sprites[i] = s
	}

	if totalArea > dimX*dimY {
		log.Fatalf("%s atlas to small", dim)
	}

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(sprites)))

	img, data, err := packAtlas(sprites, dimX, dimY, space)
	if err != nil {
		log.Fatal(err)
	}

	err = writeAtlas(outputName, img, data)
	if err != nil {
		log.Fatal(err)
	}
}

func parseDimensions(dim string) (dimX, dimY int, err error) {
	dims := strings.Split(dim, "x")
	if len(dims) != 2 {
		err = fmt.Errorf("couldn't parse dimension %s\n", dims)
		return
	}

	dimX, err = strconv.Atoi(dims[0])
	if err != nil {
		return
	}

	dimY, err = strconv.Atoi(dims[1])
	if err != nil {
		return
	}

	return
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

func packAtlas(sprites []sprite, dimX, dimY, space int) (*image.RGBA, []byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, dimX, dimY))

	n := Node{rect: image.Rect(0, 0, dimX, dimY)}

	data := []byte("[")

	for i := range sprites {
		s := &sprites[i]
		node := n.insert(s)

		// need to remove all space added around the sprite
		s.Size.X -= space
		s.Size.Y -= space

		j, err := json.MarshalIndent(s, "", " ")
		if err != nil {
			return nil, nil, err
		}
		data = append(data, '\n')
		data = append(data, j...)

		if node != nil {
			draw.Draw(img, node.rect, s.img, image.ZP, draw.Src)
		} else {
			return nil, nil, fmt.Errorf("could not place %s\n", s.Name)
		}

	}

	data = append(data, []byte("]")...)

	return img, data, nil
}

func writeAtlas(name string, img *image.RGBA, data []byte) error {
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

	n, err := dataWriter.Write(data)
	if err != nil {
		return fmt.Errorf("position: %d error: %s", n, err)
	}

	err = png.Encode(imageWriter, img)
	if err != nil {
		return err
	}

	return nil
}
