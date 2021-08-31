package avatar

import (
	"bufio"
	"bytes"
	"embed"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	defaultfontFace = "fonts/Roboto-Bold.ttf" //SourceSansVariable-Roman.ttf"
	fontSize        = 60.0
	imageWidth      = 100.0
	imageHeight     = 100.0
	dpi             = 72.0
	spacer          = 4
	textY           = 71
)

//go:embed fonts/*
var fonts embed.FS

// ToDisk saves the image to disk
func ToDisk(initials, path string) error {
	return saveToDisk(initials, path, "", "")
}

// ToDiskCustom saves the image to disk
func ToDiskCustom(initials, path, bgColor, fontColor string) error {
	return saveToDisk(initials, path, bgColor, fontColor)
}

// saveToDisk saves the image to disk
func saveToDisk(initials, path, bgColor, fontColor string) error {
	bgC, fgC := parseOrDefault(initials, bgColor, fontColor)
	rgba, err := createAvatar(initials, bgC, fgC)
	if err != nil {
		return fmt.Errorf("unable to create avatar: %w", err)
	}

	// Save image to disk
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unable to create target path: %w", err)
	}
	defer out.Close()

	b := bufio.NewWriter(out)

	err = png.Encode(b, rgba)
	if err != nil {
		return fmt.Errorf("unable to encode image: %w", err)
	}

	err = b.Flush()
	if err != nil {
		return err
	}

	return nil
}

// ToHTTP sends the image to a http.ResponseWriter (as a PNG)
func ToHTTP(initials string, w http.ResponseWriter) error {
	return saveToHTTP(initials, "", "", w)
}

// ToHTTPCustom sends the image to a http.ResponseWriter (as a PNG)
func ToHTTPCustom(initials, bgColor, fontColor string, w http.ResponseWriter) error {
	return saveToHTTP(initials, bgColor, fontColor, w)
}

// saveToHTTP sends the image to a http.ResponseWriter (as a PNG)
func saveToHTTP(initials, bgColor, fontColor string, w http.ResponseWriter) error {
	bgC, fgC := parseOrDefault(initials, bgColor, fontColor)
	rgba, err := createAvatar(initials, bgC, fgC)
	if err != nil {
		return err
	}

	b := new(bytes.Buffer)
	key := fmt.Sprintf("avatar%s", initials) // for Etag

	err = png.Encode(b, rgba)
	if err != nil {
		return fmt.Errorf("unable to encode image: %w", err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(b.Bytes())))
	w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	w.Header().Set("Etag", `"`+key+`"`)

	if _, err := w.Write(b.Bytes()); err != nil {
		return fmt.Errorf("unable to write image to response: %w", err)
	}

	return nil
}

// ToSlice simply buffers the image and returns the byte slice (as a PNG)
func ToSlice(initials string, colorSalt int) ([]byte, error) {
	bgC, fgC := saltedColor(initials, colorSalt)
	rgba, err := createAvatar(initials, bgC, fgC)
	if err != nil {
		return nil, fmt.Errorf("unable to create avatar: %w", err)
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, rgba)
	if nil == err {
		return buf.Bytes(), err
	} else {
		return nil, fmt.Errorf("unable to encode image: %w", err)
	}
}

// ToSlice simply buffers the image and returns the byte slice (as a PNG)
func ToSliceCustomColors(initials string, bgColor, fontColor string) ([]byte, error) {

	bgC, fgC := parseOrDefault(initials, bgColor, fontColor)
	rgba, err := createAvatar(initials, bgC, fgC)
	if err != nil {
		return nil, fmt.Errorf("unable to create avatar: %w", err)
	}
	buf := new(bytes.Buffer)
	err = png.Encode(buf, rgba)
	if nil == err {
		return buf.Bytes(), err
	} else {
		return nil, fmt.Errorf("unable to encode image: %w", err)
	}
}

func cleanString(incoming string) string {
	incoming = strings.TrimSpace(incoming)

	// If its something like "firstname surname" get the initials out
	split := strings.Split(incoming, " ")
	if len(split) == 2 {
		incoming = split[0][0:1] + split[1][0:1]
	}

	// Max length of 2
	if len(incoming) > 2 {
		incoming = incoming[0:2]
	}

	// To upper and trimmed
	return strings.ToUpper(strings.TrimSpace(incoming))
}

func getFont(fontPath string) (theFont *truetype.Font, err error) {
	if fontPath == "" {
		fontPath = defaultfontFace
	}
	// Read the font data.
	var fontBytes []byte

	fontBytes, err = fonts.ReadFile(fontPath)

	if err != nil {
		return nil, fmt.Errorf("unable to read font file: %w", err)
	}

	return freetype.ParseFont(fontBytes)
}

var imageCache sync.Map

func cacheKey(initials string, fg color.Color, bg color.Color) string {
	return strings.ToLower(fmt.Sprintf("%s%v%v", initials, fg, bg))
}

func getImage(initials string, fg color.Color, bg color.Color) *image.RGBA {
	value, ok := imageCache.Load(cacheKey(initials, fg, bg))

	if !ok {
		return nil
	}

	image, ok2 := value.(*image.RGBA)
	if !ok2 {
		return nil
	}
	return image
}

func setImage(initials string, fg color.Color, bg color.Color, image *image.RGBA) {
	imageCache.Store(cacheKey(initials, fg, bg), image)
}

func createAvatar(initials string, bgColor, fontColor *image.Uniform) (*image.RGBA, error) {
	// Make sure the string is OK
	text := cleanString(initials)

	// Check cache
	cachedImage := getImage(text, fontColor.C, bgColor.C)
	if cachedImage != nil {
		return cachedImage, nil
	}

	// Load and get the font
	f, err := getFont(defaultfontFace)
	if err != nil {
		return nil, fmt.Errorf("invalid font: %w", err)
	}

	rgba := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(rgba, rgba.Bounds(), bgColor, image.ZP, draw.Src)
	c := freetype.NewContext()
	c.SetDPI(dpi)
	c.SetFont(f)
	c.SetFontSize(fontSize)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(fontColor)
	c.SetHinting(font.HintingFull)

	// We need to convert the font into a "font.Face" so we can read the glyph
	// info
	to := truetype.Options{}
	to.Size = fontSize
	face := truetype.NewFace(f, &to)

	// Calculate the widths and print to image
	xPoints := []int{0, 0}
	textWidths := []int{0, 0}

	// Get the widths of the text characters
	for i, char := range text {
		width, ok := face.GlyphAdvance(rune(char))
		if !ok {
			return nil, err
		}

		textWidths[i] = int(float64(width) / 64)
	}

	// TODO need some tests for this
	if len(textWidths) == 1 {
		textWidths[1] = 0
	}

	// Get the combined width of the characters
	combinedWidth := textWidths[0] + spacer + textWidths[1]

	// Draw first character
	xPoints[0] = int((imageWidth - combinedWidth) / 2)
	xPoints[1] = int(xPoints[0] + textWidths[0] + spacer)

	for i, char := range text {
		pt := freetype.Pt(xPoints[i], textY)
		_, err := c.DrawString(string(char), pt)
		if err != nil {
			return nil, err
		}
	}

	// Cache it
	setImage(initials, fontColor.C, bgColor.C, rgba)

	return rgba, nil
}
