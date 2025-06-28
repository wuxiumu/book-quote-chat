package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestLikeStore(t *testing.T) {
	_ = os.Remove("test_likes.json")
	likeFile = "test_likes.json"

	l := model.Like{ID: "lid1", UserID: "u1", TargetType: "quote", TargetID: "q1", Created: 1234567}
	if err := SaveLikes([]model.Like{l}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadLikes("quote", "q1")
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "lid1" {
		t.Fatal("not saved or loaded")
	}
	_ = os.Remove("test_likes.json")
}
