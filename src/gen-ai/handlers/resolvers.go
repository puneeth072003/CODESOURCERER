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

func (s *server) GenerateTestFiles(_ context.Context, payload *pb.GithubContextRequest) (*pb.GeneratedTestsResponse, error) {

	ctx, client, model := models.InitializeModel()
	defer client.Close()

	res, err := GetTestsFromAI(ctx, payload, model)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetGrpcServer() *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterGenAiServiceServer(grpcServer, &server{})
	return grpcServer
}
