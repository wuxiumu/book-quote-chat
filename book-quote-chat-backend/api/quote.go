package api

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/service"
	"book-quote-chat-backend/store"
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
)

func HandleGetQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := service.GetQuotes()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(quotes)
}

func HandleAddQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
		Book string `json:"book"`
		User string `json:"user"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	q, err := service.AddQuote(req.Text, req.Book, req.User)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(q)
}

// 管理后台金句审核分页接口
// GET /api/admin/quotes?offset=0&limit=20&keyword=xxx
func HandleAdminListQuotes(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	keyword := r.URL.Query().Get("keyword")
	status := r.URL.Query().Get("status")
	topic := r.URL.Query().Get("topic")
	list, total, err := service.ListQuotesForAudit(offset, limit, keyword, status, topic)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"list":  list,
		"total": total,
	})
}

// api/quote.go
func HandleGetRandomQuotes(w http.ResponseWriter, r *http.Request) {
	quotes, err := store.LoadQuotes()
	if err != nil {
		http.Error(w, "数据加载失败", 500)
		return
	}
	// 过滤已审核通过
	var approved []model.Quote
	for _, q := range quotes {
		if q.Status == "approved" {
			approved = append(approved, q)
		}
	}
	// 随机打乱
	rand.Shuffle(len(approved), func(i, j int) { approved[i], approved[j] = approved[j], approved[i] })
	// limit=10
	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}
	if limit > len(approved) {
		limit = len(approved)
	}
	// 只返回 text、book 字段（前端只展示这两个字段）
	type simpleQuote struct {
		Content string `json:"text"`
		Book    string `json:"book"`
	}
	var result []simpleQuote
	for i := 0; i < limit; i++ {
		result = append(result, simpleQuote{Content: approved[i].Content, Book: approved[i].Book})
	}
	_ = json.NewEncoder(w).Encode(result)
}
