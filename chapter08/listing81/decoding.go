package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Post struct {
	Id       int       `json:"id"`
	Content  string    `json:"content"`
	Author   Author    `json:"author"`
	Comments []Comment `json:"comments"`
}

type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func decode(filename string) (post Post, err error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	decoder := json.NewDecoder(jsonFile)
	err = decoder.Decode(&post)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	return
}

func unmarshal(filename string) (post Post, err error) {
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening JSON file:", err)
		return
	}
	defer func(jsonFile *os.File) {
		_ = jsonFile.Close()
	}(jsonFile)

	jsonData, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading JSON data:", err)
		return
	}
	_ = json.Unmarshal(jsonData, &post)
	return
}

func main() {
	post, err := decode("post.json")
	if err != nil {
		return
	}
	fmt.Println(post)
}
