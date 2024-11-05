package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

func store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}
}

func load(data interface{}, filename string) {
	raw, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	decoder := gob.NewDecoder(buffer)
	err = decoder.Decode(data)
	if err != nil {
		panic(err)
	}
}

type Post struct {
	Id      int
	Content string
	Author  string
}

func main() {
	post := Post{Id: 1, Content: "Hello World!", Author: "John Doe"}
	store(post, "post1")

	var postRead Post
	load(&postRead, "post1")
	fmt.Println(postRead)
}
