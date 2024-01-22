package main

import (
	"context"
	"github.com/labstack/gommon/log"
	"github.com/sshilko/go-signer/internal/signaling"
	"sync"
	"time"
)

func main() {
	logger := log.New("worker")
	logger.SetHeader("${level} ")

	ctx, ctxClose := context.WithCancel(context.Background())

	var wg sync.WaitGroup

	wg.Add(1)
	t := time.NewTicker(1 * time.Second)
	go func() {
		defer wg.Done()
		defer t.Stop()
		logger.Info("Starting timed job")
		for {
			select {
			case <-ctx.Done():
				logger.Info("Stopped timed job")
				return
			case <-t.C:
				logger.Info(time.Now().UTC().String())
			}
		}
	}()

	wg.Add(1)
	go signaling.HandleSignals(ctx, func(c context.Context) {
		defer wg.Done() // indicate Done for WG no matter what
		logger.Info("Stop command received, finishing all jobs")
		ctxClose()
		logger.Info("Context closed")

		// TODO: stop worker job here
	})

	// TODO: start worker job listener here via GRPC listener

	wg.Wait()
	logger.Info("All jobs stopped")
}
