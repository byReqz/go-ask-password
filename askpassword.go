package AskPassword

import (
	"fmt"
	"github.com/mattn/go-tty"
	"unicode"
)

// Scan takes (printable) input till a newline is entered.
// The prefix is shown before the input field.
func Scan(prefix string) (string, error) {
	tty, err := tty.Open()
	if err != nil {
		return "", err
	}
	defer tty.Close()

	fmt.Print(prefix)
	var buf string
	for {
		r, err := tty.ReadRune()
		if err != nil {
			return buf, err
		}
		if r == 13 { // rune 13 == return carriage
			fmt.Print("\n")
			break
		} else {
			if unicode.IsPrint(r) {
				s := string(r)
				buf = buf + s
				fmt.Print(s)
			} else {
				return buf, fmt.Errorf("unprintable character entered")
			}
		}
	}
	return buf, nil
}

// ScanSecret takes (printable) input till a newline is entered.
// The prefix is shown before the input field.
// The substitute is what's shown instead of the entered character.
func ScanSecret(prefix string, substitute string) (string, error) {
	tty, err := tty.Open()
	if err != nil {
		return "", err
	}
	defer tty.Close()

	fmt.Print(prefix)
	var buf string
	for {
		r, err := tty.ReadRune()
		if err != nil {
			return buf, err
		}
		if r == 13 { // rune 13 == return carriage
			fmt.Print("\n")
			break
		} else {
			if unicode.IsPrint(r) {
				buf = buf + string(r)
				fmt.Print(substitute)
			} else {
				return buf, fmt.Errorf("unprintable character entered")
			}
		}
	}
	return buf, nil
}

func AskPassword(prefix string) (string, error) {
	return ScanSecret(prefix, "*")
}
