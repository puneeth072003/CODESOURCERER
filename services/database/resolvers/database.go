package resolvers

import (
	"fmt"
	"os"
)

type Database interface {
	Set(key string, val string) (bool, error)
	Get(key string) (string, error)
	Delete(key string) (bool, error)
}

func Factory() (Database, error) {

	// Maybe ill remove this once i know the Valkey URL
	databaseType := os.Getenv("DATABASE_TYPE")
	databaseUrl := os.Getenv("DATABASE_URL")

	switch databaseType {
	case "redis":
		return createRedisDatabase(databaseUrl)

	case "valkey":
		return createValkeyDatabase(databaseUrl)

	default:
		return nil, fmt.Errorf("Database not implemented")

	}
}
