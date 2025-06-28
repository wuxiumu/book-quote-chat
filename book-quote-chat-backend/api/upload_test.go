package api

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestHandleUpload(t *testing.T) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 写入一个假图片内容
	part, err := writer.CreateFormFile("file", "test.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write([]byte{0xff, 0xd8, 0xff}) // jpeg文件头

	writer.Close()

	req := httptest.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	HandleUpload(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Fatalf("上传失败，返回码:%d", resp.StatusCode)
	}
	b, _ := io.ReadAll(resp.Body)
	if !bytes.Contains(b, []byte("url")) {
		t.Fatal("未返回 url 字段")
	}
}
