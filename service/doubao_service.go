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
	c := client.GetDoubaoClient()

	messages := s.convertToDoubaoMessages(req.Messages)

	doubaoReq := model.CreateChatCompletionRequest{
		Model:    req.ModelId,
		Messages: messages,
	}

	temperature := float32(s.cfg.DefaultTemperature)
	doubaoReq.Temperature = &temperature
	maxTokens := int(s.cfg.DefaultMaxTokens)
	doubaoReq.MaxTokens = &maxTokens
	topP := float32(s.cfg.DefaultTopP)
	doubaoReq.TopP = &topP

	// 深度思考模式：使用 SDK 的 Thinking 字段（ThinkingTypeEnabled）
	if req.ExtraParams != nil {
		if v, ok := req.ExtraParams["enable_thinking"]; ok && v == "true" {
			doubaoReq.Thinking = &model.Thinking{
				Type: model.ThinkingTypeEnabled,
			}
			log.Printf("[DoubaoService] 深度思考模式已开启，模型: %s", req.ModelId)
		}
	}

	doubaoStream, err := c.CreateChatCompletionStream(ctx, doubaoReq)
	if err != nil {
		log.Printf("创建豆包流式请求失败: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}
	defer doubaoStream.Close()

	// 是否开启了深度思考
	enableThinking := req.ExtraParams != nil && req.ExtraParams["enable_thinking"] == "true"

	for {
		recv, err := doubaoStream.Recv()
		if err == io.EOF {
			endResp := &pb.StreamChatResponse{IsEnd: true}
			if err := stream.Send(endResp); err != nil {
				log.Printf("发送结束标记失败: %v", err)
				return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
			}
			return nil
		}
		if err != nil {
			log.Printf("接收豆包流式响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}

		resp := s.convertDoubaoStreamResponse(&recv, enableThinking)
		if err := stream.Send(resp); err != nil {
			log.Printf("发送豆包响应失败: %v", err)
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

		if len(msg.Content) == 1 && msg.Content[0].Type == "text" {
			doubaoMsg.Content = &model.ChatCompletionMessageContent{
				StringValue: &msg.Content[0].Text,
			}
		} else {
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
func (s *DoubaoService) convertDoubaoStreamResponse(resp *model.ChatCompletionStreamResponse, enableThinking bool) *pb.StreamChatResponse {
	result := &pb.StreamChatResponse{
		Id:      resp.ID,
		Model:   resp.Model,
		Created: resp.Created,
		IsEnd:   false,
	}

	if len(resp.Choices) > 0 {
		choices := make([]*pb.Choice, 0, len(resp.Choices))
		for _, choice := range resp.Choices {
			content := choice.Delta.Content
			// 仅在开启深度思考时，才将 reasoning_content 用 <think> 标签包裹输出
			if enableThinking && choice.Delta.ReasoningContent != nil && *choice.Delta.ReasoningContent != "" {
				content = "<think>" + *choice.Delta.ReasoningContent + "</think>" + content
			}
			c := &pb.Choice{
				Index: int32(choice.Index),
				Delta: &pb.Delta{
					Role:    string(choice.Delta.Role),
					Content: content,
				},
				FinishReason: string(choice.FinishReason),
			}
			choices = append(choices, c)
		}
		result.Choices = choices
	}

	if resp.Usage != nil {
		result.Usage = &pb.Usage{
			PromptTokens:     int32(resp.Usage.PromptTokens),
			CompletionTokens: int32(resp.Usage.CompletionTokens),
			TotalTokens:      int32(resp.Usage.TotalTokens),
		}
	}

	return result
}
