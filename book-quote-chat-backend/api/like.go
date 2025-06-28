package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"fmt"
	"net/http"
)

func HandleGetLikes(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	likes, err := service.GetLikes(targetType, targetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(likes)
}

func HandleAddLike(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"userId"`
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	ip := r.RemoteAddr
	like, err := service.AddLikeWithIPAndCheck(req.UserID, req.TargetType, req.TargetID, ip)
	if err != nil {
		http.Error(w, err.Error(), 429) // 429: Too Many Requests
		return
	}
	_ = json.NewEncoder(w).Encode(like)
}

func HandleCancelLike(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"userId"`
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.CancelLike(req.UserID, req.TargetType, req.TargetID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// 点赞数统计
func HandleCountLikes(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	count, err := service.CountLikes(targetType, targetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]int{"count": count})
}

// 用户已点赞状态
func HandleIsLiked(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	liked, err := service.IsLiked(userId, targetType, targetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]bool{"liked": liked})
}

// 分页查询
func HandleGetLikesPaged(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType")
	targetId := r.URL.Query().Get("targetId")
	offset := 0
	limit := 20
	if o := r.URL.Query().Get("offset"); o != "" {
		fmt.Sscanf(o, "%d", &offset)
	}
	if l := r.URL.Query().Get("limit"); l != "" {
		fmt.Sscanf(l, "%d", &limit)
	}
	likes, total, err := service.GetLikesPaged(targetType, targetId, offset, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  likes,
		"total": total,
	})
}

// 批量统计
func HandleBatchCountLikes(w http.ResponseWriter, r *http.Request) {
	var req []service.LikeCountReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	res, err := service.BatchCountLikes(req)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// 取消点赞，兼容GET/POST
func HandleCancelLikeCompat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		UserID     string `json:"userId"`
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if r.Method == http.MethodGet {
		req.UserID = r.URL.Query().Get("userId")
		req.TargetType = r.URL.Query().Get("targetType")
		req.TargetID = r.URL.Query().Get("targetId")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", 400)
			return
		}
	}
	if err := service.CancelLike(req.UserID, req.TargetType, req.TargetID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
