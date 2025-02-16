package main

import (
	"log"

	"github.com/codesourcerer-bot/gen-ai/handlers"
	"github.com/codesourcerer-bot/gen-ai/utils"
)

func main() {

	utils.LoadEnv()
	lis, port := utils.GetListener()

	grpcServer := handlers.GetGrpcServer()

	log.Println("GenAI gRPC Server started at PORT ", port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Unable to gRPC Server: %v", err)
	}

}
