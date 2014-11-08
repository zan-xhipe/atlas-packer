package main

import "image"

type Size struct {
	X, Y int
}

// sprite contains image and placement data for sprites
// it also has exported fields for data serialisation
type sprite struct {
	Name   string      // sprite name without extension
	Offset image.Point // offset in atlas
	Size   Size        // image size
	img    image.Image
	area   int
}

// sprites can be sorted by area for easy placement
type ByArea []sprite

func (a ByArea) Len() int           { return len(a) }
func (a ByArea) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByArea) Less(i, j int) bool { return a[i].area < a[j].area }
