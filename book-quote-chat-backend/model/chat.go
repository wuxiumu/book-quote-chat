package model

type ChatMsg struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Avatar  string `json:"avatar"`
	Text    string `json:"text"`
	Created int64  `json:"created"`
	RoomID  string `json:"roomId"` // 支持多聊天室，简单可以直接写 "main"
}
