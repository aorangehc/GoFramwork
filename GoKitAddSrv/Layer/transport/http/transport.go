package http

import (
	"context"
	"encoding/json"
	"gokitaddsrv/endpoint"
	"gokitaddsrv/service"
	"net/http"

	"github.com/gin-gonic/gin"
	httptransport "github.com/go-kit/kit/transport/http"
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

func NewHttpServer(svc service.AddService) http.Handler {
	sumHandler := httptransport.NewServer(
		endpoint.MakeAddEndpoint(svc),
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
	// r.Run(":8889")
	return r
}
