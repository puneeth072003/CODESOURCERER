package connections

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func getGrpcConnection(port string) (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return conn, nil
}
