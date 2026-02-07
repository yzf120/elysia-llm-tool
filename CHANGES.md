[//]: # (# 修改说明文档)

[//]: # ()
[//]: # (## 修改概览)

[//]: # ()
[//]: # (根据需求，对 elysia-llm-tool 项目进行了以下两个主要修改：)

[//]: # ()
[//]: # (### 1. Proto Makefile 格式统一)

[//]: # ()
[//]: # (**问题：** llm-tool 的 proto Makefile 命令格式与 backend 不一致)

[//]: # ()
[//]: # (**解决方案：**)

[//]: # (- 参考 backend 的 proto Makefile 格式进行统一)

[//]: # (- 使用 `generate-all` 目标)

[//]: # (- 添加 `--nogomod` 和 `--mock=false` 参数)

[//]: # (- 统一清理命令格式)

[//]: # ()
[//]: # (**修改文件：** `proto/Makefile`)

[//]: # ()
[//]: # (**修改后的格式：**)

[//]: # (```makefile)

[//]: # (all: generate-all)

[//]: # ()
[//]: # (# 生成所有proto文件)

[//]: # (generate-all:)

[//]: # (	trpc create \)

[//]: # (		-p llm/llm.proto \)

[//]: # (		-o llm \)

[//]: # (		--rpconly \)

[//]: # (		--nogomod \)

[//]: # (		--mock=false)

[//]: # ()
[//]: # (# 清理生成的代码)

[//]: # (clean:)

[//]: # (	rm -rf */*.pb.go)

[//]: # (	rm -rf */*.trpc.go)

[//]: # ()
[//]: # (.PHONY: all generate-all clean)

[//]: # (```)

[//]: # ()
[//]: # (### 2. 移除非流式 Chat 接口)

[//]: # ()
[//]: # (**问题：** 不需要非流式的 chat 接口)

[//]: # ()
[//]: # (**解决方案：**)

[//]: # (- 从 proto 文件中删除 `Chat` RPC 接口)

[//]: # (- 删除相关的 `ChatRequest` 和 `ChatResponse` 消息定义)

[//]: # (- 删除 `CompletionChoice` 消息定义)

[//]: # (- 更新所有相关代码)

[//]: # ()
[//]: # (**修改文件：**)

[//]: # (- `proto/llm/llm.proto`)

[//]: # (- `service/llm_service.go`（删除旧文件，创建新文件）)

[//]: # (- `service_impl/llm_service_impl.go`（新建）)

[//]: # (- `main.go`)

[//]: # (- `example/client_example.go`)

[//]: # ()
[//]: # (### 3. 模型参数配置化)

[//]: # ()
[//]: # (**问题：** 流式 chat 接口的参数（temperature, max_tokens, top_p, reasoning_effort）写死在请求中，需要改为从 env 配置读取)

[//]: # ()
[//]: # (**解决方案：**)

[//]: # ()
[//]: # (#### 3.1 简化 Proto 定义)

[//]: # ()
[//]: # (**修改前：**)

[//]: # (```protobuf)

[//]: # (message StreamChatRequest {)

[//]: # (  string model_id = 1;)

[//]: # (  repeated ChatMessage messages = 2;)

[//]: # (  float temperature = 3;)

[//]: # (  int32 max_tokens = 4;)

[//]: # (  float top_p = 5;)

[//]: # (  string reasoning_effort = 6;)

[//]: # (  map<string, string> extra_params = 7;)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**修改后：**)

[//]: # (```protobuf)

[//]: # (message StreamChatRequest {)

[//]: # (  string model_id = 1;)

[//]: # (  repeated ChatMessage messages = 2;)

[//]: # (  map<string, string> extra_params = 3;)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (#### 3.2 配置管理)

[//]: # ()
[//]: # (**新增配置项：** `config/config.go`)

[//]: # (```go)

[//]: # (type Config struct {)

[//]: # (    // ... 其他配置)

[//]: # (    )
[//]: # (    // 模型默认参数配置)

[//]: # (    DefaultTemperature     float32)

[//]: # (    DefaultMaxTokens       int32)

[//]: # (    DefaultTopP            float32)

[//]: # (    DefaultReasoningEffort string)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**环境变量：** `.env.example`)

[//]: # (```env)

[//]: # (# 模型默认参数配置（可选，不配置则使用代码中的默认值）)

[//]: # (DEFAULT_TEMPERATURE=0.7)

[//]: # (DEFAULT_MAX_TOKENS=2000)

[//]: # (DEFAULT_TOP_P=0.9)

[//]: # (DEFAULT_REASONING_EFFORT=medium)

[//]: # (```)

[//]: # ()
[//]: # (#### 3.3 Service 层使用配置)

[//]: # ()
[//]: # (在 `service/llm_service.go` 中，从配置读取参数：)

[//]: # (```go)

[//]: # (// 从配置中设置默认参数)

[//]: # (temperature := s.cfg.DefaultTemperature)

[//]: # (doubaoReq.Temperature = &temperature)

[//]: # ()
[//]: # (maxTokens := s.cfg.DefaultMaxTokens)

[//]: # (doubaoReq.MaxTokens = &maxTokens)

[//]: # ()
[//]: # (topP := s.cfg.DefaultTopP)

[//]: # (doubaoReq.TopP = &topP)

[//]: # ()
[//]: # (doubaoReq.ReasoningEffort = s.cfg.DefaultReasoningEffort)

[//]: # (```)

[//]: # ()
[//]: # (### 4. 新增 Impl 层)

[//]: # ()
[//]: # (**问题：** llm-tool 目前没有 impl 层)

[//]: # ()
[//]: # (**解决方案：**)

[//]: # (- 参考 backend 的 service_impl 层结构)

[//]: # (- 创建 `service_impl/llm_service_impl.go`)

[//]: # (- impl 层负责：)

[//]: # (  - RPC 接口实现)

[//]: # (  - 参数校验)

[//]: # (  - 错误处理)

[//]: # (  - 调用 service 层)

[//]: # (  - 返回响应给 RPC 调用方)

[//]: # ()
[//]: # (**新增文件：** `service_impl/llm_service_impl.go`)

[//]: # ()
[//]: # (**职责划分：**)

[//]: # (- **service_impl 层**：处理 RPC 接口、参数校验、错误处理)

[//]: # (- **service 层**：处理业务逻辑、调用第三方 SDK)

[//]: # (- **client 层**：管理第三方 SDK 客户端)

[//]: # ()
[//]: # (### 5. 新增错误处理包)

[//]: # ()
[//]: # (**新增文件：**)

[//]: # (- `consts/base.go` - 常量定义)

[//]: # (- `errs/errs.go` - 错误码和错误处理)

[//]: # ()
[//]: # (**错误码定义：**)

[//]: # (```go)

[//]: # (const &#40;)

[//]: # (    // 通用错误)

[//]: # (    ErrInternalServer = 500)

[//]: # (    ErrBadRequest     = 400)

[//]: # (    )
[//]: # (    // LLM相关错误 &#40;30000+&#41;)

[//]: # (    ErrModelNotFound      = 30001)

[//]: # (    ErrModelRequestFailed = 30002)

[//]: # (    ErrModelTimeout       = 30003)

[//]: # (    ErrModelStreamFailed  = 30004)

[//]: # (    ErrInvalidModelID     = 30005)

[//]: # (    ErrProviderNotSupport = 30006)

[//]: # (    )
[//]: # (    // 参数错误 &#40;40000+&#41;)

[//]: # (    ErrParamMissing  = 40001)

[//]: # (    ErrParamInvalid  = 40002)

[//]: # (    ErrMessageEmpty  = 40003)

[//]: # (    ErrMessageFormat = 40004)

[//]: # (    )
[//]: # (    // 配置错误 &#40;50000+&#41;)

[//]: # (    ErrConfigMissing = 50001)

[//]: # (    ErrAPIKeyMissing = 50002)

[//]: # (&#41;)

[//]: # (```)

[//]: # ()
[//]: # (## 项目结构变化)

[//]: # ()
[//]: # (### 修改前)

[//]: # (```)

[//]: # (elysia-llm-tool/)

[//]: # (├── client/)

[//]: # (├── config/)

[//]: # (├── proto/)

[//]: # (├── service/)

[//]: # (│   └── llm_service_impl.go  # 混合了 impl 和 service 逻辑)

[//]: # (├── main.go)

[//]: # (└── ...)

[//]: # (```)

[//]: # ()
[//]: # (### 修改后)

[//]: # (```)

[//]: # (elysia-llm-tool/)

[//]: # (├── client/              # 客户端管理层)

[//]: # (├── config/              # 配置管理)

[//]: # (├── consts/              # 常量定义（新增）)

[//]: # (│   └── base.go)

[//]: # (├── errs/                # 错误处理（新增）)

[//]: # (│   └── errs.go)

[//]: # (├── proto/               # Proto 定义)

[//]: # (├── service/             # 业务逻辑层（重构）)

[//]: # (│   └── llm_service.go)

[//]: # (├── service_impl/        # RPC 实现层（新增）)

[//]: # (│   └── llm_service_impl.go)

[//]: # (├── main.go)

[//]: # (└── ...)

[//]: # (```)

[//]: # ()
[//]: # (## 使用方式变化)

[//]: # ()
[//]: # (### 1. 环境变量配置)

[//]: # ()
[//]: # (**新增配置项：**)

[//]: # (```bash)

[//]: # (# 在 .env 文件中添加)

[//]: # (DEFAULT_TEMPERATURE=0.7)

[//]: # (DEFAULT_MAX_TOKENS=2000)

[//]: # (DEFAULT_TOP_P=0.9)

[//]: # (DEFAULT_REASONING_EFFORT=medium)

[//]: # (```)

[//]: # ()
[//]: # (### 2. 客户端调用)

[//]: # ()
[//]: # (**修改前：**)

[//]: # (```go)

[//]: # (req := &pb.StreamChatRequest{)

[//]: # (    ModelId: "doubao-seed-1-6-lite-251015",)

[//]: # (    Messages: messages,)

[//]: # (    Temperature: 0.7,      // 需要传递)

[//]: # (    MaxTokens: 100,        // 需要传递)

[//]: # (    TopP: 0.9,            // 需要传递)

[//]: # (    ReasoningEffort: "medium", // 需要传递)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (**修改后：**)

[//]: # (```go)

[//]: # (req := &pb.StreamChatRequest{)

[//]: # (    ModelId: "doubao-seed-1-6-lite-251015",)

[//]: # (    Messages: messages,)

[//]: # (    // 参数从服务端配置读取，客户端无需传递)

[//]: # (})

[//]: # (```)

[//]: # ()
[//]: # (### 3. 服务启动)

[//]: # ()
[//]: # (服务启动时会打印配置的默认参数：)

[//]: # (```)

[//]: # (配置初始化完成)

[//]: # (模型默认参数: Temperature=0.70, MaxTokens=2000, TopP=0.90, ReasoningEffort=medium)

[//]: # (LLM Tool 服务启动成功！)

[//]: # (```)

[//]: # ()
[//]: # (## 优势)

[//]: # ()
[//]: # (1. **统一规范**：Proto Makefile 格式与 backend 保持一致)

[//]: # (2. **简化接口**：移除不需要的非流式接口，减少维护成本)

[//]: # (3. **配置灵活**：模型参数可通过环境变量配置，无需修改代码)

[//]: # (4. **分层清晰**：impl 层、service 层、client 层职责明确)

[//]: # (5. **错误规范**：统一的错误码和错误处理机制)

[//]: # (6. **易于扩展**：新增模型提供商只需在 service 层添加实现)

[//]: # ()
[//]: # (## 测试建议)

[//]: # ()
[//]: # (1. **生成 Proto 代码**)

[//]: # (   ```bash)

[//]: # (   cd proto && make && cd ..)

[//]: # (   ```)

[//]: # ()
[//]: # (2. **配置环境变量**)

[//]: # (   ```bash)

[//]: # (   cp .env.example .env)

[//]: # (   # 编辑 .env 文件，填入 API Key 和参数配置)

[//]: # (   ```)

[//]: # ()
[//]: # (3. **启动服务**)

[//]: # (   ```bash)

[//]: # (   make run)

[//]: # (   ```)

[//]: # ()
[//]: # (4. **测试客户端**)

[//]: # (   ```bash)

[//]: # (   cd example)

[//]: # (   go run client_example.go)

[//]: # (   ```)

[//]: # ()
[//]: # (## 注意事项)

[//]: # ()
[//]: # (1. 需要重新生成 proto 代码（因为接口有变化）)

[//]: # (2. 旧的客户端代码需要更新（移除参数传递）)

[//]: # (3. 确保 .env 文件中配置了模型默认参数)

[//]: # (4. service_impl 层会处理所有错误，返回友好的错误信息)

[//]: # ()
[//]: # (## 后续优化建议)

[//]: # ()
[//]: # (1. 支持客户端覆盖默认参数（通过 extra_params）)

[//]: # (2. 添加参数验证（范围检查）)

[//]: # (3. 添加单元测试)

[//]: # (4. 完善日志记录)

[//]: # (5. 添加监控指标)