Gemini系列模型调用

# 1. 05-Gemini系列.md

## 1.1. 05-Gemini系列.md
本文件详细介绍了Google Gemini系列模型（如Gemini 1.5 Pro, Gemini 1.5 Flash）的调用规范。重点说明了Generative AI SDK的使用、`contents` 层级结构以及系统指令（System Instructions）的配置。

# 2. 接口调用方式

## 2.1. GenerateContent调用
Gemini主要使用 `generateContent` 接口。其消息结构使用 `contents` 数组，每个 content 包含一个 `role` (user 或 model) 和一个 `parts` 数组。
Go 语言调用示例：
```go
client, err := genai.NewClient(ctx, option.WithAPIKey("your-key"))
model := client.GenerativeModel("gemini-1.5-flash")
// 设置系统指令
model.SystemInstruction = &genai.Content{
    Parts: []genai.Part{genai.Text("You are a helpful assistant.")},
}
resp, err := model.GenerateContent(ctx, genai.Text("Hello, Gemini!"))
```

## 2.2. 流式输出（Streaming）
使用 `GenerateContentStream` 方法。Gemini 的流式响应会分批次返回完整的 `candidates`，开发者需要从中提取新增的文本片段。

# 3. 多模态与特殊功能

## 3.1. 多模态输入（图像/视频）
Gemini 原生支持大上下文多模态。除了 Base64，还可以通过 `File API` 上传大文件（如长视频）并在调用时引用。
Go 语言图像示例：
```go
imgData, _ := os.ReadFile("image.jpg")
resp, err := model.GenerateContent(ctx, 
    genai.ImageData("jpeg", imgData),
    genai.Text("What is in this image?"),
)
```

## 3.2. Context Caching
对于超大规模上下文（如整个代码库或数小时视频），Gemini 提供 Context Caching 功能，通过 TTL（生存时间）管理缓存，显著降低重复调用的成本。
