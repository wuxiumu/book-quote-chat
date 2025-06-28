package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	notifyFile = "data/notifies.json"
	notifyMu   sync.Mutex
)

func LoadNotifies() ([]model.Notify, error) {
	notifyMu.Lock()
	defer notifyMu.Unlock()
	b, err := os.ReadFile(notifyFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Notify{}, nil
	}
	var notifies []model.Notify
	err = json.Unmarshal(b, &notifies)
	return notifies, err
}

func SaveNotifies(notifies []model.Notify) error {
	notifyMu.Lock()
	defer notifyMu.Unlock()
	b, err := json.MarshalIndent(notifies, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(notifyFile, b, 0644)
}
