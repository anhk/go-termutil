package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/anhk/go-termutil"
	"github.com/creack/pty"
	"golang.org/x/term"
)

func test() error {
	c := exec.Command("bash")

	// Start the command with a pty.
	ptmx, err := pty.Start(c)
	if err != nil {
		return err
	}
	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			if err := pty.InheritSize(os.Stdin, ptmx); err != nil {
				log.Printf("error resizing pty: %s", err)
			}
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Set stdin in raw mode.
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer func() { _ = term.Restore(int(os.Stdin.Fd()), oldState) }() // Best effort.

	go func() { _, _ = io.Copy(ptmx, os.Stdin) }()
	// _, _ = io.Copy(os.Stdout, ptmx)
	copy(os.Stdout, ptmx)
	return nil
}

var terminal *termutil.Terminal

func main() {
	terminal = termutil.New(termutil.WithLogFile("./termutil.log"))

	test()
}

func copy(w io.Writer, r io.Reader) {
	buff := make([]byte, 4096)
	for {
		nr, err := r.Read(buff)
		if err != nil {
			break
		}
		terminal.Write(buff[:nr])
		w.Write(buff[:nr])
	}
}
