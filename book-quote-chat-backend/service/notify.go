package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"errors"
	"github.com/google/uuid"
	"time"
)

// 发送通知
func SendNotify(userId, typ, title, content, fromId string) (model.Notify, error) {
	notifies, err := store.LoadNotifies()
	if err != nil {
		return model.Notify{}, err
	}
	n := model.Notify{
		ID:      uuid.New().String(),
		UserID:  userId,
		Type:    typ,
		Title:   title,
		Content: content,
		FromID:  fromId,
		Read:    false,
		Created: time.Now().Unix(),
	}
	notifies = append(notifies, n)
	if err := store.SaveNotifies(notifies); err != nil {
		return model.Notify{}, err
	}
	return n, nil
}

// 查询用户通知，分页
func GetNotifies(userId string, offset, limit int) ([]model.Notify, int, error) {
	notifies, err := store.LoadNotifies()
	if err != nil {
		return nil, 0, err
	}
	var list []model.Notify
	for _, n := range notifies {
		if n.UserID == userId {
			list = append(list, n)
		}
	}
	total := len(list)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return list[offset:end], total, nil
}

// 标记已读
func MarkNotifyRead(userId, notifyId string) error {
	notifies, err := store.LoadNotifies()
	if err != nil {
		return err
	}
	changed := false
	for i, n := range notifies {
		if n.UserID == userId && n.ID == notifyId {
			notifies[i].Read = true
			changed = true
		}
	}
	if changed {
		return store.SaveNotifies(notifies)
	}
	return errors.New("通知不存在")
}

// 批量已读
func MarkNotifyBatchRead(userId string, notifyIds []string) error {
	notifies, err := store.LoadNotifies()
	if err != nil {
		return err
	}
	changed := false
	for i, n := range notifies {
		if n.UserID == userId {
			for _, id := range notifyIds {
				if n.ID == id {
					notifies[i].Read = true
					changed = true
				}
			}
		}
	}
	if changed {
		return store.SaveNotifies(notifies)
	}
	return nil
}
