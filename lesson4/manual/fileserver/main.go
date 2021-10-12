package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
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

	files, err := ioutil.ReadDir(f.UploadDir)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)

	comma := true

	fmt.Fprintf(w, "[")
	defer fmt.Fprintln(w, "]")

	q := r.URL.Query().Get("ext")

	for _, file := range files {
		if !findByQuery(q, file) {
			continue
		}
		if comma {
			comma = false
		} else {
			fmt.Fprintf(w, ",")
		}
		_ = enc.Encode(&File{
			Name: file.Name(),
			Ext:  filepath.Ext(file.Name()),
			Size: file.Size(),
		})
	}

	w.(http.Flusher).Flush()
}

func findByQuery(q string, file fs.FileInfo) bool {
	if q != "" {
		return strings.Contains(filepath.Ext(file.Name()), q)
	}
	return true
}

func main() {
	dirToServe := http.Dir("./upload")

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
