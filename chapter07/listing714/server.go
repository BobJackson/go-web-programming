package main

import (
	"encoding/json"
	"net/http"
	"path"
	"strconv"
)

const (
	contentTypeHeader = "Content-Type"
	jsonContentType   = "application/json"
)

func main() {
	server := http.Server{
		Addr: ":127.0.0.1:8080",
	}
	http.HandleFunc("/post/", handleRequest)
	_ = server.ListenAndServe()
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		err = handleGet(w, r)
	case "POST":
		err = handlePost(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&post, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set(contentTypeHeader, jsonContentType)
	_, err = w.Write(output)
	return nil
}

func handlePost(w http.ResponseWriter, r *http.Request) (err error) {
	contentLength := r.ContentLength
	body := make([]byte, contentLength)
	_, _ = r.Body.Read(body)
	var post Post
	err = json.Unmarshal(body, &post)
	if err != nil {
		return
	}
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(http.StatusCreated)
	return
}

func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	contentLength := r.ContentLength
	body := make([]byte, contentLength)
	_, _ = r.Body.Read(body)

	err = json.Unmarshal(body, &post)
	if err != nil {
		return
	}
	err = post.update()
	if err != nil {
		return
	}
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(http.StatusNoContent)
	return
}

func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	err = post.delete()
	if err != nil {
		return
	}
	w.Header().Set(contentTypeHeader, jsonContentType)
	w.WriteHeader(http.StatusOK)
	return
}
