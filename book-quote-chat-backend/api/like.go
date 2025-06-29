package api

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/service"
	"encoding/json"
	"fmt"
	"net/http"
)

// 获取点赞列表或批量查询已点赞 targetId
// 支持批量查询：如果仅传 targetType（无 targetId），则获取当前 userId 的所有该类型点赞，返回已点赞 targetId 数组
// 如果传了 targetId，则兼容原来行为，返回点赞用户列表（LikeView）
// 这样前端可用于批量已点赞判断
func HandleGetLikes(w http.ResponseWriter, r *http.Request) {
	targetType := r.URL.Query().Get("targetType") // 目标类型
	targetId := r.URL.Query().Get("targetId")     // 目标ID
	userId, _ := GetUserIDAndGroup(r)
	if targetId == "" {
		// 批量查询：获取当前 userId 对该类型点赞过的全部 targetId
		targetIds, err := service.GetLikedTargetIDsByUser(userId, targetType)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		_ = json.NewEncoder(w).Encode(targetIds) // 返回 []string
		return
	}
	// 原有逻辑：返回 targetId 下所有点赞用户
	likes, err := service.GetLikes(targetType, targetId)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var res []model.LikeView
	for _, l := range likes { // 转换成 LikeView
		res = append(res, model.LikeView{
			Like:  l,
			Liked: l.UserID == userId, // 通常是 true，仅本人能查到自己点赞
		})
	}
	_ = json.NewEncoder(w).Encode(res)
}

// 添加点赞
func HandleAddLike(w http.ResponseWriter, r *http.Request) {
	var req struct {
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("[HandleAddLike] decode error:", err)
		http.Error(w, "参数错误", 400)
		return
	}
	ip := r.RemoteAddr
	userId, _ := GetUserIDAndGroup(r)
	if userId == "" {
		fmt.Println("[HandleAddLike] userId 为空，未登录或 token 解析失败")
		http.Error(w, "未登录", 401)
		return
	}
	like, err := service.AddLikeWithIPAndCheck(userId, req.TargetType, req.TargetID, ip)
	if err != nil {
		fmt.Println("[HandleAddLike] AddLikeWithIPAndCheck err:", err)
		http.Error(w, err.Error(), 429)
		return
	}
	fmt.Println("[HandleAddLike] Like 成功:", like)
	_ = json.NewEncoder(w).Encode(like)
}

func HandleCancelLike(w http.ResponseWriter, r *http.Request) {
	// 此改法防止前端伪造 UserID，所有敏感操作只信任 token 解析结果
	var req struct {
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	userId, _ := GetUserIDAndGroup(r)
	if userId == "" {
		http.Error(w, "未登录", http.StatusUnauthorized)
		return
	}
	if err := service.CancelLike(userId, req.TargetType, req.TargetID); err != nil {
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
	// 此改法防止前端伪造 UserID，所有敏感操作只信任 token 解析结果
	var req struct {
		TargetType string `json:"targetType"`
		TargetID   string `json:"targetId"`
	}
	if r.Method == http.MethodGet {
		req.TargetType = r.URL.Query().Get("targetType")
		req.TargetID = r.URL.Query().Get("targetId")
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "参数错误", 400)
			return
		}
	}
	userId, _ := GetUserIDAndGroup(r)
	if err := service.CancelLike(userId, req.TargetType, req.TargetID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
