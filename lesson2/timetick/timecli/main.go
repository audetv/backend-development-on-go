package main

import (
	"context"
	"fmt"
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
	defer func() {
		log.Println("Close conn. Exit")
		conn.Close()
	}()
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			buf := make([]byte, 256) // создаем буфер
			for {
				_, err = conn.Read(buf)
				if err == io.EOF {
					log.Println("Server close conn")
					return
				}
				_, err := io.WriteString(os.Stdout, fmt.Sprintf("Custom output! %s", string(buf)))
				if err != nil {
					log.Println("WriteString err", err)
					return
				} // выводим измененное сообщение сервера в консоль
			}
		}
	}()

	<-ctx.Done()
}
