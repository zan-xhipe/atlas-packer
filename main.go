package main

import (
	"flag"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/zan-xhipe/atlas-packer/packer"
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

	sprites, err := packer.ReadSprites(inputDir)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("read sprite directory")
	}

	atlas := packer.Atlas{DimX: dimX, DimY: dimY, Space: space}

	img, data, err := packer.Pack(sprites, atlas)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Println("packed sprites in atlas")
	}

	err = writeData(data, dataName)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Printf("wrote data to file")
	}

	err = writeImage(img, imgName)
	if err != nil {
		log.Fatal(err)
	}
	if verbose {
		log.Printf("wrote texture to file")
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

func writeData(data []byte, filename string) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()

	n, err := writer.Write(data)
	if err != nil {
		return fmt.Errorf("position: %d error: %s", n, err)
	}

	return nil
}

func writeImage(image *image.RGBA, filename string) error {
	writer, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer writer.Close()

	err = png.Encode(writer, image)
	if err != nil {
		return err
	}

	return nil
}
