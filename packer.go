package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"sort"

	"image"
	"image/draw"
	"image/png"
)

type Size struct {
	x, y int
}

type sprite struct {
	name string
	img  image.Image
	rect image.Rectangle
	area int
	size Size
}

type ByArea []sprite

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].area < a[j].area }

func main() {
	flag.Parse()

	args := flag.Args()

	inputDir := args[0]

	files, _ := ioutil.ReadDir(inputDir)
	sprites := make([]sprite, len(files))

	totalX := 0
	totalY := 0

	for i := range sprites {
		s := readSprite(inputDir, files[i].Name())
		sprites[i] = s
		totalX += s.size.x
		totalY += s.size.y
	}

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(sprites)))

	// the final image
	dst := image.NewRGBA(image.Rect(0, 0, 1024, 1024))

	for i := range sprites {
		s := &sprites[i]
		draw.Draw(dst, s.rect, s.img, image.ZP, draw.Src)
	}

	writer, err := os.Create("test.png")
	err = png.Encode(writer, dst)
	if err != nil {
		log.Fatal(err)
	}
}

func readSprite(dir, name string) (s sprite) {
	path := path.Join(dir, name)
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	s.name = name
	s.img = img
	s.rect = img.Bounds()
	s.size = Size{s.rect.Dx(), s.rect.Dy()}
	s.area = s.size.x * s.size.y
	return
}
