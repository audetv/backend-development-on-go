package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

type FilesHandler struct {
	UploadDir string
}

func (f *FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(f.UploadDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		extension := filepath.Ext(file.Name())
		fmt.Fprintf(w, "name: %v, ext: %v, size: %v\r\n", file.Name(), extension, file.Size())
	}
}

func main() {
	dirToServe := http.Dir("./upload")

	filesHandler := &FilesHandler{
		UploadDir: "upload",
	}

	fs := &http.Server{
		Addr:         "localhost:3002",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	http.Handle("/", http.FileServer(dirToServe))
	http.Handle("/files", filesHandler)

	err := fs.ListenAndServe()
	if err != nil {
		log.Println(err)
		return
	}
}
