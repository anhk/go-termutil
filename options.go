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
