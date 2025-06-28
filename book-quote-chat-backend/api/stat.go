package api

import (
	"book-quote-chat-backend/service"
	"encoding/json"
	"net/http"
)

func HandleAdminStatOverview(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	stat, err := service.GetStatOverview(start, end)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(stat)
}
