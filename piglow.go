package main

import (
	"math"
	"sync"
	"time"

	"gobot.io/x/gobot/drivers/i2c"
)

// Array with the bytes for all the individual leds mapped according to [leg][color]
var (
	spiral0 = [6]byte{6, 7, 8, 5, 4, 9}
	spiral1 = [6]byte{17, 16, 15, 13, 11, 10}
	spiral2 = [6]byte{0, 1, 2, 3, 14, 12}
	Spirals = [3][6]byte{spiral0, spiral1, spiral2}

	mu sync.Mutex
)

// Implements the gobot.io/x/gobot/sysfs/I2cDevice interface
type PiGlowDriver struct {
	name string
	conn i2c.Connection
}

func NewPiGlow(name string, conn i2c.Connection) *PiGlowDriver {
	p := new(PiGlowDriver)
	p.name = name
	p.conn = conn
	return p
}

// Name returns the label for the Driver
func (p *PiGlowDriver) Name() string {
	return p.name
}

// SetName sets the label for the Driver
func (p *PiGlowDriver) SetName(s string) {
	p.name = s
}

// Start initiates the Driver
func (p *PiGlowDriver) Start() error {
	mu.Lock()
	defer mu.Unlock()
	err := p.turnOn()
	if err != nil {
		return err
	}
	return p.enableLeds()
}

func (p *PiGlowDriver) turnOn() error {
	return p.writeByteData(enableOutput, turnOn)
	//return nil
}

func (p *PiGlowDriver) turnOff() error {
	return p.writeByteData(enableOutput, turnOff)
	//return nil
}

func (p *PiGlowDriver) enableLeds() error {
	return p.writeBlockData(enableLeds, []byte{0x3F, 0x3F, 0x3F})
}

// Halt terminates the Driver
// Resets the SN3218 IC to its default state.
func (p *PiGlowDriver) Halt() error {
	return p.turnOff()
}

// Connection returns the Connection assiciated with the Driver
func (p *PiGlowDriver) Connection() i2c.Connection {
	return p.conn
}

func (p *PiGlowDriver) Reset() error {
	mu.Lock()
	defer mu.Unlock()
	return p.reset()
}

func (p *PiGlowDriver) reset() error {
	return p.writeByteData(reset, 0xFF)
}

func (p *PiGlowDriver) LightLed(spiral, color, intensity byte) error {
	mu.Lock()
	defer mu.Unlock()
	return p.lightLed(spiral, color, intensity)
}

func (p *PiGlowDriver) lightLed(spiral, color, intensity byte) error {
	var values = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	values[Spirals[spiral][color]] = intensity
	return p.writeValues(values)
}

func (p *PiGlowDriver) FadeLed(spiral, color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	return p.fadeLed(spiral, color, initialIntensity, finalIntensity, fadeTime)
}

func (p *PiGlowDriver) fadeLed(spiral, color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	return p.fader(initialIntensity, finalIntensity, fadeTime, func(intensity byte) error {
		spiral := spiral
		color := color
		return p.lightLed(spiral, color, intensity)
	})
}

func (p *PiGlowDriver) fader(initialIntensity, finalIntensity byte, fadeTime time.Duration, lightFunc func(i byte) error) error {
	var d time.Duration
	d = fadeTime / time.Duration(int64(math.Abs(float64(initialIntensity-finalIntensity))))

	var err error
	var intensity = initialIntensity
	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for range ticker.C {
		err = lightFunc(intensity)
		if err != nil {
			return err
		}
		if initialIntensity < finalIntensity && intensity < finalIntensity {
			intensity++
		} else if initialIntensity > finalIntensity && intensity > finalIntensity {
			intensity--
		} else {
			break
		}
	}
	return err
}

func (p *PiGlowDriver) ThrobLed(spiral, color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	d := fadeTime / 2
	err := p.fadeLed(spiral, color, initialIntensity, finalIntensity, d)
	if err != nil {
		return err
	}
	return p.fadeLed(spiral, color, finalIntensity, initialIntensity, d)
}

func (p *PiGlowDriver) FlashLed(spiral, color, intensity byte, numberOfFlashes int, totalFlashTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	return p.flasher(numberOfFlashes, totalFlashTime, func() error {
		spiral := spiral
		color := color
		return p.lightLed(spiral, color, intensity)
	})
}

func (p *PiGlowDriver) LightLeds(pattern [MaxSpirals][MaxColors]byte) error {
	mu.Lock()
	defer mu.Unlock()
	return p.lightLeds(pattern)
}

func (p *PiGlowDriver) lightLeds(pattern [MaxSpirals][MaxColors]byte) error {
	var values = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for s := 0; s < MaxSpirals; s++ {
		for c := 0; c < MaxColors; c++ {
			values[Spirals[s][c]] = pattern[s][c]
		}
	}
	return p.writeValues(values)
}

func (p *PiGlowDriver) FadeLeds(pattern [MaxSpirals][MaxColors]bool, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	return p.fadeLeds(pattern, initialIntensity, finalIntensity, fadeTime)
}

func (p *PiGlowDriver) fadeLeds(pattern [MaxSpirals][MaxColors]bool, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	return p.fader(initialIntensity, finalIntensity, fadeTime, func(intensity byte) error {
		var values [MaxSpirals][MaxColors]byte
		for s := 0; s < MaxSpirals; s++ {
			for c := 0; c < MaxColors; c++ {
				if pattern[s][c] == true {
					values[s][c] = intensity
				}
			}
		}
		return p.lightLeds(values)

	})
}

func (p *PiGlowDriver) ThrobLeds(pattern [MaxSpirals][MaxColors]bool, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	d := fadeTime / 2
	err := p.fadeLeds(pattern, initialIntensity, finalIntensity, d)
	if err != nil {
		return err
	}
	return p.fadeLeds(pattern, finalIntensity, initialIntensity, d)
}

func (p *PiGlowDriver) FlashLeds(pattern [MaxSpirals][MaxColors]byte, numberOfFlashes int, totalFlashTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	return p.flasher(numberOfFlashes, totalFlashTime, func() error {
		pattern := pattern
		return p.lightLeds(pattern)
	})
}

func (p *PiGlowDriver) flasher(numberOfFlashes int, totalFlashTime time.Duration, lightFunc func() error) error {
	defer func() error {
		return p.turnOn()
	}()

	var d time.Duration
	d = totalFlashTime / time.Duration(int64(numberOfFlashes))

	var err error
	err = lightFunc()
	if err != nil {
		return err
	}
	var count int = 1

	ticker := time.NewTicker(d / 2)
	numberOfFlashes = 2*numberOfFlashes - 1
	defer ticker.Stop()
	for range ticker.C {
		if count%2 == 0 {
			//turn on
			err = p.turnOn()
			if err != nil {
				return err
			}
		} else {
			//turn off
			err = p.turnOff()
			if err != nil {
				return err
			}
		}
		count++
		if numberOfFlashes--; numberOfFlashes <= 0 {
			break
		}
	}
	return err

}

func (p *PiGlowDriver) LightSpiral(spiral, intensity byte) error {
	mu.Lock()
	defer mu.Unlock()
	return p.lightSpiral(spiral, intensity)
}

func (p *PiGlowDriver) lightSpiral(spiral, intensity byte) error {
	var values = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	for _, pin := range Spirals[spiral] {
		values[pin] = intensity
	}
	return p.writeValues(values)
}

func (p *PiGlowDriver) FadeSpiral(spiral, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	return p.fadeSpiral(spiral, initialIntensity, finalIntensity, fadeTime)
}

func (p *PiGlowDriver) fadeSpiral(spiral, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	return p.fader(initialIntensity, finalIntensity, fadeTime, func(intensity byte) error {
		spiral := spiral
		return p.lightSpiral(spiral, intensity)
	})
}

func (p *PiGlowDriver) ThrobSpiral(spiral, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	d := fadeTime / 2
	err := p.fadeSpiral(spiral, initialIntensity, finalIntensity, d)
	if err != nil {
		return err
	}
	return p.fadeSpiral(spiral, finalIntensity, initialIntensity, d)
}

func (p *PiGlowDriver) FlashSpiral(spiral, intensity byte, numberOfFlashes int, totalFlashTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	return p.flasher(numberOfFlashes, totalFlashTime, func() error {
		spiral := spiral
		intensity := intensity
		return p.lightSpiral(spiral, intensity)
	})
}

func (p *PiGlowDriver) LightRing(color, intensity byte) error {
	mu.Lock()
	defer mu.Unlock()
	return p.lightRing(color, intensity)
}

func (p *PiGlowDriver) lightRing(color, intensity byte) error {
	var values = []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	values[spiral0[color]] = intensity
	values[spiral1[color]] = intensity
	values[spiral2[color]] = intensity

	return p.writeValues(values)
}

func (p *PiGlowDriver) FadeRing(color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	return p.fadeRing(color, initialIntensity, finalIntensity, fadeTime)
}

func (p *PiGlowDriver) fadeRing(color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	return p.fader(initialIntensity, finalIntensity, fadeTime, func(intensity byte) error {
		color := color
		return p.lightRing(color, intensity)
	})
}

func (p *PiGlowDriver) ThrobRing(color, initialIntensity, finalIntensity byte, fadeTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()
	d := fadeTime / 2
	err := p.fadeRing(color, initialIntensity, finalIntensity, d)
	if err != nil {
		return err
	}
	return p.fadeRing(color, finalIntensity, initialIntensity, d)
}

func (p *PiGlowDriver) FlashRing(color, intensity byte, numberOfFlashes int, totalFlashTime time.Duration) error {
	mu.Lock()
	defer mu.Unlock()

	return p.flasher(numberOfFlashes, totalFlashTime, func() error {
		color := color
		intensity := intensity
		return p.lightRing(color, intensity)
	})
}

func (p *PiGlowDriver) TurnAllLedsOff() error {
	return p.Reset()
}

func (p *PiGlowDriver) writeValues(values []byte) error {
	err := p.writeBlockData(setPwmValues, values)
	if err != nil {
		return err
	}
	return p.writeByteData(update, 0xFF)
}
