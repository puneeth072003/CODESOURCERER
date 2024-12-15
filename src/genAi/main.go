package main

import "genAi/services"

func main() {

	go services.StartGinServer(":3001")
	go services.StartGrpcServer(":9001")

	select {}
}
