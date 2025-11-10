package connections

import (
	"context"
	"os"
	"time"

	pb "github.com/codesourcerer-bot/proto/generated"
)

func getGenAIURL() string {
	if url := os.Getenv("GENAI_SERVICE_URL"); url != "" {
		return url
	}
	return ":9001" // default fallback
}

func GetGeneratedTestsFromGenAI(payload *pb.GithubContextRequest) (*pb.GeneratedTestsResponse, error) {
	conn, err := getGrpcConnection(getGenAIURL())
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

func GetRetriedTestsFromGenAI(payload *pb.RetryMechanismPayload) (*pb.GeneratedTestsResponse, error) {
	conn, err := getGrpcConnection(getGenAIURL())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewGenAiServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	res, err := client.GenerateRetriedTestFiles(ctx, payload)
	if err != nil {
		return nil, err
	}

	return res, nil
}
