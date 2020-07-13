package avatar

import (
	"image"
	"image/color"
	"log"
	"strings"

	"github.com/argylelabcoat/avatar/palettes"
)

// Colors for background
var (
	White = image.Uniform{color.RGBA{255, 255, 255, 255}}
	Black = image.Uniform{color.RGBA{0, 0, 0, 255}}
)

const (
	LettersCap = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Threshold  = 130000
)

var loadedPalette color.Palette

// TODO add some more colors
func defaultColor(initial string) (image.Uniform, image.Uniform) {
	if len(loadedPalette) == 0 {
		var err error
		loadedPalette, err = LoadHex(palettes.FS(false), "/downgraded-36.hex")
		if err != nil {
			log.Println("Failed to load Palette", err)
		}
	}

	numCols := len(loadedPalette)

	upperInitial := strings.ToUpper(initial)

	colorIndex := strings.Index(LettersCap, upperInitial) % numCols
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
