package model

type Like struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	TargetID   string `json:"targetId"`   // 被点赞的对象id
	TargetType string `json:"targetType"` // "quote", "msg", "chat", etc.
	Created    int64  `json:"created"`
}

// Like 仅做数据存储，LikeView 用于接口返回，便于返回“是否点赞”状态。
type LikeView struct {
	Like
	Liked bool `json:"liked"`
}
