package api

import (
	"book-quote-chat-backend/store"
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func HandleClap(w http.ResponseWriter, r *http.Request) {
	ip := getClientIP(r)
	if ip == "" {
		ip = r.RemoteAddr
	}
	count, ok, err := store.AddClap(ip)
	if err != nil {
		http.Error(w, "服务器异常", 500)
		return
	}
	if !ok {
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"count":   count,
			"message": "每个IP只能点一次掌声",
		})
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   count,
	})
}

func HandleClapCount(w http.ResponseWriter, r *http.Request) {
	count := store.GetClapCount()
	_ = json.NewEncoder(w).Encode(map[string]int{"count": count})
}

func getClientIP(r *http.Request) string {
	// 优先 X-Forwarded-For，逗号分割，取第一个
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}
	// 再用 X-Real-IP
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	// 最后 r.RemoteAddr（去端口）
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return ip
	}
	return r.RemoteAddr // 最后兜底
}
