package store

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"os"
	"sync"
	"time"
)

type LikeRateLimiter interface {
	Allow(ip, targetType, targetId string) bool
	Record(ip, targetType, targetId string)
}

// 内存实现
type memoryLikeRateLimiter struct {
	mu      sync.Mutex
	history map[string]int64
	window  int64
}

func NewMemoryLikeRateLimiter(window int64) LikeRateLimiter {
	return &memoryLikeRateLimiter{
		history: make(map[string]int64),
		window:  window,
	}
}
func (m *memoryLikeRateLimiter) Allow(ip, targetType, targetId string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := ip + "_" + targetType + "_" + targetId
	now := time.Now().Unix()
	if last, ok := m.history[key]; ok && now-last < m.window {
		return false
	}
	return true
}
func (m *memoryLikeRateLimiter) Record(ip, targetType, targetId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := ip + "_" + targetType + "_" + targetId
	m.history[key] = time.Now().Unix()
}

// 文件实现
type fileLikeRateLimiter struct {
	file   string
	window int64
	mu     sync.Mutex
}

func NewFileLikeRateLimiter(file string, window int64) LikeRateLimiter {
	return &fileLikeRateLimiter{file: file, window: window}
}
func (f *fileLikeRateLimiter) Allow(ip, targetType, targetId string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	hist := f.load()
	key := ip + "_" + targetType + "_" + targetId
	now := time.Now().Unix()
	if last, ok := hist[key]; ok && now-last < f.window {
		return false
	}
	return true
}
func (f *fileLikeRateLimiter) Record(ip, targetType, targetId string) {
	f.mu.Lock()
	defer f.mu.Unlock()
	hist := f.load()
	key := ip + "_" + targetType + "_" + targetId
	hist[key] = time.Now().Unix()
	data, _ := json.Marshal(hist)
	_ = os.WriteFile(f.file, data, 0644)
}
func (f *fileLikeRateLimiter) load() map[string]int64 {
	data, _ := os.ReadFile(f.file)
	hist := map[string]int64{}
	_ = json.Unmarshal(data, &hist)
	return hist
}

// Redis实现
type redisLikeRateLimiter struct {
	rdb    *redis.Client
	window int64
}

func NewRedisLikeRateLimiter(rdb *redis.Client, window int64) LikeRateLimiter {
	return &redisLikeRateLimiter{rdb: rdb, window: window}
}
func (r *redisLikeRateLimiter) Allow(ip, targetType, targetId string) bool {
	ctx := context.Background()
	key := "like_limit:" + ip + "_" + targetType + "_" + targetId
	val, err := r.rdb.Get(ctx, key).Int64()
	now := time.Now().Unix()
	return err == redis.Nil || now-val >= r.window
}
func (r *redisLikeRateLimiter) Record(ip, targetType, targetId string) {
	ctx := context.Background()
	key := "like_limit:" + ip + "_" + targetType + "_" + targetId
	now := time.Now().Unix()
	r.rdb.Set(ctx, key, now, time.Duration(r.window)*time.Second)
}

// 数据库实现
type dbLikeRateLimiter struct {
	db     *sql.DB
	window int64
}

func NewDBLikeRateLimiter(db *sql.DB, window int64) LikeRateLimiter {
	return &dbLikeRateLimiter{db: db, window: window}
}
func (d *dbLikeRateLimiter) Allow(ip, targetType, targetId string) bool {
	key := ip + "_" + targetType + "_" + targetId
	var lastTime int64
	err := d.db.QueryRow("SELECT last_time FROM like_limit WHERE `key` = ?", key).Scan(&lastTime)
	now := time.Now().Unix()
	return err == sql.ErrNoRows || now-lastTime >= d.window
}
func (d *dbLikeRateLimiter) Record(ip, targetType, targetId string) {
	key := ip + "_" + targetType + "_" + targetId
	now := time.Now().Unix()
	_, err := d.db.Exec("REPLACE INTO like_limit (`key`, last_time) VALUES (?,?)", key, now)
	_ = err
}

// 热加载配置
type RateLimitConfig struct {
	Mode   string `json:"mode"`   // memory/file/redis/db
	Window int64  `json:"window"` // 秒
	File   string `json:"file"`
	Redis  string `json:"redis"` // redis连接串
	MySQL  string `json:"mysql"` // mysql连接串
}

func LoadRateLimitConfig(configPath string) (RateLimitConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return RateLimitConfig{}, err
	}
	var cfg RateLimitConfig
	err = json.Unmarshal(data, &cfg)
	return cfg, err
}

func NewLikeRateLimiterFromConfig(cfg RateLimitConfig) LikeRateLimiter {
	switch cfg.Mode {
	case "memory":
		return NewMemoryLikeRateLimiter(cfg.Window)
	case "file":
		return NewFileLikeRateLimiter(cfg.File, cfg.Window)
	case "redis":
		rdb := redis.NewClient(&redis.Options{Addr: cfg.Redis})
		return NewRedisLikeRateLimiter(rdb, cfg.Window)
	case "db":
		db, err := sql.Open("mysql", cfg.MySQL)
		if err != nil {
			panic(err)
		}
		return NewDBLikeRateLimiter(db, cfg.Window)
	default:
		return NewMemoryLikeRateLimiter(cfg.Window)
	}
}
