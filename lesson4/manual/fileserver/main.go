package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	dirToServe := http.Dir("./upload")

	fs := &http.Server{
		Addr:         "localhost:3002",
		Handler:      http.FileServer(dirToServe),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	err := fs.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
