package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	commentFile = "data/comments.json"
	commentMu   sync.Mutex
)

func LoadComments() ([]model.Comment, error) {
	commentMu.Lock()
	defer commentMu.Unlock()
	b, err := os.ReadFile(commentFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Comment{}, nil
	}
	var comments []model.Comment
	err = json.Unmarshal(b, &comments)
	return comments, err
}

func SaveComments(comments []model.Comment) error {
	commentMu.Lock()
	defer commentMu.Unlock()
	b, err := json.MarshalIndent(comments, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(commentFile, b, 0644)
}
