package service

import (
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

// LLMService LLM服务层（通用）
type LLMService struct {
	doubaoService         *DoubaoService
	qwenService           *QwenService
	inferenceUsageService *InferenceUsageService
}

// NewLLMService 创建LLM服务
func NewLLMService() *LLMService {
	return &LLMService{
		doubaoService:         NewDoubaoService(),
		qwenService:           NewQwenService(),
		inferenceUsageService: NewInferenceUsageService(),
	}
}

// GetDoubaoService 获取豆包服务
func (s *LLMService) GetDoubaoService() *DoubaoService {
	return s.doubaoService
}

// GetQwenService 获取通义千问服务
func (s *LLMService) GetQwenService() *QwenService {
	return s.qwenService
}

// GetInferenceUsageService 获取推理用量查询服务
func (s *LLMService) GetInferenceUsageService() *InferenceUsageService {
	return s.inferenceUsageService
}

// ListModels 获取支持的模型列表
func (s *LLMService) ListModels(provider string) ([]*pb.ModelInfo, error) {
	models := []*pb.ModelInfo{
		// 豆包模型
		{
			ModelId:       "doubao-seed-2-0-lite-260215",
			ModelName:     "Doubao-Seed-2.0-lite",
			Provider:      "doubao",
			Description:   "多模态模型，支持深度思考，适合快速响应场景",
			SupportStream: true,
			SupportVision: true,
		},
		// 通义千问模型
		{
			ModelId:       "qwen3-omni-flash",
			ModelName:     "Qwen3-Omni-Flash",
			Provider:      "qwen",
			Description:   "全模态模型，Thinker–Talker 架构，支持深度思考",
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
	if len(modelID) >= 4 && modelID[:4] == "qwen" {
		return "qwen"
	}
	// 默认使用豆包
	return "doubao"
}
