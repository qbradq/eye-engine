package main

import (
	"bytes"
	"flag"
	"image"
	"image/png"
	"os"

	"github.com/qbradq/eye-engine/engine"
)

func main() {
	// Flags
	paletteName := flag.String(
		"palette",
		"default",
		"Name of the palette to dump.",
	)
	outputName := flag.String(
		"outfile",
		"default.png",
		"Name of the PNG file to dump to.",
	)
	flag.Parse()
	// Auto-set the output file name
	if *outputName == "default.png" {
		*outputName = *paletteName + ".png"
	}
	// Load the palette
	p, err := engine.LoadPalette(*paletteName)
	if err != nil {
		panic(err)
	}
	// Generate the palette image
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	for y := range 16 {
		for x := range 16 {
			img.SetRGBA(x, y, p[y*16+x])
		}
	}
	// Encode the image
	buf := bytes.NewBuffer(nil)
	if err := png.Encode(buf, img); err != nil {
		panic(err)
	}
	// Write the image file
	if err := os.WriteFile(*outputName, buf.Bytes(), 0666); err != nil {
		panic(err)
	}
}
