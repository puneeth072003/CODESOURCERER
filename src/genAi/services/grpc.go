package services

import (
	"context"
	"fmt"
	"genAi/handlers"
	"genAi/models"
	"log"
	"net"
	pb "protobuf/generated"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGenAiServiceServer
}

func (s *server) GenerateTestFiles(_ context.Context, payload *pb.GithubContextRequest) (*pb.GeneratedTestsResponse, error) {

	ctx, client, model := models.InitializeModel()
	defer client.Close()

	res, err := handlers.ProcessAI(ctx, payload, model)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func StartGrpcServer(addr string) {

	grpcServer := grpc.NewServer()

	pb.RegisterGenAiServiceServer(grpcServer, &server{})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not Listen at Port %s: %v", addr, err)
	}

	fmt.Println("gRPC Server started at PORT ", addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to gRPC Server: %v", err)
	}

}
