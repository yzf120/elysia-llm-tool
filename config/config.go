package config

import (
	"os"
	"strconv"
)

// Config 全局配置
type Config struct {
	// 豆包配置
	DoubaoAPIKey  string
	DoubaoBaseURL string

	// 混元配置
	HunyuanSecretID  string
	HunyuanSecretKey string

	// DeepSeek配置
	DeepSeekAPIKey  string
	DeepSeekBaseURL string

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
		// 豆包配置
		DoubaoAPIKey:  getEnv("DOUBAO_API_KEY", ""),
		DoubaoBaseURL: getEnv("DOUBAO_BASE_URL", "https://ark.cn-beijing.volces.com/api/v3"),

		// 混元配置
		HunyuanSecretID:  getEnv("HUNYUAN_SECRET_ID", ""),
		HunyuanSecretKey: getEnv("HUNYUAN_SECRET_KEY", ""),

		// DeepSeek配置
		DeepSeekAPIKey:  getEnv("DEEPSEEK_API_KEY", ""),
		DeepSeekBaseURL: getEnv("DEEPSEEK_BASE_URL", "https://api.deepseek.com"),

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
