// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type CostMiddleware struct {
}

func NewCostMiddleware() *CostMiddleware {
	return &CostMiddleware{}
}

// 中间件的功能
func (m *CostMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO generate middleware implement function, delete after code implementation
		now := time.Now()
		// Passthrough to next handler if need
		next(w, r)
		cost := time.Since(now)
		logx.Infof("请求 %s 耗时: %v", r.RequestURI, cost)
		fmt.Println("请求", r.RequestURI, "耗时:", cost)
	}
}
