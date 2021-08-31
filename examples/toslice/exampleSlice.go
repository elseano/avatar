package main

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/elseano/avatar"
)

const (
	LettersCap = "ABCDEFGHIJKLMNOPQRSTUVWXYZABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func main() {
	outDir := "examples/output/"
	_ = os.Mkdir(outDir, 0700)
	for index, letter := range LettersCap {
		otherindex := (index + 1) % 26
		//fmt.Println(otherindex)
		initials := fmt.Sprintf("%c%c", letter, LettersCap[otherindex])
		//fmt.Println(initials)
		salt := rand.Int()
		fname := fmt.Sprintf("%s%s%v.png", outDir, initials, salt)
		byteslice, err := avatar.ToSlice(initials, salt)
		if nil != err {
			panic(err)
		}
		file, err := os.OpenFile(fname, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
		if nil != err {
			panic(err)
		}
		defer file.Close()

		written, err := file.Write(byteslice)
		fmt.Printf("Wrote %v bytes to %v.\n", written, fname)
	}

}
