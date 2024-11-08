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
	Content string `json:"content"`
	Author  string `json:"author"`
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
				Content: "Thanks!",
				Author:  "John Smith",
			},
		},
	}

	output, err := json.MarshalIndent(&post, "", "\t\t")
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}
	err = os.WriteFile("post.json", output, 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
}
