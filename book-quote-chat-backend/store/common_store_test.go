package store

import (
	"book-quote-chat-backend/model"
	"os"
	"testing"
)

func TestCommentStore(t *testing.T) {
	_ = os.Remove("test_comments.json")
	commentFile = "test_comments.json"
	c := model.Comment{
		ID:         "c1",
		UserID:     "u1",
		UserName:   "测试用户",
		Avatar:     "",
		TargetType: "quote",
		TargetID:   "q1",
		Content:    "棒极了！",
		ParentID:   "",
		Created:    12345,
	}
	if err := SaveComments([]model.Comment{c}); err != nil {
		t.Fatal(err)
	}
	list, err := LoadComments()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 || list[0].ID != "c1" {
		t.Fatal("评论存取失败")
	}
	_ = os.Remove("test_comments.json")
}
