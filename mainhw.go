// +build !sim

package main

import (
	"fmt"

	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	defer r.Finalize()
	err := r.Connect()
	if err != nil {
		fmt.Printf("Could not connect to raspi %v\n", err)
	}
	conn, err := r.GetConnection(address, r.GetDefaultBus())
	if err != nil {
		fmt.Printf("Failed to get connection: %v\n", err)
	}
	defer conn.Close()
	p := NewPiGlow("PiGlow", conn)

	DoMain(p)

}
