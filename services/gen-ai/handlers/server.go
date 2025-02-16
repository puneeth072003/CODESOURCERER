package handlers

import (
	"context"

	pb "github.com/codesourcerer-bot/proto/generated"

	"github.com/codesourcerer-bot/gen-ai/models"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGenAiServiceServer
}

func GetGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterGenAiServiceServer(grpcServer, &server{})
	return grpcServer
}

func (s *server) GenerateTestFiles(_ context.Context, payload *pb.GithubContextRequest) (*pb.GeneratedTestsResponse, error) {

	ctx, client, model := models.InitializeGeneratorModel()
	defer client.Close()

	res, err := getTestsFromAI(ctx, payload, model)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *server) GenerateRetriedTestFiles(_ context.Context, payload *pb.RetryMechanismPayload) (*pb.GeneratedTestsResponse, error) {
	ctx, client, model := models.InitializeParserModel()
	defer client.Close()

	logsPart, err := getParsedLogsFromAI(ctx, payload.GetLogs(), model)
	if err != nil {
		return nil, err
	}

	ctx, client, model = models.InitializeRetryModel()
	defer client.Close()

	res, err := generateRetriedTestsFromAI(ctx, logsPart, payload.GetCache(), model)
	if err != nil {
		return nil, err
	}

	return res, nil

}
