package main

import (
	"fmt"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	pb "gokitaddsrv/proto"
	"gokitaddsrv/service"
	httpserver "gokitaddsrv/transport/http"
	"gokitaddsrv/transport/rpc"

	metricsMinddle "gokitaddsrv/minddleware/metrics"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
)

func main() {
	// 初始化指标追踪中间件 - 使用空结构，服务层会进行防御性检查
	svc := service.NewAddService()

	// instrumentation
	fieldKeys := []string{"method", "error"}
	requestCount := kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, fieldKeys)
	requestLatency := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, fieldKeys)
	countResult := kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "my_group",
		Subsystem: "string_service",
		Name:      "count_result",
		Help:      "The result of each count method.",
	}, []string{}) // no fields here

	svc = metricsMinddle.NewInstrumentingMiddleware(
		requestCount,
		requestLatency,
		countResult,
	)(svc)

	var g errgroup.Group

	// http
	g.Go(func() error {
		listener, err := net.Listen("tcp", ":8889")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
			return err
		}
		defer listener.Close()
		fmt.Println("HTTP Server is listening on :8889...")

		// 初始化go-kit logger
		logger := log.NewLogfmtLogger(os.Stderr)

		httpHandler := httpserver.NewHttpServer(svc, logger)

		return http.Serve(listener, httpHandler)
	})

	// rpc
	g.Go(func() error {
		listener, err := net.Listen("tcp", ":8888")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to listen: %v\n", err)
			return err
		}
		defer listener.Close()
		fmt.Println("gRPC Server is listening on :8888...")

		s := grpc.NewServer()
		pb.RegisterAddSrvServer(s, rpc.NewAddServer(svc))

		return s.Serve(listener)
	})

	// 等待所有服务完成
	if err := g.Wait(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
	}
}
