package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"github.com/google/uuid"
	"time"
)

func AddMessage(from, to, avatar, text string) (model.Message, error) {
	msgs, err := store.LoadMessages(from, to)
	if err != nil {
		return model.Message{}, err
	}
	m := model.Message{
		ID:       uuid.New().String(),
		FromUser: from,
		ToUser:   to,
		Avatar:   avatar,
		Text:     text,
		Created:  time.Now().Unix(),
		Read:     false,
		Revoke:   false,
	}
	msgs = append(msgs, m)
	err = store.SaveMessages(msgs)
	return m, err
}

func GetMessages(user1, user2 string) ([]model.Message, error) {
	return store.LoadMessages(user1, user2)
}
