package oprds

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func SetLoginToken(username string, token string) {
	// redis保存此字符串，用户名：token, 有效时间为24小时
	_, err := rds.Set(context.Background(), username, token, 86400*time.Second).Result()
	if err != nil {
		log.Println(err)
	}
}

func SetZsetKey(key string, score int, value string) {
	rds.ZAdd(context.Background(), key, &redis.Z{
		Score:  float64(score),
		Member: value,
	})
}
func SetHashKey(key string, k1, v1 string) {
	rds.HSet(context.Background(), key, k1, v1)
}

func IncZsetKey(key, value string) {
	rds.ZIncrBy(context.Background(), key, 1, value)
}
