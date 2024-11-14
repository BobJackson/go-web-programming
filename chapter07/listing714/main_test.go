package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var mux *http.ServeMux
var writer *httptest.ResponseRecorder

func TestMain(m *testing.M) {
	setUp()
	code := m.Run()
	os.Exit(code)
}

func setUp() {
	mux = http.NewServeMux()
	mux.HandleFunc("/post/", handleRequest)
	writer = httptest.NewRecorder()
}

func TestHandleRequest(t *testing.T) {
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
	}
	var post Post
	_ = json.Unmarshal(writer.Body.Bytes(), &post)

	if post.Id != 1 {
		t.Errorf("Expected post id %d, got %d", 1, post.Id)
	}
}

func TestHandlePut(t *testing.T) {
	body := strings.NewReader(`{"content": "Updated post", "author": "Sau Sheong"}`)
	request, _ := http.NewRequest("PUT", "/post/1", body)
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
	}
}
