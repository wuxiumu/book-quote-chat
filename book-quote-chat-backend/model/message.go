package model

type Message struct {
	ID       string `json:"id"`
	FromUser string `json:"fromUser"`
	ToUser   string `json:"toUser"`
	Avatar   string `json:"avatar"`
	Text     string `json:"text"`
	Created  int64  `json:"created"`
	Read     bool   `json:"read"`
	Revoke   bool   `json:"revoke"`
}
