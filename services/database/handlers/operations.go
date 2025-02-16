package handlers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/codesourcerer-bot/database/resolvers"
	pb "github.com/codesourcerer-bot/proto/generated"
)

func getContextAndTests(db resolvers.Database, key string) (*pb.CachedContents, error) {

	val, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	retries, err := db.Get(key + "/retries")
	if err != nil {
		return nil, err
	}

	intRetries, err := strconv.Atoi(retries)
	if err != nil {
		return nil, fmt.Errorf("cannot convert retries to int: %v", err)
	}

	reducedRetries := intRetries - 1

	ok, err := db.Set(key, strconv.Itoa(reducedRetries))
	if err != nil || !ok {
		return nil, err
	}

	var value pb.CachedContents

	if err = json.Unmarshal([]byte(val), &value); err != nil {
		return nil, fmt.Errorf("unable to unmarshal json")
	}

	return &value, nil
}

func setContextAndTests(db resolvers.Database, key string, value *pb.CachedContents) (*pb.ResultType, error) {

	valBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json")
	}
	ok, err := db.Set(key, string(valBytes))
	if err != nil {
		return nil, err
	}

	ok, err = db.Set(key+"/retries", strconv.Itoa(3))
	if err != nil {
		return nil, err
	}

	return &pb.ResultType{Result: ok}, nil
}

func deleteContextAndTests(db resolvers.Database, key string) (*pb.ResultType, error) {
	ok, err := db.Delete(key)
	if err != nil {
		return nil, err
	}

	ok, err = db.Delete(key + "/retries")
	if err != nil {
		return nil, err
	}

	return &pb.ResultType{Result: ok}, nil
}

func isRetriesExhauted(db resolvers.Database, key string) (*pb.ResultType, error) {
	retries, err := db.Get(key + "/retries")
	if err != nil {
		return nil, err
	}

	intRetries, err := strconv.Atoi(retries)
	if err != nil {
		return nil, fmt.Errorf("cannot convert retries to int: %v", err)
	}

	return &pb.ResultType{Result: !(intRetries > 0)}, nil

}
