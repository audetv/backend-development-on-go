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
	conn, err := dialer.DialContext(ctx, "tcp", "[::1]:9000")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			text, _ := io.Copy(os.Stdout, conn)
			if text == 0 {
				break
			}
			log.Println(text)
		}
	}()

	<-ctx.Done()
}
