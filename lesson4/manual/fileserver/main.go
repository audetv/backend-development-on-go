package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type FilesHandler struct {
	UploadDir string
}

type File struct {
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

func (f *FilesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "bad method", http.StatusMethodNotAllowed)
		return
	}

	q := r.URL.Query().Get("ext")
	log.Println(q)

	files, err := ioutil.ReadDir(f.UploadDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")

	enc := json.NewEncoder(w)

	comma := true
	fmt.Fprintf(w, "[")
	defer fmt.Fprintln(w, "]")

	for _, file := range files {
		f := &File{
			Name: file.Name(),
			Ext:  filepath.Ext(file.Name()),
			Size: file.Size(),
		}
		if comma {
			comma = false
		} else {
			fmt.Fprintf(w, ",")
		}
		if q != "" {
			if strings.Contains(filepath.Ext(file.Name()), q) {
				log.Println(strings.Contains(filepath.Ext(file.Name()), q))
				_ = enc.Encode(f)
				w.(http.Flusher).Flush()
			}
		} else {
			_ = enc.Encode(f)
			w.(http.Flusher).Flush()
		}
	}
}

func main() {
	dirToServe := http.Dir("/Windows/system32")

	filesHandler := &FilesHandler{
		UploadDir: string(dirToServe),
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
