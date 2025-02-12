package handlers

import (
	"encoding/json"
	"fmt"

	"github.com/codesourcerer-bot/database/resolvers"
	pb "github.com/codesourcerer-bot/proto/generated"
)

func getContextAndTests(db resolvers.Database, key string) (*pb.ValueType, error) {

	val, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var value pb.ValueType

	if err = json.Unmarshal([]byte(val), &value); err != nil {
		return nil, fmt.Errorf("unable to unmarshal json")
	}

	return &value, nil
}

func setContextAndTests(db resolvers.Database, key string, value *pb.ValueType) (*pb.ResultType, error) {

	valBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal json")
	}
	ok, err := db.Set(key, string(valBytes))
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

	return &pb.ResultType{Result: ok}, nil
}
