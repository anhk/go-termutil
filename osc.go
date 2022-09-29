package termutil

func (t *Terminal) handleOSC(reader *Reader) (renderRequired bool) {

READ:
	for {
		b := reader.ReadRune()
		if t.isOSCTerminator(b.Rune) {
			break READ
		}
		if b.Rune == ';' {
			continue
		} else {
			return false
		}
	}

	return false
}

func (t *Terminal) isOSCTerminator(r rune) bool {
	for _, terminator := range oscTerminators {
		if terminator == r {
			return true
		}
	}
	return false
}
