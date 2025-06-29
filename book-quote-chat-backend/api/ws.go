package api

import (
	"book-quote-chat-backend/service"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := os.Getenv("ALLOWED_ORIGIN")
		if allowed == "" {
			return true // 未配置则全部允许
		}
		// 支持多个域名
		for _, item := range strings.Split(allowed, ",") {
			if origin == strings.TrimSpace(item) {
				return true
			}
		}
		return false
	},
}
var (
	userConnMap   = make(map[string][]*websocket.Conn) // userId=>[]conn
	userConnMu    sync.Mutex
	broadcastChan = make(chan map[string]interface{}, 128)
)

func init() {
	go func() {
		for msg := range broadcastChan {
			userConnMu.Lock()
			for _, conns := range userConnMap {
				for _, conn := range conns {
					_ = conn.WriteJSON(msg)
				}
			}
			userConnMu.Unlock()
		}
	}()

	// 新增定时广播在线人数
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			count := 0
			userConnMu.Lock()
			for _, conns := range userConnMap {
				count += len(conns)
			}
			userConnMu.Unlock()
			broadcastChan <- map[string]interface{}{
				"type":  "online",
				"count": count,
			}
		}
	}()
}

func HandleWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	//log.Println("WS连接请求，token:", token)
	userId, _, err := service.ParseJWT(token)
	if err != nil {
		//log.Println("WS鉴权失败：", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//log.Println("WS鉴权通过，userId:", userId)
	conn, err := upgrader.Upgrade(w, r, nil) // 升级协议
	if err != nil {
		log.Println("WS Upgrade Error:", err)
		return
	}
	userConnMu.Lock()
	userConnMap[userId] = append(userConnMap[userId], conn)
	userConnMu.Unlock()

	// 启动 goroutine 读取消息，收到 type == "chat" 的消息写入 broadcastChan
	go func() { // 启动 goroutine 读消息
		for {
			var msg map[string]interface{}
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Println("WS Read Error:", err)
				break
			}
			//log.Println("收到前端消息：", msg)
			if msgType, ok := msg["type"].(string); ok && msgType == "chat" {
				broadcastChan <- msg
			}
		}
		// 连接断开，删除 conn
		userConnMu.Lock()
		conns := userConnMap[userId]
		for i, c := range conns {
			if c == conn {
				userConnMap[userId] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(userConnMap[userId]) == 0 {
			delete(userConnMap, userId)
		}
		userConnMu.Unlock()
		conn.Close()
	}()
	// 必须阻塞 main goroutine，不能 return
	// 保持连接，不断读取避免超时（原代码中循环读取已由上面 goroutine 代替，此处不再阻塞）
	select {} // 阻塞主线程
}

// 后端推送
func PushNotifyWS(userId, title, content string) {
	userConnMu.Lock()
	defer userConnMu.Unlock()
	if conns, ok := userConnMap[userId]; ok {
		for _, conn := range conns {
			_ = conn.WriteJSON(map[string]interface{}{
				"type":    "notify",
				"title":   title,
				"content": content,
			})
		}
	}
}
