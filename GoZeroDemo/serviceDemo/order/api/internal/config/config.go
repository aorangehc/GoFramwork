// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	// 配置上 RPC 的客户端
	UserRpc zrpc.RpcClientConf

	Mysql struct {
		DataSource string // $user:$password@tcp($host:$port)/$dbname
	}
	CacheRedis cache.CacheConf
}
