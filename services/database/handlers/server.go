package handlers

import (
	"context"

	"github.com/codesourcerer-bot/database/resolvers"
	pb "github.com/codesourcerer-bot/proto/generated"

	"google.golang.org/grpc"
)

type server struct {
	db resolvers.Database
	pb.UnimplementedDatabaseServiceServer
}

func GetGrpcServer(db resolvers.Database) *grpc.Server {
	grpcServer := grpc.NewServer()
	pb.RegisterDatabaseServiceServer(grpcServer, &server{db: db})
	return grpcServer
}

func (s *server) Set(_ context.Context, payload *pb.KeyValType) (*pb.ResultType, error) {
	return setContextAndTests(s.db, payload.Key, payload.GetValue())
}

func (s *server) Get(_ context.Context, payload *pb.KeyType) (*pb.CachedContents, error) {
	return getContextAndTests(s.db, payload.Key)
}

func (s *server) Delete(_ context.Context, payload *pb.KeyType) (*pb.ResultType, error) {
	return deleteContextAndTests(s.db, payload.Key)
}

func (s *server) IsRetriesExhauted(_ context.Context, payload *pb.KeyType) (*pb.ResultType, error) {
	return isRetriesExhauted(s.db, payload.Key)
}
