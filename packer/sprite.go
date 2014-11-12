package packer

import "image"

// Size of a sprite
type Size struct {
	X, Y int
}

// Sprite contains image and placement data for sprites
// it also has exported fields for data serialisation
type Sprite struct {
	Name   string      // sprite name without extension
	Offset image.Point // offset in atlas
	Size   Size        // image size
	img    image.Image
	area   int
}

// ByArea to help sort sprites for easy placement
type ByArea []Sprite

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].area < a[j].area }
