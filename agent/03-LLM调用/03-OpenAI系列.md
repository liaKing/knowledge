OpenAI系列模型调用

# 1. 03-OpenAI系列.md

## 1.1. 03-OpenAI系列.md
本文件详细介绍了OpenAI系列模型（如GPT-4o, GPT-4 Turbo, GPT-3.5 Turbo）的调用规范。涵盖了标准OpenAI API的使用方法、核心入参出参说明、多模态图像输入技巧以及Prompt Caching缓存机制的应用。

# 2. 接口调用方式

## 2.1. 标准API调用
使用OpenAI官方SDK或HTTP请求调用 `/v1/chat/completions` 接口。核心入参包括 `model` (模型ID), `messages` (对话列表), `temperature` (随机性), `max_tokens` (最大长度) 等。
Go 语言调用示例：
```go
client := openai.NewClient("your-token")
resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model:    openai.GPT4o,
        Messages: []openai.ChatCompletionMessage{{Role: "user", Content: "Hello!"}},
    },
)
```

## 2.2. 流式输出（Streaming）
通过设置 `stream: true`，模型将以Server-Sent Events (SSE) 形式逐字返回响应。适用于需要实时交互激的UI场景，能显著提升用户感知速度。注意：流式输出默认不返回 `usage` 统计，需显式设置 `stream_options: {"include_usage": true}`。
Go 语言流式示例（含 Token 统计）：
```go
stream, err := client.CreateChatCompletionStream(
    context.Background(),
    openai.ChatCompletionRequest{
        Model:    openai.GPT4o,
        Messages: []openai.ChatCompletionMessage{{Role: "user", Content: "Hello!"}},
        Stream:   true,
        // 开启后，最后一个数据块将包含 usage 统计
        StreamOptions: &openai.StreamOptions{
            IncludeUsage: true,
        },
    },
)
defer stream.Close()
for {
    response, err := stream.Recv()
    if errors.Is(err, io.EOF) { break }
    if response.Usage != nil {
        fmt.Printf("\nTokens: %d", response.Usage.TotalTokens)
    }
    if len(response.Choices) > 0 {
        fmt.Print(response.Choices[0].Delta.Content)
    }
}
```

# 3. 多模态图像调用

## 3.1. 图像输入方式
GPT-4o等模型支持图像输入。可以通过 `image_url` 传递公开可访问的图片链接，或者将图片转为 `Base64` 编码字符串并按 `data:image/jpeg;base64,{base64_string}` 格式嵌入请求。
Go 语言多模态示例：
```go
resp, err := client.CreateChatCompletion(
    context.Background(),
    openai.ChatCompletionRequest{
        Model: openai.GPT4o,
        Messages: []openai.ChatCompletionMessage{
            {
                Role: openai.ChatMessageRoleUser,
                MultiContent: []openai.ChatMessagePart{
                    {
                        Type: openai.ChatMessagePartTypeImageURL,
                        ImageURL: &openai.ChatMessageImageURL{URL: "https://url/to/image.jpg"},
                    },
                },
            },
        },
    },
)
```

## 3.2. 分辨率控制
通过 `detail` 参数（low, high, auto）控制模型对图像的解析精度。`high` 模式会将大图切分为 512x512 的网格进行详细分析，消耗更多 Token。

# 4. 缓存与入参出参

## 4.1. Prompt Caching
OpenAI 自动为超过 1024 个 Token 且最近被使用过的相同 Prompt 前缀启用缓存。缓存命中可享受大幅折扣（通常为 50%），并显著降低首字延迟。

## 4.2. 入参出参结构
入参以 JSON 格式发送，包含 `messages` 数组（Role 分为 system, user, assistant）。出参包含 `choices` 数组（包含 `message` 内容和 `finish_reason`）。关于 Token 统计：
- **标准响应**：默认在 `usage` 对象中返回 `prompt_tokens` 和 `completion_tokens`。
- **流式响应**：默认 `usage` 为空。必须在入参中设置 `stream_options: {"include_usage": true}`，OpenAI 才会增加最后一个数据块来返回 `usage` 信息。
