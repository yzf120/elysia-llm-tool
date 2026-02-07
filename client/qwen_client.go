package client

import (
	"sync"

	"github.com/sashabaranov/go-openai"
	"github.com/yzf120/elysia-llm-tool/config"
)

var (
	qwenClient     *openai.Client
	qwenClientOnce sync.Once
)

// GetQwenClient 获取通义千问客户端（单例模式）
func GetQwenClient() *openai.Client {
	qwenClientOnce.Do(func() {
		cfg := config.GetConfig()
		clientConfig := openai.DefaultConfig(cfg.QwenAPIKey)
		clientConfig.BaseURL = cfg.QwenBaseURL
		qwenClient = openai.NewClientWithConfig(clientConfig)
	})
	return qwenClient
}
