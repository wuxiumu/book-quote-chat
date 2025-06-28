package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"errors"
	"os"
	"sync"
)

var (
	adminLogFile = "data/admin_logs.json"
	adminLogMu   sync.Mutex
)

func LoadAdminLogs() ([]model.AdminLog, error) {
	adminLogMu.Lock()
	defer adminLogMu.Unlock()
	b, err := os.ReadFile(adminLogFile)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	if len(b) == 0 {
		return []model.AdminLog{}, nil
	}
	var logs []model.AdminLog
	err = json.Unmarshal(b, &logs)
	return logs, err
}

func SaveAdminLogs(logs []model.AdminLog) error {
	adminLogMu.Lock()
	defer adminLogMu.Unlock()
	b, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(adminLogFile, b, 0644)
}
