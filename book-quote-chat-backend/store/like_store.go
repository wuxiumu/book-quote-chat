package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	likeFile = "data/likes.json"
	likeMu   sync.Mutex
)

func LoadLikes(targetType, targetId string) ([]model.Like, error) {
	likeMu.Lock()
	defer likeMu.Unlock()
	b, err := os.ReadFile(likeFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Like{}, nil
	}
	var likes []model.Like
	if err := json.Unmarshal(b, &likes); err != nil {
		return nil, err
	}
	// 过滤目标
	if targetType == "" && targetId == "" {
		return likes, nil
	}
	out := make([]model.Like, 0, len(likes))
	for _, l := range likes {
		if l.TargetType == targetType && l.TargetID == targetId {
			out = append(out, l)
		}
	}
	return out, nil
}

func SaveLikes(likes []model.Like) error {
	likeMu.Lock()
	defer likeMu.Unlock()
	b, err := json.MarshalIndent(likes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(likeFile, b, 0644)
}
