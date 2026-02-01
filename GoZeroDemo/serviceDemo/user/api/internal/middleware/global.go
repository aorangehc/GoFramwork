package middleware

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
)

type bodyCopy struct {
	// 结构体能调用接口中实现的方法
	http.ResponseWriter               // 结构体嵌入接口类型
	body                *bytes.Buffer // 用于存储响应体内容
}

func NewBodyCopy(w http.ResponseWriter) *bodyCopy {
	return &bodyCopy{
		ResponseWriter: w,
		body:           bytes.NewBuffer([]byte{}),
	}
}

func (bc *bodyCopy) Write(b []byte) (int, error) {
	// 先存储
	bc.body.Write(b)
	// 往 HTTP 中 png 写入内容
	return bc.ResponseWriter.Write(b)
}

func CopyResp(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 处理请求前
		bc := NewBodyCopy(w) // 用我们自定义的 bodyCopy 替换原有的 ResponseWriter
		next(bc, r)          // 处实际路由请求函数,换成自己的
		// 处理请求后
		logx.Infof("req.url: %s, 响应体内容: %s", r.URL.Path, bc.body.String())
		fmt.Println("req.url:", r.URL.Path, "响应体内容:", bc.body.String())
	}
}
