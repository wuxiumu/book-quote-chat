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
)

func AddLike(userId, targetType, targetId string) (model.Like, error) {
	// 并发安全：加锁保证原子性
	addLikeMu.Lock()         // 加锁
	defer addLikeMu.Unlock() // 解锁

	likes, err := store.LoadLikes("", "")
	if err != nil {
		return model.Like{}, err
	}
	// 检查是否已点赞（幂等）
	for _, l := range likes {
		if l.UserID == userId && l.TargetType == targetType && l.TargetID == targetId {
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
	likes = append(likes, like)
	err = store.SaveLikes(likes)
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
	// 校验目标
	switch targetType {
	case "quote":
		if !ExistsQuote(targetId) {
			return model.Like{}, fmt.Errorf("点赞目标不存在")
		}
	}
	// 防刷
	key := fmt.Sprintf("%s_%s_%s", ip, targetType, targetId)
	ipLikeHistoryMu.Lock()
	defer ipLikeHistoryMu.Unlock()
	now := time.Now().Unix()
	if last, ok := ipLikeHistory[key]; ok && now-last < rateLimitWindow {
		return model.Like{}, fmt.Errorf("点赞过于频繁，请稍后再试")
	}
	ipLikeHistory[key] = now
	// 并发锁
	addLikeMu.Lock()
	defer addLikeMu.Unlock()
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
