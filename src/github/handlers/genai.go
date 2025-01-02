package handlers

import (
	"context"
	pb "protobuf/generated"
	"time"
)

func GetGeneratedTestsFromGenAI(payload *pb.GithubContextRequest) (*pb.GeneratedTestsResponse, error) {
	conn, err := getConnection(":9001")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewGenAiServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.GenerateTestFiles(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}
