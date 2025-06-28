package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// 新增评论
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"userId"`
		UserName   string `json:"userName"`
		Avatar     string `json:"avatar"`
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
		Content    string `json:"content"`
		ParentID   string `json:"parentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	comment, err := service.AddComment(req.UserID, req.UserName, req.Avatar, req.TargetType, req.TargetID, req.Content, req.ParentID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(comment)
}

// 查询评论
func HandleGetComments(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	comments, total, err := service.GetComments(targetType, targetId, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := map[string]interface{}{
		"list":  comments,
		"total": total,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// 查询单条评论
func HandleGetCommentByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	comment, err := service.GetCommentByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(comment)
}
