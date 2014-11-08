package main

import (
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
)

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

	atlas := atlas{dimX: dimX, dimY: dimY, space: space}

	err = atlas.readDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}

	err = atlas.pack()
	if err != nil {
		log.Fatal(err)
	}

	err = atlas.write(outputName)
	if err != nil {
		log.Fatal(err)
	}
}

// parseDimensions takes a string of two numbers separated with an x
// and returns two ints representing a width and a height
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
