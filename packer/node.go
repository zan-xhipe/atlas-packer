package packer

import (
	"fmt"
	"image"
)

// Node in a binary tree that hold sprites
type node struct {
	child [2]*node
	rect  image.Rectangle
	img   *Sprite
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

// insert a sprite in the tree. Returns nil if it can't fit
// returns the node that it does fit in.
// will split the node if it is too big
func (n *node) insert(img *Sprite) *node {
	if !n.isLeaf() {
		node := n.child[0].insert(img)

		if node != nil {
			return node
		}
		return n.child[1].insert(img)
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

// split the node along one axis. The first child is the right
// shape so that it should fit perfectly within two splits
func (n *node) split(img *Sprite) {
	dx := n.rect.Dx() - img.Size.X
	dy := n.rect.Dy() - img.Size.Y

	// which axis to split on
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
