package main

import (
	"fmt"
	"time"

	"github.com/macroblock/rawin"
)

var quit = false

func main() {
	rawin.AddAction(rawin.PreFilter, func(r rune) bool {
		fmt.Printf("key: %q, %v\n", r, int(r))
		return false
	})
	rawin.AddAction('t', func(r rune) bool {
		fmt.Printf("-> Test <-\n")
		return true
	})
	rawin.AddAction('q', func(r rune) bool {
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
