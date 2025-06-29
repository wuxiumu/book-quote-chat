package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"` // 不直接返回密码，只返回加密后的密码哈希值
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Status   string `json:"status"` // "normal" or "banned"
	Group    string `json:"group"`  // "admin", "user"
	Created  int64  `json:"created"`
}
