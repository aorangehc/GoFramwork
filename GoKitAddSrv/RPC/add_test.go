package main

import (
	"context"
	pb "gokitaddsrv/proto"
	"log"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
)

// 使用bufconn进行测试，模拟构建测试链接，不用实际端口启动项目
const bufSize = 1024 * 1024

var bufListener *bufconn.Listener

func init() {
	bufListener = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	gs := NewAddServer(&addService{})
	pb.RegisterAddSrvServer(s, gs)

	go func() {
		if err := s.Serve(bufListener); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return bufListener.Dial()
}

// 创建测试客户端连接的辅助函数
func createTestClient(t *testing.T) pb.AddSrvClient {
	conn, err := grpc.DialContext(
		context.Background(),
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	assert.NoError(t, err, "Failed to create gRPC connection")
	t.Cleanup(func() {
		conn.Close()
	})

	return pb.NewAddSrvClient(conn)
}

// TestServiceLayer 直接测试Service层功能
func TestServiceLayer(t *testing.T) {
	svc := &addService{}
	ctx := context.Background()

	// 测试Add方法
	t.Run("Add", func(t *testing.T) {
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

		// 测试大数值
		result, err = svc.Add(ctx, 1000000, 2000000)
		assert.NoError(t, err)
		assert.Equal(t, 3000000, result)
	})

	// 测试Concat方法
	t.Run("Concat", func(t *testing.T) {
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

		// 检查错误是gRPC状态错误
		errStatus, ok := status.FromError(err)
		assert.True(t, ok, "Expected gRPC status error")
		assert.Equal(t, codes.InvalidArgument, errStatus.Code())
		assert.Equal(t, "两个参数不能同时为空", errStatus.Message())
		assert.Equal(t, "", result)

		// 测试包含特殊字符
		result, err = svc.Concat(ctx, "test", "!@#$%")
		assert.NoError(t, err)
		assert.Equal(t, "test!@#$%", result)
	})
}

// TestGRPCSum 测试Sum gRPC接口
func TestGRPCSum(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// 定义测试用例
	testCases := []struct {
		name     string
		a        int32
		b        int32
		expected int32
		wantErr  bool
	}{
		{"正常求和", 10, 2, 12, false},
		{"零值求和", 0, 0, 0, false},
		{"负数求和", -5, 3, -2, false},
		{"大数值求和", 1000000, 2000000, 3000000, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Add(ctx, &pb.AddReq{
				A: tc.a,
				B: tc.b,
			})

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expected, resp.Sum)
		})
	}
}

// TestGRPCConcat 测试Concat gRPC接口
func TestGRPCConcat(t *testing.T) {
	client := createTestClient(t)
	ctx := context.Background()

	// 定义测试用例
	testCases := []struct {
		name     string
		a        string
		b        string
		expected string
		wantErr  bool
		errCode  codes.Code
	}{
		{"正常连接", "hello", "world", "helloworld", false, codes.OK},
		{"左侧空字符串", "", "world", "world", false, codes.OK},
		{"右侧空字符串", "hello", "", "hello", false, codes.OK},
		{"包含特殊字符", "test", "!@#$%", "test!@#$%", false, codes.OK},
		{"两个都为空", "", "", "", true, codes.InvalidArgument},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := client.Concat(ctx, &pb.ConcatReq{
				A: tc.a,
				B: tc.b,
			})

			if tc.wantErr {
				assert.Error(t, err)

				// 检查错误码
				errStatus, ok := status.FromError(err)
				assert.True(t, ok, "Expected gRPC status error")
				assert.Equal(t, tc.errCode, errStatus.Code())

				// 检查错误信息
				if tc.errCode == codes.InvalidArgument {
					assert.Contains(t, errStatus.Message(), "两个参数不能同时为空")
				}

				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, resp)
			assert.Equal(t, tc.expected, resp.Res)
		})
	}
}

// TestEndpointLayer 测试Endpoint层
func TestEndpointLayer(t *testing.T) {
	svc := &addService{}
	ctx := context.Background()

	// 测试Add Endpoint
	t.Run("AddEndpoint", func(t *testing.T) {
		addEndpoint := makeAddEndpoint(svc)

		req := AddReq{A: 5, B: 3}
		resp, err := addEndpoint(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		addResp, ok := resp.(AddResp)
		assert.True(t, ok, "Expected AddResp type")
		assert.Equal(t, 8, addResp.Result)
	})

	// 测试Concat Endpoint
	t.Run("ConcatEndpoint", func(t *testing.T) {
		concatEndpoint := makeConcatEndpoint(svc)

		// 正常情况
		req := ConcatReq{A: "test", B: "123"}
		resp, err := concatEndpoint(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, resp)

		concatResp, ok := resp.(ConcatResp)
		assert.True(t, ok, "Expected ConcatResp type")
		assert.Equal(t, "test123", concatResp.Result)
	})
}
