package main

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	_ "os"
	"strings"
	"time"
)

type OriginQuote struct {
	ID      string `json:"id"`
	Text    string `json:"text"`
	Book    string `json:"book"`
	User    string `json:"user"`
	Created string `json:"created"`
}

type Quote struct {
	ID           string   `json:"id"`
	UserId       string   `json:"userId"`
	Text         string   `json:"text"`
	Book         string   `json:"book"`
	User         string   `json:"user"`
	Source       string   `json:"source"`
	Tags         []string `json:"tags"`
	Image        string   `json:"image"`
	AuditBy      string   `json:"auditBy"`
	AuditTime    int64    `json:"auditTime"`
	RejectReason string   `json:"rejectReason"`
	LikeCount    int      `json:"likeCount"`
	CommentCount int      `json:"commentCount"`
	ReportCount  int      `json:"reportCount"`
	Created      int64    `json:"created"`
	Status       string   `json:"status"`
}

var tagList = [][]string{
	{"情感", "哲理", "成长"},
	{"治愈", "温暖"},
	{"成长", "人生"},
	{"青春", "自省"},
	{"友情", "正能量"},
	{"爱情"},
	{"哲理"},
	{"治愈"},
}

var statusList = []string{"approved", "pending", "rejected"}

func randomTags() []string {
	return tagList[rand.Intn(len(tagList))]
}

func randomStatus() string {
	return statusList[rand.Intn(len(statusList))]
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// 1. 读取原始文件（只保留数组部分，没有md标记）
	data, err := ioutil.ReadFile("data/jingju.json")
	if err != nil {
		fmt.Println("读取文件失败:", err)
		return
	}
	// 打印原始数据

	var originQuotes []OriginQuote
	err = json.Unmarshal(data, &originQuotes)
	if err != nil {
		fmt.Println("JSON解析失败:", err)
		return
	}
	fmt.Printf("原始数据共 %d 条\n", len(originQuotes))
	var newQuotes []Quote
	for _, oq := range originQuotes {
		created, _ := time.Parse(time.RFC3339, oq.Created)
		q := Quote{
			ID:           oq.ID,
			UserId:       fmt.Sprintf("u%d", rand.Intn(10)+1),
			Text:         oq.Text,
			Book:         strings.Trim(oq.Book, "《》"),
			User:         extractUserName(oq.User),
			Source:       "",
			Tags:         randomTags(),
			Image:        "",
			AuditBy:      "",
			AuditTime:    created.Unix(),
			RejectReason: "",
			LikeCount:    rand.Intn(200),
			CommentCount: rand.Intn(15),
			ReportCount:  rand.Intn(5),
			Created:      created.Unix(),
			Status:       randomStatus(),
		}
		newQuotes = append(newQuotes, q)
	}
	//
	out, _ := json.MarshalIndent(newQuotes, "", "  ")
	err = os.WriteFile("test.json", out, 0644)
	if err != nil {
		fmt.Println("写入 test.json 失败:", err)
		return
	}
	fmt.Printf("生成成功，共 %d 条，保存在 test.json\n", len(newQuotes))
}

func extractUserName(u string) string {
	// 支持 "豆瓣网友・xxx" => xxx  或 "云村・无语" => 无语  或 "情感网友・xxx"
	if i := strings.Index(u, "・"); i >= 0 {
		return u[i+len("・"):]
	}
	return u
}
