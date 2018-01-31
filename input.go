package rawin

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/ssh/terminal"
)

var (
	sin      = os.Stdin
	state, _ = terminal.GetState(int(sin.Fd()))
	buff     chan rune
	actions  map[rune]func()
)

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
	//terminal.Restore(int(sin.Fd()), state)
}

// Read -
func Read() (rune, error) {
	// fd := int(sin.Fd())
	// state, err := terminal.MakeRaw(fd)
	// if err != nil {
	// 	return 0, fmt.Errorf("cannot set raw mode")
	// }
	// defer terminal.Restore(fd, state)

	// var buf [1]byte
	// _, err = sin.Read(buf[:])
	// key := rune(buf[0])
	// return key, err
	if buff == nil {
		return 0, fmt.Errorf("raw input mode thread is not running yet")
	}
	if len(buff) == 0 {
		return 0, fmt.Errorf("no events")
	}
	return <-buff, nil
}

// Start -
func Start() error {
	if buff != nil {
		return fmt.Errorf("raw input mode thread is already running")
	}
	buff = make(chan rune)

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

// AddAction -
func AddAction(r rune, fn func()) {
	if actions == nil {
		actions = make(map[rune]func())
	}
	actions[r] = fn
}

// ClearActions -
func ClearActions() {
	actions = nil
}

func readKeys() {
	for buff != nil {
		var buf [1]byte
		_, err := sin.Read(buf[:])
		key := rune(buf[0])
		if err != nil {
			break
		}
		// for len(buff) > 0 {
		// 	<-buff
		// }

		if fn, ok := actions[key]; ok {
			fn()
		}
		//buff <- key
	}
}

// Size -
func Size() (int64, error) {
	info, err := sin.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), err
}
