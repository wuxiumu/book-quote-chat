// service/ratelimiter_registry.go
package service

import (
	"book-quote-chat-backend/store"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

type RateLimiterRegistry struct {
	Like    store.LikeRateLimiter
	Comment store.LikeRateLimiter
	Report  store.LikeRateLimiter
}

var RateLimiters RateLimiterRegistry

func InitRateLimitersFromEnv() {
	_ = godotenv.Load(".env")
	RateLimiters.Like = newLimiterFromEnv("LIKE_LIMIT")
	RateLimiters.Comment = newLimiterFromEnv("COMMENT_LIMIT")
	RateLimiters.Report = newLimiterFromEnv("REPORT_LIMIT")
}

func newLimiterFromEnv(prefix string) store.LikeRateLimiter {
	mode := os.Getenv(prefix + "_MODE")
	window, _ := strconv.ParseInt(os.Getenv(prefix+"_WINDOW"), 10, 64)
	if window == 0 {
		window = 60
	}
	file := os.Getenv(prefix + "_FILE")
	redisAddr := os.Getenv(prefix + "_REDIS")
	mysqlDSN := os.Getenv(prefix + "_MYSQL")
	switch mode {
	case "memory":
		return store.NewMemoryLikeRateLimiter(window)
	case "file":
		return store.NewFileLikeRateLimiter(file, window)
	case "redis":
		rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
		return store.NewRedisLikeRateLimiter(rdb, window)
	case "db":
		db, err := sql.Open("mysql", mysqlDSN)
		if err != nil {
			panic(err)
		}
		return store.NewDBLikeRateLimiter(db, window)
	default:
		return store.NewMemoryLikeRateLimiter(window)
	}
}
