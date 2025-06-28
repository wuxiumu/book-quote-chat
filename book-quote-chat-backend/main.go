package main

import (
	"book-quote-chat-backend/api"
	"net/http"
)

// 在 main.go 所有 HandleFunc 前增加
func withCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next(w, r)
	}
}

// 其它路由同理

func main() {
	http.HandleFunc("/api/quotes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetQuotes(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddQuote(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetFriends(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddFriend(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/chats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetChats(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddChat(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetMessages(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddMessage(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/likes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetLikes(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddLike(w, r)
		} else if r.Method == http.MethodDelete {
			api.HandleCancelLike(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/likes/count", api.HandleCountLikes)
	http.HandleFunc("/api/likes/status", api.HandleIsLiked)
	http.HandleFunc("/api/likes/paged", api.HandleGetLikesPaged)
	http.HandleFunc("/api/likes/batch-count", api.HandleBatchCountLikes)
	http.HandleFunc("/api/likes/cancel", api.HandleCancelLikeCompat)

	http.HandleFunc("/api/register", api.HandleRegister)
	http.HandleFunc("/api/login", api.HandleLogin)
	http.HandleFunc("/api/user", api.HandleGetUser)           // 必须带Bearer Token
	http.HandleFunc("/api/user_by_id", api.HandleGetUserByID) // 无鉴权版

	http.HandleFunc("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetComments(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddComment(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	http.HandleFunc("/api/comment", api.HandleGetCommentByID)

	http.HandleFunc("/api/upload", api.HandleUpload)
	// 静态文件服务，支持 /uploads 访问
	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	http.HandleFunc("/api/notify/send", api.HandleSendNotify)
	http.HandleFunc("/api/notifies", api.HandleGetNotifies)
	http.HandleFunc("/api/notify/read", api.HandleMarkNotifyRead)

	http.HandleFunc("/ws", api.HandleWS)

	http.HandleFunc("/api/notify/batchread", api.HandleMarkNotifyBatchRead)

	// 用法示例：
	http.HandleFunc("/api/admin/comments", withCORS(api.HandleAdminListComments))

	http.HandleFunc("/api/admin/quotes", withCORS(api.HandleAdminListQuotes))        // GET
	http.HandleFunc("/api/admin/quote/add", withCORS(api.HandleAdminAddQuote))       // POST
	http.HandleFunc("/api/admin/quote/edit", withCORS(api.HandleAdminEditQuote))     // POST
	http.HandleFunc("/api/admin/quote/delete", withCORS(api.HandleAdminDeleteQuote)) // POST
	http.HandleFunc("/api/admin/quote/audit", withCORS(api.HandleAdminAuditQuote))   // POST

	http.HandleFunc("/api/admin/audit_comment", withCORS(api.HandleAdminAuditComment))
	http.HandleFunc("/api/admin/logs", withCORS(api.HandleAdminListLogs))
	http.HandleFunc("/api/admin/create_user", withCORS(api.HandleAdminCreateUser))
	http.HandleFunc("/api/admin/users", withCORS(api.HandleAdminListUsers))
	http.HandleFunc("/api/admin/ban_user", withCORS(api.HandleAdminBanUser))
	http.HandleFunc("/api/admin/unban_user", withCORS(api.HandleAdminUnbanUser))

	http.HandleFunc("/api/admin/config", withCORS(api.HandleAdminListConfig))          // GET 查询全部
	http.HandleFunc("/api/admin/config/add", withCORS(api.HandleAdminAddConfig))       // POST 新增
	http.HandleFunc("/api/admin/config/update", withCORS(api.HandleAdminUpdateConfig)) // POST 更新
	http.HandleFunc("/api/admin/config/delete", withCORS(api.HandleAdminDeleteConfig)) // GET ?id= 删除
	http.HandleFunc("/api/admin/config/import", withCORS(api.HandleAdminImportConfig)) // POST 批量

	http.HandleFunc("/api/admin/stat_overview", withCORS(api.HandleAdminStatOverview)) // 数据统

	http.ListenAndServe(":8080", nil)
}
