package main

import (
	"fmt"
	"sync"
	"time"
)

func DoMain(p *PiGlowDriver) {
	err := p.Start()
	if err != nil {
		fmt.Printf("Could not start piglow %v\n", err)
	}

	var wg sync.WaitGroup

	// Faders
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p
		err := p.FadeLed(Spiral1, White, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not fade led %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p
		err := p.FadeRing(Green, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not fade ring %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p
		err := p.FadeSpiral(Spiral0, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not fade spiral %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var pattern [MaxSpirals][MaxColors]bool
		pattern[Spiral0][Green] = true
		pattern[Spiral1][White] = true
		pattern[Spiral2][Orange] = true // should panic if these are different!
		p := p
		err := p.FadeLeds(pattern, 128, 0, time.Second*2)
		if err != nil {
			fmt.Printf("Could not fade leds %v\n", err)
		}
	}()

	// Flashers
	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashLed(Spiral0, White, 10, 3, time.Second)
		if err != nil {
			fmt.Printf("Could not flash led %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashSpiral(Spiral2, 32, 15, time.Second*3)
		if err != nil {
			fmt.Printf("Could not flash spiral %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashRing(Red, 128, 3, time.Second*10)
		if err != nil {
			fmt.Printf("Could not flash ring %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var pattern [MaxSpirals][MaxColors]byte
		pattern[Spiral0][Green] = 64
		pattern[Spiral1][White] = 64
		pattern[Spiral2][Orange] = 64 // should panic if these are different!
		p := p
		err := p.FlashLeds(pattern, 3, time.Second*2)
		if err != nil {
			fmt.Printf("Could not flash leds %v\n", err)
		}
	}()

	// Throbbers
	wg.Add(1)
	go func() {
		defer wg.Done()
		var pattern [MaxSpirals][MaxColors]bool
		pattern[Spiral0][Orange] = true
		pattern[Spiral0][Blue] = true
		pattern[Spiral1][Orange] = true
		pattern[Spiral1][Blue] = true
		pattern[Spiral2][Orange] = true
		pattern[Spiral2][Blue] = true // should panic if these are different!
		p := p
		err := p.ThrobLeds(pattern, 64, 0, time.Second*5)
		if err != nil {
			fmt.Printf("Could not throb leds %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.ThrobRing(Orange, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not throb ring %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.ThrobSpiral(Spiral1, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not throb spiral %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.ThrobLed(Spiral1, Yellow, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not throb led %v\n", err)
		}
	}()

	// simple lighters - each lights the lied(s) for one second
	var pattern [MaxSpirals][MaxColors]byte
	pattern[Spiral0][Red] = 64
	pattern[Spiral1][White] = 64
	pattern[Spiral2][Blue] = 64
	err = p.LightLeds(pattern)
	if err != nil {
		fmt.Printf("Could not light pattern %v\n", err)
	}
	time.Sleep(time.Second)

	err = p.LightRing(Yellow, 128)
	if err != nil {
		fmt.Printf("Could not light ring %v\n", err)
	}
	time.Sleep(time.Second)

	err = p.LightSpiral(Spiral2, 128)
	if err != nil {
		fmt.Printf("Could not light spiral %v\n", err)
	}
	time.Sleep(time.Second)

	err = p.LightLed(Spiral2, Blue, 128)
	if err != nil {
		fmt.Printf("Could not light led %v\n", err)
	}
	time.Sleep(time.Second)

	wg.Wait()

	err = p.Reset()
	if err != nil {
		fmt.Printf("Could not halt %v\n", err)
	}
}
