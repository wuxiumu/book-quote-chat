package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestFriendStore(t *testing.T) {
	_ = os.Remove("test_friends.json")
	friendFile = "test_friends.json"

	f := model.Friend{ID: "fid1", Name: "Aric", Avatar: "a.png", Group: "同事", Remark: "大佬", Created: 1234567}
	if err := SaveFriends([]model.Friend{f}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadFriends()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "fid1" {
		t.Fatal("not saved or loaded")
	}
	_ = os.Remove("test_friends.json")
}
