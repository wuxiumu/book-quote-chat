package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestMessageStore(t *testing.T) {
	_ = os.Remove("test_messages.json")
	messageFile = "test_messages.json"

	m := model.Message{ID: "mid1", FromUser: "Aric", ToUser: "Lina", Avatar: "a.png", Text: "hi", Created: 1234567}
	if err := SaveMessages([]model.Message{m}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadMessages("Aric", "Lina")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "mid1" {
		t.Fatal("not saved or loaded")
	}
	_ = os.Remove("test_messages.json")
}
