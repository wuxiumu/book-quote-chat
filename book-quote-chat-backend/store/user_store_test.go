package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestUserStore(t *testing.T) {
	_ = os.Remove("test_users.json")
	userFile = "test_users.json"
	u := model.User{ID: "u1", Name: "test", Password: "123", Email: "t@x.com", Avatar: "", Created: 12345}
	if err := SaveUsers([]model.User{u}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadUsers()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "u1" {
		t.Fatal("保存或读取失败")
	}
	_ = os.Remove("test_users.json")
}
