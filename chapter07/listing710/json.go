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
	Id      string `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

func main() {
	jsonFile, err := os.Open("post.json")
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
	var post Post
	_ = json.Unmarshal(jsonData, &post)
	fmt.Println(post)
}
