package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Post struct {
	Id         int
	Content    string
	AuthorName string `db:"author"`
}

var Db *sqlx.DB

func init() {
	var err error
	Db, err = sqlx.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func GetPost(id int) (post Post, err error) {
	post = Post{}
	err = Db.QueryRowx("select id, content, author from posts where id=$1", id).StructScan(&post)
	return
}

func (post *Post) Create() (err error) {
	err = Db.QueryRowx("insert into posts (content, author) values ($1, $2) returning id", post.Content, post.AuthorName).Scan(&post.Id)
	return
}

func main() {
	post := Post{Content: "Hello World", AuthorName: "Joey Bloggs"}
	_ = post.Create()

	readPost, err := GetPost(post.Id)
	if err != nil {
		panic(err)
	}
	fmt.Println(readPost)
}
