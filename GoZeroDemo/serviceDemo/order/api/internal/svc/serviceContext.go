// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gozeroservicedemo/order/api/internal/config"
	"gozeroservicedemo/order/sql/model"
	"gozeroservicedemo/user/rpc/user"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config     config.Config
	UserRpc    user.UserServiceClient
	OrderModel model.OrderModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)
	userConn := zrpc.MustNewClient(c.UserRpc).Conn()
	return &ServiceContext{
		Config:     c,
		UserRpc:    user.NewUserServiceClient(userConn),
		OrderModel: model.NewOrderModel(sqlxConn),
	}
}
