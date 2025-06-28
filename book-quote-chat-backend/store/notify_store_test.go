package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestNotifyStore(t *testing.T) {
	_ = os.Remove("test_notifies.json")
	notifyFile = "test_notifies.json"
	n := model.Notify{ID: "n1", UserID: "u1", Type: "like", Title: "有新点赞", Content: "你的金句被点赞", FromID: "u2", Read: false, Created: 12345}
	if err := SaveNotifies([]model.Notify{n}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadNotifies()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "n1" {
		t.Fatal("通知存取失败")
	}
	_ = os.Remove("test_notifies.json")
}
