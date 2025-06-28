package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"github.com/google/uuid"
	"time"
)

func AddFriend(name, avatar, group, remark string) (model.Friend, error) {
	friends, err := store.LoadFriends()
	if err != nil {
		return model.Friend{}, err
	}
	f := model.Friend{
		ID:      uuid.New().String(),
		Name:    name,
		Avatar:  avatar,
		Group:   group,
		Remark:  remark,
		Created: time.Now().Unix(),
	}
	friends = append(friends, f)
	err = store.SaveFriends(friends)
	return f, err
}

func GetFriends() ([]model.Friend, error) {
	return store.LoadFriends()
}
