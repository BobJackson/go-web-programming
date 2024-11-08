package main

import (
	"encoding/json"
	"fmt"
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
	Author  string `json:"author"`
	Content string `json:"content"`
}

func main() {
	post := Post{
		Id:      1,
		Content: "Hello, world!",
		Author: Author{
			Id:   1,
			Name: "John Doe",
		},
		Comments: []Comment{
			{
				Id:      1,
				Content: "Have a great day!",
				Author:  "Jane Doe",
			},
			{
				Id:      2,
				Content: "How are you today?",
				Author:  "John Smith",
			},
		},
	}

	jsonFile, err := os.Create("post.json")
	if err != nil {
		fmt.Println("Error creating JSON file:", err)
		return
	}
	encoder := json.NewEncoder(jsonFile)
	err = encoder.Encode(&post)
	if err != nil {
		fmt.Println("Error encoding JSON to file:", err)
		return
	}
}
