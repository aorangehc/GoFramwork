package logic

import (
	"context"

	"gozerodemo/source2go/userRpc/example/proto/greet"
	"gozerodemo/source2go/userRpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendMessageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMessageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMessageLogic {
	return &SendMessageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 定义客户端流式 rpc
func (l *SendMessageLogic) SendMessage(stream greet.Greet_SendMessageServer) error {
	// todo: add your logic here and delete this line

	return nil
}
