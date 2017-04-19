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

	var wg sync.WaitGroup

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
		err := p.FadeSpiral(Spiral0, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not fade led %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		var pattern [MaxSpirals][MaxColors]bool
		pattern[Spiral0][Red] = true
		pattern[Spiral0][Green] = true
		pattern[Spiral1][Red] = true
		pattern[Spiral1][Green] = true
		pattern[Spiral2][Red] = true
		pattern[Spiral2][Green] = true // should panic if these are different!
		p := p
		err := p.FadeLeds(pattern, 64, 0, time.Second*2)
		if err != nil {
			fmt.Printf("Could not fade leds %v\n", err)
		}
	}()

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
			fmt.Printf("Could not fade leds %v\n", err)
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
			fmt.Printf("Could not fade leds %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p
		err := p.FadeRing(Green, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not fade led %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.ThrobRing(Orange, 128, 0, time.Second*3)
		if err != nil {
			fmt.Printf("Could not throb led %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashRing(Red, 128, 3, time.Second*10)
		if err != nil {
			fmt.Printf("Could not light spiral %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashLed(Spiral0, White, 128, 10, time.Second*10)
		if err != nil {
			fmt.Printf("Could not light spiral %v\n", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		p := p

		err := p.FlashSpiral(Spiral2, 32, 15, time.Second*3)
		if err != nil {
			fmt.Printf("Could not light spiral %v\n", err)
		}
	}()

	wg.Wait()
	// pause - required for the sim
	// TODO  replace this quit a CTRL-C or q to quit in the sim
	time.Sleep(5 * time.Second)
	err = p.Reset()
	if err != nil {
		fmt.Printf("Could not halt %v\n", err)
	}
}
