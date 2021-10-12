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
		name       string
	}{
		{
			dirToServe: "testdata/upload",
			url:        "/files",
			expected:   `[{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.jpeg","ext":".jpeg","size":36},{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.txt","ext":".txt","size":36},{"name":"8848ae1f-a31c-42b8-b512-4ec57906f065.testfile2.png","ext":".png","size":0}]`,
			name:       "ListOfFiles",
		},
		{
			dirToServe: "testdata/upload",
			url:        "/files?ext=gif",
			expected:   "[]",
			name:       "NoMatchesFound",
		},
		{
			dirToServe: "testdata/upload",
			url:        "/files?ext=png",
			expected:   `[{"name":"8848ae1f-a31c-42b8-b512-4ec57906f065.testfile2.png","ext":".png","size":0}]`,
			name:       "PngMatchesFound",
		},
		{
			dirToServe: "testdata/upload",
			url:        "/files?ext=jpeg",
			expected:   `[{"name":"34cf8110-d6f7-4c72-8126-e2b33917f6ae.testfile.jpeg","ext":".jpeg","size":36}]`,
			name:       "JpegMatchesFound",
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
