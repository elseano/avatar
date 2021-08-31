package avatar

import (
	"embed"
	"errors"
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
var errInvalidFormat = errors.New("invalid format")

func parseOrDefault(initial, bgColor, fontColor string) (bgC *image.Uniform, fgC *image.Uniform) {
	if len(bgColor) == 0 || len(fontColor) == 0 {
		bgC, fgC = defaultColor(initial)
	}
	if len(bgColor) != 0 {
		bgC = colorFromHex(bgColor)
	}
	if len(fontColor) != 0 {
		fgC = colorFromHex(fontColor)
	}
	return bgC, fgC
}

// parseHexColorFast was found here:
// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func parseHexColorFast(s string) (c color.RGBA, err error) {
	c.A = 0xff

	if s[0] != '#' {
		return c, errInvalidFormat
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errInvalidFormat
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	case 4:
		c.R = hexToByte(s[1]) * 17
		c.G = hexToByte(s[2]) * 17
		c.B = hexToByte(s[3]) * 17
	default:
		err = errInvalidFormat
	}
	return
}

func colorFromHex(s string) (colImg *image.Uniform) {
	c, err := parseHexColorFast(s)
	if err == nil {
		return &image.Uniform{c}
	} else {
		return &White
	}
}

// TODO add some more colors
func defaultColor(initial string) (bgC *image.Uniform, fgC *image.Uniform) {
	return saltedColor(initial, 0)

}

// TODO add some more colors
func saltedColor(initial string, salt int) (bgC *image.Uniform, fgC *image.Uniform) {
	if len(loadedPalette) == 0 {
		var err error
		loadedPalette, err = LoadHex(palettes, "palettes/downgraded-36.hex")
		if err != nil {
			log.Println("Failed to load Palette", err)
		}
	}

	numCols := len(loadedPalette)

	upperInitial := strings.ToUpper(initial)

	colorIndex := (strings.Index(LettersCap, upperInitial) + salt) % numCols
	if colorIndex < 0 {
		colorIndex = rand.Intn(numCols)
	}
	colorFromIndex := loadedPalette[colorIndex]
	bgC = &image.Uniform{colorFromIndex}

	r, g, b, _ := colorFromIndex.RGBA()
	bgsum := r + g + b
	log.Println(upperInitial, bgsum, r, g, b)

	if bgsum > Threshold {
		return bgC, &Black
	} else {
		return bgC, &White
	}
}
