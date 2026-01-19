package main

import (
	"log"
	"net"
	"net/http"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "gokitaddsrv/proto"
	"gokitaddsrv/service"
	httpserver "gokitaddsrv/transport/http"
	"gokitaddsrv/transport/rpc"
)

func main() {
	svc := service.NewAddService()

	var g errgroup.Group

	// http
	g.Go(func() error {
		listener, err := net.Listen("tcp", ":8889")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
			return err
		}
		defer listener.Close()
		log.Println("HTTP Server is listening on :8889...")

		httpHandler := httpserver.NewHttpServer(svc)

		return http.Serve(listener, httpHandler)
	})

	// rpc
	g.Go(func() error {
		listener, err := net.Listen("tcp", ":8888")
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
			return err
		}
		defer listener.Close()
		log.Println("gRPC Server is listening on :8888...")

		s := grpc.NewServer()
		pb.RegisterAddSrvServer(s, rpc.NewAddServer(svc))

		return s.Serve(listener)
	})

	// 等待所有服务完成
	if err := g.Wait(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
