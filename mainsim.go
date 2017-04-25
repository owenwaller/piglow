// +build sim

package main

import (
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
)

var (
	screen tcell.Screen
	quit   chan struct{}
)

func main() {
	initUi()
	defer closeUi()
	// make the quit channel
	quit = make(chan struct{})
	initI2CBus()
	// just create a PiGlow
	p := NewPiGlow("PiGlow", nil) // the connection is repalced in the sim so nil is safe

	// start the sim
	go runSimulator()

	// run the demo code
	go DoMain(p)

	// make sure we quit if CTRL-C or Escape or Enter are pressed
	go pressKeyToExit()

	// wait for the qiuit channel to close - the defered closeUi will cleanup
	<-quit
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
	panic("End Of Program - only main should remain")
}

func pressKeyToExit() {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape, tcell.KeyEnter, tcell.KeyCtrlC:
				close(quit)
				return
			}
		}
	}
}
