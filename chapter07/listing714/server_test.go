package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandleGet(t *testing.T) {
	// 创建HTTP路由处理器
	mux := http.NewServeMux()
	mux.HandleFunc("/post/", handleRequest)

	// 创建响应记录器
	writer := httptest.NewRecorder()

	// 设置请求超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 创建HTTP GET请求
	request, err := http.NewRequestWithContext(ctx, "GET", "/post/1", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 发送请求并记录响应
	done := make(chan struct{})
	go func() {
		mux.ServeHTTP(writer, request)
		close(done)
	}()

	select {
	case <-done:
		// 检查响应状态码
		if writer.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, writer.Code)
		}

		// 解析响应体为Post结构体
		var post Post
		if err := json.Unmarshal(writer.Body.Bytes(), &post); err != nil {
			t.Errorf("Failed to unmarshal response body: %v", err)
		}

		// 检查Post ID
		if post.Id != 1 {
			t.Errorf("Expected post ID %d, got %d", 1, post.Id)
		}
	case <-ctx.Done():
		t.Errorf("Request timed out: %v", ctx.Err())
	}
}
