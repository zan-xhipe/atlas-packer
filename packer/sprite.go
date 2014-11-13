package packer

import "image"

// size of a sprite
type size struct {
	X, Y int
}

// Sprite contains image and placement data for sprites
// it also has exported fields for data serialisation
type Sprite struct {
	Name   string      // sprite name without extension
	Offset image.Point // offset in atlas
	Size   size        // image size
	img    image.Image
	area   int
}

// byArea to help sort sprites for easy placement
type byArea []Sprite

func (a byArea) Len() int           { return len(a) }
func (a byArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byArea) Less(i, j int) bool { return a[i].area < a[j].area }
