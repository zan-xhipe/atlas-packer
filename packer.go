package main

import (
	"errors"
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

func (n *Node) insert(img *sprite) (image.Rectangle, error) {
	// returned when no valid rectangle is found
	nop := image.Rect(0, 0, 0, 0)

	// there is already an image in this node
	if n.img != nil {
		fmt.Println("already contains an image")
		return nop, errors.New("already contains an image")
	}

	// try to insert into either of the nodes children
	if n.child[0] != nil {
		fmt.Println("has a 0 child")
		rc, err := n.child[0].insert(img)

		if err != nil {
			fmt.Println("has a 1 child")
			return n.child[1].insert(img)
		} else {
			return rc, nil
		}
	}

	if n.rect.Dx() < img.size.x || n.rect.Dy() < img.size.y {
		fmt.Println("space too small")
		return nop, errors.New("space too small")
	}

	if n.rect.Dx() == img.size.x && n.rect.Dy() == img.size.y {
		fmt.Println("prefect fit")
		n.img = img
		return n.rect, nil
	}

	if n.rect.Dx() >= img.size.x && n.rect.Dy() >= img.size.y {
		fmt.Println("Split")
		n.split(img)
	}

	return n.insert(img)

}

func (n *Node) split(img *sprite) {
	var tl0 image.Point
	var br0 image.Point

	var tl1 image.Point
	var br1 image.Point

	dx := n.size().x - img.size.x
	dy := n.size().y - img.size.y

	rc := n.rect

	tl0 = rc.Min
	br1 = rc.Max

	if dx > dy {
		fmt.Println("split on x")
		br0 = image.Point{rc.Min.X + img.size.x, rc.Dy()}
		tl1 = image.Point{rc.Min.X + img.size.x - 1, rc.Min.Y}
	} else {
		fmt.Println("split on y")
		br0 = image.Point{rc.Dx(), rc.Min.Y + img.size.y}
		tl1 = image.Point{rc.Min.X, rc.Min.Y + img.size.y - 1}
	}

	rect0 := image.Rectangle{tl0, br0}
	n.child[0] = &Node{rect: rect0}

	rect1 := image.Rectangle{tl1, br1}
	n.child[1] = &Node{rect: rect1}
}

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
	dst := image.NewRGBA(image.Rect(0, 0, 2048, 2048))

	n := Node{rect: image.Rect(0, 0, 1024, 1024)}

	for i := range sprites {
		s := &sprites[i]
		fmt.Printf("inserting %s\n", s.name)
		rect, err := n.insert(s)
		if err == nil {
			draw.Draw(dst, rect, s.img, image.ZP, draw.Src)
		}
	}

	n.print()

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
