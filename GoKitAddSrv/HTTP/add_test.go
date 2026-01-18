package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/stretchr/testify/assert"
)

// TestService层功能
func TestAddService_Add(t *testing.T) {
	svc := &addService{}
	ctx := context.Background()

	// 测试正常情况
	result, err := svc.Add(ctx, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, 5, result)

	// 测试负数
	result, err = svc.Add(ctx, -1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 0, result)

	// 测试零值
	result, err = svc.Add(ctx, 0, 0)
	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}

func TestAddService_Concat(t *testing.T) {
	svc := &addService{}
	ctx := context.Background()

	// 测试正常情况
	result, err := svc.Concat(ctx, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, "helloworld", result)

	// 测试其中一个参数为空
	result, err = svc.Concat(ctx, "hello", "")
	assert.NoError(t, err)
	assert.Equal(t, "hello", result)

	result, err = svc.Concat(ctx, "", "world")
	assert.NoError(t, err)
	assert.Equal(t, "world", result)

	// 测试两个参数都为空（应该返回错误）
	result, err = svc.Concat(ctx, "", "")
	assert.Error(t, err)
	assert.Equal(t, "两个参数都为空字符串", err.Error())
	assert.Equal(t, "", result)
}

// TestEndpoint层功能
func TestSumEndpoint(t *testing.T) {
	svc := &addService{}
	endpoint := makeSumEndpoint(svc)
	ctx := context.Background()

	// 测试正常情况
	req := SumRequest{A: 10, B: 20}
	resp, err := endpoint(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	sumResp, ok := resp.(*SumResponse)
	assert.True(t, ok)
	assert.Equal(t, 30, sumResp.Result)
	assert.Empty(t, sumResp.Error)
}

func TestConcatEndpoint(t *testing.T) {
	svc := &addService{}
	endpoint := makeConcatEndpoint(svc)
	ctx := context.Background()

	// 测试正常情况
	req := ConcatRequest{A: "test", B: "123"}
	resp, err := endpoint(ctx, req)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	concatResp, ok := resp.(*ConcatResponse)
	assert.True(t, ok)
	assert.Equal(t, "test123", concatResp.Result)
	assert.Empty(t, concatResp.Error)

	// 测试两个参数都为空
	req2 := ConcatRequest{A: "", B: ""}
	resp2, err := endpoint(ctx, req2)
	assert.NoError(t, err)
	assert.NotNil(t, resp2)

	concatResp2, ok := resp2.(*ConcatResponse)
	assert.True(t, ok)
	assert.Equal(t, "", concatResp2.Result)
	assert.NotEmpty(t, concatResp2.Error)
	assert.Equal(t, "两个参数都为空字符串", concatResp2.Error)
}

// TestHTTPHandler层功能
func TestHTTPHandlers(t *testing.T) {
	svc := &addService{}
	sumEndpoint := makeSumEndpoint(svc)
	concatEndpoint := makeConcatEndpoint(svc)

	// 创建HTTP handlers
	sumHandler := httptransport.NewServer(
		sumEndpoint,
		decodeSumRequest,
		encodeResponse,
	)

	concatHandler := httptransport.NewServer(
		concatEndpoint,
		decodeConcatRequest,
		encodeResponse,
	)

	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/sum":
			sumHandler.ServeHTTP(w, r)
		case "/concat":
			concatHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	// 测试Sum endpoint
	t.Run("SumHandler", func(t *testing.T) {
		// 正常请求
		reqBody := bytes.NewBufferString(`{"a": 5, "b": 3}`)
		req, err := http.NewRequest("POST", server.URL+"/sum", reqBody)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var sumResp SumResponse
		err = json.NewDecoder(resp.Body).Decode(&sumResp)
		assert.NoError(t, err)
		assert.Equal(t, 8, sumResp.Result)
		assert.Empty(t, sumResp.Error)

		// 无效的JSON请求
		reqBody = bytes.NewBufferString(`{invalid json}`)
		req, err = http.NewRequest("POST", server.URL+"/sum", reqBody)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		// 注意：Go kit的HTTP transport在解码失败时返回500错误
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	// 测试Concat endpoint
	t.Run("ConcatHandler", func(t *testing.T) {
		// 正常请求
		reqBody := bytes.NewBufferString(`{"a": "hello", "b": "kit"}`)
		req, err := http.NewRequest("POST", server.URL+"/concat", reqBody)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var concatResp ConcatResponse
		err = json.NewDecoder(resp.Body).Decode(&concatResp)
		assert.NoError(t, err)
		assert.Equal(t, "hellokit", concatResp.Result)
		assert.Empty(t, concatResp.Error)

		// 测试两个参数都为空的情况
		reqBody = bytes.NewBufferString(`{"a": "", "b": ""}`)
		req, err = http.NewRequest("POST", server.URL+"/concat", reqBody)
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		resp, err = client.Do(req)
		assert.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		err = json.NewDecoder(resp.Body).Decode(&concatResp)
		assert.NoError(t, err)
		assert.Equal(t, "", concatResp.Result)
		assert.Equal(t, "两个参数都为空字符串", concatResp.Error)
	})
}
