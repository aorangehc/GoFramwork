package rpc

import (
	"context"
	"gokitaddsrv/endpoint"
	pb "gokitaddsrv/proto"
	"gokitaddsrv/service"

	grpctransport "github.com/go-kit/kit/transport/grpc"
)

// 3. Transport 层 (Protobuf <-> Go Struct)

type grpcServer struct {
	pb.UnimplementedAddSrvServer
	add    grpctransport.Handler
	concat grpctransport.Handler
}

// DecodeRequest: 将 Proto 对象转为 Go 结构体
func decodeAddRequest(_ context.Context, grpcreq interface{}) (interface{}, error) {
	req := grpcreq.(*pb.AddReq)
	return endpoint.AddReq{A: int(req.A), B: int(req.B)}, nil
}

func decodeConcatRequest(_ context.Context, grpcreq interface{}) (interface{}, error) {
	req := grpcreq.(*pb.ConcatReq)
	return endpoint.ConcatReq{A: req.A, B: req.B}, nil
}

// EncodeResponse: 将 Go 结构体转为 Proto 对象
func encodeAddResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.AddResp)
	return &pb.AddResp{
		Sum: int32(resp.Result),
	}, nil
}

func encodeConcatResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(endpoint.ConcatResp)
	return &pb.ConcatResp{
		Res: resp.Result,
	}, nil
}

// 构造 gRPC Handler
func NewAddServer(svc service.AddService) pb.AddSrvServer {
	return &grpcServer{
		add: grpctransport.NewServer(
			endpoint.MakeAddEndpoint(svc),
			decodeAddRequest,
			encodeAddResponse,
		),
		concat: grpctransport.NewServer(
			endpoint.MakeConcatEndpoint(svc),
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
