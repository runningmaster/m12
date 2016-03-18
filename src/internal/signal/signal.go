package signal

import (
	"os"
	"os/signal"
	"syscall"
)

// WaitToExit helps wait for the event to exit
func WaitToExit() <-chan os.Signal {
	return waitFor(syscall.SIGINT, syscall.SIGTERM)
}

func waitFor(sig ...os.Signal) <-chan os.Signal {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, sig...)
	return ch
}
