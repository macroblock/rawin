package main

import (
	"fmt"
	"time"

	"github.com/macroblock/rawin"
)

var quit = false

func main() {
	rawin.SetAction(rawin.PreFilter, func(r rune) bool {
		fmt.Printf("key: %q, %v\n", r, int(r))
		return false
	})
	rawin.SetAction('t', func(r rune) bool {
		fmt.Printf("-> Test <-\n")
		return true
	})
	rawin.SetAction('q', func(r rune) bool {
		fmt.Printf("Quit\n")
		quit = true
		return true
	})

	err := rawin.Start()
	fmt.Printf("start err: %v\n", err)
	defer rawin.Stop()

	for !quit {
		time.Sleep(250 * time.Millisecond)
		fmt.Printf("ping\n")
	}
}
