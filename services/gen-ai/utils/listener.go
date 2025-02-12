package utils

import (
	"fmt"
	"log"
	"net"
	"os"
)

func GetListener() (net.Listener, string) {
	port := os.Getenv("PORT")
	addr := fmt.Sprintf("127.0.0.1:%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Could not Listen at Port %s: %v", port, err)
	}
	return lis, port
}
