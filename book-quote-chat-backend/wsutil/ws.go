package wsutil

import (
	"github.com/gorilla/websocket"
	"sync"
)

var (
	userConnMap = make(map[string]*websocket.Conn)
	userConnMu  sync.Mutex
)

// 你的 ws 相关管理和推送代码...

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
