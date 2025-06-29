package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
	"time"
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

func AuditComment(c model.Comment) error {
	comments, err := LoadComments()
	if err != nil {
		return err
	}
	found := false
	for i := range comments {
		if comments[i].ID == c.ID {
			comments[i].Status = c.Status
			comments[i].AuditBy = c.AuditBy
			comments[i].AuditTime = c.AuditTime
			comments[i].RejectReason = c.RejectReason
			found = true
			break
		}
	}
	if !found {
		return errors.New("comment not found")
	}
	return SaveComments(comments)
}

func DeleteCommentByID(id string) error {
	comments, err := LoadComments()
	if err != nil {
		return err
	}
	idx := -1
	for i := range comments {
		if comments[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("comment not found")
	}
	comments = append(comments[:idx], comments[idx+1:]...)
	return SaveComments(comments)
}

func BatchAuditComment(ids []string, status, by, reason string) error {
	comments, err := LoadComments()
	if err != nil {
		return err
	}
	idSet := map[string]struct{}{}
	for _, id := range ids {
		idSet[id] = struct{}{}
	}
	for i := range comments {
		if _, ok := idSet[comments[i].ID]; ok {
			comments[i].Status = status
			comments[i].AuditBy = by
			comments[i].AuditTime = time.Now().Unix()
			comments[i].RejectReason = reason
		}
	}
	return SaveComments(comments)
}
