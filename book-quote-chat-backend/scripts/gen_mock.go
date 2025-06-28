package main

// 命令行写入数据脚本（推荐长期维护）写代码自动生成 mock 数据写入 data/comments.json，以后可以不断扩充。
// 运行方式：go run scripts/gen_mock.go
import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type Comment struct {
	ID         string `json:"id"`
	TargetType string `json:"targetType"` // "comment" 或 "quote"
	Content    string `json:"content"`
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
	Book       string `json:"book,omitempty"` // 金句时有
	Created    int64  `json:"created"`
}

var commentsText = []string{
	"这本书看得我热泪盈眶！",
	"有些人注定只能陪你走一段路。",
	"内容很触动我，推荐大家看看。",
	"字字珠玑，太真实了。",
	"总有人教会你成长然后转身离开。",
}

var quotesText = []string{
	"人生就像一场马拉松，贵在坚持。",
	"选择比努力更重要。",
	"幸福不是得到想要的一切，而是享受你拥有的一切。",
	"生活就像海洋，只有意志坚强的人才能到达彼岸。",
	"不经历风雨，怎能见彩虹。",
}

var users = []struct {
	Id   string
	Name string
}{
	{"u1", "Aric"},
	{"u2", "Coder"},
	{"u3", "云淡风轻"},
	{"u4", "Geek"},
	{"u5", "小明"},
}

var books = []string{
	"活着", "自控力", "人性的弱点", "平凡的世界", "三体",
}

func main() {
	rand.Seed(time.Now().UnixNano())

	var comments []Comment
	// 生成 10 条评论
	for i := 0; i < 10; i++ {
		user := users[rand.Intn(len(users))]
		comments = append(comments, Comment{
			ID:         fmt.Sprintf("cmt%d", i+1),
			TargetType: "comment",
			Content:    commentsText[rand.Intn(len(commentsText))],
			UserId:     user.Id,
			UserName:   user.Name,
			Created:    time.Now().Unix() - int64(rand.Intn(100000)),
		})
	}
	// 生成 10 条金句
	for i := 0; i < 10; i++ {
		user := users[rand.Intn(len(users))]
		comments = append(comments, Comment{
			ID:         fmt.Sprintf("qt%d", i+1),
			TargetType: "quote",
			Content:    quotesText[rand.Intn(len(quotesText))],
			UserId:     user.Id,
			UserName:   user.Name,
			Book:       books[rand.Intn(len(books))],
			Created:    time.Now().Unix() - int64(rand.Intn(100000)),
		})
	}
	// 确保 data 目录存在
	_ = os.MkdirAll("data", 0755)
	f, err := os.Create("data/comments.json")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(comments); err != nil {
		panic(err)
	}
	fmt.Printf("生成 %d 条评论/金句到 data/comments.json\n", len(comments))
}
