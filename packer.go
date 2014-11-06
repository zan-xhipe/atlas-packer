package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"image"
	"image/draw"
	"image/png"
)

type Size struct {
	X, Y int
}

type sprite struct {
	Name   string
	Offset image.Point
	Size   Size
	img    image.Image
	area   int
}

type ByArea []sprite

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].area < a[j].area }

type Node struct {
	child [2]*Node
	rect  image.Rectangle
	img   *sprite
}

func (n *Node) size() Size {
	return Size{n.rect.Dx(), n.rect.Dy()}
}

func (n *Node) print() {
	fmt.Println(n)
	if n.child[0] != nil {
		n.child[0].print()
	}
	if n.child[1] != nil {
		n.child[1].print()
	}
}

func (n *Node) isLeaf() bool {
	return n.child[0] == nil || n.child[1] == nil
}

func (n *Node) insert(img *sprite) *Node {
	if !n.isLeaf() {
		node := n.child[0].insert(img)

		if node != nil {
			return node
		} else {
			return n.child[1].insert(img)
		}
	}

	// there is already an image in this node
	if n.img != nil {
		return nil
	}

	// space too small
	if n.rect.Dx() < img.Size.X || n.rect.Dy() < img.Size.Y {
		return nil
	}

	// just right
	if n.rect.Dx() == img.Size.X && n.rect.Dy() == img.Size.Y {
		img.Offset = n.rect.Min
		n.img = img
		return n
	}

	// the space that is left will be large enough to split
	n.split(img)

	return n.child[0].insert(img)
}

func (n *Node) split(img *sprite) {
	dx := n.rect.Dx() - img.Size.X
	dy := n.rect.Dy() - img.Size.Y

	if dx > dy {
		n.child[0] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Min.X+img.Size.X,
			n.rect.Max.Y)}
		n.child[1] = &Node{rect: image.Rect(
			n.rect.Min.X+img.Size.X,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	} else {
		n.child[0] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Min.Y+img.Size.Y)}
		n.child[1] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y+img.Size.Y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	}

}

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

	// the final image
	dst := image.NewRGBA(image.Rect(0, 0, dimX, dimY))

	n := Node{rect: image.Rect(0, 0, dimX, dimY)}

	for i := range sprites {
		s := &sprites[i]
		node := n.insert(s)
		if node != nil {
			draw.Draw(dst, node.rect, s.img, image.ZP, draw.Src)
		} else {
			log.Fatalf("could not place %s\n", s.Name)
		}

	}

	err = writeAtlas(outputName, dst)
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

func writeAtlas(filename string, data *image.RGBA) (err error) {
	writer, err := os.Create(filename)
	if err != nil {
		return
	}
	defer writer.Close()

	err = png.Encode(writer, data)
	if err != nil {
		return
	}

	return
}
