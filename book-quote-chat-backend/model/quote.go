package model

type Quote struct {
	ID           string   `json:"id"`           // 金句ID
	UserID       string   `json:"userId"`       // 发布者ID
	Content      string   `json:"text"`         // 金句内容
	Book         string   `json:"book"`         // 书名
	User         string   `json:"user"`         // 发布者
	Source       string   `json:"source"`       // 金句详细来源
	Tags         []string `json:"tags"`         // 主题标签
	Image        string   `json:"image"`        // 可选，配图地址
	AuditBy      string   `json:"auditBy"`      // 审核人ID或名称
	AuditTime    int64    `json:"auditTime"`    // 审核时间
	RejectReason string   `json:"rejectReason"` // 驳回理由
	LikeCount    int      `json:"likeCount"`    // 点赞数
	CommentCount int      `json:"commentCount"` // 评论数
	ReportCount  int      `json:"reportCount"`  // 举报数
	Created      int64    `json:"created"`      // 创建时间
	Status       string   `json:"status"`       // 状态，"pending"/"approved"/"rejected"
}
