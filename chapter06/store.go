package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

var PostById map[int]*Post
var PostByAuthor map[string][]*Post

func store(post Post) {
	PostById[post.Id] = &post
	PostByAuthor[post.Author] = append(PostByAuthor[post.Author], &post)
}

func main() {

	//storeInMemory()

	//storeToFile()

	storeToCsv()
}

func storeToCsv() {
	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}
	defer func(csvFile *os.File) {
		err := csvFile.Close()
		if err != nil {
			panic(err)
		}
	}(csvFile)

	allPosts := []Post{
		Post{Id: 1, Content: "Hello World", Author: "John"},
		Post{Id: 2, Content: "Bonjour Monde", Author: "Pierre"},
		Post{Id: 3, Content: "Hola Mundo", Author: "Pedro"},
		Post{Id: 4, Content: "Ciao mondo", Author: "Marco"},
	}

	write := csv.NewWriter(csvFile)

	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		err := write.Write(line)
		if err != nil {
			panic(err)
		}
	}
	write.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var posts []Post
	for _, item := range records {
		id, _ := strconv.Atoi(item[0])
		posts = append(posts, Post{Id: id, Content: item[1], Author: item[2]})
	}
	fmt.Println(posts[0].Id)
	fmt.Println(posts[0].Content)
	fmt.Println(posts[0].Author)
}

func storeToFile() {
	data := []byte("Hello World!\n")
	err := os.WriteFile("data1", data, 0644)
	if err != nil {
		panic(err)
	}

	read1, _ := os.ReadFile("data1")
	fmt.Println(string(read1))

	file1, _ := os.Create("data2")
	defer func(file1 *os.File) {
		err := file1.Close()
		if err != nil {
			panic(err)
		}
	}(file1)

	bytes, _ := file1.Write(data)
	fmt.Printf("Wrote %d bytes to file\n", bytes)

	file2, _ := os.Open("data2")
	defer func(file2 *os.File) {
		err := file2.Close()
		if err != nil {
			panic(err)
		}
	}(file2)

	read2 := make([]byte, len(data))
	bytes, _ = file2.Read(read2)

	fmt.Printf("Read %d bytes from file\n", bytes)
	fmt.Println(string(read2))
}

func storeInMemory() {
	PostById = make(map[int]*Post)
	PostByAuthor = make(map[string][]*Post)

	post1 := Post{Id: 1, Content: "Hello World!", Author: "John Doe"}
	post2 := Post{Id: 2, Content: "Bonjour Monde!", Author: "Pierre"}
	post3 := Post{Id: 3, Content: "Hola Mundo!", Author: "Pedro"}
	post4 := Post{Id: 4, Content: "Greetings Earthlings!", Author: "Sau Sheong"}

	store(post1)
	store(post2)
	store(post3)
	store(post4)

	for _, post := range PostByAuthor["Sau Sheong"] {
		fmt.Println(post)
	}

	for _, post := range PostByAuthor["Pedro"] {
		fmt.Println(post)
	}
}
