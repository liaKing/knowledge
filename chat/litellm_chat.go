package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

func keepLastNMessages(msgs []openai.ChatCompletionMessageParamUnion, n int) []openai.ChatCompletionMessageParamUnion {
	if n <= 0 || len(msgs) <= n {
		return msgs
	}
	return msgs[len(msgs)-n:]
}

func main() {
	in := bufio.NewReader(os.Stdin)

	fmt.Print("请输入 LiteLLM Base URL (例如 https://api-dev.joycloud.ai/v1): ")
	baseURL, _ := in.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" {
		fmt.Println("base url 不能为空")
		os.Exit(1)
	}
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		fmt.Println("base url 必须以 http:// 或 https:// 开头（你可能把 api key 输到这里了）")
		os.Exit(1)
	}

	apiKey := os.Getenv("LITELLM_API_KEY")
	if apiKey == "" {
		fmt.Print("请输入 LiteLLM API Key (或先设置环境变量 LITELLM_API_KEY): ")
		apiKey, _ = in.ReadString('\n')
		apiKey = strings.TrimSpace(apiKey)
		if apiKey == "" {
			fmt.Println("api key 不能为空")
			os.Exit(1)
		}
	}

	client := openai.NewClient(
		option.WithAPIKey(apiKey),
		option.WithBaseURL(baseURL),
	)

	modelName := "claude-sonnet-4-6"
	ctx := context.Background()

	systemMsg := openai.SystemMessage("你是一个有帮助的 AI 助手。")
	history := []openai.ChatCompletionMessageParamUnion{systemMsg}

	fmt.Println("开始聊天，输入 exit 退出。")
	for {
		fmt.Print("\nYou: ")
		userText, _ := in.ReadString('\n')
		userText = strings.TrimSpace(userText)
		if userText == "" {
			continue
		}
		if userText == "exit" {
			return
		}

		history = append(history, openai.UserMessage(userText))

		// 只携带最近10条对话消息（不含 system），并确保 system 永远在最前面
		msgs := append([]openai.ChatCompletionMessageParamUnion{systemMsg}, keepLastNMessages(history[1:], 10)...)

		fmt.Print("Assistant: ")
		stream := client.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
			Model:    openai.ChatModel(modelName),
			Messages: msgs,
			MaxTokens: openai.Int(2000),
		})

		var assistant strings.Builder
		for stream.Next() {
			chunk := stream.Current()
			if len(chunk.Choices) == 0 {
				continue
			}
			delta := chunk.Choices[0].Delta
			if delta.Content != "" {
				fmt.Print(delta.Content)
				assistant.WriteString(delta.Content)
			}
		}
		if err := stream.Err(); err != nil {
			fmt.Printf("\n(stream error: %v)\n", err)
			continue
		}
		fmt.Println()

		history = append(history, openai.AssistantMessage(assistant.String()))
	}
}
