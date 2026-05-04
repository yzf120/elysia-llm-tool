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
	c := client.GetQwenClient()

	messages := s.convertToQwenMessages(req.Messages)

	qwenReq := openai.ChatCompletionRequest{
		Model:    req.ModelId,
		Messages: messages,
		Stream:   true,
	}

	temperature := float32(s.cfg.DefaultTemperature)
	qwenReq.Temperature = temperature
	maxTokens := int(s.cfg.DefaultMaxTokens)
	qwenReq.MaxTokens = maxTokens
	topP := float32(s.cfg.DefaultTopP)
	qwenReq.TopP = topP

	// 处理深度思考模式
	// 千问深度思考通过 ChatTemplateKwargs 传递 enable_thinking: true/false
	if req.ExtraParams != nil {
		if v, ok := req.ExtraParams["enable_thinking"]; ok {
			if v == "true" {
				qwenReq.ChatTemplateKwargs = map[string]any{
					"enable_thinking": true,
				}
				log.Printf("[QwenService] 已开启深度思考模式，模型: %s", req.ModelId)
			} else if v == "false" {
				qwenReq.ChatTemplateKwargs = map[string]any{
					"enable_thinking": false,
				}
				log.Printf("[QwenService] 已禁用深度思考模式（加速），模型: %s", req.ModelId)
			}
		}
	}

	qwenStream, err := c.CreateChatCompletionStream(ctx, qwenReq)
	if err != nil {
		log.Printf("创建通义千问流式请求失败: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}
	defer qwenStream.Close()

	for {
		recv, err := qwenStream.Recv()
		if err == io.EOF {
			endResp := &pb.StreamChatResponse{IsEnd: true}
			if err := stream.Send(endResp); err != nil {
				log.Printf("发送结束标记失败: %v", err)
				return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
			}
			return nil
		}
		if err != nil {
			log.Printf("接收通义千问流式响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}

		// 过滤无意义的空 chunk（content 为空、无 finish_reason、无 usage）
		if s.isEmptyChunk(&recv) {
			continue
		}

		resp := s.convertQwenStreamResponse(&recv)
		if err := stream.Send(resp); err != nil {
			log.Printf("发送通义千问响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}
	}
}

// isEmptyChunk 判断是否为无意义的空 chunk（应跳过不发送）
func (s *QwenService) isEmptyChunk(resp *openai.ChatCompletionStreamResponse) bool {
	if len(resp.Choices) == 0 {
		return resp.Usage == nil
	}

	choice := resp.Choices[0]

	// 有 finish_reason 的 chunk 是有意义的
	if choice.FinishReason != "" {
		return false
	}

	// 有实际内容的 chunk 是有意义的
	if choice.Delta.Content != "" {
		return false
	}

	// 有思考内容的 chunk 是有意义的
	if choice.Delta.ReasoningContent != "" {
		return false
	}

	// 有 usage 信息的 chunk 是有意义的
	if resp.Usage != nil && (resp.Usage.PromptTokens > 0 || resp.Usage.CompletionTokens > 0 || resp.Usage.TotalTokens > 0) {
		return false
	}

	return true
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
			content := choice.Delta.Content
			// 只要有思考内容就用 <think> 标签包裹输出（由模型自行决定是否思考）
			if choice.Delta.ReasoningContent != "" {
				content = "<think>" + choice.Delta.ReasoningContent + "</think>" + content
			}
			c := &pb.Choice{
				Index: int32(choice.Index),
				Delta: &pb.Delta{
					Role:    choice.Delta.Role,
					Content: content,
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
