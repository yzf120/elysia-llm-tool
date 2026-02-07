package service

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/sashabaranov/go-openai"
	"github.com/yzf120/elysia-llm-tool/client"
	"github.com/yzf120/elysia-llm-tool/config"
	"github.com/yzf120/elysia-llm-tool/errs"
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

// QwenService 通义千问服务
type QwenService struct {
	cfg *config.Config
}

// NewQwenService 创建通义千问服务
func NewQwenService() *QwenService {
	return &QwenService{
		cfg: config.GetConfig(),
	}
}

// StreamChat 通义千问流式对话
func (s *QwenService) StreamChat(ctx context.Context, req *pb.StreamChatRequest, stream pb.LLMService_StreamChatServer) error {
	client := client.GetQwenClient()

	// 转换消息格式
	messages := s.convertToQwenMessages(req.Messages)

	// 构建请求
	qwenReq := openai.ChatCompletionRequest{
		Model:    req.ModelId,
		Messages: messages,
		Stream:   true,
	}

	// 从配置中设置默认参数
	temperature := float32(s.cfg.DefaultTemperature)
	qwenReq.Temperature = temperature

	maxTokens := int(s.cfg.DefaultMaxTokens)
	qwenReq.MaxTokens = maxTokens

	topP := float32(s.cfg.DefaultTopP)
	qwenReq.TopP = topP

	// 创建流式请求
	qwenStream, err := client.CreateChatCompletionStream(ctx, qwenReq)
	if err != nil {
		log.Printf("创建通义千问流式请求失败: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}
	defer qwenStream.Close()

	// 读取流式响应
	for {
		recv, err := qwenStream.Recv()
		if err == io.EOF {
			// 发送结束标记
			endResp := &pb.StreamChatResponse{
				IsEnd: true,
			}
			if err := stream.Send(endResp); err != nil {
				log.Printf("发送结束标记失败: %v", err)
				return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
			}
			return nil
		}
		if err != nil {
			log.Printf("接收流式响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}

		// 转换并发送响应
		resp := s.convertQwenStreamResponse(&recv)
		if err := stream.Send(resp); err != nil {
			log.Printf("发送响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}
	}
}

// convertToQwenMessages 转换消息格式到通义千问格式
func (s *QwenService) convertToQwenMessages(messages []*pb.ChatMessage) []openai.ChatCompletionMessage {
	result := make([]openai.ChatCompletionMessage, 0, len(messages))

	for _, msg := range messages {
		qwenMsg := openai.ChatCompletionMessage{
			Role: msg.Role,
		}

		// 处理内容
		if len(msg.Content) == 1 && msg.Content[0].Type == "text" {
			// 纯文本消息
			qwenMsg.Content = msg.Content[0].Text
		} else {
			// 多模态消息
			parts := make([]openai.ChatMessagePart, 0, len(msg.Content))
			for _, part := range msg.Content {
				if part.Type == "text" {
					parts = append(parts, openai.ChatMessagePart{
						Type: openai.ChatMessagePartTypeText,
						Text: part.Text,
					})
				} else if part.Type == "image_url" && part.ImageUrl != nil {
					parts = append(parts, openai.ChatMessagePart{
						Type: openai.ChatMessagePartTypeImageURL,
						ImageURL: &openai.ChatMessageImageURL{
							URL: part.ImageUrl.Url,
						},
					})
				}
			}
			qwenMsg.MultiContent = parts
		}

		result = append(result, qwenMsg)
	}

	return result
}

// convertQwenStreamResponse 转换通义千问流式响应
func (s *QwenService) convertQwenStreamResponse(resp *openai.ChatCompletionStreamResponse) *pb.StreamChatResponse {
	result := &pb.StreamChatResponse{
		Id:      resp.ID,
		Model:   resp.Model,
		Created: resp.Created,
		IsEnd:   false,
	}

	// 转换choices
	if len(resp.Choices) > 0 {
		choices := make([]*pb.Choice, 0, len(resp.Choices))
		for _, choice := range resp.Choices {
			c := &pb.Choice{
				Index: int32(choice.Index),
				Delta: &pb.Delta{
					Role:    choice.Delta.Role,
					Content: choice.Delta.Content,
				},
				FinishReason: string(choice.FinishReason),
			}
			choices = append(choices, c)
		}
		result.Choices = choices
	}

	// 转换usage（如果有）
	if resp.Usage != nil {
		result.Usage = &pb.Usage{
			PromptTokens:     int32(resp.Usage.PromptTokens),
			CompletionTokens: int32(resp.Usage.CompletionTokens),
			TotalTokens:      int32(resp.Usage.TotalTokens),
		}
	}

	return result
}
