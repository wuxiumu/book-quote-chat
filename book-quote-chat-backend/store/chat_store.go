package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	chatFile = "data/chats.json"
	chatMu   sync.Mutex
)

// Load all messages (optionally by room)
func LoadChats(roomId string) ([]model.ChatMsg, error) {
	chatMu.Lock()
	defer chatMu.Unlock()
	b, err := os.ReadFile(chatFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.ChatMsg{}, nil
	}
	var msgs []model.ChatMsg
	err = json.Unmarshal(b, &msgs)
	if err != nil {
		return nil, err
	}
	// Filter by roomId
	if roomId == "" {
		return msgs, nil
	}
	out := make([]model.ChatMsg, 0, len(msgs))
	for _, m := range msgs {
		if m.RoomID == roomId {
			out = append(out, m)
		}
	}
	return out, nil
}

func SaveChats(msgs []model.ChatMsg) error {
	chatMu.Lock()
	defer chatMu.Unlock()
	b, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(chatFile, b, 0644)
}
