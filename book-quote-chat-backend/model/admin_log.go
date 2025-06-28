package model

type AdminLog struct {
	ID       string `json:"id"`
	AdminID  string `json:"adminId"`
	Action   string `json:"action"` // "approve", "reject", "delete", "login"
	Target   string `json:"target"` // "comment", "quote", "user"
	TargetID string `json:"targetId"`
	Detail   string `json:"detail"` // 备注说明
	Created  int64  `json:"created"`
}
