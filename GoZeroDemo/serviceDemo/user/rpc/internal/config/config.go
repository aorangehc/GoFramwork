package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	// MySQL
	Mysql struct {
		DataSource string // $user:$password@tcp($host:$port)/$dbname
	}

	// Redis
	CacheRedis cache.CacheConf
}
