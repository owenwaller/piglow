// +build !sim

package main

func (p *PiGlowDriver) writeByteData(address, value byte) error {
	return p.conn.WriteByteData(address, value)
}

func (p *PiGlowDriver) writeBlockData(address byte, values []byte) error {
	return p.conn.WriteBlockData(address, values)
}
