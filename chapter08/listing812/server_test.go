package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type FakePost struct {
	Id      int
	Content string
	Author  string
}

func (post *FakePost) fetch(id int) (err error) {
	post.Id = id
	return
}

func (post *FakePost) create() (err error) {
	return
}

func (post *FakePost) update() (err error) {
	return
}

func (post *FakePost) delete() (err error) {
	return
}

func TestHandleGet(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/post/", handleRequest(&FakePost{}))

	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)

	if writer.Code != http.StatusOK {
		t.Error("Expected status code 200, got ", writer.Code)
	}
	var post Post
	_ = json.Unmarshal(writer.Body.Bytes(), &post)
	if post.Id != 1 {
		t.Error("Expected post id 1, got ", post.Id)
	}
}
