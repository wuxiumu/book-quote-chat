// service/like_service_test.go
package service

import (
	"book-quote-chat-backend/model"
	"fmt"
	"os"
	"testing"
)

func TestLikeStoreBatchAndPaged(t *testing.T) {
	_ = os.Remove("test_likes.json")
	likeFile = "test_likes.json"
	// mock 100 条
	likes := []model.Like{}
	for i := 0; i < 100; i++ {
		likes = append(likes, model.Like{
			ID:         fmt.Sprintf("id%d", i),
			UserID:     fmt.Sprintf("u%d", i%10),
			TargetType: "quote",
			TargetID:   fmt.Sprintf("q%d", i%5),
			Created:    int64(10000 + i),
		})
	}
	if err := SaveLikes(likes); err != nil {
		t.Fatal(err)
	}
	// 分页
	list, total, err := service.GetLikesPaged("quote", "q1", 0, 10)
	if err != nil {
		t.Fatal(err)
	}
	if total == 0 || len(list) == 0 {
		t.Fatal("分页异常")
	}
	_ = os.Remove("test_likes.json")
}
