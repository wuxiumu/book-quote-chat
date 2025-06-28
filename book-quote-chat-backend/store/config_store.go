package store

import (
	"book-quote-chat-backend/model"
	"encoding/json"
	"os"
)

func LoadConfig() ([]model.ConfigItem, error) {
	f, err := os.Open("data/config.json")
	if err != nil {
		if os.IsNotExist(err) {
			return []model.ConfigItem{}, nil
		}
		return nil, err
	}
	defer f.Close()
	var list []model.ConfigItem
	if err := json.NewDecoder(f).Decode(&list); err != nil {
		return nil, err
	}
	return list, nil
}

func SaveConfig(list []model.ConfigItem) error {
	f, err := os.Create("data/config.json")
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(list)
}
