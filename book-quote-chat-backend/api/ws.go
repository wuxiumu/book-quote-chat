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
	userConnMap    = make(map[string][]*websocket.Conn)     // userId=>[]conn
	sessionConnMap = make(map[string][]*websocket.Conn)     // sessionKey => []conn
	userConnMu     sync.Mutex                               // 读写 userConnMap 锁
	broadcastChan  = make(chan map[string]interface{}, 128) // 广播消息通道
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
	sessionKey := r.URL.Query().Get("sessionKey") // 新增参数 sessionKey
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
	userConnMap[userId] = append(userConnMap[userId], conn) //   userId 映射
	if sessionKey != "" {
		sessionConnMap[sessionKey] = append(sessionConnMap[sessionKey], conn) //   sessionKey 映射
	}
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
			if msgType, ok := msg["type"].(string); ok && msgType == "chat" {
				msg["from"] = userId // 自动带上发送者
				msg["sessionKey"] = sessionKey
				if sessionKey != "" {
					// 只发给同 sessionKey 的用户（私聊）
					userConnMu.Lock()
					for _, c := range sessionConnMap[sessionKey] {
						_ = c.WriteJSON(msg)
					}
					userConnMu.Unlock()
				} else {
					broadcastChan <- msg // 群聊
				}
			} else if msgType == "rtc-offer" || msgType == "rtc-answer" || msgType == "rtc-ice" || msgType == "rtc-hangup" {
				// ------- WebRTC 信令分发 -------
				// 打印日志
				log.Printf("WS RTC %s: %v", msgType, msg)
				if to, ok := msg["to"].(string); ok && to != "" {
					userConnMu.Lock()
					if conns, ok := userConnMap[to]; ok {
						for _, c := range conns {
							_ = c.WriteJSON(msg)
						}
					}
					userConnMu.Unlock()
				}
			}

		}
		// 连接断开，删除 conn
		userConnMu.Lock()            // 清理连接
		conns := userConnMap[userId] // 群聊
		for i, c := range conns {
			if c == conn {
				userConnMap[userId] = append(conns[:i], conns[i+1:]...)
				break
			}
		}
		if len(userConnMap[userId]) == 0 {
			delete(userConnMap, userId)
		}
		// 私聊
		if sessionKey != "" {
			sConns := sessionConnMap[sessionKey] // 私聊
			for i, c := range sConns {           // 遍历私聊连接
				if c == conn { // 找到当前连接
					sessionConnMap[sessionKey] = append(sConns[:i], sConns[i+1:]...) // 删除连接
					break
				}
			}
			if len(sessionConnMap[sessionKey]) == 0 { // 私聊连接全部断开
				delete(sessionConnMap, sessionKey) // 私聊连接全部断开，删除 sessionKey
			}
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
