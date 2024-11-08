package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
)

type Post struct {
	XMLName  xml.Name  `xml:"post"`
	Id       string    `xml:"id,attr"`
	Author   string    `xml:"author"`
	Xml      string    `xml:",innerxml"`
	Comments []Comment `xml:"comments>comment"`
}

type Author struct {
	Id   string `xml:"id,attr"`
	Name string `xml:",chardata"`
}

type Comment struct {
	Id      string `xml:"user,attr"`
	Content string `xml:"content"`
	Author  Author `xml:"author"`
}

func main() {
	xmlFile, err := os.Open("post.xml") // post.xml 文件需要放在最外层，而不是和xml.go 文件同级
	if err != nil {
		fmt.Println("Error opening XML file:", err)
		return
	}
	defer func(xmlFile *os.File) {
		_ = xmlFile.Close()
	}(xmlFile)
	xmlData, err := io.ReadAll(xmlFile)
	if err != nil {
		fmt.Println("Error reading XML data:", err)
		return
	}

	var post Post
	_ = xml.Unmarshal(xmlData, &post)
	fmt.Println(post)
}
