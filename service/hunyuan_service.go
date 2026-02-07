package service

import (
	"context"
	"fmt"
	"log"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	hunyuan "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan/v20230901"
	"github.com/yzf120/elysia-llm-tool/client"
	"github.com/yzf120/elysia-llm-tool/config"
	"github.com/yzf120/elysia-llm-tool/errs"
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
)

// HunyuanService 混元服务
type HunyuanService struct {
	cfg *config.Config
}

// NewHunyuanService 创建混元服务
func NewHunyuanService() *HunyuanService {
	return &HunyuanService{
		cfg: config.GetConfig(),
	}
}

// StreamChat 混元流式对话
func (s *HunyuanService) StreamChat(ctx context.Context, req *pb.StreamChatRequest, stream pb.LLMService_StreamChatServer) error {
	client := client.GetHunyuanClient()

	// 转换消息格式
	messages := s.convertToHunyuanMessages(req.Messages)

	// 实例化一个请求对象
	request := hunyuan.NewChatCompletionsRequest()
	request.Model = common.StringPtr(req.ModelId)
	request.Messages = messages

	// 设置流式响应
	request.Stream = common.BoolPtr(true)

	// 从配置中设置默认参数
	temperature := float64(s.cfg.DefaultTemperature)
	request.Temperature = common.Float64Ptr(temperature)

	topP := float64(s.cfg.DefaultTopP)
	request.TopP = common.Float64Ptr(topP)

	// 发起请求
	response, err := client.ChatCompletions(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		log.Printf("腾讯云 API 错误: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}
	if err != nil {
		log.Printf("创建混元流式请求失败: %v", err)
		return fmt.Errorf("[%d]%s", errs.ErrModelRequestFailed, err.Error())
	}

	// 处理响应
	if response.Response != nil {
		// 非流式响应（不应该发生，因为我们设置了 Stream=true）
		log.Printf("收到非流式响应")
		resp := s.convertHunyuanResponse(response)
		if err := stream.Send(resp); err != nil {
			log.Printf("发送响应失败: %v", err)
			return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
		}
		// 发送结束标记
		endResp := &pb.StreamChatResponse{
			IsEnd: true,
		}
		return stream.Send(endResp)
	} else if response.Events != nil {
		// 流式响应
		for event := range response.Events {
			if event.Data != nil {
				resp := s.convertHunyuanStreamEvent(event.Data)
				if err := stream.Send(resp); err != nil {
					log.Printf("发送流式响应失败: %v", err)
					return fmt.Errorf("[%d]%s", errs.ErrModelStreamFailed, err.Error())
				}
			}
		}
		// 发送结束标记
		endResp := &pb.StreamChatResponse{
			IsEnd: true,
		}
		return stream.Send(endResp)
	} else {
		// 没有响应数据
		log.Printf("没有收到响应数据")
		return fmt.Errorf("[%d]没有收到响应数据", errs.ErrModelRequestFailed)
	}
}

// convertToHunyuanMessages 转换消息格式到混元格式
func (s *HunyuanService) convertToHunyuanMessages(messages []*pb.ChatMessage) []*hunyuan.Message {
	result := make([]*hunyuan.Message, 0, len(messages))

	for _, msg := range messages {
		hunyuanMsg := &hunyuan.Message{
			Role: common.StringPtr(msg.Role),
		}

		// 处理内容
		if len(msg.Content) == 1 && msg.Content[0].Type == "text" {
			// 纯文本消息
			hunyuanMsg.Content = common.StringPtr(msg.Content[0].Text)
		} else {
			// 多模态消息
			contents := make([]*hunyuan.Content, 0, len(msg.Content))
			for _, part := range msg.Content {
				hunyuanContent := &hunyuan.Content{
					Type: common.StringPtr(part.Type),
				}
				if part.Type == "text" {
					hunyuanContent.Text = common.StringPtr(part.Text)
				} else if part.Type == "image_url" && part.ImageUrl != nil {
					hunyuanContent.ImageUrl = &hunyuan.ImageUrl{
						Url: common.StringPtr(part.ImageUrl.Url),
					}
				}
				contents = append(contents, hunyuanContent)
			}
			hunyuanMsg.Contents = contents
		}

		result = append(result, hunyuanMsg)
	}

	return result
}

// convertHunyuanResponse 转换混元非流式响应
func (s *HunyuanService) convertHunyuanResponse(resp *hunyuan.ChatCompletionsResponse) *pb.StreamChatResponse {
	result := &pb.StreamChatResponse{
		Id:      *resp.Response.Id,
		Model:   "",
		Created: *resp.Response.Created,
		IsEnd:   false,
	}

	// 转换choices
	if len(resp.Response.Choices) > 0 {
		choices := make([]*pb.Choice, 0, len(resp.Response.Choices))
		for _, choice := range resp.Response.Choices {
			c := &pb.Choice{
				Index: int32(*choice.Index),
				Delta: &pb.Delta{
					Role:    *choice.Message.Role,
					Content: *choice.Message.Content,
				},
				FinishReason: *choice.FinishReason,
			}
			choices = append(choices, c)
		}
		result.Choices = choices
	}

	// 转换usage
	if resp.Response.Usage != nil {
		result.Usage = &pb.Usage{
			PromptTokens:     int32(*resp.Response.Usage.PromptTokens),
			CompletionTokens: int32(*resp.Response.Usage.CompletionTokens),
			TotalTokens:      int32(*resp.Response.Usage.TotalTokens),
		}
	}

	return result
}

// convertHunyuanStreamEvent 转换混元流式事件
func (s *HunyuanService) convertHunyuanStreamEvent(event []byte) *pb.StreamChatResponse {
	// 解析事件数据
	// 注意：这里需要根据实际的事件数据格式进行解析
	// 腾讯云的流式响应格式是 SSE (Server-Sent Events)
	result := &pb.StreamChatResponse{
		Id:      "",
		Model:   "",
		Created: 0,
		IsEnd:   false,
	}

	// 简单处理：将事件数据作为内容返回
	if event != nil {
		result.Choices = []*pb.Choice{
			{
				Index: 0,
				Delta: &pb.Delta{
					Role:    "assistant",
					Content: string(event),
				},
				FinishReason: "",
			},
		}
	}

	return result
}
