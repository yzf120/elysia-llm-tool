package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

// Config 全局配置
type Config struct {
	// 豆包配置
	DoubaoAPIKey  string
	DoubaoBaseURL string

	// 火山开放平台配置（用于调用 GetInferenceUsage 等管理API）
	VolcAccessKeyID     string
	VolcSecretAccessKey string
	VolcRegion          string

	// 混元配置
	HunyuanSecretID  string
	HunyuanSecretKey string

	// 通义千问配置
	QwenAPIKey  string
	QwenBaseURL string

	// 服务配置
	ServerPort int

	// 模型默认参数配置
	DefaultTemperature     float32
	DefaultMaxTokens       int32
	DefaultTopP            float32
	DefaultReasoningEffort string
}

var globalConfig *Config

// InitConfig 初始化配置
func InitConfig() *Config {
	globalConfig = &Config{
		// 豆包配置（强制使用北京节点，深圳节点 ark.cn-shenzhen.volces.com 不可用）
		DoubaoAPIKey:  getEnv("DOUBAO_API_KEY", ""),
		DoubaoBaseURL: getDoubaoBaseURL(),

		// 火山开放平台配置（调用推理用量等管理API）
		VolcAccessKeyID:     getEnv("VOLC_ACCESS_KEY_ID", ""),
		VolcSecretAccessKey: getEnv("VOLC_SECRET_ACCESS_KEY", ""),
		VolcRegion:          getEnv("VOLC_REGION", "cn-beijing"),

		// 混元配置
		HunyuanSecretID:  getEnv("HUNYUAN_SECRET_ID", ""),
		HunyuanSecretKey: getEnv("HUNYUAN_SECRET_KEY", ""),

		// 通义千问配置
		QwenAPIKey:  getEnv("DASHSCOPE_API_KEY", ""),
		QwenBaseURL: getEnv("QWEN_BASE_URL", "https://dashscope.aliyuncs.com/compatible-mode/v1"),

		// 服务配置
		ServerPort: 9001,

		// 模型默认参数配置
		DefaultTemperature:     getEnvFloat32("DEFAULT_TEMPERATURE", 0.7),
		DefaultMaxTokens:       getEnvInt32("DEFAULT_MAX_TOKENS", 2000),
		DefaultTopP:            getEnvFloat32("DEFAULT_TOP_P", 0.9),
		DefaultReasoningEffort: getEnv("DEFAULT_REASONING_EFFORT", "medium"),
	}
	return globalConfig
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if globalConfig == nil {
		return InitConfig()
	}
	return globalConfig
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvFloat32 获取 float32 类型的环境变量
func getEnvFloat32(key string, defaultValue float32) float32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	f, err := strconv.ParseFloat(value, 32)
	if err != nil {
		return defaultValue
	}
	return float32(f)
}

// getDoubaoBaseURL 获取豆包BaseURL，强制使用北京节点（深圳节点不可用）
func getDoubaoBaseURL() string {
	url := os.Getenv("DOUBAO_BASE_URL")
	if url == "" {
		return "https://ark.cn-beijing.volces.com/api/v3"
	}
	// 如果环境变量设置了深圳节点，强制替换为北京节点并打印警告
	if strings.Contains(url, "shenzhen") {
		log.Printf("[Config] 警告：DOUBAO_BASE_URL 设置了深圳节点 (%s)，已自动切换为北京节点", url)
		return "https://ark.cn-beijing.volces.com/api/v3"
	}
	return url
}

// getEnvInt32 获取 int32 类型的环境变量
func getEnvInt32(key string, defaultValue int32) int32 {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return defaultValue
	}
	return int32(i)
}
