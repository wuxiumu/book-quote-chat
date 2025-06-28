package model

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"` // 不直接返回
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Status   string `json:"status"` // "normal" or "banned"
	Group    string `json:"group"`  // "admin", "user"
	Created  int64  `json:"created"`
}
