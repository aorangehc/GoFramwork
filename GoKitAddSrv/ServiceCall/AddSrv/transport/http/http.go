package http

import (
	"context"
	"encoding/json"
	"gokitaddsrv/endpoint"
	"gokitaddsrv/service"
	"net/http"

	loggerMinddle "gokitaddsrv/minddleware/logger"
	"gokitaddsrv/minddleware/ratelimit"

	"github.com/baidubce/bce-sdk-go/rate"
	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func DecodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.AddReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func DecodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoint.ConcatReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func EncodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func NewHttpServer(svc service.AddService, logger log.Logger) http.Handler {
	sum := endpoint.MakeAddEndpoint(svc)

	// 为sum方法添加日志中间件
	sum = loggerMinddle.LoggingMiddleware(log.With(logger, "method", "sum"))(sum)
	// 为sum方法添加限流中间件
	sum = ratelimit.RateLimitMiddleware(rate.NewLimiter(1, 1))(sum)

	sumHandler := httptransport.NewServer(
		sum,
		DecodeSumRequest,
		EncodeResponse,
	)

	concatHandler := httptransport.NewServer(
		endpoint.MakeConcatEndpoint(svc),
		DecodeConcatRequest,
		EncodeResponse,
	)

	r := gin.Default()
	r.POST("/sum", gin.WrapH(sumHandler))
	r.POST("/concat", gin.WrapH(concatHandler))
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	// r.Run(":8889")
	return r
}
