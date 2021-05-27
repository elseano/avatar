package avatar

import (
	"image"
	"image/draw"

	"github.com/anthonynsimon/bild/adjust"
	"github.com/anthonynsimon/bild/blend"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/noise"
)

type BGMethod int

const (
	Fast BGMethod = iota
	LinenLike
	Drops
)

var defaultBGType = Fast

func SetDefaultBGType(bgtype BGMethod) {
	defaultBGType = bgtype
}

func GetDefaultBG(bg image.Uniform) *image.RGBA {
	return GetBG(bg, defaultBGType)
}

func GetBG(bg image.Uniform, method BGMethod) *image.RGBA {
	switch method {
	case Fast:
		return FastBG(bg)
	case LinenLike:
		return LinenLikeBG(bg)
	case Drops:
		return DropsBG(bg)
	default:
		return FastBG(bg)
	}
}

func FastBG(bg image.Uniform) *image.RGBA {
	bgc := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(bgc, bgc.Bounds(), &bg, image.Point{}, draw.Src)
	return bgc
}

func DropsBG(bg image.Uniform) *image.RGBA {
	bgc := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(bgc, bgc.Bounds(), &bg, image.Point{}, draw.Src)
	bgi := noise.Generate(imageWidth, imageHeight, &noise.Options{Monochrome: true, NoiseFn: noise.Gaussian})
	bgi = effect.Dilate(bgi, 6)

	bgi = adjust.Contrast(bgi, -0.3)
	bgi = adjust.Brightness(bgi, 0.1)

	return blend.Multiply(bgc, bgi)

}

func LinenLikeBG(bg image.Uniform) *image.RGBA {
	bgc := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(bgc, bgc.Bounds(), &bg, image.Point{}, draw.Src)
	bgi := noise.Generate(imageWidth, imageHeight, &noise.Options{Monochrome: true, NoiseFn: noise.Binary})
	bgi = effect.Median(bgi, 10.0)

	bgi = adjust.Contrast(bgi, -0.95)
	bgi = adjust.Brightness(bgi, 0.5)

	return blend.Multiply(bgc, bgi)
}
