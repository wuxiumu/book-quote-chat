package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
)

func HandleGetFriends(w http.ResponseWriter, r *http.Request) {
	friends, err := service.GetFriends()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(friends)
}

func HandleAddFriend(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
		Group  string `json:"group"`
		Remark string `json:"remark"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	f, err := service.AddFriend(req.Name, req.Avatar, req.Group, req.Remark)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(f)
}
