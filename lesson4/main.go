package main

import (
	"fmt"
	"net/http"
	"time"
)

type helloHandler struct {
	subject string
}

func (h *helloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello, %s!", h.subject)
	switch r.Method {
	case http.MethodGet:
		fmt.Fprintln(w, "some GET-related logic")
	case http.MethodPost:
		fmt.Fprintln(w, "some POST-related logic")
	}
}

func main() {
	worldHandler := helloHandler{"World"}
	roomHandler := helloHandler{"Mark"}

	http.Handle("/world", &worldHandler)
	http.Handle("/room", &roomHandler)

	srv := &http.Server{
		Addr:         "localhost:80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		return
	}
}
