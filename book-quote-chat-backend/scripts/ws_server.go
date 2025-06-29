package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// 单条消息结构体，和你的前端 Msg 完全对应
type Msg struct {
	ID      string `json:"id"`
	User    string `json:"user"`
	Avatar  string `json:"avatar"`
	Text    string `json:"text"`
	Created string `json:"created"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	broadcast = make(chan Msg, 128)
	msgMutex  sync.Mutex
	msgList   []Msg
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
)

func main() {
	http.HandleFunc("/ws/chat", handleWS)
	log.Println("WebSocket server at ws://localhost:8080/ws/chat")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 新连接加到 clients
	clients[conn] = true
	defer delete(clients, conn)

	// 首次连接推送历史消息（最新 50 条）
	msgMutex.Lock()
	if len(msgList) > 0 {
		history := msgList
		if len(history) > 50 {
			history = history[len(history)-50:]
		}
		conn.WriteJSON(history)
	}
	msgMutex.Unlock()

	// 读协程
	go func() {
		for {
			var msg Msg
			err := conn.ReadJSON(&msg)
			if err != nil {
				break
			}
			// 保存历史
			msgMutex.Lock()
			msgList = append(msgList, msg)
			if len(msgList) > 200 {
				msgList = msgList[len(msgList)-200:]
			}
			msgMutex.Unlock()
			// 广播
			broadcast <- msg
		}
		conn.Close()
	}()

	// 写协程
	for msg := range broadcast {
		for c := range clients {
			if err := c.WriteJSON(msg); err != nil {
				c.Close()
				delete(clients, c)
			}
		}
	}
}
