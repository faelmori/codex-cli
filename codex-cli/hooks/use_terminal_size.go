package hooks

import (
	"os"
	"os/signal"
	"syscall"
)

const TERMINAL_PADDING_X = 8

type TerminalSize struct {
	Columns int
	Rows    int
}

func GetTerminalSize() TerminalSize {
	return TerminalSize{
		Columns: (os.Stdout().Columns() || 60) - TERMINAL_PADDING_X,
		Rows:    os.Stdout().Rows() || 20,
	}
}

func WatchTerminalSize(callback func(size TerminalSize)) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGWINCH)

	go func() {
		for range sig {
			callback(GetTerminalSize())
		}
	}()
}
