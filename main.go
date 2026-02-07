package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/yzf120/elysia-llm-tool/config"
	pb "github.com/yzf120/elysia-llm-tool/proto/llm"
	"github.com/yzf120/elysia-llm-tool/service_impl"
	"trpc.group/trpc-go/trpc-go"
)

func main() {
	// 加载环境变量
	err := godotenv.Load()
	if err != nil {
		log.Println("未找到.env文件，使用系统环境变量")
	}

	// 初始化配置
	cfg := config.InitConfig()
	log.Println("配置初始化完成")

	// 验证必要的配置
	if cfg.DoubaoAPIKey == "" {
		log.Println("警告: DOUBAO_API_KEY 未设置，豆包模型将无法使用")
	}
	
	// 打印模型默认参数配置
	log.Printf("模型默认参数: Temperature=%.2f, MaxTokens=%d, TopP=%.2f, ReasoningEffort=%s",
		cfg.DefaultTemperature, cfg.DefaultMaxTokens, cfg.DefaultTopP, cfg.DefaultReasoningEffort)

	// 创建trpc服务器
	s := trpc.NewServer()

	// 注册LLM服务（使用service_impl层）
	pb.RegisterLLMServiceService(s, service_impl.NewLLMServiceImpl())

	log.Println("LLM Tool 服务启动成功！")
	log.Printf("服务地址: 127.0.0.1:%d", cfg.ServerPort)
	log.Println("支持的接口:")
	log.Println("  - StreamChat: 流式对话")
	log.Println("  - ListModels: 获取模型列表")

	// 启动服务器
	if err := s.Serve(); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
