Claude系列模型调用

# 1. 04-Claude系列.md

## 1.1. 04-Claude系列.md
本文件详细介绍了Anthropic公司Claude系列模型（如Claude 3.5 Sonnet, Claude 3 Opus）的调用规范。重点说明了Messages API的使用、系统提示词（System Prompt）的独立设置以及Prompt Caching缓存机制。

# 2. 接口调用方式

## 2.1. Messages API调用
Claude使用 `/v1/messages` 接口。与OpenAI不同，Claude将 `system` 提示词作为一个顶层独立参数，而不是放在 `messages` 数组中。核心入参包括 `model`, `messages`, `system`, `max_tokens` 等。
Go 语言调用示例：
```go
// 使用官方 anthropic-sdk-go (或第三方库)
client := anthropic.NewClient(option.WithAPIKey("your-key"))
resp, err := client.Messages.New(context.Background(), anthropic.MessageNewParams{
    Model:     anthropic.F(anthropic.ModelClaude3_5Sonnet),
    MaxTokens: anthropic.F(int64(1024)),
    System:    anthropic.F([]anthropic.TextBlockParam{{
        Type: anthropic.F(anthropic.TextBlockParamTypeText),
        Text: anthropic.F("You are a helpful assistant."),
    }}),
    Messages: anthropic.F([]anthropic.MessageParam{{
        Role: anthropic.F(anthropic.MessageParamRoleUser),
        Content: anthropic.F([]anthropic.MessagePartContentUnionParam{
            anthropic.TextBlockParam{
                Type: anthropic.F(anthropic.TextBlockParamTypeText),
                Text: anthropic.F("Hello, Claude!"),
            },
        }),
    }}),
})
```

## 2.2. 流式输出（Streaming）
设置 `stream: true` 后，Claude 通过 SSE 返回事件。常见的事件类型包括 `message_start`, `content_block_delta` (包含文本片段) 和 `message_stop`。

# 3. 多模态与缓存

## 3.1. 图像输入
支持通过 Base64 传递图像。格式要求为 `image/jpeg`, `image/png`, `image/gif` 或 `image/webp`。
Go 语言图像示例：
```go
Content: anthropic.F([]anthropic.MessagePartContentUnionParam{
    anthropic.ImageBlockParam{
        Type: anthropic.F(anthropic.ImageBlockParamTypeImage),
        Source: anthropic.F(anthropic.ImageBlockParamSource{
            Type:    anthropic.F(anthropic.ImageBlockParamSourceTypeBase64),
            MediaType: anthropic.F(anthropic.ImageBlockParamSourceMediaTypeImageJPEG),
            Data:    anthropic.F("base64_encoded_data"),
        }),
    },
    anthropic.TextBlockParam{
        Type: anthropic.F(anthropic.TextBlockParamTypeText),
        Text: anthropic.F("Describe this image."),
    },
}),
```

## 3.2. Prompt Caching
Claude 支持显式指定缓存点。通过在 `messages` 或 `system` 的内容块中添加 `"cache_control": {"type": "ephemeral"}`，可以缓存该点之前的内容。适用于长文档问答。
