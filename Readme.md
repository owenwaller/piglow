# GoBot PiGlow driver

## Work in porgress
This code is **not** yet a full GoBot driver. The code will need to be refactored first.
This is something I need help with. The GoBot documentation doesn't really describe the
driver API. If anyone understands the GoBot driver architecture I would appreciate
both pull requests and advice to help refactor the code.
I have an open GoBot issue [#377](https://github.com/hybridgroup/gobot/issues/377) to track this.

## The PiGlow Driver
The [PiGlow](https://shop.pimoroni.com/products/piglow) is a small add-on board for the [RaspberryPi](https://www.raspberrypi.org/) that provides 18 controllable LEDs.
The board is based around the [SN3218 microcontroller](https://github.com/pimoroni/piglow/tree/master/datasheets) that is used to control the LEDs

There are already a number of other PiGlow drivers, written in Go, as well as Python and Scratch but no GoBot driver as far as I could see. This project aims to fix that.

## Features
The current feature list is

*    light any single led in any spiral
*    light a single arm/spiral, so that's all 6 leds in a arm
*    light a single ring, so that's the same colour lit in each of the three spiral arms
*    light any random pattern of leds
*    fade up-down or down-up over a period any single led, spiral arm or ring
*    throb so, that's fade up-down-up or down-up-down, over a period for any led, spiral arm or ring
*    flash n times over a period for any led, spiral arm or ring
*    Fade, throb, or flash any random pattern of leds
*    Turn all of the leds off

## ToDo
- [x] remove the 5 second sleep at the end of the demo. This needs to be replaced with a "CTRL-C" or "q" to quit in the simulator.
- [ ] Refactor the code into a GoBot metal driver
- [ ] Add more tests
- [ ] Check the test coverage
- [ ] Add Gamma correction
- [ ] Create  a CLI application for the scripting community
- [ ] Vendor the simulators dependency

## Installing
Use go get
```
go get -u -v github.com/owenwaller/piglow
```
To use the simulator, which you can use if you do not have a PiGlow, you will also need the `tcell` package
```
go get -u -v github.com/gdamore/tcell
```

## Building
To build for a Raspberry Pi For ARM
```
GOARCH=arm GOARM=6 go build -v
```
To build the simulator
```
go install -tags=sim
```
Building with the race detector
```
go install -race -tags=sim

```
*Running with the race detector is currently only supported on platforms that support the go race detector. Currently this is AMD64 only. As such you can only run the race detector against the
simulator and not the actual hardware. However if you are race free on the simulator you will
also be race free on the (RaspberryPi/ARM) hardware.*

## Running
Without the race detector
```
piglow
```
With the race detector
```
GORACE="log_path=/tmp/" piglow
```

The `piglow` program will run through a demo of the drivers features.

### The Simulator
The simulator provides a terminal ui based simulation of the PiGlow's underlying SN3218 IC and the LEDs.

It displays the SN3218's register status and the intensity of each of the 18 leds.

Each arm has 6 leds - red, orange, yellow, green blue and white. The intensity value (in hex) is printed in the same colour as the LED.

If the value is printed in reverse then the LED would be lit with the given intensity on the real hardware. If the value is printed normally the led would be off (i.e. intensity zero)

The 4 key registers are also printed. These are:

* The shutdown register - shown as `EnableOutput` in the simulator (at address 0x00). This is used to enable and disable the LEDs. Setting it to zero switches any lit LEDs off - regardless of their intensity. Setting it to one restores the LED's at the previous intensity. This is used to create the a flashing effect by repeatedly cycling between shutdown and awake modes. See the SN32218 datasheet for the details.

* The LED control register - shown as `EnableLeds` in the simulator (starting at address 0x13). This register and the following two (at addresses 0x14 and 0x15) enable the LEDs in each arm.
The bottom 6 bits of each register enable or disable the individual LEDs. All LEDs are enabled by default.

* The update register updates the LED intensity. Any non-zero write into the register caused the
new LED intensity values to be flushed to the underlying PWM registers in the hardware, so updating the LEDs.

* The reset register resets all of the registers back to their default state. Any non-zero write will cause a reset. See the SN32218 datasheet for the details.

![The PiGlow Simualtor](https://github.com/owenwaller/piglow/blob/master/PiGlowSimulator.png)

To quit the simulator press CTRL-C or Enter or Escape.

### Platform Support
The simulator should run pretty much anywhere that Go runs. I've only tested it on Linux/AMD64 but Darwin/AMD64 and Windows/AMD64 should also work. The limiting factor will be the TCell terminal Ui library that is used to create the terminal Ui.

The code should run on any Raspberry Pi that GoBot supports (Pi 1 Model A and B/B+, Pi 2 and Pi 3 Model B) that can connect to a PiGlow.

I've only tested it on a Pi 1 Model B. I would welcome reports for others to confirm that the driver runs on ohter Raspberry Pi platforms.

### Dependancies
The driver requires GoBot v1.3.0 or higher. The driver requires the fix for issue [#372](https://github.com/hybridgroup/gobot/issues/372)

The simulator uses the [tcell](https://github.com/gdamore/tcell) package At the minute this dependency not vendored.
