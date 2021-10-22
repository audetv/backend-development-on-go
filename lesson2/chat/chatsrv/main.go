package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan<- string

type chatChannels struct {
	entering chan client
	leaving  chan client
	messages chan string
}

func main() {
	// In order to pass the gochecknoglobals linter check
	chatChan := &chatChannels{
		entering: make(chan client),
		leaving:  make(chan client),
		messages: make(chan string),
	}

	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(chatChan)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn, chatChan)
	}
}

func handleConn(conn net.Conn, chatChan *chatChannels) {
	ch := make(chan string)
	go clientWriter(conn, ch)

	remoteAddr := conn.RemoteAddr().String()
	ch <- "You remote addr is " + remoteAddr + ". Please enter you nickname:"

	nickname := awaitInputNickname(conn)

	chatChan.entering <- ch
	chatChan.messages <- "Nickname: " + nickname + " has arrived"
	ch <- "Welcome, " + nickname + "!"

	log.Println(remoteAddr + "@" + nickname + " has arrived")

	input := bufio.NewScanner(conn)
	for input.Scan() {
		chatChan.messages <- nickname + ": " + input.Text()
	}
	chatChan.leaving <- ch
	chatChan.messages <- nickname + " has left"
	conn.Close()
}

func awaitInputNickname(conn net.Conn) string {
	var nickname string
	// Created the scanner and await input
	nickScan := bufio.NewScanner(conn)
	count := 0
	for nickScan.Scan() {
		count++
		nickname = nickScan.Text()
		if nickname != "" && count > 0 {
			break
		}
	}
	return nickname
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func broadcaster(chatChan *chatChannels) {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-chatChan.messages:
			for cli := range clients {
				cli <- msg
			}

		case cli := <-chatChan.entering:
			clients[cli] = true
		case cli := <-chatChan.leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}
