package avatar

import (
	"bufio"
	"encoding/hex"
	"errors"
	"image/color"
	"net/http"
	"strconv"
	"strings"
)

func HexToColor(hexColor string) (color.Color, error) {
	b, err := hex.DecodeString(hexColor)
	if err != nil {
		return nil, err
	}
	if len(b) == 4 {
		return color.RGBA{b[0], b[1], b[2], b[3]}, nil
	} else if len(b) == 3 {
		return color.RGBA{b[0], b[1], b[2], 255}, nil
	}
	return nil, errors.New("Not valid RGB(A) Hex Color")
}

func LoadHex(fs http.FileSystem, fname string) (color.Palette, error) {

	var err error
	var colors color.Palette

	file, err := fs.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			c, err := HexToColor(line)

			if err != nil {
				return colors, err
			}
			colors = append(colors, c)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return colors, err
}

func LoadGPL(fs http.FileSystem, fname string) (color.Palette, error) {

	var err error
	var colors color.Palette

	file, err := fs.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) > 0 {
			if line[0] != '#' {
				if !strings.HasPrefix(line, "Name:") &&
					!strings.HasPrefix(line, "Columns:") {
					pieces := strings.Split(line, "\t")
					if len(pieces) >= 3 {
						rb, err := strconv.ParseUint(pieces[0], 10, 8)
						if err != nil {
							return colors, err
						}
						gb, err := strconv.ParseUint(pieces[1], 10, 8)
						if err != nil {
							return colors, err
						}
						bb, err := strconv.ParseUint(pieces[2], 10, 8)
						if err != nil {
							return colors, err
						}

						c := color.RGBA{byte(rb), byte(gb), byte(bb), 255}
						colors = append(colors, c)

					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return colors, err
}
