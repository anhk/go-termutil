package termutil

func (t *Terminal) handleOSC(readChan chan MeasuredRune) (renderRequired bool) {

READ:
	for {
		select {
		case b := <-readChan:
			if t.isOSCTerminator(b.Rune) {
				break READ
			}
			if b.Rune == ';' {
				continue
			}
		default:
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
