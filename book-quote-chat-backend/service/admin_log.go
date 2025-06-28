package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"errors"
	"github.com/google/uuid"
	"time"
)

// 查询全部评论/金句（支持分页/过滤）
func ListComments(offset, limit int) ([]model.Comment, int, error) {
	all, err := store.LoadComments()
	if err != nil {
		return nil, 0, err
	}
	total := len(all)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return all[offset:end], total, nil
}

// 审核评论（通过/驳回/删除）
func AuditComment(adminId, commentId, action, detail string) error {
	comments, err := store.LoadComments()
	if err != nil {
		return err
	}
	found := false
	for i, c := range comments {
		if c.ID == commentId {
			if action == "delete" {
				comments = append(comments[:i], comments[i+1:]...)
			}
			// 可补充其它审核操作
			found = true
			break
		}
	}
	if !found {
		return errors.New("评论不存在")
	}
	if err := store.SaveComments(comments); err != nil {
		return err
	}
	_ = AddAdminLog(adminId, action, "comment", commentId, detail)
	return nil
}

// 管理日志
func AddAdminLog(adminId, action, target, targetId, detail string) error {
	logs, _ := store.LoadAdminLogs()
	log := model.AdminLog{
		ID:       uuid.New().String(),
		AdminID:  adminId,
		Action:   action,
		Target:   target,
		TargetID: targetId,
		Detail:   detail,
		Created:  time.Now().Unix(),
	}
	logs = append(logs, log)
	return store.SaveAdminLogs(logs)
}

func ListAdminLogs(offset, limit int) ([]model.AdminLog, int, error) {
	logs, err := store.LoadAdminLogs()
	if err != nil {
		return nil, 0, err
	}
	total := len(logs)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return logs[offset:end], total, nil
}
