// +build sim

package main

import "fmt"

const (
	singleByte = iota
	multiByte
)

type i2cData struct {
	payload byte
	address byte
	value   byte
	values  []byte
}

var (
	i2cBus    chan i2cData
	done      chan struct{}
	busClosed bool
)

func newI2CDataByte(address, value byte) i2cData {
	return i2cData{
		payload: singleByte,
		address: address,
		value:   value}
	// values defaults to nil
}

func newI2CDataBlock(address byte, values []byte) i2cData {
	return i2cData{
		payload: multiByte,
		address: address,
		values:  values}
	// value defaults to zero
}

func initI2CBus() {
	i2cBus = make(chan i2cData)
}

func cancelled() bool {
	select {
	case <-quit:
		if !busClosed {
			close(i2cBus)
			busClosed = true // make sure we don't call close twice
		}
		return true
	default:
		return false
	}
}

func (p *PiGlowDriver) writeByteData(address, value byte) error {
	if cancelled() {
		return fmt.Errorf("program quiting - write cancelled and bus closed")
	}
	i2cBus <- newI2CDataByte(address, value)
	return nil
}

func (p *PiGlowDriver) writeBlockData(address byte, values []byte) error {
	if cancelled() {
		return fmt.Errorf("program quiting - write cancelled and bus closed")
	}
	i2cBus <- newI2CDataBlock(address, values)
	return nil
}
