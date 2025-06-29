package model

type Comment struct {
	ID           string `json:"id"`           // 评论ID
	UserID       string `json:"userId"`       // 发表评论的用户ID
	UserName     string `json:"userName"`     // 发表评论的用户名
	Avatar       string `json:"avatar"`       // 发表评论的用户头像
	TargetType   string `json:"targetType"`   // "quote"、"msg"等
	TargetID     string `json:"targetId"`     // 评论目标ID，比如说回复哪条消息
	Content      string `json:"content"`      //评论内容
	ParentID     string `json:"parentId"`     // 回复父评论ID，无则""
	Created      int64  `json:"created"`      // 评论创建时间
	Status       string `json:"status"`       // "pending", "approved", "rejected"
	RejectReason string `json:"rejectReason"` // 审核不通过的原因
	AuditBy      string `json:"auditBy"`      // 审核人ID
	AuditTime    int64  `json:"auditTime"`    // 审核时间
}
