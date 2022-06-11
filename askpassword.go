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

func Fillerstring(prelen int, buflen int, filler string) string {
	var space string
	var i int
	for i < (prelen + buflen) {
		space = space + filler
		i++
	}
	return space
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

	var buf string
	var toggled bool
	for {
		if len(buf) == 0 && !toggled {
			fmt.Print(prefix, "(press TAB for no echo)")
		}
		r, err := tty.ReadRune()
		if err != nil {
			return buf, err
		}
		if r == 13 { // rune 13 == return carriage
			fmt.Print("\n")
			break
		} else if r == 9 { // rune 9 == tab
			if !toggled && len(buf) == 0 {
				toggled = !toggled
			} else if toggled {
				space := Fillerstring(len(prefix), len(buf), " ")
				mask := Fillerstring(0, len(buf), "*")
				fmt.Print("\r", space, "\r", prefix, mask)
				toggled = !toggled
			} else {
				space := Fillerstring(len(prefix), len(buf), " ")
				fmt.Print("\r", space, "\r", prefix, buf)
				toggled = !toggled
			}
		} else if r == 127 { // rune 127 == backspace
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
		} else {
			if unicode.IsPrint(r) {
				if len(buf) == 0 {
					space := Fillerstring(len(prefix), 24, " ")
					fmt.Print("\r", space, "\r", prefix)
				}
				buf = buf + string(r)
				if toggled {
					fmt.Print(string(r))
				} else {
					fmt.Print(substitute)
				}
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
