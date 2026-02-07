[//]: # (# LLM Tool - 统一大模型接入服务)

[//]: # ()
[//]: # (## 项目简介)

[//]: # ()
[//]: # (llm-tool 是一个统一的大模型能力接入服务，为上层服务提供统一的 StreamChat 方法（通过 tRPC 的 RPC 范式调用）。)

[//]: # ()
[//]: # (目前支持的模型提供商：)

[//]: # (- ✅ **豆包（Doubao）** - 字节跳动火山引擎)

[//]: # (- ✅ **混元（Hunyuan）** - 腾讯云)

[//]: # (- ✅ **通义千问（Qwen）** - 阿里云)

[//]: # ()
[//]: # (## 项目架构)

[//]: # ()
[//]: # (```)

[//]: # (llm-tool/)

[//]: # (├── client/              # 第三方 SDK 客户端封装)

[//]: # (│   ├── doubao_client.go    # 豆包客户端)

[//]: # (│   ├── hunyuan_client.go   # 混元客户端)

[//]: # (│   └── qwen_client.go      # 通义千问客户端)

[//]: # (├── service/             # 业务逻辑层)

[//]: # (│   ├── llm_service.go      # 通用服务（路由、模型列表）)

[//]: # (│   ├── doubao_service.go   # 豆包业务逻辑)

[//]: # (│   ├── hunyuan_service.go  # 混元业务逻辑)

[//]: # (│   └── qwen_service.go     # 通义千问业务逻辑)

[//]: # (├── service_impl/        # RPC 实现层)

[//]: # (│   └── llm_service_impl.go # RPC 接口实现)

[//]: # (├── proto/               # Protocol Buffers 定义)

[//]: # (│   └── llm/)

[//]: # (│       └── llm.proto       # 服务接口定义)

[//]: # (├── config/              # 配置管理)

[//]: # (│   └── config.go)

[//]: # (├── consts/              # 常量定义)

[//]: # (│   └── base.go)

[//]: # (├── errs/                # 错误处理)

[//]: # (│   └── errs.go)

[//]: # (└── main.go              # 服务入口)

[//]: # (```)

[//]: # ()
[//]: # (### 分层说明)

[//]: # ()
[//]: # (1. **service_impl 层**：接收 RPC 请求，根据 modelId 路由到对应的 service)

[//]: # (2. **service 层**：处理业务逻辑，调用第三方 SDK，进行参数转换和错误处理)

[//]: # (3. **client 层**：封装第三方 SDK 客户端，提供单例模式)

[//]: # ()
[//]: # (## 快速开始)

[//]: # ()
[//]: # (### 1. 环境准备)

[//]: # ()
[//]: # (```bash)

[//]: # (# 安装 Go 1.21+)

[//]: # (go version)

[//]: # ()
[//]: # (# 克隆项目)

[//]: # (cd elysia-llm-tool)

[//]: # ()
[//]: # (# 安装依赖)

[//]: # (go mod tidy)

[//]: # (```)

[//]: # ()
[//]: # (### 2. 配置环境变量)

[//]: # ()
[//]: # (复制 `.env.example` 为 `.env` 并填写配置：)

[//]: # ()
[//]: # (```bash)

[//]: # (cp .env.example .env)

[//]: # (```)

[//]: # ()
[//]: # (编辑 `.env` 文件：)

[//]: # ()
[//]: # (```bash)

[//]: # (# 豆包配置)

[//]: # (DOUBAO_API_KEY=your_doubao_api_key_here)

[//]: # (DOUBAO_BASE_URL=https://ark.cn-beijing.volces.com/api/v3)

[//]: # ()
[//]: # (# 腾讯云混元配置)

[//]: # (TENCENTCLOUD_SECRET_ID=your_tencentcloud_secret_id_here)

[//]: # (TENCENTCLOUD_SECRET_KEY=your_tencentcloud_secret_key_here)

[//]: # ()
[//]: # (# 阿里通义千问配置)

[//]: # (DASHSCOPE_API_KEY=your_dashscope_api_key_here)

[//]: # (QWEN_BASE_URL=https://dashscope.aliyuncs.com/compatible-mode/v1)

[//]: # ()
[//]: # (# 模型默认参数配置)

[//]: # (DEFAULT_TEMPERATURE=0.7)

[//]: # (DEFAULT_MAX_TOKENS=2000)

[//]: # (DEFAULT_TOP_P=0.9)

[//]: # (```)

[//]: # ()
[//]: # (### 3. 生成 Proto 代码)

[//]: # ()
[//]: # (```bash)

[//]: # (cd proto)

[//]: # (make)

[//]: # (```)

[//]: # ()
[//]: # (### 4. 编译运行)

[//]: # ()
[//]: # (```bash)

[//]: # (# 编译)

[//]: # (go build -o llm-tool .)

[//]: # ()
[//]: # (# 运行)

[//]: # (./llm-tool)

[//]: # (```)

[//]: # ()
[//]: # (或者使用 Makefile：)

[//]: # ()
[//]: # (```bash)

[//]: # (make run)

[//]: # (```)

[//]: # ()
[//]: # (## API 使用)

[//]: # ()
[//]: # (### 1. 流式对话（StreamChat）)

[//]: # ()
[//]: # (**请求参数：**)

[//]: # ()
[//]: # (```protobuf)

[//]: # (message StreamChatRequest {)

[//]: # (  string model_id = 1;                    // 模型ID（必填）)

[//]: # (  repeated ChatMessage messages = 2;       // 对话消息列表)

[//]: # (})

[//]: # ()
[//]: # (message ChatMessage {)

[//]: # (  string role = 1;                         // 角色：user/assistant/system)

[//]: # (  repeated ContentPart content = 2;        // 消息内容（支持多模态）)

[//]: # (})

[//]: # ()
[//]: # (message ContentPart {)

[//]: # (  string type = 1;                         // 类型：text/image_url)

[//]: # (  string text = 2;                         // 文本内容)

[//]: # (  ImageUrl image_url = 3;                  // 图片URL)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**响应：**)

[//]: # ()
[//]: # (```protobuf)

[//]: # (message StreamChatResponse {)

[//]: # (  string id = 1;                           // 响应ID)

[//]: # (  string model = 2;                        // 模型名称)

[//]: # (  int64 created = 3;                       // 创建时间戳)

[//]: # (  repeated Choice choices = 4;             // 选择列表)

[//]: # (  Usage usage = 5;                         // Token使用情况)

[//]: # (  bool is_end = 6;                         // 是否结束)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (### 2. 获取模型列表（ListModels）)

[//]: # ()
[//]: # (**请求参数：**)

[//]: # ()
[//]: # (```protobuf)

[//]: # (message ListModelsRequest {)

[//]: # (  string provider = 1;  // 提供商过滤（可选）：doubao/hunyuan/qwen)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**响应：**)

[//]: # ()
[//]: # (```protobuf)

[//]: # (message ListModelsResponse {)

[//]: # (  repeated ModelInfo models = 1;)

[//]: # (})

[//]: # ()
[//]: # (message ModelInfo {)

[//]: # (  string model_id = 1;        // 模型ID)

[//]: # (  string model_name = 2;      // 模型名称)

[//]: # (  string provider = 3;        // 提供商)

[//]: # (  string description = 4;     // 描述)

[//]: # (  bool support_stream = 5;    // 是否支持流式)

[//]: # (  bool support_vision = 6;    // 是否支持视觉)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (## 支持的模型)

[//]: # ()
[//]: # (### 豆包（Doubao）)

[//]: # ()
[//]: # (| 模型ID | 模型名称 | 描述 | 流式 | 多模态 |)

[//]: # (|--------|---------|------|------|--------|)

[//]: # (| doubao-seed-1-6-lite-251015 | 豆包 Seed 1.6 Lite | 轻量级模型，适合快速响应 | ✅ | ✅ |)

[//]: # ()
[//]: # (### 混元（Hunyuan）)

[//]: # ()
[//]: # (| 模型ID | 模型名称 | 描述 | 流式 | 多模态 |)

[//]: # (|--------|---------|------|------|--------|)

[//]: # (| hunyuan-turbos-latest | 混元 Turbos Latest | 腾讯混元最新模型 | ✅ | ✅ |)

[//]: # ()
[//]: # (### 通义千问（Qwen）)

[//]: # ()
[//]: # (| 模型ID | 模型名称 | 描述 | 流式 | 多模态 |)

[//]: # (|--------|---------|------|------|--------|)

[//]: # (| qwen3-omni-flash | 通义千问 Qwen3 Omni Flash | 阿里通义千问最新模型 | ✅ | ✅ |)

[//]: # ()
[//]: # (## 开发指南)

[//]: # ()
[//]: # (### 添加新的模型提供商)

[//]: # ()
[//]: # (1. **创建 client**：在 `client/` 目录下创建新的客户端文件)

[//]: # (2. **创建 service**：在 `service/` 目录下创建新的服务文件)

[//]: # (3. **更新路由**：在 `service/llm_service.go` 中更新 `GetProviderFromModelID` 方法)

[//]: # (4. **更新 impl**：在 `service_impl/llm_service_impl.go` 中添加新的 case)

[//]: # (5. **更新模型列表**：在 `service/llm_service.go` 的 `ListModels` 方法中添加模型信息)

[//]: # ()
[//]: # (### 示例：添加新模型)

[//]: # ()
[//]: # (```go)

[//]: # (// 1. client/newmodel_client.go)

[//]: # (func GetNewModelClient&#40;&#41; *newmodel.Client {)

[//]: # (    // 实现客户端单例)

[//]: # (})

[//]: # ()
[//]: # (// 2. service/newmodel_service.go)

[//]: # (type NewModelService struct {)

[//]: # (    cfg *config.Config)

[//]: # (})

[//]: # ()
[//]: # (func &#40;s *NewModelService&#41; StreamChat&#40;ctx context.Context, req *pb.StreamChatRequest, stream pb.LLMService_StreamChatServer&#41; error {)

[//]: # (    // 实现流式对话逻辑)

[//]: # (})

[//]: # ()
[//]: # (// 3. service/llm_service.go)

[//]: # (func &#40;s *LLMService&#41; GetNewModelService&#40;&#41; *NewModelService {)

[//]: # (    return s.newmodelService)

[//]: # (})

[//]: # ()
[//]: # (// 4. service_impl/llm_service_impl.go)

[//]: # (case "newmodel":)

[//]: # (    return s.llmService.GetNewModelService&#40;&#41;.StreamChat&#40;ctx, req, stream&#41;)

[//]: # (```)

[//]: # ()
[//]: # (## 配置说明)

[//]: # ()
[//]: # (### 环境变量)

[//]: # ()
[//]: # (| 变量名 | 说明 | 默认值 |)

[//]: # (|--------|------|--------|)

[//]: # (| DOUBAO_API_KEY | 豆包 API Key | - |)

[//]: # (| DOUBAO_BASE_URL | 豆包 API 地址 | https://ark.cn-beijing.volces.com/api/v3 |)

[//]: # (| TENCENTCLOUD_SECRET_ID | 腾讯云 Secret ID | - |)

[//]: # (| TENCENTCLOUD_SECRET_KEY | 腾讯云 Secret Key | - |)

[//]: # (| DASHSCOPE_API_KEY | 阿里云 DashScope API Key | - |)

[//]: # (| QWEN_BASE_URL | 通义千问 API 地址 | https://dashscope.aliyuncs.com/compatible-mode/v1 |)

[//]: # (| DEFAULT_TEMPERATURE | 默认温度参数 | 0.7 |)

[//]: # (| DEFAULT_MAX_TOKENS | 默认最大 Token 数 | 2000 |)

[//]: # (| DEFAULT_TOP_P | 默认 Top P 参数 | 0.9 |)

[//]: # ()
[//]: # (### tRPC 配置)

[//]: # ()
[//]: # (服务配置在 `trpc_go.yaml` 中，默认监听端口：)

[//]: # (- RPC 端口：8000)

[//]: # (- HTTP 端口：8080)

[//]: # ()
[//]: # (## 错误处理)

[//]: # ()
[//]: # (错误码定义在 `consts/base.go` 和 `errs/errs.go` 中：)

[//]: # ()
[//]: # (| 错误码 | 说明 |)

[//]: # (|--------|------|)

[//]: # (| 10001 | 模型ID无效 |)

[//]: # (| 10002 | 消息列表为空 |)

[//]: # (| 10003 | 模型提供商不支持 |)

[//]: # (| 10004 | 模型请求失败 |)

[//]: # (| 10005 | 模型流式响应失败 |)

[//]: # ()
[//]: # (## 测试)

[//]: # ()
[//]: # (```bash)

[//]: # (# 运行测试)

[//]: # (go test ./...)

[//]: # ()
[//]: # (# 运行特定包的测试)

[//]: # (go test ./service/...)

[//]: # (```)

[//]: # ()
[//]: # (## 部署)

[//]: # ()
[//]: # (### Docker 部署)

[//]: # ()
[//]: # (```bash)

[//]: # (# 构建镜像)

[//]: # (docker build -t llm-tool:latest .)

[//]: # ()
[//]: # (# 运行容器)

[//]: # (docker run -d \)

[//]: # (  --name llm-tool \)

[//]: # (  -p 8000:8000 \)

[//]: # (  -p 8080:8080 \)

[//]: # (  --env-file .env \)

[//]: # (  llm-tool:latest)

[//]: # (```)

[//]: # ()
[//]: # (### Kubernetes 部署)

[//]: # ()
[//]: # (```bash)

[//]: # (# 应用配置)

[//]: # (kubectl apply -f k8s/deployment.yaml)

[//]: # (kubectl apply -f k8s/service.yaml)

[//]: # (```)

[//]: # ()
[//]: # (## 使用示例)

[//]: # ()
[//]: # (### 豆包模型)

[//]: # ()
[//]: # (```go)

[//]: # (req := &pb.StreamChatRequest{)

[//]: # (    ModelId: "doubao-seed-1-6-lite-251015",)

[//]: # (    Messages: []*pb.ChatMessage{)

[//]: # (        {)

[//]: # (            Role: "user",)

[//]: # (            Content: []*pb.ContentPart{)

[//]: # (                {Type: "text", Text: "你好"},)

[//]: # (            },)

[//]: # (        },)

[//]: # (    },)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (### 混元模型)

[//]: # ()
[//]: # (```go)

[//]: # (req := &pb.StreamChatRequest{)

[//]: # (    ModelId: "hunyuan-turbos-latest",)

[//]: # (    Messages: []*pb.ChatMessage{)

[//]: # (        {)

[//]: # (            Role: "user",)

[//]: # (            Content: []*pb.ContentPart{)

[//]: # (                {Type: "text", Text: "你好"},)

[//]: # (            },)

[//]: # (        },)

[//]: # (    },)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (### 通义千问模型)

[//]: # ()
[//]: # (```go)

[//]: # (req := &pb.StreamChatRequest{)

[//]: # (    ModelId: "qwen3-omni-flash",)

[//]: # (    Messages: []*pb.ChatMessage{)

[//]: # (        {)

[//]: # (            Role: "user",)

[//]: # (            Content: []*pb.ContentPart{)

[//]: # (                {Type: "text", Text: "你好"},)

[//]: # (            },)

[//]: # (        },)

[//]: # (    },)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (## 许可证)

[//]: # ()
[//]: # (MIT License)

[//]: # ()
[//]: # (## 贡献)

[//]: # ()
[//]: # (欢迎提交 Issue 和 Pull Request！)
