package model

type ConfigItem struct {
	ID    string `json:"id"`
	Key   string `json:"key"`   // 如 friendlink
	Name  string `json:"name"`  // 展示名
	Value string `json:"value"` // 配置值（如友链url等）
	Desc  string `json:"desc"`  // 备注/描述
}
