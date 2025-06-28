package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"github.com/google/uuid"
	"strings"
	"time"
)

func AddQuote(text, book, user string) (model.Quote, error) {
	quotes, err := store.LoadQuotes()
	if err != nil {
		return model.Quote{}, err
	}
	q := model.Quote{
		ID:      uuid.New().String(),
		Content: text,
		Book:    book,
		User:    user,
		Created: time.Now().Unix(),
	}
	quotes = append(quotes, q)
	err = store.SaveQuotes(quotes)
	return q, err
}

func GetQuotes() ([]model.Quote, error) {
	return store.LoadQuotes()
}

// 如未有目标校验，则返回true
func ExistsQuote(id string) bool {
	quotes, err := store.LoadQuotes()
	if err != nil {
		return false
	}
	for _, q := range quotes {
		if q.ID == id {
			return true
		}
	}
	return false
}

// ListQuotesForAudit 分页、关键词筛选待审核金句 status 可选值 "" | "pending" | "approved" | "rejected"
// ListQuotesForAudit 多条件分页筛选
func ListQuotesForAudit(offset, limit int, keyword, status, topic string) ([]model.Quote, int, error) {
	list, err := store.LoadQuotes()
	if err != nil {
		return nil, 0, err
	}
	var filtered []model.Quote
	for _, q := range list {
		// 状态筛选
		if status != "" && q.Status != status {
			continue
		}
		// 专题/标签筛选
		if topic != "" {
			hasTopic := false
			for _, tag := range q.Tags {
				if tag == topic {
					hasTopic = true
					break
				}
			}
			if !hasTopic {
				continue
			}
		}
		// 关键词筛选（内容、用户、书名、标签）
		if keyword != "" &&
			!strings.Contains(q.Content, keyword) &&
			!strings.Contains(q.User, keyword) &&
			!strings.Contains(q.Book, keyword) {
			found := false
			for _, tag := range q.Tags {
				if strings.Contains(tag, keyword) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		filtered = append(filtered, q)
	}
	total := len(filtered)
	start := offset
	end := offset + limit
	if start > total {
		start = total
	}
	if end > total {
		end = total
	}
	return filtered[start:end], total, nil
}
