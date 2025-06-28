package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	userFile = "data/users.json"
	userMu   sync.Mutex
)

func LoadUsers() ([]model.User, error) {
	userMu.Lock()
	defer userMu.Unlock()
	b, err := os.ReadFile(userFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.User{}, nil
	}
	var users []model.User
	err = json.Unmarshal(b, &users)
	return users, err
}

func SaveUsers(users []model.User) error {
	userMu.Lock()
	defer userMu.Unlock()
	b, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(userFile, b, 0644)
}
