// +build sim

package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var screen tcell.Screen

func main() {
	initUi()
	defer closeUi()
	// just create a PiGlow
	p := NewPiGlow("PiGlow", nil) // the connection is repalced in the sim so nil is safe
	DoMain(p)
}

func initUi() {
	var err error
	screen, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}
	encoding.Register()
	screen.Clear()
	initScreen()
}

func closeUi() {
	screen.Fini()
}
