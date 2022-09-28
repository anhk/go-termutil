package termutil

import (
	"image"
)

type Sixel struct {
	X      uint16
	Y      uint64 // raw line
	Width  uint64
	Height uint64
	Image  image.Image
}

type VisibleSixel struct {
	ViewLineOffset int
	Sixel          Sixel
}

func (b *Buffer) clearSixelsAtRawLine(rawLine uint64) {
	var filtered []Sixel

	for _, sixelImage := range b.sixels {
		if sixelImage.Y+sixelImage.Height-1 >= rawLine && sixelImage.Y <= rawLine {
			continue
		}

		filtered = append(filtered, sixelImage)
	}

	b.sixels = filtered
}

func (b *Buffer) GetVisibleSixels() []VisibleSixel {

	firstLine := b.convertViewLineToRawLine(0)
	lastLine := b.convertViewLineToRawLine(b.viewHeight - 1)

	var visible []VisibleSixel

	for _, sixelImage := range b.sixels {
		if sixelImage.Y+sixelImage.Height-1 < firstLine {
			continue
		}
		if sixelImage.Y > lastLine {
			continue
		}

		visible = append(visible, VisibleSixel{
			ViewLineOffset: int(sixelImage.Y) - int(firstLine),
			Sixel:          sixelImage,
		})
	}

	return visible
}

func (t *Terminal) handleSixel(readChan chan MeasuredRune) (renderRequired bool) {

	var inEscape bool

	for {
		r := <-readChan

		switch r.Rune {
		case 0x1b:
			inEscape = true
			continue
		case 0x5c:
			if inEscape {
				return true
			}
		}

		inEscape = false
	}
}
