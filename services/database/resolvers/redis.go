package resolvers

import (
	"fmt"

	"github.com/go-redis/redis"
)

type redisDatabase struct {
	client *redis.Client
}

func (r *redisDatabase) Set(key string, value string) (bool, error) {
	if _, err := r.client.Set(key, value, 0).Result(); err != nil {
		return false, fmt.Errorf("unable to set value: %v", err)
	}
	return true, nil
}

func (r *redisDatabase) Get(key string) (string, error) {
	if value, err := r.client.Get(key).Result(); err != nil {
		return "", fmt.Errorf("unable to get value: %v", err)
	} else {
		return value, nil
	}
}

func (r *redisDatabase) Delete(key string) (bool, error) {
	if _, err := r.client.Del(key).Result(); err != nil {
		return false, fmt.Errorf("unable to delete key: %v", err)
	}
	return true, nil
}
