package connections

import (
	"context"
	"os"
	"time"

	pb "github.com/codesourcerer-bot/proto/generated"
)

func getDatabaseURL() string {
	if url := os.Getenv("DATABASE_SERVICE_URL"); url != "" {
		return url
	}
	return ":9002" // default fallback
}

func GetContextAndTestsFromDatabase(key string) (*pb.CachedContents, error) {
	conn, err := getGrpcConnection(getDatabaseURL())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewDatabaseServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.Get(ctx, &pb.KeyType{Key: key})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func SetContextAndTestsToDatabase(key string, ctx []*pb.SourceFilePayload, tests []*pb.TestFilePayload) (bool, error) {
	conn, err := getGrpcConnection(getDatabaseURL())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewDatabaseServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	val := &pb.CachedContents{Contexts: ctx, Tests: tests}

	res, err := client.Set(c, &pb.KeyValType{Key: key, Value: val})
	if err != nil {
		return false, err
	}

	return res.Result, nil
}

func DeleteContextAndTestsToDatabase(key string) (bool, error) {
	conn, err := getGrpcConnection(getDatabaseURL())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewDatabaseServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.Delete(c, &pb.KeyType{Key: key})
	if err != nil {
		return false, err
	}

	return res.Result, nil
}

func GetRetryExhaustionStatus(key string) (bool, error) {
	conn, err := getGrpcConnection(getDatabaseURL())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	client := pb.NewDatabaseServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.IsRetriesExhauted(c, &pb.KeyType{Key: key})
	if err != nil {
		return false, err
	}

	return res.Result, nil
}
