package api

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/service"
	"book-quote-chat-backend/store"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"strconv"
	"time"
)

// 评论列表
func HandleAdminListComments(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	list, total, err := service.ListComments(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := map[string]interface{}{"list": list, "total": total}
	_ = json.NewEncoder(w).Encode(resp)
}

// 审核评论
func HandleAdminAuditComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AdminID   string `json:"adminId"`
		CommentID string `json:"commentId"`
		Action    string `json:"action"` // "delete"...
		Detail    string `json:"detail"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.AuditComment(req.AdminID, req.CommentID, req.Action, req.Detail); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 管理日志
func HandleAdminListLogs(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	list, total, err := service.ListAdminLogs(offset, limit)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := map[string]interface{}{"list": list, "total": total}
	_ = json.NewEncoder(w).Encode(resp)
}

// 管理员新建用户
func HandleAdminCreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	id, err := service.RegisterUser(req.Name, req.Password, req.Email, "", "user")
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"ok": true, "id": id})
}

// 用户列表
func HandleAdminListUsers(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit == 0 {
		limit = 20
	}
	id := r.URL.Query().Get("id")
	nickname := r.URL.Query().Get("nickname")

	list, total, err := service.ListUsers(offset, limit, id, nickname)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	resp := map[string]interface{}{"list": list, "total": total}
	_ = json.NewEncoder(w).Encode(resp)
}

// 封禁用户
func HandleAdminBanUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.BanUser(req.ID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 解封用户
func HandleAdminUnbanUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.UnbanUser(req.ID); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 列表
func HandleAdminListConfig(w http.ResponseWriter, r *http.Request) {
	list, err := service.ListConfig()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(list)
}

// 新增
func HandleAdminAddConfig(w http.ResponseWriter, r *http.Request) {
	var item model.ConfigItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.AddConfigItem(item); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 更新
func HandleAdminUpdateConfig(w http.ResponseWriter, r *http.Request) {
	var item model.ConfigItem
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.UpdateConfigItem(item); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 删除
func HandleAdminDeleteConfig(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if err := service.DeleteConfigItem(id); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 批量导入
func HandleAdminImportConfig(w http.ResponseWriter, r *http.Request) {
	var list []model.ConfigItem
	if err := json.NewDecoder(r.Body).Decode(&list); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if err := service.ImportConfig(list); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// POST /api/admin/quote/add
func HandleAdminAddQuote(w http.ResponseWriter, r *http.Request) {
	var q model.Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	q.ID = uuid.New().String()
	q.Created = time.Now().Unix()
	q.Status = "pending"
	if err := store.AddQuote(q); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(q)
}

// 编辑金句（支持内容、书名、标签等字段编辑）
func HandleAdminEditQuote(w http.ResponseWriter, r *http.Request) {
	var q model.Quote
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if q.ID == "" {
		http.Error(w, "ID必填", 400)
		return
	}
	if err := store.UpdateQuote(q); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 删除金句
func HandleAdminDeleteQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if req.ID == "" {
		http.Error(w, "ID必填", 400)
		return
	}
	if err := store.DeleteQuoteByID(req.ID); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// 审核金句（通过/驳回）
func HandleAdminAuditQuote(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string `json:"id"`
		Status string `json:"status"` // "approved" or "rejected"
		By     string `json:"by"`
		Reason string `json:"reason"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "参数错误", 400)
		return
	}
	if req.ID == "" || (req.Status != "approved" && req.Status != "rejected") {
		http.Error(w, "参数错误", 400)
		return
	}
	// 打印 req
	fmt.Println(req)
 
	if err := store.AuditQuote(model.Quote{
		ID:           req.ID,
		Status:       req.Status,
		AuditBy:      req.By,
		AuditTime:    time.Now().Unix(),
		RejectReason: req.Reason,
	}); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
