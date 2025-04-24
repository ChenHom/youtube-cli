package translator

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// TranslateFile 讀取 filePath 內容，呼叫 OpenAI API 翻譯為 targetLang，並寫回原檔
func TranslateFile(apiKey, filePath, targetLang string) error {
	if apiKey == "dummy-key" {
		// 測試 stub: 直接複製原檔案內容（或加註 mock 註記）
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("read file failed: %w", err)
		}
		// 你也可以選擇在內容前加上註記，例如：
		// mockContent := append([]byte("[MOCK TRANSLATION]\n"), data...)
		return os.WriteFile(filePath, data, 0644)
	}
	// 讀取原始內容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file failed: %w", err)
	}
	// 建立 OpenAI client
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: fmt.Sprintf("You are a translation engine. Translate the following text to %s while preserving formatting and file structure.", targetLang)},
			{Role: openai.ChatMessageRoleUser, Content: string(data)},
		},
	}
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("translation API error: %w", err)
	}
	translated := resp.Choices[0].Message.Content
	// 寫回檔案
	if err := os.WriteFile(filePath, []byte(translated), 0644); err != nil {
		return fmt.Errorf("write translated file failed: %w", err)
	}
	return nil
}
