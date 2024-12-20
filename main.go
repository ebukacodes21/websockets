package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"websockets/servers"
	"websockets/socket"

	"golang.org/x/sync/errgroup"
)

var signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), signals...)
	defer stop()

	group, ctx := errgroup.WithContext(ctx)
	manager := socket.NewSocketManager()

	servers.RunWebsocket(ctx, group, manager)

	err := group.Wait()
	if err != nil {
		log.Fatal(err)
	}
}
