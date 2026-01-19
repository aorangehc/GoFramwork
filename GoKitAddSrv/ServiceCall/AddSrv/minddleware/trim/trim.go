package trim

import (
	"context"
	"gokitaddsrv/service"
	"gokittrimsrv/pb"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

type withTrimMiddleware struct {
	next        service.AddService
	trimService endpoint.Endpoint // trim 交给这个endpoint处理
}

func NewServiceWithTrim(trimEndpoint endpoint.Endpoint, svc service.AddService) service.AddService {
	return &withTrimMiddleware{
		trimService: trimEndpoint,
		next:        svc,
	}
}

func (mw withTrimMiddleware) Add(ctx context.Context, a, b int) (res int, err error) {
	return mw.next.Add(ctx, a, b) // 与之前一致
}

// Concat 方法需要先发起gRPC调用外部trim service服务
func (mw withTrimMiddleware) Concat(ctx context.Context, a, b string) (res string, err error) {
	// 先调用trim服务，去除字符串中可能存在的空格
	respA, err := mw.trimService(ctx, trimRequest{s: a}) // 请求trim服务处理a
	if err != nil {
		return "", err
	}
	respB, err := mw.trimService(ctx, trimRequest{s: b}) // 请求trim服务处理b
	if err != nil {
		return "", err
	}
	trimA := respA.(trimResponse)
	trimB := respB.(trimResponse)
	return mw.next.Concat(ctx, trimA.s, trimB.s)
}

type trimRequest struct {
	s string
}

type trimResponse struct {
	s string
}

func MakeTrimEndpoint(conn *grpc.ClientConn) endpoint.Endpoint {
	return grpctransport.NewClient(
		conn,
		"pb.Trim",
		"TrimSpace",
		encodeTrimRequest,
		decodeTrimResponse,
		pb.TrimResponse{},
	).Endpoint()
}

// encodeTrimRequest 将内部使用的数据编码为proto
func encodeTrimRequest(_ context.Context, response interface{}) (request interface{}, err error) {
	resp := response.(trimRequest)
	return &pb.TrimRequest{S: resp.s}, nil
}

// decodeTrimResponse 解析pb消息
func decodeTrimResponse(_ context.Context, in interface{}) (interface{}, error) {
	resp := in.(*pb.TrimResponse)
	return trimResponse{s: resp.S}, nil
}
