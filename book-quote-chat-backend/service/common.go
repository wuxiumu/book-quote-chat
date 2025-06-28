package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"book-quote-chat-backend/wsutil"
	"errors"
	"github.com/google/uuid"
	"time"
)

// 新增评论
func AddComment(userId, userName, avatar, targetType, targetId, content, parentId string) (model.Comment, error) {
	if content == "" {
		return model.Comment{}, errors.New("评论内容不能为空")
	}
	comments, err := store.LoadComments()
	if err != nil {
		return model.Comment{}, err
	}
	comment := model.Comment{
		ID:         uuid.New().String(),
		UserID:     userId,
		UserName:   userName,
		Avatar:     avatar,
		TargetType: targetType,
		TargetID:   targetId,
		Content:    content,
		ParentID:   parentId,
		Created:    time.Now().Unix(),
	}
	comments = append(comments, comment)
	if err := store.SaveComments(comments); err != nil {
		return model.Comment{}, err
	}
	// 发通知
	targetOwnerId := GetTargetOwnerId(targetType, targetId)
	if targetOwnerId != "" && targetOwnerId != userId {
		_, _ = SendNotify(targetOwnerId, "comment", "新评论", "你的内容收到新评论", userId)
		wsutil.PushNotifyWS(targetOwnerId, "新评论", "你的内容收到新评论")
	}
	return comment, nil
}

// 获取某目标下评论，支持分页
func GetComments(targetType, targetId string, offset, limit int) ([]model.Comment, int, error) {
	comments, err := store.LoadComments()
	if err != nil {
		return nil, 0, err
	}
	list := make([]model.Comment, 0)
	for _, c := range comments {
		if c.TargetType == targetType && c.TargetID == targetId {
			list = append(list, c)
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

// 获取单个评论
func GetCommentByID(id string) (model.Comment, error) {
	comments, err := store.LoadComments()
	if err != nil {
		return model.Comment{}, err
	}
	for _, c := range comments {
		if c.ID == id {
			return c, nil
		}
	}
	return model.Comment{}, errors.New("评论不存在")
}
