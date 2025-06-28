package model

type Notify struct {
	ID      string `json:"id"`
	UserID  string `json:"userId"`  // 被通知用户
	Type    string `json:"type"`    // "like", "comment", "friend", "sys"
	Title   string `json:"title"`   // 简短描述
	Content string `json:"content"` // 详细内容
	FromID  string `json:"fromId"`  // 触发通知的人/对象
	Read    bool   `json:"read"`    // 已读未读
	Created int64  `json:"created"`
}
