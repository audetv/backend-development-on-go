package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandler_ServeHTTP(t *testing.T) {
	// Создаем запрос с указанием нашего хендлера. Так как мы тестируем GET-эндпоинт
	// то нам не нужно передавать тело, поэтому третьим аргументом передаем nil
	req, err := http.NewRequest(http.MethodGet, "/?name=Jhon", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
	// и используем его для получения ответа
	rr := httptest.NewRecorder()
	handler := &Handler{}

	// Наш хендлер соответствует интерфейсу http.Handler, а значит
	// мы можем использовать ServeHTTP и напрямую указать
	// Request и ResponseRecorder
	handler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expect %v", status, http.StatusOK)
	}

	// Проверяем тело ответа
	expected := `Parsed query-param with key "name": Jhon`
	if rr.Body.String() != expected {
		t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body.String(), expected)
	}
}

func TestUploadHandler_ServeHTTP(t *testing.T) {
	// открываем файл, который хотим отправить
	file, err := os.Open("testfile")
	if err != nil {
		log.Println(err)
	}
	defer file.Close()

	// действия, необходимые для того, чтобы засунуть файл в запрос
	// в качестве мультипарт-формы
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	// опять создаем запрос, теперь уже на /upload эндпоинт
	req, _ := http.NewRequest(http.MethodPost, "/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// создаем ResponseRecorder
	rr := httptest.NewRecorder()

	// создаем заглушку файлового сервера. Для прохождения тестов
	// нам достаточно чтобы он возвращал 200 статус
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "ok!")
	}))
	defer ts.Close()

	uploadHandler := &UploadHandler{
		UploadDir: "upload",
		// таким образом мы подменим адрес файлового сервера
		// и вместо реального, хэндлер будет стучаться на заглушку
		// которая всегда будет возвращать 200 статус, что нам и нужно
		HostAddr: ts.URL,
	}

	// опять же, вызываем ServeHTTP у тестируемого обработчика
	uploadHandler.ServeHTTP(rr, req)

	// Проверяем статус-код ответа
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v expext %v", status, http.StatusOK)
	}

	// Мы не будем тестировать создание UUID и добавление его к имени файла,
	// так как uuid каждый раз при старте теста будет другой
	// поэтому можно оставить условие без изменения: если в строке ответа не содержится имя файла testfile
	expected := "testfile"
	if !strings.Contains(rr.Body.String(), expected) {
		t.Errorf("Handler returned unexpected body: got %v expect %v", rr.Body.String(), expected)
	}
}
