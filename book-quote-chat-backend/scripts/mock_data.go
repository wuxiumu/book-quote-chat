package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"time"
)

type Comment struct {
	ID         string `json:"id"`
	UserID     string `json:"userId"`
	UserName   string `json:"userName"`
	Avatar     string `json:"avatar"`
	TargetType string `json:"targetType"`
	TargetID   string `json:"targetId"`
	Content    string `json:"content"`
	ParentID   string `json:"parentId"`
	Created    int64  `json:"created"`
}

type Quote struct {
	ID           string   `json:"id"`
	UserID       string   `json:"userId"`
	Content      string   `json:"text"`
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

func randFrom(slice []string) string {
	return slice[rand.Intn(len(slice))]
}

func main() {
	rand.Seed(time.Now().UnixNano())
	users := []struct{ ID, Name, Avatar string }{
		{"u1", "Aric", "https://api.dicebear.com/7.x/bottts/svg?seed=aric"},
		{"u2", "Luna", "https://api.dicebear.com/7.x/bottts/svg?seed=luna"},
		{"u3", "Mars", "https://api.dicebear.com/7.x/bottts/svg?seed=mars"},
		{"u4", "小乐", "https://api.dicebear.com/7.x/bottts/svg?seed=xiaole"},
		{"u5", "Niko", "https://api.dicebear.com/7.x/bottts/svg?seed=niko"},
	}
	books := []string{"活着", "百年孤独", "追风筝的人", "解忧杂货店", "嫌疑人X的献身"}
	quotes := []string{
		"世界上最遥远的距离，是鱼与飞鸟的距离。",
		"成年人的崩溃，常常只在一瞬间，但没人会在意你的眼泪。",
		"最怕你一生碌碌无为，还安慰自己平凡可贵。",
		"后来我明白，很多事说再多也没有用。",
		"遗憾和年轻总绑在一起。",
		"你不主动，我们就真的没关系了。",
		"幸福不是得到你想要的一切，而是享受你拥有的一切。",
		"时间会告诉你，越是平淡越长久。",
	}
	tagsPool := [][]string{
		{"情感"}, {"哲理"}, {"成长", "励志"}, {"人生"}, {"孤独", "青春"}, {"生活"}, {"成长", "情感"}, {"回忆"},
	}

	var quoteList []Quote
	for i := 0; i < 20; i++ {
		u := users[rand.Intn(len(users))]
		idx := rand.Intn(len(quotes))
		q := Quote{
			ID:           "q" + randomString(5),
			UserID:       u.ID,
			Content:      quotes[idx],
			Book:         randFrom(books),
			User:         u.Name,
			Source:       "",
			Tags:         tagsPool[idx%len(tagsPool)],
			Image:        "",
			AuditBy:      "",
			AuditTime:    0,
			RejectReason: "",
			LikeCount:    rand.Intn(200),
			CommentCount: rand.Intn(30),
			ReportCount:  rand.Intn(5),
			Created:      time.Now().Add(time.Duration(-rand.Intn(10*24)) * time.Hour).Unix(),
			Status:       []string{"pending", "approved", "rejected"}[rand.Intn(3)],
		}
		quoteList = append(quoteList, q)
	}

	// 生成评论 mock
	var commentList []Comment
	commentSamples := []string{
		"说得太真实了。",
		"每一句都戳心！",
		"人生啊，总要走点弯路。",
		"收藏了，感同身受。",
		"真是金句，想哭。",
		"我也有同样的感受。",
		"写进心里了。",
	}
	for i := 0; i < 40; i++ {
		u := users[rand.Intn(len(users))]
		c := Comment{
			ID:         "c" + randomString(6),
			UserID:     u.ID,
			UserName:   u.Name,
			Avatar:     u.Avatar,
			TargetType: "quote",
			TargetID:   quoteList[rand.Intn(len(quoteList))].ID,
			Content:    randFrom(commentSamples),
			ParentID:   "",
			Created:    time.Now().Add(time.Duration(-rand.Intn(10*24)) * time.Hour).Unix(),
		}
		commentList = append(commentList, c)
	}

	// 写入文件
	_ = writeJSON("data/mock_quotes.json", quoteList)
	_ = writeJSON("data/mock_comments.json", commentList)
}

func randomString(n int) string {
	letters := []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func writeJSON(filename string, data interface{}) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(data)
}
