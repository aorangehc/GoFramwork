// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"

	"gozeroservicedemo/order/api/internal/svc"
	"gozeroservicedemo/order/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLogic) Delete(req *types.DeleteReq) (resp *types.DeleteResp, err error) {
	// 检查参数
	if req.OrderSn == "" || req.UserId <= 0 {
		return &types.DeleteResp{
			Success: false,
			Message: "参数错误",
		}, fmt.Errorf("参数错误")
	}

	if _, err := l.svcCtx.OrderModel.FindOne(l.ctx, req.OrderSn); err == sqlx.ErrNotFound {
		return &types.DeleteResp{
			Success: false,
			Message: "订单不存在",
		}, fmt.Errorf("订单不存在")
	} else if err != nil {
		return &types.DeleteResp{
			Success: false,
			Message: "订单查询失败",
		}, fmt.Errorf("订单查询失败：%v", err)
	}

	// 删除订单
	err = l.svcCtx.OrderModel.Delete(l.ctx, req.OrderSn)
	if err != nil {
		return &types.DeleteResp{
			Success: false,
			Message: "订单删除失败",
		}, fmt.Errorf("订单删除失败：%v", err)
	}

	return &types.DeleteResp{
		Success: true,
		Message: "订单删除成功",
	}, nil
}
