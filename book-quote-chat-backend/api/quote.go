package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
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
