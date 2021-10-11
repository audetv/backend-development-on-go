package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFileServer(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	dirToServe := http.Dir("upload")
	handler := http.FileServer(dirToServe)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expect %v", status, http.StatusOK)
	}

	expected := `<pre><a href="34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt">34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt</a><a href="3ad59372-cc02-4b0a-8e80-1f7a0e6a541e.testfile1.txt">3ad59372-cc02-4b0a-8e80-1f7a0e6a541e.testfile1.txt</a></pre>`
	if strings.ReplaceAll(rr.Body.String(), "\n", "") != expected {
		t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body, expected)
	}
}

// Тестируем фильтр ext=jpg, который должен вернуть пустой массив
func TestFilesHandler_ExtEmptyResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/files?ext=jpg", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	dirToServe := http.Dir("upload")
	handler := FilesHandler{
		UploadDir: string(dirToServe),
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expect %v", status, http.StatusOK)
	}

	expected := "[]"
	if strings.ReplaceAll(rr.Body.String(), "\n", "") != expected {
		t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body, expected)
	}
}

// Тестируем фильтр файлов ext=txt который должен вернуть массив json из 2-х объектов
func TestFilesHandler_ExtWithJsonResponse(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/files?ext=txt", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	dirToServe := http.Dir("upload")
	handler := FilesHandler{
		UploadDir: string(dirToServe),
	}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expect %v", status, http.StatusOK)
	}

	expected := `[{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt","ext":".txt","size":36},{"name":"3ad59372-cc02-4b0a-8e80-1f7a0e6a541e.testfile1.txt","ext":".txt","size":0}]`
	if strings.ReplaceAll(rr.Body.String(), "\n", "") != expected {
		t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body, expected)
	}
}
