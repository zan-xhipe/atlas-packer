package main

import "image"

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
