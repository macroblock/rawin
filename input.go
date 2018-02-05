package rawin

import (
	"fmt"
	"os"
	"time"
	"unicode/utf8"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	// RuneError -
	RuneError = utf8.RuneError
	// PreFilter -
	PreFilter = utf8.MaxRune + 1
	// PostFilter -
	PostFilter = utf8.MaxRune + 2
)

var (
	sin     *os.File
	state   *terminal.State
	buff    chan rune
	actions map[rune]func(r rune) bool
)

func init() {
	sin = os.Stdin
	state, _ = terminal.GetState(int(sin.Fd()))
}

// Init -
func Init(f *os.File) error {
	err := terminal.Restore(int(sin.Fd()), state)
	if err != nil {
		return err
	}
	fd := int(f.Fd())
	if !terminal.IsTerminal(fd) {
		return fmt.Errorf("it's not a terminal descriptor")
	}
	st, err := terminal.GetState(fd)
	if err != nil {
		return err
	}
	sin = f
	state = st
	return nil
}

// Close -
func Close() {
	Stop()
}

// Read -
func Read() (rune, error) {
	if buff == nil {
		return RuneError, fmt.Errorf("raw input mode thread is not running yet")
	}
	if len(buff) == 0 {
		return RuneError, fmt.Errorf("no events")
	}
	return <-buff, nil
}

// Start -
func Start() error {
	if buff != nil {
		return fmt.Errorf("raw input mode thread is already running")
	}
	buff = make(chan rune, 8) // TODO: bad constant

	fd := int(sin.Fd())
	st, err := terminal.MakeRaw(fd)
	if err != nil {
		return fmt.Errorf("cannot set raw mode")
	}
	state = st
	go readKeys()

	return nil
}

// Stop -
func Stop() {
	terminal.Restore(int(sin.Fd()), state)
	close(buff)
	for len(buff) > 0 {
		<-buff
		time.Sleep(10) // TODO: bad constant
	}
	buff = nil
}

// SetAction -
func SetAction(r rune, fn func(r rune) bool) {
	if actions == nil {
		actions = make(map[rune]func(r rune) bool)
	}
	if fn == nil {
		delete(actions, r)
		return
	}
	actions[r] = fn
}

// ClearActions -
func ClearActions() {
	actions = nil
}

func readKeys() {
	for buff != nil {
		key := RuneError
		b := [1]byte{}
		_, err := sin.Read(b[:])
		if err != nil {
			break
		}
		key = rune(b[0])

		for len(buff) > 0 {
			<-buff
		}
		fn, ok := actions[PreFilter]
		if ok && fn(key) {
			continue
		}
		fn, ok = actions[key]
		if ok && fn(key) {
			continue
		}
		fn, ok = actions[PostFilter]
		if ok && fn(key) {
			continue
		}
		buff <- key
	}
}
