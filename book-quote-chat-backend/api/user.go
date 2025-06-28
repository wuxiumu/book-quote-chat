package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
	"strings"
)

// 注册
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
		Email    string `json:"email"`
		Avatar   string `json:"avatar"`
		Group    string `json:"group"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	u, err := service.RegisterUser(req.Name, req.Password, req.Email, req.Avatar, req.Group)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(u)
}

// 登录
func HandleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	token, user, err := service.LoginUser(req.Name, req.Password)
	if err != nil {
		http.Error(w, err.Error(), 401)
		return
	}
	resp := map[string]interface{}{
		"token": token,
		"user":  user,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// 鉴权用例
func HandleGetUser(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, "请登录", 401)
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	userId, group, err := service.ParseJWT(token)
	if err != nil {
		http.Error(w, "token无效", 401)
		return
	}
	// 可根据需要再细分 group 权限判断
	u, err := service.GetUserByID(userId)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	resp := map[string]interface{}{
		"user":  u,
		"group": group,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// 查询用户
func HandleGetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	u, err := service.GetUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), 404)
		return
	}
	_ = json.NewEncoder(w).Encode(u)
}
