package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"github.com/google/uuid"
	"time"
)

func AddChatMsg(user, avatar, text, roomId string) (model.ChatMsg, error) {
	msgs, err := store.LoadChats(roomId)
	if err != nil {
		return model.ChatMsg{}, err
	}
	m := model.ChatMsg{
		ID:      uuid.New().String(),
		User:    user,
		Avatar:  avatar,
		Text:    text,
		Created: time.Now().Unix(),
		RoomID:  roomId,
	}
	msgs = append(msgs, m)
	err = store.SaveChats(msgs)
	return m, err
}

func GetChatMsgs(roomId string) ([]model.ChatMsg, error) {
	return store.LoadChats(roomId)
}
