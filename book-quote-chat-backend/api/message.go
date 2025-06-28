package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
)

func HandleGetMessages(w http.ResponseWriter, r *http.Request) {
	user1 := r.URL.Query().Get("user1")
	user2 := r.URL.Query().Get("user2")
	msgs, err := service.GetMessages(user1, user2)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(msgs)
}

func HandleAddMessage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		FromUser string `json:"fromUser"`
		ToUser   string `json:"toUser"`
		Avatar   string `json:"avatar"`
		Text     string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	m, err := service.AddMessage(req.FromUser, req.ToUser, req.Avatar, req.Text)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(m)
}
