// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"gozeroservicedemo/order/api/internal/interceptor"
	"gozeroservicedemo/order/api/internal/svc"
	"gozeroservicedemo/order/api/internal/types"
	"gozeroservicedemo/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.DetailReq) (resp *types.DetailResp, err error) {
	// 根据订单号找到数据
	if len(req.OrderSn) == 0 {
		return nil, fmt.Errorf("参数错误")
	}
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, req.OrderSn)
	if err == sqlx.ErrNotFound {
		return nil, fmt.Errorf("订单不存在")
	}
	if err != nil {
		return nil, fmt.Errorf("订单查询失败：%v", err)
	}
	l.ctx = context.WithValue(l.ctx, interceptor.CtxKeyAdminID, "123456")
	// 根据其中的 user_id 调用 user 服务，获取用户信息
	userInfo, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: order.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("获取用户信息失败：%v", err)
	}

	fmt.Printf("订单用户信息：%+v\n", userInfo)

	// 拼接返回结果

	return &types.DetailResp{
		OrderSn:  order.OrderSn,
		UserId:   order.UserId,
		GoodsId:  order.GoodsId,
		Num:      order.Num,
		Amount:   order.Amount,
		Status:   order.Status,
		CreateAt: order.CreateAt.String(),
	}, nil
}
