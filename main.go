package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	var space int
	var dim string
	var verbose bool
	var forceRebuild bool

	flag.IntVar(&space, "space", 1, "space added between images")
	flag.StringVar(&dim, "dimensions", "1024x1024", "atlas size")
	flag.BoolVar(&verbose, "v", false, "verbose output")
	flag.BoolVar(&forceRebuild, "force", false, "force rebuild")

	flag.Parse()

	args := flag.Args()

	inputDir := args[0]
	outputName := args[1]

	imgName := outputName + ".png"
	dataName := outputName + ".json"

	// check if it is necessary to rebuild the atlas
	if !forceRebuild {
		imgInfo, imgErr := os.Stat(imgName)
		datInfo, datErr := os.Stat(dataName)
		dirInfo, dirErr := os.Stat(inputDir)

		if imgErr == nil && datErr == nil && dirErr == nil &&
			dirInfo.ModTime().Before(imgInfo.ModTime()) &&
			dirInfo.ModTime().Before(datInfo.ModTime()) {

			if verbose {
				log.Println("already up to date")
			}

			os.Exit(0)
		}
	}

	dimX, dimY, err := parseDimensions(dim)
	if err != nil {
		log.Fatal(err)
	}

	atlas := atlas{dimX: dimX, dimY: dimY, space: space}

	err = atlas.readDir(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("read sprite directory")
	}

	err = atlas.pack()
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("packed sprites in atlas")
	}

	err = atlas.write(outputName)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Printf("wrote atlas to file")
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
