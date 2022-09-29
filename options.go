package termutil

import (
	"os"
)

type Option func(t *Terminal)

func WithLogFile(path string) Option {
	return func(t *Terminal) {
		if path == "-" {
			t.logFile = os.Stdout
			return
		}
		t.logFile, _ = os.Create(path)
	}
}

func WithRequestRender(f func()) Option {
	return func(t *Terminal) {
		t.RequestRender = f
	}
}
