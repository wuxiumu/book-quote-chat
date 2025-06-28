package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 上传图片（头像/聊天图片通用）
func HandleUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		http.Error(w, "解析表单失败", 400)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "未选择文件", 400)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" {
		http.Error(w, "不支持的文件类型", 400)
		return
	}

	// 生成唯一文件名
	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), RandString(8), ext)
	savePath := filepath.Join("uploads", filename)
	out, err := os.Create(savePath)
	if err != nil {
		http.Error(w, "保存失败", 500)
		return
	}
	defer out.Close()
	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "写入文件失败", 500)
		return
	}

	// 拼接可访问的 URL
	imgUrl := "/uploads/" + filename
	resp := map[string]interface{}{
		"success": true,
		"url":     imgUrl,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

// 生成随机字符串
func RandString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1)
	}
	return string(b)
}
