package db_test

import (
	"go-nuxt-blogs/models"
	"testing"

	"github.com/fzzp/gotk"
	_ "github.com/mattn/go-sqlite3"
)

func createUser(t *testing.T) int64 {
	user := models.User{
		ID:       int64(gotk.RandomInt(100, 10000)),
		Email:    gotk.RandomString(6) + "@qq.com",
		Username: gotk.RandomString(4),
		Avatar:   "http://" + gotk.RandomString(8),
		Password: gotk.RandomString(6),
	}

	if _, err := user.Hash(); err != nil {
		t.Error(err)
	}
	id, err := testRepo.Users.Create(&user)
	if err != nil {
		t.Error(err)
	}
	if id <= 0 {
		t.Errorf("0 <= %d ", id)
	}

	return id
}

func TestAddUser(t *testing.T) {
	createUser(t)
}

func TestGetParamByKey(t *testing.T) {
	id := createUser(t)
	var data = map[string]interface{}{"id": id}
	user, err := testRepo.Users.GetByUnique(data)
	if err != nil {
		t.Error(err)
	}
	if user.ID != id {
		t.Errorf("%d != %d", user.ID, id)
	}
}

func TestUpdateUser(t *testing.T) {
	id := createUser(t)
	var data = map[string]interface{}{"id": id}
	user, err := testRepo.Users.GetByUnique(data)
	if err != nil {
		t.Error(err)
	}
	if user.ID != id {
		t.Errorf("%d != %d", user.ID, id)
	}
	userName := gotk.RandomString(4)
	avatar := gotk.RandomString(10)

	user.Username = userName
	user.Avatar = avatar

	err = testRepo.Users.Update(user)
	if err != nil {
		t.Error(err)
	}

	user2, err := testRepo.Users.GetByUnique(data)
	if err != nil {
		t.Error(err)
	}

	if user2.Username != userName {
		t.Error("uerName no update")
	}
	if user2.Avatar != avatar {
		t.Error("avatar no update")
	}
}
