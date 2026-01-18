package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-kit/kit/endpoint"
)

// Service层
type AddService interface {
	Add(ctx context.Context, a, b int) (int, error)
	Concat(ctx context.Context, a, b string) (string, error)
}

type addService struct{}

func (s *addService) Add(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}

func (s *addService) Concat(_ context.Context, a, b string) (string, error) {
	if a == "" && b == "" {
		return "", errors.New("两个参数都为空字符串")
	}
	return a + b, nil
}

// Endpoint层
type SumRequest struct {
	A int `json:"a"`
	B int `json:"b"`
}

type SumResponse struct {
	Result int    `json:"result"`
	Error  string `json:"error"`
}

type ConcatRequest struct {
	A string `json:"a"`
	B string `json:"b"`
}

type ConcatResponse struct {
	Result string `json:"result"`
	Error  string `json:"error"`
}

func makeSumEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(SumRequest)
		result, err := svc.Add(ctx, req.A, req.B)
		if err != nil {
			return &SumResponse{Error: err.Error()}, nil
		}
		return &SumResponse{Result: result}, nil
	}
}

func makeConcatEndpoint(svc AddService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ConcatRequest)
		result, err := svc.Concat(ctx, req.A, req.B)
		if err != nil {
			return &ConcatResponse{Error: err.Error()}, nil
		}
		return &ConcatResponse{Result: result}, nil
	}
}

// Transport层
func decodeSumRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req SumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeConcatRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req ConcatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}
