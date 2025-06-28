package api

import (
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var (
	upgrader    = websocket.Upgrader{}
	userConnMap = make(map[string]*websocket.Conn) // userId=>conn
	userConnMu  sync.Mutex
)

func HandleWS(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	if userId == "" {
		http.Error(w, "缺少userId", 400)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	userConnMu.Lock()
	userConnMap[userId] = conn
	userConnMu.Unlock()

	// 保持连接，不断读取避免超时
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
	userConnMu.Lock()
	delete(userConnMap, userId)
	userConnMu.Unlock()
	conn.Close()
}

// 后端推送
func PushNotifyWS(userId, title, content string) {
	userConnMu.Lock()
	defer userConnMu.Unlock()
	if conn, ok := userConnMap[userId]; ok {
		_ = conn.WriteJSON(map[string]interface{}{
			"type":    "notify",
			"title":   title,
			"content": content,
		})
	}
}
