// +build !sim

package main

func (p *PiGlowDriver) writeByteData(address, value byte) error {
	/*
		b := []byte{address, value}
		w, err := p.conn.Write(b)
		if err != nil {
			return err
		}
		if w != len(b) {
			return fmt.Errorf("failed to write the expected number of bytes. Wrote %d but expected %d", w, len(b))
		}
		return nil
	*/
	return p.conn.WriteByteData(address, value)
}

func (p *PiGlowDriver) writeBlockData(address byte, values []byte) error {
	/*
		b := []byte{address}
		b = append(b, values...)
		w, err := p.conn.Write(b)
		if err != nil {
			return err
		}
		if w != len(b) {
			return fmt.Errorf("failed to write the expected number of bytes. Wrote %d but expected %d", w, len(b))
		}
		return nil
	*/
	return p.conn.WriteBlockData(address, values)
}
