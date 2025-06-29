package service

import (
	"book-quote-chat-backend/model"
	"book-quote-chat-backend/store"
	"book-quote-chat-backend/wsutil"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

// 新增全局锁和限流数据
var (
	addLikeMu       sync.Mutex
	ipLikeHistory   = make(map[string]int64) // key: ip+targetType+targetId, value: last like unix
	ipLikeHistoryMu sync.Mutex
	rateLimitWindow int64 = 60 // 秒
	likeRateLimiter store.LikeRateLimiter
)

func AddLike(userId, targetType, targetId string) (model.Like, error) {
	fmt.Println("AddLike called with userId:", userId, "targetType:", targetType, "targetId:", targetId)
	// 并发安全：加锁保证原子性
	addLikeMu.Lock()         // 加锁
	defer addLikeMu.Unlock() // 解锁

	likes, err := store.LoadLikes("", "")
	if err != nil {
		fmt.Println("Error loading likes:", err)
		return model.Like{}, err
	}
	fmt.Println("Current likes count:", len(likes))
	// 检查是否已点赞（幂等）
	for _, l := range likes {
		if l.UserID == userId && l.TargetType == targetType && l.TargetID == targetId {
			fmt.Println("Already liked found:", l)
			return l, nil // 已点直接返回
		}
	}
	like := model.Like{
		ID:         uuid.New().String(),
		UserID:     userId,
		TargetType: targetType,
		TargetID:   targetId,
		Created:    time.Now().Unix(),
	}
	fmt.Println("New like created:", like)
	likes = append(likes, like)
	err = store.SaveLikes(likes)
	fmt.Println("SaveLikes error:", err)
	return like, err
}

func CancelLike(userId, targetType, targetId string) error {
	likes, err := store.LoadLikes("", "")
	if err != nil {
		return err
	}
	out := make([]model.Like, 0, len(likes))
	for _, l := range likes {
		if !(l.UserID == userId && l.TargetType == targetType && l.TargetID == targetId) {
			out = append(out, l)
		}
	}
	return store.SaveLikes(out)
}

func GetLikes(targetType, targetId string) ([]model.Like, error) {
	return store.LoadLikes(targetType, targetId)
}

func CountLikes(targetType, targetId string) (int, error) {
	likes, err := GetLikes(targetType, targetId)
	return len(likes), err
}

// 分页获取点赞
func GetLikesPaged(targetType, targetId string, offset, limit int) ([]model.Like, int, error) {
	likes, err := store.LoadLikes(targetType, targetId)
	if err != nil {
		return nil, 0, err
	}
	total := len(likes)
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return likes[offset:end], total, nil
}

// 用户是否已点赞
func IsLiked(userId, targetType, targetId string) (bool, error) {
	likes, err := store.LoadLikes(targetType, targetId)
	if err != nil {
		return false, err
	}
	for _, l := range likes {
		if l.UserID == userId {
			return true, nil
		}
	}
	return false, nil
}

// 批量计数
type LikeCountReq struct {
	TargetType string `json:"targetType"`
	TargetID   string `json:"targetId"`
}
type LikeCountResp struct {
	TargetID string `json:"targetId"`
	Count    int    `json:"count"`
}

func BatchCountLikes(items []LikeCountReq) ([]LikeCountResp, error) {
	likes, err := store.LoadLikes("", "")
	if err != nil {
		return nil, err
	}
	res := make([]LikeCountResp, 0, len(items))
	for _, item := range items {
		count := 0
		for _, l := range likes {
			if l.TargetType == item.TargetType && l.TargetID == item.TargetID {
				count++
			}
		}
		res = append(res, LikeCountResp{TargetID: item.TargetID, Count: count})
	}
	return res, nil
}

// ip 参数需传入
func AddLikeWithIP(userId, targetType, targetId, ip string) (model.Like, error) {
	key := fmt.Sprintf("%s_%s_%s", ip, targetType, targetId)
	ipLikeHistoryMu.Lock()
	defer ipLikeHistoryMu.Unlock()
	now := time.Now().Unix()
	if last, ok := ipLikeHistory[key]; ok && now-last < rateLimitWindow {
		return model.Like{}, fmt.Errorf("点赞过于频繁，请稍后再试")
	}
	ipLikeHistory[key] = now
	return AddLike(userId, targetType, targetId)
}

func AddLikeWithCheck(userId, targetType, targetId string) (model.Like, error) {
	// 支持多种目标
	switch targetType {
	case "quote":
		if !ExistsQuote(targetId) {
			return model.Like{}, fmt.Errorf("点赞目标不存在")
		}
		// 其它类型可补充
	}
	return AddLike(userId, targetType, targetId)
}

// 集成点赞：IP限流+目标校验+并发锁
func AddLikeWithIPAndCheck(userId, targetType, targetId, ip string) (model.Like, error) {
	// 检查目标是否存在
	if targetType == "quote" && !ExistsQuote(targetId) {
		return model.Like{}, fmt.Errorf("点赞目标不存在")
	}
	// IP限流
	//if !RateLimiters.Like.Allow(ip, targetType, targetId) {
	//	return model.Like{}, fmt.Errorf("点赞过于频繁，请稍后再试")
	//}
	//RateLimiters.Like.Record(ip, targetType, targetId)

	// 并发锁、写入点赞
	//addLikeMu.Lock()
	//defer addLikeMu.Unlock()
	return AddLike(userId, targetType, targetId)
}

func AddLikeWithNotify(userId, targetType, targetId, ip string) (model.Like, error) {
	// ...原有AddLikeWithIPAndCheck逻辑
	like, err := AddLikeWithIPAndCheck(userId, targetType, targetId, ip)
	if err == nil {
		// 获取被点赞对象owner
		ownerId := GetTargetOwnerId(targetType, targetId) // 你需实现
		if ownerId != "" && ownerId != userId {
			title := "收到点赞"
			content := fmt.Sprintf("你的%s被点赞", targetType)
			_, _ = SendNotify(ownerId, "like", title, content, userId)
			// WebSocket实时推送
			wsutil.PushNotifyWS(ownerId, title, content)
		}
	}
	return like, err
}

// 假设有一个获取目标owner的方法
func GetTargetOwnerId(targetType, targetId string) string {
	// 这里只做demo，实际需查数据库/文件
	if targetType == "quote" && targetId == "q1" {
		return "u_owner" // mock
	}
	return ""
}

// GetLikedTargetIDsByUser 查询某用户点赞过的目标ID map，可选 targetType 过滤
// 返回 map[targetId]count，一般只需要 key 判断是否已点赞
// 高性能 map 方案用于接口批量判断“是否点赞”，如前端批量展示点赞状态。
func GetLikedTargetIDsByUser(userId string, targetType string) (map[string]int, error) {
	likes, err := store.LoadLikes("", "")
	if err != nil {
		return nil, err
	}
	m := make(map[string]int)
	for _, l := range likes {
		if l.UserID == userId && (targetType == "" || l.TargetType == targetType) {
			m[l.TargetID]++
		}
	}
	return m, nil
}

// GetLikedTargetIDListByUser 查询某用户点赞过的所有目标ID列表，可选 targetType 过滤
// 高性能 map 方案用于接口批量判断“是否点赞”，如前端批量展示点赞状态。
func GetLikedTargetIDListByUser(userId string, targetType string) ([]string, error) {
	m, err := GetLikedTargetIDsByUser(userId, targetType)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(m))
	for id := range m {
		ids = append(ids, id)
	}
	return ids, nil
}
