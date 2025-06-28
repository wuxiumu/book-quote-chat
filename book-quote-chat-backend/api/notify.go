package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// 推送通知
func HandleSendNotify(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID  string `json:"userId"`
		Type    string `json:"type"`
		Title   string `json:"title"`
		Content string `json:"content"`
		FromID  string `json:"fromId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	n, err := service.SendNotify(req.UserID, req.Type, req.Title, req.Content, req.FromID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(n)
}

// 查询通知
func HandleGetNotifies(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	list, total, err := service.GetNotifies(userId, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := map[string]interface{}{
		"list":  list,
		"total": total,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// 标记已读
func HandleMarkNotifyRead(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID   string `json:"userId"`
		NotifyID string `json:"notifyId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.MarkNotifyRead(req.UserID, req.NotifyID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func HandleMarkNotifyBatchRead(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID    string   `json:"userId"`
		NotifyIDs []string `json:"notifyIds"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.MarkNotifyBatchRead(req.UserID, req.NotifyIDs); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
