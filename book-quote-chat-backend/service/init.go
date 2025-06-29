package service

import (
	"book-quote-chat-backend/store"
	"database/sql"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
)

func InitLikeRateLimiterFromEnv() store.LikeRateLimiter {
	_ = godotenv.Load(".env") // 加载.env到os.Environ
	mode := os.Getenv("LIKE_LIMIT_MODE")
	window, _ := strconv.ParseInt(os.Getenv("LIKE_LIMIT_WINDOW"), 10, 64)
	if window == 0 {
		window = 60
	}
	file := os.Getenv("LIKE_LIMIT_FILE")
	redisAddr := os.Getenv("LIKE_LIMIT_REDIS")
	mysqlDSN := os.Getenv("LIKE_LIMIT_MYSQL")

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
