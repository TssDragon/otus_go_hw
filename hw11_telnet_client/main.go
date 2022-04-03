package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

const DefaultTimeout = 10

func init() {
	flag.DurationVar(&timeout, "timeout", DefaultTimeout*time.Second, "connection timeout")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		log.Fatal("Address & port must be type")
	}

	host := args[0]
	port := args[1]
	socket := net.JoinHostPort(host, port)

	telnetClient := NewTelnetClient(socket, timeout, os.Stdin, os.Stdout)
	if err := telnetClient.Connect(); err != nil {
		log.Fatalf("Error while open connection: %s", err)
	}
	defer telnetClient.Close()

	ctx, cancelFunc := context.WithCancel(context.Background())
	go writer(telnetClient, cancelFunc)
	go reader(telnetClient, cancelFunc)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	select {
	case <-signalChannel:
		cancelFunc()
		signal.Stop(signalChannel)
		return

	case <-ctx.Done():
		close(signalChannel)
		return
	}
}

func writer(client TelnetClient, cancelFunction context.CancelFunc) {
	if err := client.Send(); err != nil {
		cancelFunction()
	}
}

func reader(client TelnetClient, cancelFunction context.CancelFunc) {
	if err := client.Receive(); err != nil {
		cancelFunction()
	}
}
