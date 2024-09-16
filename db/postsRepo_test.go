package db_test

import (
	"encoding/json"
	"fmt"
	"go-nuxt-blogs/db"
	"go-nuxt-blogs/models"
	"testing"

	"github.com/fzzp/gotk"
)

func createPosts(t *testing.T) int64 {
	uid := createUser(t)
	p := models.Posts{
		ID:       int64(gotk.RandomInt(10, 100000)),
		Title:    gotk.RandomString(10),
		Content:  gotk.RandomString(20),
		AuthorID: uid,
	}

	pid, err := testRepo.Posts.Create(&p)
	if err != nil {
		t.Error(err)
	}
	if pid <= 0 {
		t.Errorf("pid: %d", pid)
	}

	return pid
}

func TestAddPosts(t *testing.T) {
	createPosts(t)
}

func TestGetPosts(t *testing.T) {
	pid := createPosts(t)
	posts, err := testRepo.Posts.Get(pid)
	if err != nil {
		t.Error(err)
	}
	if posts.ID != pid {
		t.Error("id != pid")
	}
	if posts.AttrID <= 0 {
		t.Error("posts.AttrID <= 0")
	}
	if posts.Title == "" {
		t.Error("title blank")
	}
	if posts.Content == "" {
		t.Error("content blank")
	}
	if posts.AuthorID <= 0 {
		t.Error("authorId blank")
	}
}

func TestUpdatePosts(t *testing.T) {
	pid := createPosts(t)
	posts, err := testRepo.Posts.Get(pid)
	if err != nil {
		t.Error(err)
	}
	content := gotk.RandomString(30)
	title := gotk.RandomString(10)
	posts.Content = content
	posts.Title = title

	if err := testRepo.Posts.Update(posts); err != nil {
		t.Error(err)
	}

	posts2, _ := testRepo.Posts.Get(pid)
	if posts2.Title != title {
		t.Error("title")
	}

	if posts2.Content != content {
		t.Error("content")
	}
}

func TestGetListPosts(t *testing.T) {
	createPosts(t)
	f := db.Filter{}
	list, _, err := testRepo.Posts.List(f)
	if err != nil {
		t.Error(err)
	}
	if len(list) <= 0 {
		t.Errorf("len = %d", len(list))
	}
	buf, _ := json.MarshalIndent(list, " ", "")
	fmt.Println(string(buf))
}
