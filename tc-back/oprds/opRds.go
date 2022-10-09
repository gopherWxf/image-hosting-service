package oprds

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"tc-back/config"
)

var rds *redis.Client

//从配置文件中读取数据库的配置信息并连接数据库
func InitRedisConn() (err error) {
	//读取配置文件内容
	redisCfg := config.LoadRedisConfig()
	log.Printf("redis : %s:%d\n", redisCfg.Host, redisCfg.Port)
	rds = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password: redisCfg.Pwd,
		DB:       0,
	})
	_, err = rds.Ping(context.Background()).Result()
	if err != nil {
		return err
	}
	return nil
}
