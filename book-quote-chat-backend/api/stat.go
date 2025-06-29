package api

import (
	"book-quote-chat-backend/service"
	"book-quote-chat-backend/store"
	"encoding/json"
	"net/http"
	"time"
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

func HandleAdminAuditStatOverview(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	quotes, _ := store.LoadQuotes()
	comments, _ := store.LoadComments()

	// 统计今日、累计、待审核、昨日通过等
	quotesToday, quotesTotal, quotesPending, quotesApprovedYesterday := 0, 0, 0, 0
	commentsToday, commentsTotal, commentsPending, commentsApprovedYesterday := 0, 0, 0, 0

	// 近7天数据
	days := 7
	chartDates := []string{}
	chartQuotes := make([]int, days)
	chartComments := make([]int, days)
	for i := days - 1; i >= 0; i-- {
		day := now.AddDate(0, 0, -i)
		chartDates = append(chartDates, day.Format("01-02"))
	}

	quotesTotal = len(quotes)
	commentsTotal = len(comments)

	for _, q := range quotes {
		t := time.Unix(q.Created, 0)
		if t.Format("2006-01-02") == now.Format("2006-01-02") {
			quotesToday++
		}
		if q.Status == "pending" {
			quotesPending++
		}
		if q.Status == "approved" && t.Format("2006-01-02") == now.AddDate(0, 0, -1).Format("2006-01-02") {
			quotesApprovedYesterday++
		}
		// 近7天趋势
		for i := range chartDates {
			day := now.AddDate(0, 0, -(days - 1 - i)).Format("01-02")
			if t.Format("01-02") == day {
				chartQuotes[i]++
			}
		}
	}
	for _, c := range comments {
		t := time.Unix(c.Created, 0)
		if t.Format("2006-01-02") == now.Format("2006-01-02") {
			commentsToday++
		}
		if c.Status == "pending" {
			commentsPending++
		}
		if c.Status == "approved" && t.Format("2006-01-02") == now.AddDate(0, 0, -1).Format("2006-01-02") {
			commentsApprovedYesterday++
		}
		for i := range chartDates {
			day := now.AddDate(0, 0, -(days - 1 - i)).Format("01-02")
			if t.Format("01-02") == day {
				chartComments[i]++
			}
		}
	}

	res := map[string]interface{}{
		"quotesToday":               quotesToday,
		"quotesPending":             quotesPending,
		"quotesTotal":               quotesTotal,
		"quotesApprovedYesterday":   quotesApprovedYesterday,
		"commentsToday":             commentsToday,
		"commentsPending":           commentsPending,
		"commentsTotal":             commentsTotal,
		"commentsApprovedYesterday": commentsApprovedYesterday,
		"chartDates":                chartDates,
		"chartQuotes":               chartQuotes,
		"chartComments":             chartComments,
	}
	_ = json.NewEncoder(w).Encode(res)
}
