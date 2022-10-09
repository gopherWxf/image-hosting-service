package oprds

import (
	"context"
	"errors"
)

func CheckToekn(username string, token string) error {
	result, _ := rds.Get(context.Background(), username).Result()
	if result != token {
		return errors.New("token faild")
	}
	return nil
}

func CheckZsetHasFile(key, file string) (bool, error) {
	temp, err := rds.ZRange(context.Background(), key, 0, -1).Result()
	if err != nil {
		return false, err
	}
	for _, t := range temp {
		if t == file {
			return true, nil
		}
	}
	return false, err
}
