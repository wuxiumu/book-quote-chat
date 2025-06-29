package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
	"strconv"
)

// 评论
func HandleAddComment(w http.ResponseWriter, r *http.Request) {
	userId, _ := GetUserIDAndGroup(r)
	if userId == "" {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	var req struct {
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
		Content    string `json:"content"`
		ParentID   string `json:"parentId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	comment, err := service.AddComment(userId, "", "", req.TargetType, req.TargetID, req.Content, req.ParentID)
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
		limit = 2
	}
	order := r.URL.Query().Get("order") // 新增
	if order == "" {
		order = "desc"
	}
	comments, total, err := service.GetComments(targetType, targetId, offset, limit, order)
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

// 友情链接
func HandleGetLinks(w http.ResponseWriter, r *http.Request) {
	links, err := service.GetLinks()
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(links)
}

// GetUserFromRequest 从请求上下文获取用户信息
func GetUserFromRequest(r *http.Request) (userId, userName, avatar string) {
	// 假设你用 JWT 中间件或统一在 r.Context() 里塞了用户信息
	if val := r.Context().Value("userId"); val != nil {
		userId, _ = val.(string)
	}
	if val := r.Context().Value("userName"); val != nil {
		userName, _ = val.(string)
	}
	if val := r.Context().Value("avatar"); val != nil {
		avatar, _ = val.(string)
	}
	return
}
