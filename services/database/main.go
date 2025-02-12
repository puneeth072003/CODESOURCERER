package main

import (
	"fmt"
	"log"

	"github.com/codesourcerer-bot/database/handlers"
	"github.com/codesourcerer-bot/database/resolvers"
	"github.com/codesourcerer-bot/database/utils"
)

func main() {
	utils.LoadEnv()
	lis, port := utils.GetListener()

	db, err := resolvers.Factory()
	if err != nil {
		log.Fatalf("Unable to initiate database: %v", err)
	}

	grpcServer := handlers.GetGrpcServer(db)

	fmt.Println("Database gRPC Server started at PORT ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to gRPC Server: %v", err)
	}

}
