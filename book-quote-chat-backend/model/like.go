package model

type Like struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	TargetID   string `json:"targetId"`   // 被点赞的对象id
	TargetType string `json:"targetType"` // "quote", "msg", "chat", etc.
	Created    int64  `json:"created"`
}
