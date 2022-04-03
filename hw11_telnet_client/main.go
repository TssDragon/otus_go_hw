package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

var timeout time.Duration

const DefaultTimeout = 10

func main() {
	flag.DurationVar(&timeout, "timeout", DefaultTimeout*time.Second, "connection timeout")
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

	ctx, cancelFunc := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancelFunc()

	go writer(telnetClient, cancelFunc)
	go reader(telnetClient, cancelFunc)

	<-ctx.Done()
	cancelFunc()
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
