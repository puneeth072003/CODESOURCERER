package resolvers

import (
	"os"
)

type Database interface {
	Set(key string, val string) (bool, error)
	Get(key string) (string, error)
	Delete(key string) (bool, error)
}

func Factory() (Database, error) {

	databaseUrl := os.Getenv("DATABASE_URL")
	return createRedisDatabase(databaseUrl)

}
