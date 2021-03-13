package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anitta/eguchi-wedding-bot/pkg/di"
)

func main() {
	log.Println("Started eguchi-wedding-bot.")

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 10)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		defer close(sig)
		<-sig
		cancel()
	}()

	err := di.Start(ctx)
	if err != nil {
		panic(err)
	}
}
