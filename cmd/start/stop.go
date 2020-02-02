package start

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func ContextWithInterrupt() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cancel()
	}()
	return ctx
}
