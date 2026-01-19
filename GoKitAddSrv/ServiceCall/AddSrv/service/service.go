package service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// 1. Service 层 (纯业务逻辑)
// 数据库相关操作可以放在这一层
type AddService interface {
	Add(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

type addService struct {
	// 数据库相关操作可以放在这一层
}

func (s *addService) Add(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func (s *addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", status.Error(codes.InvalidArgument, "两个参数不能同时为空")
	}
	return a + b, nil
}

// 构造函数
func NewAddService() AddService {
	return &addService{}
}
