package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	cases := []struct {
		dirToServe string
		url        string
		expected   string
		returnErr  bool
		name       string
	}{
		{
			dirToServe: "testdata/upload",
			url:        "/files",
			expected:   `[{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt","ext":".txt","size":36},{"name":"3ad59372-cc02-4b0a-8e80-1f7a0e6a541e.testfile1.txt","ext":".txt","size":0},{"name":"8848ae1f-a31c-42b8-b512-4ec57906f065.testfile2","ext":".testfile2","size":0}]`,
			name:       "ListOfFiles",
		},
		{
			dirToServe: "testdata/upload",
			url:        "/files?ext=jpg",
			expected:   "[]",
			name:       "NoMatchesFound",
		},
		{
			dirToServe: "testdata/upload",
			url:        "/files?ext=txt",
			expected:   `[{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt","ext":".txt","size":36},{"name":"3ad59372-cc02-4b0a-8e80-1f7a0e6a541e.testfile1.txt","ext":".txt","size":0}]`,
			name:       "MatchesFound",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tc.url, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			dirToServe := http.Dir(tc.dirToServe)
			handler := FilesHandler{
				UploadDir: string(dirToServe),
			}
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != http.StatusOK {
				t.Errorf("Handler returned wrong status code: got %v expect %v", status, http.StatusOK)
			}

			expected := tc.expected
			if strings.ReplaceAll(rr.Body.String(), "\n", "") != expected {
				t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body, expected)
			}
		})
	}
}
