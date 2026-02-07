package service

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/yzf120/elysia-llm-tool/client"
	"github.com/yzf120/elysia-llm-tool/config"
	"github.com/yzf120/elysia-llm-tool/errs"
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

// DoubaoService 豆包服务
type DoubaoService struct {
	cfg *config.Config
}

// NewDoubaoService 创建豆包服务
func NewDoubaoService() *DoubaoService {
	return &DoubaoService{
		cfg: config.GetConfig(),
	}
}

// StreamChat 豆包流式对话
func (s *DoubaoService) StreamChat(ctx context.Context, req *pb.StreamChatRequest, stream pb.LLMService_StreamChatServer) error {
	client := client.GetDoubaoClient()

	// 转换消息格式
	messages := s.convertToDoubaoMessages(req.Messages)

	// 构建请求
	doubaoReq := model.CreateChatCompletionRequest{
		Model:    req.ModelId,
		Messages: messages,
	}

	// 从配置中设置默认参数
	temperature := float32(s.cfg.DefaultTemperature)
	doubaoReq.Temperature = &temperature

	maxTokens := int(s.cfg.DefaultMaxTokens)
	doubaoReq.MaxTokens = &maxTokens

	topP := float32(s.cfg.DefaultTopP)
	doubaoReq.TopP = &topP

	// 创建流式请求
	doubaoStream, err := client.CreateChatCompletionStream(ctx, doubaoReq)
	if err != nil {
		log.Printf("创建豆包流式请求失败: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}
	defer doubaoStream.Close()

	// 读取流式响应
	for {
		recv, err := doubaoStream.Recv()
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
		resp := s.convertDoubaoStreamResponse(&recv)
		if err := stream.Send(resp); err != nil {
			log.Printf("发送响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}
	}
}

// convertToDoubaoMessages 转换消息格式到豆包格式
func (s *DoubaoService) convertToDoubaoMessages(messages []*pb.ChatMessage) []*model.ChatCompletionMessage {
	result := make([]*model.ChatCompletionMessage, 0, len(messages))

	for _, msg := range messages {
		doubaoMsg := &model.ChatCompletionMessage{
			Role: msg.Role,
		}

		// 处理内容
		if len(msg.Content) == 1 && msg.Content[0].Type == "text" {
			// 纯文本消息
			doubaoMsg.Content = &model.ChatCompletionMessageContent{
				StringValue: &msg.Content[0].Text,
			}
		} else {
			// 多模态消息
			parts := make([]*model.ChatCompletionMessageContentPart, 0, len(msg.Content))
			for _, part := range msg.Content {
				doubaPart := &model.ChatCompletionMessageContentPart{
					Type: model.ChatCompletionMessageContentPartType(part.Type),
				}
				if part.Type == "text" {
					doubaPart.Text = part.Text
				} else if part.Type == "image_url" && part.ImageUrl != nil {
					doubaPart.ImageURL = &model.ChatMessageImageURL{
						URL: part.ImageUrl.Url,
					}
				}
				parts = append(parts, doubaPart)
			}
			doubaoMsg.Content = &model.ChatCompletionMessageContent{
				ListValue: parts,
			}
		}

		result = append(result, doubaoMsg)
	}

	return result
}

// convertDoubaoStreamResponse 转换豆包流式响应
func (s *DoubaoService) convertDoubaoStreamResponse(resp *model.ChatCompletionStreamResponse) *pb.StreamChatResponse {
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
					Role:    string(choice.Delta.Role),
					Content: choice.Delta.Content,
				},
				FinishReason: string(choice.FinishReason),
			}
			choices = append(choices, c)
		}
		result.Choices = choices
	}

	// 转换usage
	if resp.Usage != nil {
		result.Usage = &pb.Usage{
			PromptTokens:     int32(resp.Usage.PromptTokens),
			CompletionTokens: int32(resp.Usage.CompletionTokens),
			TotalTokens:      int32(resp.Usage.TotalTokens),
		}
	}

	return result
}
