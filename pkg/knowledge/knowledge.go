package knowledge

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// DeepDive 對指定章節做深入分析，輸出到 outputPath
func DeepDive(apiKey, transcriptPath string, chapter int, outputPath string) error {
	// 測試 stub
	if apiKey == "dummy-key" {
		content := fmt.Sprintf("Deep dive placeholder for chapter %d of %s\n", chapter, transcriptPath)
		return os.WriteFile(outputPath, []byte(content), 0644)
	}
	// 讀取全文
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return err
	}
	// 呼叫 GPT 生成深入分析
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an assistant for detailed analysis of transcripts."},
			{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Provide a deep dive analysis for chapter %d of the following transcript:\n%s", chapter, string(data))},
		},
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(resp.Choices[0].Message.Content), 0644)
}

// Related 列出與指定章節相關的知識主題，輸出到 outputPath
func Related(apiKey, transcriptPath string, chapter int, outputPath string) error {
	// 測試 stub
	if apiKey == "dummy-key" {
		content := fmt.Sprintf("Related topics placeholder for chapter %d of %s\n", chapter, transcriptPath)
		return os.WriteFile(outputPath, []byte(content), 0644)
	}
	// 讀取全文
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return err
	}
	// 呼叫 GPT 生成相關主題列表
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are an assistant for suggesting related knowledge topics."},
			{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("List 5 related knowledge topics for chapter %d of the following transcript:\n%s", chapter, string(data))},
		},
	}
	resp, err := client.CreateChatCompletion(context.Background(), req)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(resp.Choices[0].Message.Content), 0644)
}
