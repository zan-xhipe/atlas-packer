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
	_ "image/png"
)

type sprite struct {
	name   string
	img    image.Image
	bounds image.Rectangle
	area   int
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

	for i := range sprites {
		s := readSprite(inputDir, files[i].Name())
		sprites[i] = s
	}

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(sprites)))

	for _, s := range sprites {
		fmt.Println(s)
	}
}

func readSprite(dir, name string) (s sprite) {
	path := path.Join(dir, name)
	reader, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	img, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	s.img = img
	s.bounds = img.Bounds()
	s.area = s.bounds.Dx() * s.bounds.Dy()
	return
}
