package resolvers

import (
	"fmt"

	"github.com/go-redis/redis"
)

func createRedisDatabase(databaseUrl string) (Database, error) {

	redisOptions, err := redis.ParseURL(databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("Unable to parse redis URL: %v", err)
	}

	client := redis.NewClient(redisOptions)

	if _, err = client.Ping().Result(); err != nil {
		return nil, fmt.Errorf("unable to connect to redis url: %v", err)
	}
	return &redisDatabase{client: client}, nil
}
