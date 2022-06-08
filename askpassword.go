package AskPassword

import (
	"fmt"
	"github.com/mattn/go-tty"
	"log"
)

func Scanln(prefix string, substitute string) (string, error) {
	tty, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer tty.Close()

	fmt.Print(prefix)
	var buf string
	for {
		r, err := tty.ReadRune()
		if err != nil {
			log.Fatal(err)
		}
		if string(r) == "\r" || string(r) == "\n" {
			fmt.Print("\n")
			break
		} else {
			buf = buf + string(r)
			fmt.Print(substitute)
		}
	}
	return buf, nil
}

func AskPassword(prefix string) (string, error) {
	return Scanln(prefix, "*")
}
