package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type filesHandler struct {
	UploadDir string
}

type File struct {
	Name string `json:"name"`
	Ext  string `json:"ext"`
	Size int64  `json:"size"`
}

func (f File) findByQuery(q string) bool {
	if q != "" {
		return strings.Contains(f.Name, q)
	}
	return true
}

type Files []File

// Обычный handler без flusher
func (f *filesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	q := r.URL.Query().Get("ext")

	var filesJson Files

	for _, file := range files {
		fileJson := &File{
			Name: file.Name(),
			Ext:  filepath.Ext(file.Name()),
			Size: file.Size(),
		}
		if !fileJson.findByQuery(q) {
			continue
		}
		filesJson = append(filesJson, *fileJson)
	}
	err = enc.Encode(filesJson)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	dirToServe := http.Dir("./upload")

	filesHandler := &filesHandler{
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
