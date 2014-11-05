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
	"strings"

	"image"
	"image/draw"
	"image/png"
)

type Size struct {
	x, y int
}

func (a Size) Equal(b Size) bool {
	return a.x == b.x && a.y == b.y
}

func (a Size) Larger(b Size) bool {
	return b.x < a.x && b.y < a.y
}

func (a Size) Smaller(b Size) bool {
	return !a.Larger(b)
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
	if n.rect.Dx() < img.size.x || n.rect.Dy() < img.size.y {
		return nil
	}

	// just right
	if n.rect.Dx() == img.size.x && n.rect.Dy() == img.size.y {
		n.img = img
		return n
	}

	// the space that is left will be large enough to split
	n.split(img)

	return n.child[0].insert(img)
}

func (n *Node) split(img *sprite) {
	dx := n.rect.Dx() - img.size.x
	dy := n.rect.Dy() - img.size.y

	if dx > dy {
		n.child[0] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Min.X+img.size.x,
			n.rect.Max.Y)}
		n.child[1] = &Node{rect: image.Rect(
			n.rect.Min.X+img.size.x,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	} else {
		n.child[0] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Min.Y+img.size.y)}
		n.child[1] = &Node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y+img.size.y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	}

}

func main() {
	var space int

	flag.IntVar(&space, "space", 1, "space added between images")

	flag.Parse()

	args := flag.Args()

	inputDir := args[0]
	outputName := args[1]

	files, _ := ioutil.ReadDir(inputDir)
	sprites := make([]sprite, len(files))

	for i := range sprites {
		s := readSprite(inputDir, files[i].Name(), space)
		sprites[i] = s
	}

	// we want to place the largest sprite first
	sort.Sort(sort.Reverse(ByArea(sprites)))

	// the final image
	dst := image.NewRGBA(image.Rect(0, 0, 2048, 2048))

	n := Node{rect: image.Rect(0, 0, 2048, 2048)}

	for i := range sprites {
		s := &sprites[i]
		node := n.insert(s)
		if node != nil {
			draw.Draw(dst, node.rect, s.img, image.ZP, draw.Src)
		} else {
			log.Fatalf("could not place %s\n", s.name)
		}

	}

	writer, err := os.Create(outputName)
	err = png.Encode(writer, dst)
	if err != nil {
		log.Fatal(err)
	}
}

func readSprite(dir, name string, space int) (s sprite) {
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

	s.name = strings.TrimSuffix(name, filepath.Ext(name))
	s.img = img
	s.rect = img.Bounds()
	s.size = Size{s.rect.Dx() + space, s.rect.Dy() + space}
	s.area = s.size.x * s.size.y
	return
}
