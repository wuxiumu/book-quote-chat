package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
)

func HandleGetChats(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("roomId")
	msgs, err := service.GetChatMsgs(room)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(msgs)
}

func HandleAddChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User   string `json:"user"`
		Avatar string `json:"avatar"`
		Text   string `json:"text"`
		RoomID string `json:"roomId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	m, err := service.AddChatMsg(req.User, req.Avatar, req.Text, req.RoomID)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(m)
}
