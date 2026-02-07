package client

import (
	"sync"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/yzf120/elysia-llm-tool/config"
)

var (
	doubaoClient     *arkruntime.Client
	doubaoClientOnce sync.Once
)

// GetDoubaoClient 获取豆包客户端（单例模式）
func GetDoubaoClient() *arkruntime.Client {
	doubaoClientOnce.Do(func() {
		cfg := config.GetConfig()
		doubaoClient = arkruntime.NewClientWithApiKey(
			cfg.DoubaoAPIKey,
			arkruntime.WithBaseUrl(cfg.DoubaoBaseURL),
		)
	})
	return doubaoClient
}
