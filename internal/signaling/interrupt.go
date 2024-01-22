package signaling

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type signalCallback = func(ctx context.Context)

const signalTimeout = 10

func HandleSignals(ctx context.Context, callback signalCallback) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// blocks until context is cancelled or signal received .
	select {
	case <-ctx.Done():
		close(sigs)
		callback(ctx)
	case <-sigs:
		close(sigs)
		timeoutContext, cancel := context.WithTimeout(ctx, signalTimeout*time.Second)
		defer cancel()
		callback(timeoutContext)
	}
}
