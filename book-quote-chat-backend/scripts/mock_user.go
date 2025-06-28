package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Password string `json:"-"` // 不直接返回
	Avatar   string `json:"avatar"`
	Email    string `json:"email"`
	Status   string `json:"status"` // "normal" or "banned"
	Group    string `json:"group"`  // "admin", "user"
	Created  int64  `json:"created"`
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var users []User
	for i := 1; i <= 100; i++ {
		id := fmt.Sprintf("u%03d", i)
		name := fmt.Sprintf("用户%03d", i)
		email := fmt.Sprintf("user%03d@example.com", i)
		status := "normal"
		if rand.Float64() < 0.1 { // 10%概率为 banned
			status = "banned"
		}
		group := "user"
		if i <= 3 {
			group = "admin" // 前3个为管理员
		}
		avatar := fmt.Sprintf("https://api.dicebear.com/7.x/bottts/svg?seed=%s", randomString(6))
		users = append(users, User{
			ID:      id,
			Name:    name,
			Avatar:  avatar,
			Email:   email,
			Status:  status,
			Group:   group,
			Created: time.Now().Add(-time.Duration(rand.Intn(365*24)) * time.Hour).Unix(),
		})
	}
	f, _ := os.Create("data/mock_users.json")
	defer f.Close()
	_ = json.NewEncoder(f).Encode(users)
}
