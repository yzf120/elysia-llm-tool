package service_impl

import (
	"context"
	"log"

	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
	"github.com/yzf120/elysia-llm-tool/service"
)

// LLMServiceImpl LLM服务实现层
type LLMServiceImpl struct {
	llmService *service.LLMService
}

// NewLLMServiceImpl 创建LLM服务实现
func NewLLMServiceImpl() *LLMServiceImpl {
	return &LLMServiceImpl{
		llmService: service.NewLLMService(),
	}
}

// StreamChat 流式对话实现
func (s *LLMServiceImpl) StreamChat(req *pb.StreamChatRequest, stream pb.LLMService_StreamChatServer) error {
	log.Printf("收到流式对话请求，模型: %s", req.ModelId)

	// 根据模型ID判断使用哪个提供商
	provider := service.GetProviderFromModelID(req.ModelId)

	// 使用 stream 的 context
	ctx := stream.Context()

	switch provider {
	case "doubao":
		// 调用豆包 service 处理
		return s.llmService.GetDoubaoService().StreamChat(ctx, req, stream)

	case "hunyuan":
		// 调用混元 service 处理
		return s.llmService.GetHunyuanService().StreamChat(ctx, req, stream)

	case "qwen":
		// 调用通义千问 service 处理
		return s.llmService.GetQwenService().StreamChat(ctx, req, stream)

	default:
		log.Printf("未知的模型提供商: %s，使用默认豆包", provider)
		return s.llmService.GetDoubaoService().StreamChat(ctx, req, stream)
	}
}

// ListModels 获取支持的模型列表
func (s *LLMServiceImpl) ListModels(ctx context.Context, req *pb.ListModelsRequest) (*pb.ListModelsResponse, error) {
	log.Printf("收到获取模型列表请求，提供商: %s", req.Provider)

	// 调用service层处理业务逻辑
	models, err := s.llmService.ListModels(req.Provider)
	if err != nil {
		log.Printf("获取模型列表失败: %v", err)
		return &pb.ListModelsResponse{
			Models: nil,
		}, err
	}

	log.Printf("成功获取 %d 个模型", len(models))
	return &pb.ListModelsResponse{
		Models: models,
	}, nil
}
