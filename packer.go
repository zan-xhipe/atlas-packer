package main

import (
	"fmt"
	"image"
)

type node struct {
	child [2]*node
	rect  image.Rectangle
	img   *sprite
}

func (n *node) print() {
	fmt.Println(n)
	if n.child[0] != nil {
		n.child[0].print()
	}
	if n.child[1] != nil {
		n.child[1].print()
	}
}

func (n *node) isLeaf() bool {
	return n.child[0] == nil || n.child[1] == nil
}

func (n *node) insert(img *sprite) *node {
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

func (n *node) split(img *sprite) {
	dx := n.rect.Dx() - img.Size.X
	dy := n.rect.Dy() - img.Size.Y

	if dx > dy {
		n.child[0] = &node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Min.X+img.Size.X,
			n.rect.Max.Y)}
		n.child[1] = &node{rect: image.Rect(
			n.rect.Min.X+img.Size.X,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	} else {
		n.child[0] = &node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y,
			n.rect.Max.X,
			n.rect.Min.Y+img.Size.Y)}
		n.child[1] = &node{rect: image.Rect(
			n.rect.Min.X,
			n.rect.Min.Y+img.Size.Y,
			n.rect.Max.X,
			n.rect.Max.Y)}
	}

}
