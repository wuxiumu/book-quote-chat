package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	friendFile = "data/friends.json"
	friendMu   sync.Mutex
)

func LoadFriends() ([]model.Friend, error) {
	friendMu.Lock()
	defer friendMu.Unlock()
	b, err := os.ReadFile(friendFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Friend{}, nil
	}
	var friends []model.Friend
	err = json.Unmarshal(b, &friends)
	return friends, err
}

func SaveFriends(friends []model.Friend) error {
	friendMu.Lock()
	defer friendMu.Unlock()
	b, err := json.MarshalIndent(friends, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(friendFile, b, 0644)
}
