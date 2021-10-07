package main

import (
	"net/http"
	"net/http/httptest"
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
