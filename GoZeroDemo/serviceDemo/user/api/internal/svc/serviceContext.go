// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"gozeroservicedemo/user/api/internal/config"
	"gozeroservicedemo/user/api/internal/middleware"
	"gozeroservicedemo/user/sql/cachemodel"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/rest"
)

type ServiceContext struct {
	Config    config.Config
	Cost      rest.Middleware // 自定义路由中间件,要与 api文件中声明的一致
	UserModel cachemodel.UserModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlxConn := sqlx.NewMysql(c.Mysql.DataSource)

	return &ServiceContext{
		Config:    c,
		UserModel: cachemodel.NewUserModel(sqlxConn, c.CacheRedis),
		Cost:      middleware.NewCostMiddleware().Handle,
	}
}
