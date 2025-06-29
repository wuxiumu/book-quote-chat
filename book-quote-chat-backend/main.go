package main

import (
	"book-quote-chat-backend/api"
	"book-quote-chat-backend/service"
	"context"
	"net/http"
	"strings"
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
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "请登录", 401)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		userId, group, err := service.ParseJWT(token)
		if err != nil {
			http.Error(w, "token无效", 401)
			return
		}
		ctx := context.WithValue(r.Context(), api.CtxUserID, userId)
		ctx = context.WithValue(ctx, api.CtxGroup, group)
		next(w, r.WithContext(ctx))
	}
}
func registerAPI(pattern string, handler http.HandlerFunc) {
	http.HandleFunc(pattern, withCORS(handler))
}

func registerProtectedAPI(pattern string, handler http.HandlerFunc) {
	http.HandleFunc(pattern, withCORS(RequireAuth(handler)))
}

// 其它路由同理

func main() {
	service.InitRateLimitersFromEnv()
	service.LoadUploadExts()

	registerAPI("/api/quotes/random", api.HandleGetRandomQuotes)
	registerAPI("/api/friend-links", api.HandleGetLinks)
	registerAPI("/api/clap", api.HandleClap)
	registerAPI("/api/clap/count", api.HandleClapCount)

	registerAPI("/api/quotes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetQuotes(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddQuote(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerAPI("/api/friends", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetFriends(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddFriend(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerAPI("/api/chats", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetChats(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddChat(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerAPI("/api/messages", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetMessages(w, r)
		} else if r.Method == http.MethodPost {
			api.HandleAddMessage(w, r)
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerProtectedAPI("/api/likes", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetLikes(w, r) // 获取点赞列表
		} else if r.Method == http.MethodPost {
			api.HandleAddLike(w, r) // 点赞
		} else if r.Method == http.MethodDelete {
			api.HandleCancelLike(w, r) // 取消点赞
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerAPI("/api/likes/count", api.HandleCountLikes)
	registerAPI("/api/likes/status", api.HandleIsLiked)
	registerAPI("/api/likes/paged", api.HandleGetLikesPaged)
	registerAPI("/api/likes/batch-count", api.HandleBatchCountLikes)
	registerAPI("/api/likes/cancel", api.HandleCancelLikeCompat)

	registerAPI("/api/register", api.HandleRegister)
	registerAPI("/api/login", api.HandleLogin)
	registerProtectedAPI("/api/user", api.HandleGetUser)  // 必须带Bearer Token
	registerAPI("/api/user_by_id", api.HandleGetUserByID) // 无鉴权版

	registerProtectedAPI("/api/comments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			api.HandleGetComments(w, r) // 获取评论列表
		} else if r.Method == http.MethodPost {
			api.HandleAddComment(w, r) // 新增评论
		} else {
			http.Error(w, "Method Not Allowed", 405)
		}
	})
	registerAPI("/api/comment", api.HandleGetCommentByID)

	http.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))
	registerProtectedAPI("/api/upload/", api.HandleUpload)

	registerProtectedAPI("/api/notify/send", api.HandleSendNotify)
	registerProtectedAPI("/api/notifies", api.HandleGetNotifies)
	registerProtectedAPI("/api/notify/read", api.HandleMarkNotifyRead)

	registerAPI("/ws", api.HandleWS)

	registerProtectedAPI("/api/notify/batchread", api.HandleMarkNotifyBatchRead)

	// 用法示例：
	registerProtectedAPI("/api/admin/comments", api.HandleAdminListComments)
	registerProtectedAPI("/api/admin/comment/audit", api.HandleAdminAuditComment)
	registerProtectedAPI("/api/admin/comment/delete", api.HandleAdminDeleteComment)
	registerProtectedAPI("/api/admin/comment/audit_batch", api.HandleAdminBatchAuditComment)

	registerProtectedAPI("/api/admin/quotes", api.HandleAdminListQuotes)                 // GET
	registerProtectedAPI("/api/admin/quote/add", api.HandleAdminAddQuote)                // POST
	registerProtectedAPI("/api/admin/quote/edit", api.HandleAdminEditQuote)              // POST
	registerProtectedAPI("/api/admin/quote/delete", api.HandleAdminDeleteQuote)          // POST
	registerProtectedAPI("/api/admin/quote/audit", api.HandleAdminAuditQuote)            // POST
	registerProtectedAPI("/api/admin/quote/audit_batch", api.HandleAdminBatchAuditQuote) // POST

	registerProtectedAPI("/api/admin/audit_comment", api.HandleAdminAuditComment)
	registerProtectedAPI("/api/admin/logs", api.HandleAdminListLogs)
	registerProtectedAPI("/api/admin/create_user", api.HandleAdminCreateUser)
	registerProtectedAPI("/api/admin/users", api.HandleAdminListUsers)
	registerProtectedAPI("/api/admin/ban_user", api.HandleAdminBanUser)
	registerProtectedAPI("/api/admin/unban_user", api.HandleAdminUnbanUser)

	registerProtectedAPI("/api/admin/config", api.HandleAdminListConfig)          // GET 查询全部
	registerProtectedAPI("/api/admin/config/add", api.HandleAdminAddConfig)       // POST 新增
	registerProtectedAPI("/api/admin/config/update", api.HandleAdminUpdateConfig) // POST 更新
	registerProtectedAPI("/api/admin/config/delete", api.HandleAdminDeleteConfig) // GET ?id= 删除
	registerProtectedAPI("/api/admin/config/import", api.HandleAdminImportConfig) // POST 批量

	registerProtectedAPI("/api/admin/stat_overview", api.HandleAdminStatOverview)            // 数据统
	registerProtectedAPI("/api/admin/stat_audit_overview", api.HandleAdminAuditStatOverview) // 数据统

	//http.ListenAndServe(":8080", nil)
	http.ListenAndServe("0.0.0.0:8080", nil)
}
