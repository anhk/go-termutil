package termutil

import (
	"fmt"
	"os"
	"sync"
)

const (
	MainBuffer     uint8 = 0
	AltBuffer      uint8 = 1
	InternalBuffer uint8 = 2
)

// Terminal communicates with the underlying terminal
type Terminal struct {
	mu           sync.Mutex
	buffers      []*Buffer
	activeBuffer *Buffer
	mouseMode    MouseMode
	mouseExtMode MouseExtMode
	logFile      *os.File
	theme        *Theme

	RequestRender func()
}

// NewTerminal creates a new terminal instance
func New(options ...Option) *Terminal {
	term := &Terminal{
		theme: &Theme{},
	}
	for _, opt := range options {
		opt(term)
	}
	fg := term.theme.DefaultForeground()
	bg := term.theme.DefaultBackground()
	term.buffers = []*Buffer{
		NewBuffer(4096, 80, 0xffff, fg, bg),
		NewBuffer(4096, 80, 0xffff, fg, bg),
		NewBuffer(4096, 80, 0xffff, fg, bg),
	}
	term.activeBuffer = term.buffers[0]

	return term
}

func (t *Terminal) log(line string, params ...interface{}) {
	if t.logFile != nil {
		_, _ = fmt.Fprintf(t.logFile, line+"\n", params...)
	}
}

func (t *Terminal) reset() {
	fg := t.theme.DefaultForeground()
	bg := t.theme.DefaultBackground()
	t.buffers = []*Buffer{
		NewBuffer(4096, 80, 0xffff, fg, bg),
		NewBuffer(4096, 80, 0xffff, fg, bg),
		NewBuffer(4096, 80, 0xffff, fg, bg),
	}
	t.useMainBuffer()
}

func (t *Terminal) Theme() *Theme {
	return t.theme
}

// // write takes data from StdOut of the child shell and processes it
// func (t *Terminal) Write(data []byte) (n int, err error) {
// 	reader := bufio.NewReader(bytes.NewBuffer(data))
// 	for {
// 		r, size, err := reader.ReadRune()
// 		if err == io.EOF {
// 			break
// 		}
// 		t.processChan <- MeasuredRune{Rune: r, Width: size}
// 	}
// 	return len(data), nil
// }

func (t *Terminal) requestRender() {
	if t.RequestRender != nil {
		t.RequestRender()
	}
}

func (t *Terminal) processSequence(mr MeasuredRune, reader *Reader) {
	if mr.Rune == 0x1b {
		t.handleANSI(reader)
	} else {
		t.processRunes(mr)
	}
}

func (t *Terminal) Process(data []byte) {
	reader := NewReader(data)
	for {
		if mr := reader.ReadRune(); mr.Empty() {
			return
		} else {
			t.processSequence(mr, reader)
		}
	}
}

func (t *Terminal) processRunes(runes ...MeasuredRune) bool {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, r := range runes {
		switch r.Rune {
		case 0x05: //enq
			continue
		case 0x07: //bell
			continue
		case 0x8: //backspace
			t.activeBuffer.backspace()
		case 0x9: //tab
			t.activeBuffer.tab()
		case 0xa, 0xc: //newLine/form feed
			t.activeBuffer.newLine()
		case 0xb: //vertical tab
			t.activeBuffer.verticalTab()
		case 0xd: //carriageReturn
			t.activeBuffer.carriageReturn()
			t.requestRender()
		case 0xe: //shiftOut
			t.activeBuffer.currentCharset = 1
		case 0xf: //shiftIn
			t.activeBuffer.currentCharset = 0
		default:
			if r.Rune < 0x20 {
				// handle any other control chars here?
				continue
			}

			t.activeBuffer.write(t.translateRune(r))
		}
	}

	return false
}

func (t *Terminal) translateRune(b MeasuredRune) MeasuredRune {
	table := t.activeBuffer.charsets[t.activeBuffer.currentCharset]
	if table == nil {
		return b
	}
	chr, ok := (*table)[b.Rune]
	if ok {
		return MeasuredRune{Rune: chr, Width: 1}
	}
	return b
}

func (t *Terminal) switchBuffer(index uint8) {
	var carrySize bool
	var w, h uint16
	if t.activeBuffer != nil {
		w, h = t.activeBuffer.viewWidth, t.activeBuffer.viewHeight
		carrySize = true
	}
	t.activeBuffer = t.buffers[index]
	if carrySize {
		t.activeBuffer.resizeView(w, h)
	}
}

func (t *Terminal) GetMouseMode() MouseMode {
	return t.mouseMode
}

func (t *Terminal) GetMouseExtMode() MouseExtMode {
	return t.mouseExtMode
}

func (t *Terminal) GetActiveBuffer() *Buffer {
	return t.activeBuffer
}

func (t *Terminal) useMainBuffer() {
	t.switchBuffer(MainBuffer)
}

func (t *Terminal) useAltBuffer() {
	t.switchBuffer(AltBuffer)
}

func (t *Terminal) Lock() {
	t.mu.Lock()
}

func (t *Terminal) Unlock() {
	t.mu.Unlock()
}
