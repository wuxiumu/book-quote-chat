package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	messageFile = "data/messages.json"
	messageMu   sync.Mutex
)

func LoadMessages(user1, user2 string) ([]model.Message, error) {
	messageMu.Lock()
	defer messageMu.Unlock()
	b, err := os.ReadFile(messageFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Message{}, nil
	}
	var msgs []model.Message
	if err := json.Unmarshal(b, &msgs); err != nil {
		return nil, err
	}
	// 双方的私聊消息
	out := make([]model.Message, 0, len(msgs))
	for _, m := range msgs {
		if (m.FromUser == user1 && m.ToUser == user2) || (m.FromUser == user2 && m.ToUser == user1) {
			out = append(out, m)
		}
	}
	return out, nil
}

func SaveMessages(msgs []model.Message) error {
	messageMu.Lock()
	defer messageMu.Unlock()
	b, err := json.MarshalIndent(msgs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(messageFile, b, 0644)
}
