package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

var (
	quoteFile = "data/quotes.json"
	quotesMu  sync.Mutex
)

func LoadQuotes() ([]model.Quote, error) {
	quotesMu.Lock()
	defer quotesMu.Unlock()
	b, err := os.ReadFile(quoteFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.Quote{}, nil
	}
	var quotes []model.Quote
	err = json.Unmarshal(b, &quotes)
	return quotes, err
}

func SaveQuotes(quotes []model.Quote) error {
	fmt.Printf("SaveQuotes called, total: %d, first: %+v\n", len(quotes), quotes[0])
	quotesMu.Lock()
	defer quotesMu.Unlock()
	b, err := json.MarshalIndent(quotes, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(quoteFile, b, 0644)
}

func AddQuote(q model.Quote) error {
	// 先读取原始 quotes 列表，append，写回
	quotes, err := LoadQuotes()
	if err != nil {
		quotes = []model.Quote{}
	}
	quotes = append(quotes, q)
	return SaveQuotes(quotes)
}

func UpdateQuote(q model.Quote) error {
	quotes, err := LoadQuotes()
	if err != nil {
		return err
	}
	found := false
	for i := range quotes {
		if quotes[i].ID == q.ID {
			quotes[i] = q
			found = true
			break
		}
	}
	if !found {
		return errors.New("quote not found")
	}
	return SaveQuotes(quotes)
}

func DeleteQuote(q model.Quote) error {
	quotes, err := LoadQuotes()
	if err != nil {
		return err
	}
	idx := -1
	for i := range quotes {
		if quotes[i].ID == q.ID {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("quote not found")
	}
	quotes = append(quotes[:idx], quotes[idx+1:]...)
	return SaveQuotes(quotes)
}

func AuditQuote(q model.Quote) error {
	fmt.Printf("UpdateQuote called: %+v\n", q)
	quotes, err := LoadQuotes()
	if err != nil {
		return err
	}
	found := false
	for i := range quotes {
		if quotes[i].ID == q.ID {
			quotes[i].Status = q.Status
			quotes[i].AuditBy = q.AuditBy
			quotes[i].AuditTime = q.AuditTime
			quotes[i].RejectReason = q.RejectReason
			found = true
			break
		}
	}
	if !found {
		return errors.New("quote not found")
	}
	return SaveQuotes(quotes)
}

func DeleteQuoteByID(id string) error {
	quotes, err := LoadQuotes()
	if err != nil {
		return err
	}
	idx := -1
	for i := range quotes {
		if quotes[i].ID == id {
			idx = i
			break
		}
	}
	if idx == -1 {
		return errors.New("quote not found")
	}
	quotes = append(quotes[:idx], quotes[idx+1:]...)
	return SaveQuotes(quotes)
}
