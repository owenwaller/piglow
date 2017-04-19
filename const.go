package main

const (
	address      = 0x54
	enableOutput = 0x00
	enableLeds   = 0x13
	setPwmValues = 0x01
	update       = 0x16
	reset        = 0x17

	turnOn  = 0x01
	turnOff = 0x00

	MaxColors  = 6
	MaxSpirals = 3
)

// All the possible colors
const (
	Red byte = iota
	Orange
	Yellow
	Green
	Blue
	White
)

const (
	Spiral0 byte = iota
	Spiral1
	Spiral2
)
