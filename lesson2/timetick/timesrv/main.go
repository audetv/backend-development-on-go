package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

type client chan<- string

type events struct {
	entering chan client
	leaving  chan client
	messages chan string
}

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	event := &events{
		entering: make(chan client),
		leaving:  make(chan client),
		messages: make(chan string),
	}

	cfg := net.ListenConfig{
		KeepAlive: time.Minute,
	}

	listener, err := cfg.Listen(ctx, "tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	wg := &sync.WaitGroup{}

	log.Println("Im started!")

	go broadcaster(event)
	go scanServerInput(event)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
			}
			conn, err := listener.Accept()
			if err != nil {
				log.Println(err)
			} else {
				wg.Add(1)
				createNewChannel(conn, event)
				go handleConn(ctx, conn, wg)
				// event.leaving <- ch
				// log.Println(remoteAddr + " has left")
				// conn.Close()
			}
		}
	}()

	<-ctx.Done()

	log.Println("done")
	// after ctx.Done, got accept tcp [::]:9000: use of closed network connection
	listener.Close()
	wg.Wait()
	log.Println("exit")
}

func createNewChannel(conn net.Conn, event *events) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	remoteAddr := conn.RemoteAddr().String()
	ch <- "Welcome, " + remoteAddr
	event.entering <- ch

	log.Println(remoteAddr + " has arrived")
}

func scanServerInput(event *events) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		event.messages <- "Server message: " + scanner.Text()
	}
}

func handleConn(ctx context.Context, conn net.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	defer conn.Close()

	tck := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case t := <-tck.C:
			fmt.Fprintf(conn, "now: %s\n", t)
		}
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster(events *events) {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-events.messages:
			for cli := range clients {
				cli <- msg
			}
		case cli := <-events.entering:
			clients[cli] = true
		case cli := <-events.leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
