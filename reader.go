package termutil

import (
	"bufio"
	"bytes"
)

type Reader struct {
	r *bufio.Reader
}

func NewReader(data []byte) *Reader {
	reader := bufio.NewReader(bytes.NewBuffer(data))
	return &Reader{r: reader}
}

func (reader *Reader) ReadRune() MeasuredRune {
	r, size, err := reader.r.ReadRune()
	if err != nil {
		return MeasuredRune{}
	}
	return MeasuredRune{Rune: r, Width: size}
}
