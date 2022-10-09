package oprds

import "context"

func GetShareFileNum(key string) int {
	result, _ := rds.ZCard(context.Background(), key).Result()
	return int(result)
}

func GetZsetZrevrange(key string, start, end int) ([]string, error) {
	return rds.ZRevRange(context.Background(), key, int64(start), int64(end)).Result()
}

func GetHashFilename(key string, md5file string) (string, error) {
	return rds.HGet(context.Background(), key, md5file).Result()
}

func GetZsetScore(key string, value string) (float64, error) {
	return rds.ZScore(context.Background(), key, value).Result()
}
