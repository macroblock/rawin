package main

import (
	"fmt"

	"github.com/macroblock/rawin"
)

var quit = false

func main() {
	// key, err := rawterm.Read()
	// fmt.Printf("key: %v, %q\nerror: %v\n", int(key), key, err)
	rawin.AddAction('t', func() {
		fmt.Printf("\nTEST\n")
	})
	rawin.AddAction('q', func() {
		fmt.Printf("\nQuit\n")
		quit = true
	})
	err := rawin.Start()
	fmt.Printf("err: %v\n", err)
	defer rawin.Stop()
	for !quit {
	}
}
