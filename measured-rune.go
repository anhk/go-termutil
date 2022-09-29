package termutil

type MeasuredRune struct {
	Rune  rune
	Width int
}

func (m MeasuredRune) Empty() bool {
	return m.Width == 0
}
