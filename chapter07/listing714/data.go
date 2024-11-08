package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var Db *sql.DB

func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=postgres password=postgres sslmode=disable")
	if err != nil {
		panic(err)
	}
}

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

func retrieve(id int) (post Post, err error) {
	post = Post{}
	_ = Db.QueryRow("select id, content, author from posts where id=$1", id).Scan(&post.Id, &post.Content, &post.Author)
	return
}

func (post *Post) create(id int) (err error) {
	statement := "insert into posts (content, author) values ($1,$2) returning id"
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	err = stmt.QueryRow(post.Content, post.Author).Scan(&post.Id)
	return
}

func (post *Post) update() (err error) {
	_, err = Db.Exec("update posts set content=$2, author=$3 where id=$1", post.Id, post.Content, post.Author)
	return
}

func (post *Post) delete() (err error) {
	_, err = Db.Exec("delete from posts where id=$1", post.Id)
	return
}
