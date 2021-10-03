package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	go func() {
		_, err := io.Copy(os.Stdout, conn)
		if err != nil {
			return
		}
	}()
	_, err = io.Copy(conn, os.Stdin)
	if err != nil {
		log.Println(err)
	} // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
