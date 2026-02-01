package intereptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func MyInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// 自定义拦截器逻辑
	fmt.Printf("服务端拦截器 in\n")

	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Printf("收到的元数据：%v\n", md)
		adminID := md["adminID"]
		fmt.Printf("adminID: %v\n", adminID)
	} else {
		fmt.Printf("没有收到元数据\n")
		return nil, status.Errorf(codes.InvalidArgument, "没有收到元数据")
	}

	if md["token"][0] != "aorangeToken" {
		return nil, status.Errorf(codes.Unauthenticated, "token 验证失败")
	}

	m, err := handler(ctx, req) // 实际的 rpc方法

	fmt.Printf("服务端拦截器 out\n")

	return m, err
}
