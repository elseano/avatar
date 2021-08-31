package main

import (
	"fmt"

	"github.com/elseano/avatar"
)

const (
	LettersCap = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func main() {
	outDir := "examples/output/"

	for index, letter := range LettersCap {
		otherindex := (index + 1) % 26
		fmt.Println(otherindex)
		initials := fmt.Sprintf("%c%c", letter, LettersCap[otherindex])
		fmt.Println(initials)
		fname := fmt.Sprintf("%s%s.png", outDir, initials)
		avatar.ToDisk(initials, fname)
	}

}
