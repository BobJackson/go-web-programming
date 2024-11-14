package main

import (
	"database/sql"
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
	var err error
	db, err := sql.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr: ":127.0.0.1:8080",
	}
	http.HandleFunc("/post/", handleRequest(&Post{Db: db}))
	_ = server.ListenAndServe()
}

func handleRequest(t Text) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		switch r.Method {
		case "GET":
			err = handleGet(w, r, t)
		case "POST":
			err = handlePost(w, r, t)
		case "PUT":
			err = handlePut(w, r, t)
		case "DELETE":
			err = handleDelete(w, r, t)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func handleGet(w http.ResponseWriter, r *http.Request, post Text) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	err = post.fetch(id)
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

func handlePost(w http.ResponseWriter, r *http.Request, text Text) (err error) {
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

func handlePut(w http.ResponseWriter, r *http.Request, text Text) (err error) {
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

func handleDelete(w http.ResponseWriter, r *http.Request, text Text) (err error) {
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
