package askpassword

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"unicode"
	"unicode/utf8"

	"github.com/fatih/color"
	"github.com/mattn/go-isatty"
	"github.com/mattn/go-tty"
)

const (
	keyReturn        = 13  // rune for return carriage
	keyTab           = 9   // rune for tab
	keyBackspace     = 127 // rune for backspace
	keyCtrlBackspace = 23  // rune for ctrl + backspace
	//keyCtrlD         = 4   // rune for ctrl + d
)

type Options struct {
	Plain     bool // does not hide the password if true
	PullStdin bool // starts the prompt with the contents from stdin

	Prefix      string // whats shown in front of the password field
	Substitute  string // character/string shown instead of the password input
	Placeholder string // whats shown while there has not been user input yet

	BreakAt           rune // character that closes/exits the prompt (def: 13 / enter)
	ExitOnUnprintable bool // exit with an error if an unprintable character is entered, will ignore them if false
	MaxLength         int  // maximum input length (def: none)
}

// fillerstring constructs the string used to hide secrets
func fillerstring(prelen int, buflen int, filler string) string {
	var space string
	var i int
	for i < (prelen + buflen) {
		space = space + filler
		i++
	}
	return space
}

// readStdin will read the contents of /dev/stdin and return it in the buffer format used by the prompt
func readStdin() (stdin []string) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return
	}
	infile, err := os.Open(os.Stdin.Name())
	if err != nil {
		return
	}
	buf, err := io.ReadAll(infile)
	if err != nil {
		return
	}
	for _, c := range string(buf) {
		stdin = append(stdin, string(c))
	}
	return
}

// NewScan starts a new scanner with the given options.
func NewScan(opts Options) (string, error) {
	var inbuf []string // stdin contents before initializing the tty
	if opts.PullStdin {
		inbuf = readStdin()
	}

	t, err := tty.Open()
	if err != nil {
		return "", err
	}
	defer func(t *tty.TTY) {
		_ = t.Close()
	}(t)

	// handle interrupts (i.e. ctrl-c)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)
	go func() {
		sn := <-sigchan
		if sn == os.Interrupt || sn == syscall.SIGTERM {
			_ = t.Close()
			os.Exit(1)
		}
		signal.Stop(sigchan)
	}()
	defer close(sigchan)

	if opts.BreakAt == 0 { // set enter as breakpoint if none is set
		opts.BreakAt = keyReturn
	}

	var buf []string  // main buffer for the input prompt, each char is one string
	var revealed bool // if the contents of the buffer are revealed to the user
	fmt.Print(opts.Prefix, color.HiBlackString(opts.Placeholder))
	for {
		if opts.PullStdin {
			for _, si := range inbuf {
				t.Input().WriteString(si)
			}
		}

		r, err := t.ReadRune()
		if err != nil {
			return "", err
		}

		if r == opts.BreakAt { // def: keyReturn == return carriage
			fmt.Print("\n")
			break
		} else if r == keyTab { // rune 9 == tab
			if opts.Plain {
				continue
			}
			space := fillerstring(utf8.RuneCountInString(opts.Prefix), len(buf), " ")
			mask := fillerstring(0, len(buf), opts.Substitute)
			if !revealed {
				space = fmt.Sprint("\r", fillerstring(utf8.RuneCountInString(opts.Prefix), utf8.RuneCountInString(opts.Placeholder), " "), "\r", opts.Prefix)
				mask = strings.Join(buf, "")
			}
			fmt.Print("\r", space, "\r", opts.Prefix, mask)
			revealed = !revealed
			continue
		} else if r == keyBackspace { // rune 127 == backspace, "|| r == 8" used to be here but i cant find 8 on the keyboard anymore :|
			if len(buf) > 0 {
				buf = buf[:len(buf)-1]
				fmt.Print("\b \b")
			}
			continue
		} else if r == keyCtrlBackspace { // delete all input on ctrl + backspace
			fmt.Print(strings.Repeat("\b", len(buf)))
			fmt.Print(strings.Repeat(" ", len(buf)))
			fmt.Print(strings.Repeat("\b", len(buf)))
			buf = []string{}
			continue
		} else if opts.MaxLength > 0 && len(buf) == opts.MaxLength {
			continue
		} else if r == 0 { // drop NULLs (see issue #1, not necessary for organic input)
			continue
		}

		if !unicode.IsPrint(r) {
			if opts.ExitOnUnprintable {
				return strings.Join(buf, ""), fmt.Errorf("unprintable character entered")
			}
			continue
		}

		if len(buf) == 0 {
			space := fillerstring(utf8.RuneCountInString(opts.Prefix), 24, " ")
			fmt.Print("\r", space, "\r", opts.Prefix)
		}
		buf = append(buf, string(r))
		if revealed || opts.Plain {
			fmt.Print(string(r))
		} else {
			fmt.Print(opts.Substitute)
		}
	}
	return strings.Join(buf, ""), nil
}

// Scan takes (printable) input until a newline is entered.
func Scan(prefix string) (string, error) {
	return NewScan(Options{
		PullStdin: true,
		Plain:     true,
		Prefix:    prefix,
	})
}

// ScanSecret takes (printable) input till a newline is entered but hides the content by default.
func ScanSecret(prefix string, substitute string, placeholder string) (string, error) {
	return NewScan(Options{
		Prefix:      prefix,
		Substitute:  substitute,
		Placeholder: placeholder,
	})
}

// AskPassword is an opinionated default Password prompt like systemd-ask-password
func AskPassword(prefix string) (string, error) {
	return ScanSecret(color.New(color.Bold, color.FgHiWhite).Sprint("🔐"+prefix), "*", "(press TAB for echo)")
}

// AskUser is an opinionated default Username prompt
func AskUser(prefix string) (string, error) {
	return Scan(color.New(color.Bold, color.FgHiWhite).Sprint("👤" + prefix))
}

// AskKey is an opinionated default Password prompt like systemd-ask-password
func AskKey(prefix string) (string, error) {
	return ScanSecret(color.New(color.Bold, color.FgHiWhite).Sprint("🔑"+prefix), "*", "(press TAB for echo)")
}
