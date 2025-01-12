package main

import (
	"fmt"
	"log"
	"net"

	"github.com/codesourcerer-bot/gen-ai/handlers"
)

func main() {

	addr := ":9001"

	grpcServer := handlers.GetGrpcServer()

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not Listen at Port %s: %v", addr, err)
	}

	fmt.Println("gRPC Server started at PORT ", addr)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to gRPC Server: %v", err)
	}

}
