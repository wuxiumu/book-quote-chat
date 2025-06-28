package service

import (
	"book-quote-chat-backend/store"
	"time"
)

type TrendPoint struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type StatOverview struct {
	UserCount       int          `json:"userCount"`
	QuoteCount      int          `json:"quoteCount"`
	CommentCount    int          `json:"commentCount"`
	UserTrend       []TrendPoint `json:"userTrend"`
	CommentTrend    []TrendPoint `json:"commentTrend"`
	QuoteTrend      []TrendPoint `json:"quoteTrend"`
	ActiveUserTrend []TrendPoint `json:"activeUserTrend"`
}

func GetStatOverview(startStr, endStr string) (StatOverview, error) {
	users, _ := store.LoadUsers()
	quotes, _ := store.LoadQuotes()
	comments, _ := store.LoadComments()

	// 解析日期
	const layout = "2006-01-02"
	var (
		start time.Time
		end   time.Time
		err   error
	)
	if startStr != "" && endStr != "" {
		start, err = time.Parse(layout, startStr)
		if err != nil {
			return StatOverview{}, err
		}
		end, err = time.Parse(layout, endStr)
		if err != nil {
			return StatOverview{}, err
		}
	} else {
		end = time.Now()
		start = end.AddDate(0, 0, -6) // 默认近7天
	}

	// 生成所有天
	days := []string{}
	for t := start; !t.After(end); t = t.AddDate(0, 0, 1) {
		days = append(days, t.Format(layout))
	}

	userTrend := make([]TrendPoint, 0, len(days))
	commentTrend := make([]TrendPoint, 0, len(days))
	quoteTrend := make([]TrendPoint, 0, len(days))
	activeUserTrend := make([]TrendPoint, 0, len(days))

	for _, day := range days {
		// 用户注册数
		uCount := 0
		for _, u := range users {
			if time.Unix(u.Created, 0).Format(layout) == day {
				uCount++
			}
		}
		userTrend = append(userTrend, TrendPoint{Date: day, Count: uCount})

		// 评论数、活跃用户数
		cCount := 0
		activeUsers := make(map[string]struct{})
		for _, c := range comments {
			if time.Unix(c.Created, 0).Format(layout) == day {
				cCount++
				activeUsers[c.UserID] = struct{}{}
			}
		}
		commentTrend = append(commentTrend, TrendPoint{Date: day, Count: cCount})
		activeUserTrend = append(activeUserTrend, TrendPoint{Date: day, Count: len(activeUsers)})

		// 金句数
		qCount := 0
		for _, q := range quotes {
			if time.Unix(q.Created, 0).Format(layout) == day {
				qCount++
			}
		}
		quoteTrend = append(quoteTrend, TrendPoint{Date: day, Count: qCount})
	}

	return StatOverview{
		UserCount:       len(users),
		QuoteCount:      len(quotes),
		CommentCount:    len(comments),
		UserTrend:       userTrend,
		CommentTrend:    commentTrend,
		QuoteTrend:      quoteTrend,
		ActiveUserTrend: activeUserTrend,
	}, nil
}
