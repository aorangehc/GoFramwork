package endpoint

import (
	"context"

	"gokitaddsrv/service"

	"github.com/go-kit/kit/endpoint"
)

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

func MakeAddEndpoint(svc service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(AddReq)
		result, err := svc.Add(ctx, req.A, req.B)
		if err != nil {
			return nil, err
		}
		return AddResp{Result: result}, nil
	}
}

func MakeConcatEndpoint(svc service.AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatReq)
		result, err := svc.Concat(ctx, req.A, req.B)
		if err != nil {
			return nil, err
		}
		return ConcatResp{Result: result}, nil
	}
}
