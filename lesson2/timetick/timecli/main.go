package main

import (
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"time"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	dialer := net.Dialer{
		Timeout:   time.Second,
		KeepAlive: time.Minute,
	}
	conn, err := dialer.DialContext(ctx, "tcp", "127.127.127:9000")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(io.Copy(os.Stdout, conn))
}
