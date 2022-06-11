package AskPassword

import (
	"fmt"
	"github.com/mattn/go-tty"
	"os"
	"os/signal"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Scan takes (printable) input till a newline is entered.
// The prefix is shown before the input field.
func Scan(prefix string) (string, error) {
	tty, err := tty.Open()
	if err != nil {
		return "", err
	}
	defer tty.Close()

	// handle interrupts (i.e. ctrl-c)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		tty.Close()
		os.Exit(1)
		return
	}()
	defer func() { sigchan <- os.Interrupt }() // kill goroutine after function has ended

	fmt.Print(prefix)
	var buf []string
	for {
		r, err := tty.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 13 { // rune 13 == return carriage
			fmt.Print("\n")
			break
		} else if r == 127 || r == 8 { // rune 127 == backspace
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
		} else {
			if unicode.IsPrint(r) {
				s := string(r)
				buf = append(buf, s)
				fmt.Print(s)
			} else {
				return strings.Join(buf, ""), fmt.Errorf("unprintable character entered")
			}
		}
	}
	return strings.Join(buf, ""), nil
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

	// handle interrupts (i.e. ctrl-c)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	go func() {
		<-sigchan
		tty.Close()
		os.Exit(1)
		return
	}()
	defer func() { sigchan <- os.Interrupt }() // kill goroutine after function has ended

	var buf []string
	var toggled bool
	fmt.Print(prefix, "(press TAB for no echo)")
	for {
		if len(buf) == 0 && toggled {
			fmt.Print("\r", Fillerstring(utf8.RuneCountInString(prefix), 24, " "), "\r", prefix)
		}
		r, err := tty.ReadRune()
		if err != nil {
			return "", err
		}
		if r == 13 { // rune 13 == return carriage
			fmt.Print("\n")
			break
		} else if r == 9 { // rune 9 == tab
			if !toggled && len(buf) == 0 {
				toggled = !toggled
			} else if toggled {
				space := Fillerstring(utf8.RuneCountInString(prefix), len(buf), " ")
				mask := Fillerstring(0, len(buf), "*")
				fmt.Print("\r", space, "\r", prefix, mask)
				toggled = !toggled
			} else {
				space := Fillerstring(utf8.RuneCountInString(prefix), len(buf), " ")
				fmt.Print("\r", space, "\r", prefix, strings.Join(buf, ""))
				toggled = !toggled
			}
		} else if r == 127 || r == 8 { // rune 127 == backspace
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
		} else {
			if unicode.IsPrint(r) {
				if len(buf) == 0 {
					space := Fillerstring(utf8.RuneCountInString(prefix), 24, " ")
					fmt.Print("\r", space, "\r", prefix)
				}
				buf = append(buf, string(r))
				if toggled {
					fmt.Print(string(r))
				} else {
					fmt.Print(substitute)
				}
			} else {
				return strings.Join(buf, ""), fmt.Errorf("unprintable character entered")
			}
		}
	}
	return strings.Join(buf, ""), nil
}

func AskPassword(prefix string) (string, error) {
	return ScanSecret(prefix, "*")
}
