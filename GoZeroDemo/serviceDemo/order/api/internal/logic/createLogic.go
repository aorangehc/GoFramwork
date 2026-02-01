// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package logic

import (
	"context"
	"fmt"
	"time"

	"gozeroservicedemo/order/api/internal/svc"
	"gozeroservicedemo/order/api/internal/types"
	"gozeroservicedemo/order/sql/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateLogic) Create(req *types.CreateReq) (resp *types.CreateResp, err error) {
	// 检查参数
	if req.UserId <= 0 || req.GoodsId <= 0 || req.Num <= 0 || req.Amount <= 0 {
		return nil, fmt.Errorf("参数错误")
	}
	// 生成订单号，写入数据库
	orderSn := fmt.Sprintf("ORDER_%d_%d_%s", req.UserId, req.GoodsId, time.Now().Format("20060102150405"))
	_, err = l.svcCtx.OrderModel.Insert(l.ctx, &model.Order{
		OrderSn: orderSn,
		UserId:  req.UserId,
		GoodsId: req.GoodsId,
		Num:     req.Num,
		Amount:  req.Amount,
		Status:  1, // 默认状态为1（待支付）
	})
	if err != nil {
		return nil, fmt.Errorf("订单创建失败：%v", err)
	}
	return &types.CreateResp{
		OrderSn: orderSn,
	}, nil
}
