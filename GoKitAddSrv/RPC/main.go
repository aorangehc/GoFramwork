package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "gokitaddsrv/proto"
)

func main() {
	svc := &addService{}

	server := NewAddServer(svc)

	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
		return
	}
	log.Println("gRPC Server is listening on :8888...")

	s := grpc.NewServer()
	pb.RegisterAddSrvServer(s, server)

	err = s.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve: %v", err)
		return
	}
}
