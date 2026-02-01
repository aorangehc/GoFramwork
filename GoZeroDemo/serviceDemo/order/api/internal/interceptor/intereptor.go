package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type CtxKey string

const (
	// 使用自定义类型声明 context 中存储的key，避免冲突
	CtxKeyAdminID CtxKey = "adminID"
)

func AorangeInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	// 在这里可以添加自定义的拦截逻辑，例如日志记录、熔断等
	fmt.Printf("客户端拦截器 in")

	adminID := ctx.Value(CtxKeyAdminID)
	md := metadata.Pairs(
		"key1", "val1",
		"key2", "val1-2",
		"requestID", "12345",
		"token", "aorangeToken",
		"adminID", adminID.(string),
	)

	ctx = metadata.NewOutgoingContext(ctx, md) // 将元数据添加到上下文中

	err := invoker(ctx, method, req, reply, cc, opts...)

	fmt.Printf("客户端拦截器 out")
	return err
}
