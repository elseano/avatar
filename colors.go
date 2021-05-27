package avatar

import (
	"embed"
	"image"
	"image/color"
	"log"
	"math/rand"
	"strings"
)

// Colors for background
var (
	White = image.Uniform{color.RGBA{255, 255, 255, 255}}
	Black = image.Uniform{color.RGBA{0, 0, 0, 255}}
)

//go:embed palettes/*
var palettes embed.FS

const (
	LettersCap = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Threshold  = 130000
)

var loadedPalette color.Palette

// TODO add some more colors
func defaultColor(initial string) (image.Uniform, image.Uniform) {
	if len(loadedPalette) == 0 {
		var err error
		loadedPalette, err = LoadHex(palettes, "palettes/downgraded-36.hex")
		if err != nil {
			log.Println("Failed to load Palette", err)
		}
	}

	numCols := len(loadedPalette)

	upperInitial := strings.ToUpper(initial)

	colorIndex := strings.Index(LettersCap, upperInitial) % numCols
	if colorIndex < 0 {
		colorIndex = rand.Intn(numCols)
	}
	colorFromIndex := loadedPalette[colorIndex]
	bgC := image.Uniform{colorFromIndex}

	r, g, b, _ := colorFromIndex.RGBA()
	bgsum := r + g + b
	log.Println(upperInitial, bgsum, r, g, b)
	if bgsum > Threshold {
		return bgC, Black
	} else {
		return bgC, White
	}

}
