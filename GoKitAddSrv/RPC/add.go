package main

import (
	"context"

	pb "gokitaddsrv/proto"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 1. Service 层 (纯业务逻辑)
type AddService interface {
	Add(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

type addService struct{}

func (s *addService) Add(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func (s *addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", status.Error(codes.InvalidArgument, "两个参数不能同时为空")
	}
	return a + b, nil
}

// 2. Endpoint 层 (适配器)

// 请求/响应结构体 (去除了 Error 字段)
type AddReq struct {
	A int `json:"a"`
	B int `json:"b"`
}
type AddResp struct {
	Result int `json:"result"`
}

type ConcatReq struct {
	A string `json:"a"`
	B string `json:"b"`
}
type ConcatResp struct {
	Result string `json:"result"`
}

func makeAddEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddReq)
		result, err := svc.Add(ctx, req.A, req.B)
		if err != nil {
			return nil, err
		}
		return AddResp{Result: result}, nil
	}
}

func makeConcatEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatReq)
		result, err := svc.Concat(ctx, req.A, req.B)
		if err != nil {
			return nil, err
		}
		return ConcatResp{Result: result}, nil
	}
}

// 3. Transport 层 (Protobuf <-> Go Struct)

type grpcServer struct {
	pb.UnimplementedAddSrvServer
	add    grpctransport.Handler
	concat grpctransport.Handler
}

// DecodeRequest: 将 Proto 对象转为 Go 结构体
func decodeAddRequest(_ context.Context, grpcreq interface{}) (interface{}, error) {
	req := grpcreq.(*pb.AddReq)
	return AddReq{A: int(req.A), B: int(req.B)}, nil
}

func decodeConcatRequest(_ context.Context, grpcreq interface{}) (interface{}, error) {
	req := grpcreq.(*pb.ConcatReq)
	return ConcatReq{A: req.A, B: req.B}, nil
}

// EncodeResponse: 将 Go 结构体转为 Proto 对象
func encodeAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(AddResp)
	return &pb.AddResp{
		Sum: int32(resp.Result),
	}, nil
}

func encodeConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(ConcatResp)
	return &pb.ConcatResp{
		Res: resp.Result,
	}, nil
}

// 构造 gRPC Handler
func NewAddServer(svc AddService) pb.AddSrvServer {
	return &grpcServer{
		add: grpctransport.NewServer(
			makeAddEndpoint(svc),
			decodeAddRequest,
			encodeAddResponse,
		),
		concat: grpctransport.NewServer(
			makeConcatEndpoint(svc),
			decodeConcatRequest,
			encodeConcatResponse,
		),
	}
}

// 实现 Protobuf 接口，转交给 Go-Kit Handler
func (s *grpcServer) Add(ctx context.Context, req *pb.AddReq) (*pb.AddResp, error) {
	_, resp, err := s.add.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.AddResp), nil
}

func (s *grpcServer) Concat(ctx context.Context, req *pb.ConcatReq) (*pb.ConcatResp, error) {
	_, resp, err := s.concat.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.ConcatResp), nil
}
