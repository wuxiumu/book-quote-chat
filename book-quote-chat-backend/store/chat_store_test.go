package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestChatStore(t *testing.T) {
	_ = os.Remove("test_chats.json")
	chatFile = "test_chats.json"

	m := model.ChatMsg{ID: "mid1", User: "Aric", Avatar: "a.png", Text: "hi", Created: 1234567, RoomID: "main"}
	if err := SaveChats([]model.ChatMsg{m}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadChats("main")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "mid1" {
		t.Fatal("not saved or loaded")
	}
	_ = os.Remove("test_chats.json")
}
