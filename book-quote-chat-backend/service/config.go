package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"errors"
	"github.com/google/uuid"
	"os"
	"strings"
)

var AllowedUploadExts []string

func LoadUploadExts() {
	exts := os.Getenv("UPLOAD_ACCEPT")
	if exts == "" {
		exts = "jpg,jpeg,png,gif,webp" // 默认
	}
	arr := strings.Split(exts, ",")
	AllowedUploadExts = make([]string, 0, len(arr))
	for _, ext := range arr {
		e := strings.ToLower(strings.TrimSpace(ext))
		if !strings.HasPrefix(e, ".") {
			e = "." + e
		}
		AllowedUploadExts = append(AllowedUploadExts, e)
	}
	//log.Printf("AllowedUploadExts: %v", AllowedUploadExts)
}

// 查
func ListConfig() ([]model.ConfigItem, error) {
	return store.LoadConfig()
}

// 增
func AddConfigItem(item model.ConfigItem) error {
	list, err := store.LoadConfig()
	if err != nil {
		return err
	}
	item.ID = uuid.New().String()
	list = append(list, item)
	return store.SaveConfig(list)
}

// 改
func UpdateConfigItem(item model.ConfigItem) error {
	list, err := store.LoadConfig()
	if err != nil {
		return err
	}
	updated := false
	for i := range list {
		if list[i].ID == item.ID {
			list[i] = item
			updated = true
			break
		}
	}
	if !updated {
		return errors.New("配置不存在")
	}
	return store.SaveConfig(list)
}

// 删
func DeleteConfigItem(id string) error {
	list, err := store.LoadConfig()
	if err != nil {
		return err
	}
	newList := make([]model.ConfigItem, 0, len(list))
	for _, v := range list {
		if v.ID != id {
			newList = append(newList, v)
		}
	}
	return store.SaveConfig(newList)
}

// 批量导入
func ImportConfig(list []model.ConfigItem) error {
	return store.SaveConfig(list)
}

// GetLinks
func GetLinks() ([]model.ConfigItem, error) {
	return store.LoadConfig()
}
