// +build sim

package main

import (
	"fmt"

	"github.com/gdamore/tcell"
)

var (
	defaultStyle       = tcell.StyleDefault
	redStyle           = tcell.StyleDefault.Foreground(tcell.ColorRed)
	orangeStyle        = tcell.StyleDefault.Foreground(tcell.ColorOrange)
	yellowStyle        = tcell.StyleDefault.Foreground(tcell.ColorYellow)
	greenStyle         = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	blueStyle          = tcell.StyleDefault.Foreground(tcell.ColorBlue)
	whiteStyle         = tcell.StyleDefault.Foreground(tcell.ColorWhite)
	redReverseStyle    = tcell.StyleDefault.Foreground(tcell.ColorRed).Reverse(true)
	orangeReverseStyle = tcell.StyleDefault.Foreground(tcell.ColorOrange).Reverse(true)
	yellowReverseStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow).Reverse(true)
	greenReverseStyle  = tcell.StyleDefault.Foreground(tcell.ColorGreen).Reverse(true)
	blueReverseStyle   = tcell.StyleDefault.Foreground(tcell.ColorBlue).Reverse(true)
	whiteReverseStyle  = tcell.StyleDefault.Foreground(tcell.ColorWhite).Reverse(true)
)

const (
	dummy = iota
	titleRow
	emptyRowOne // place holder for an empty row
	registersRow
	emptyRowTwo // place holder for an empty row
	enableOutputRegisterRow
	enableLedsRegisterRow
	updateRegisterRow
	resetRegisterRow
	ringColoursRow
	arm0Row
	arm1Row
	arm2Row
)

var (
	enableOutputRegister byte
	enableLedsRegister   [3]byte
	updateRegister       byte
	resetRegister        byte
	pwmRegisters         [18]byte
)

const (
	title          = "PiGlow (sn3218 IC) Simulator"
	registersTitle = "sn3218 Registers"
	arm0Title      = "Arm 0: "
	arm1Title      = "Arm 1: "
	arm2Title      = "Arm 2: "
)

func (p *PiGlowDriver) writeByteData(address, value byte) error {
	//b := []byte{address, value}

	switch address {
	case enableOutput:
		updateEnableOutputRegister(value)
	case update:
		updateUpdateRegister(value)
	case reset:
		updateResetRegister(value)
	default:
		s := fmt.Sprintf("Unknown address: %x\n", address)
		panic(s)
	}
	return nil
}

func (p *PiGlowDriver) writeBlockData(address byte, values []byte) error {
	//	b := []byte{address}
	//	b := append(b, values...
	switch address {
	case setPwmValues:
		updatePwmRegisters(values)
	case enableLeds:
		updateEnableLedsRegister(values[:])
	default:
		s := fmt.Sprintf("Unknown address: %x\n", address)
		panic(s)
	}

	return nil
}

func initScreen() {
	displayTitle()
	displayEnableOutputRegister()
	displayEnableLedsRegister()
	displayUpdateRegister()
	displayResetRegister()
}

func updateEnableOutputRegister(value byte) {
	enableOutputRegister = value
	displayEnableOutputRegister()
	// if this has changed the leds need redrawn too
	displayArms()
}

func displayTitle() {
	putln(screen, titleRow, title)
	putln(screen, registersRow, registersTitle)
}

func displayEnableOutputRegister() {
	var str string
	str = fmt.Sprintf("EnableOutput (address:0x%02X) 0x%02X[%03d]", enableOutput, enableOutputRegister, enableOutputRegister)
	//fmt.Fprintf(os.Stderr, str)
	putln(screen, enableOutputRegisterRow, str)
}

func displayArmTitle(title string, row int) int {
	putln(screen, row, title)
	return len(title)
}

func updateEnableLedsRegister(values []byte) {
	copy(enableLedsRegister[:], values)
	//enableLedsRegister = value
	displayEnableLedsRegister()
}

func displayEnableLedsRegister() {
	var str string
	str = fmt.Sprintf("EnableLeds   (address:0x%02X) 0x%02X 0x%02X 0x%02X [%03d %03d %03d]", enableLeds, enableLedsRegister[0], enableLedsRegister[1], enableLedsRegister[2], enableLedsRegister[0], enableLedsRegister[1], enableLedsRegister[2])
	putln(screen, enableLedsRegisterRow, str)
}

func updateUpdateRegister(value byte) {
	updateRegister = value
	displayUpdateRegister()
	// if this has changed the leds ned redrawn too
	displayArms()
}

func displayUpdateRegister() {
	var str string
	str = fmt.Sprintf("Update       (address:0x%02X) 0x%02X[%03d]", update, updateRegister, updateRegister)
	putln(screen, updateRegisterRow, str)
}

func displayArms() {
	c := displayArmTitle(arm0Title, arm0Row)
	displayPwmRegistersForArm(spiral0, arm0Row, c)
	c = displayArmTitle(arm1Title, arm1Row)
	displayPwmRegistersForArm(spiral1, arm1Row, c)
	c = displayArmTitle(arm2Title, arm2Row)
	displayPwmRegistersForArm(spiral2, arm2Row, c)
}

func updateResetRegister(value byte) {
	resetRegister = value
	displayResetRegister()
}

func displayResetRegister() {
	var str string
	str = fmt.Sprintf("Reset        (address:0x%02X) 0x%02X[%03d]", reset, resetRegister, resetRegister)
	putln(screen, resetRegisterRow, str)
}

func updatePwmRegisters(values []byte) {
	copy(pwmRegisters[:], values)
	displayArms()
}

func displayPwmRegistersForArm(spiral [6]byte, row, column int) {
	redIntensity := pwmRegisters[spiral[Red]]
	orangeIntensity := pwmRegisters[spiral[Orange]]
	yellowIntensity := pwmRegisters[spiral[Yellow]]
	greenIntensity := pwmRegisters[spiral[Green]]
	blueIntensity := pwmRegisters[spiral[Blue]]
	whiteIntensity := pwmRegisters[spiral[White]]

	l := column
	l = l + printLed(redIntensity, l, row, redStyle, redReverseStyle)
	l = l + printLed(orangeIntensity, l, row, orangeStyle, orangeReverseStyle)
	l = l + printLed(yellowIntensity, l, row, yellowStyle, yellowReverseStyle)
	l = l + printLed(greenIntensity, l, row, greenStyle, greenReverseStyle)
	l = l + printLed(blueIntensity, l, row, blueStyle, blueReverseStyle)
	printLed(whiteIntensity, l, row, whiteStyle, whiteReverseStyle)
}

func printLed(intensity byte, line, row int, style, reverseStyle tcell.Style) int {
	str := fmt.Sprintf("0x%02X[%03d] ", intensity, intensity)
	if intensity == 0 || enableOutputRegister == 0 {
		putlnrowcol(screen, style, line, row, str)
	} else {
		putlnrowcol(screen, reverseStyle, line, row, str)
	}
	return len(str)
}

func putln(s tcell.Screen, row int, str string) {
	//	var ns tcell.Style = tcell.StyleDefault
	//	ns.Foreground(tcell.ColorRed)

	puts(s, defaultStyle, 0, row, str)
	s.Show()
}

func putlnrowcol(s tcell.Screen, sty tcell.Style, col, row int, str string) {
	puts(s, sty, col, row, str)
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	//fmt.Println(str)
	for _, r := range str {
		s.SetContent(x, y, r, nil, style)
		x += 1
	}
	s.Show()
}
