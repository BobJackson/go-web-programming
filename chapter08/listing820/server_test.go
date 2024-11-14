package main

import (
	"encoding/json"
	. "gopkg.in/check.v1"
	"net/http"
	"net/http/httptest"
	"testing"
)

type PostTestSuite struct {
}

type FakePost struct {
	Id      int
	Content string
	Author  string
}

func (post *FakePost) fetch(id int) (err error) {
	post.Id = id
	return
}

func (post *FakePost) create() (err error) {
	return
}

func (post *FakePost) update() (err error) {
	return
}

func (post *FakePost) delete() (err error) {
	return
}

func init() {
	Suite(&PostTestSuite{})
}

func Test(t *testing.T) {
	TestingT(t)
}

func (s *PostTestSuite) TestHandleGet(c *C) {
	mux := http.NewServeMux()
	mux.HandleFunc("/post/", handleRequest(&FakePost{}))
	writer := httptest.NewRecorder()
	request, _ := http.NewRequest("GET", "/post/1", nil)
	mux.ServeHTTP(writer, request)

	c.Check(writer.Code, Equals, 200)
	var post Post
	_ = json.Unmarshal(writer.Body.Bytes(), &post)
	c.Check(post.Id, Equals, 1)
}
