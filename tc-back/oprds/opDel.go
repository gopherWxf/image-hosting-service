package oprds

import "context"

func DelZsetKey(key string, value string) {
	rds.ZRem(context.Background(), key, value)
}
func DelHashKey(key string, value string) {
	rds.HDel(context.Background(), key, value)
}

func DelKey(key string) {
	rds.Del(context.Background(), key)
}
