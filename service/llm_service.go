package service

import (
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

// LLMService LLM服务层（通用）
type LLMService struct {
	doubaoService  *DoubaoService
	hunyuanService *HunyuanService
	qwenService    *QwenService
}

// NewLLMService 创建LLM服务
func NewLLMService() *LLMService {
	return &LLMService{
		doubaoService:  NewDoubaoService(),
		hunyuanService: NewHunyuanService(),
		qwenService:    NewQwenService(),
	}
}

// GetDoubaoService 获取豆包服务
func (s *LLMService) GetDoubaoService() *DoubaoService {
	return s.doubaoService
}

// GetHunyuanService 获取混元服务
func (s *LLMService) GetHunyuanService() *HunyuanService {
	return s.hunyuanService
}

// GetQwenService 获取通义千问服务
func (s *LLMService) GetQwenService() *QwenService {
	return s.qwenService
}

// ListModels 获取支持的模型列表
func (s *LLMService) ListModels(provider string) ([]*pb.ModelInfo, error) {
	models := []*pb.ModelInfo{
		// 豆包模型
		{
			ModelId:       "doubao-seed-1-6-lite-251015",
			ModelName:     "豆包 Seed 1.6 Lite",
			Provider:      "doubao",
			Description:   "豆包轻量级模型，适合快速响应场景",
			SupportStream: true,
			SupportVision: true,
		},
		// 混元模型
		{
			ModelId:       "hunyuan-turbos-latest",
			ModelName:     "混元 Turbos Latest",
			Provider:      "hunyuan",
			Description:   "腾讯混元最新模型，支持多模态",
			SupportStream: true,
			SupportVision: true,
		},
		// 通义千问模型
		{
			ModelId:       "qwen3-omni-flash",
			ModelName:     "通义千问 Qwen3 Omni Flash",
			Provider:      "qwen",
			Description:   "阿里通义千问最新模型，支持多模态",
			SupportStream: true,
			SupportVision: true,
		},
	}
	// 根据provider过滤
	if provider != "" {
		filtered := make([]*pb.ModelInfo, 0)
		for _, m := range models {
			if m.Provider == provider {
				filtered = append(filtered, m)
			}
		}
		return filtered, nil
	}

	return models, nil
}

// GetProviderFromModelID 从模型ID获取提供商
func GetProviderFromModelID(modelID string) string {
	// 简单的前缀匹配
	if len(modelID) >= 6 && modelID[:6] == "doubao" {
		return "doubao"
	}
	if len(modelID) >= 7 && modelID[:7] == "hunyuan" {
		return "hunyuan"
	}
	if len(modelID) >= 4 && modelID[:4] == "qwen" {
		return "qwen"
	}
	// 默认使用豆包
	return "doubao"
}
