package summarizer

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

// Overview 讀取 transcriptPath，將首行作為章節概要輸出到 outputPath
func Overview(apiKey, transcriptPath, outputPath string) error {
	// 測試 stub: dummy-key 時只輸出首行文字
	if apiKey == "dummy-key" {
		f, err := os.Open(transcriptPath)
		if err != nil {
			return err
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		// skip header
		for scanner.Scan() {
			line := scanner.Text()
			if strings.TrimSpace(line) != "" && !strings.HasPrefix(line, "WEBVTT") {
				content := fmt.Sprintf("Overview: %s\n", line)
				return os.WriteFile(outputPath, []byte(content), 0644)
			}
		}
		// 若無其他內容，寫入預設空 overview
		return os.WriteFile(outputPath, []byte("Overview:\n"), 0644)
	}

	// 讀取所有文字
	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return err
	}
	// 準備 GPT 請求
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful assistant for generating overviews."},
			{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Provide an overview of the following transcript:\n%s", string(data))},
		},
	}
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}
	output := fmt.Sprintf("Overview: %s\n", resp.Choices[0].Message.Content)
	return os.WriteFile(outputPath, []byte(output), 0644)
}

// Summarize 生成影片摘要，輸出到 outputPath
func Summarize(apiKey, transcriptPath, outputPath string) error {
	// 測試 stub: dummy-key 時直接輸出 placeholder
	if apiKey == "dummy-key" {
		return os.WriteFile(outputPath, []byte("Summary placeholder\n"), 0644)
	}

	data, err := os.ReadFile(transcriptPath)
	if err != nil {
		return err
	}
	client := openai.NewClient(apiKey)
	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{Role: openai.ChatMessageRoleSystem, Content: "You are a helpful assistant for summarization."},
			{Role: openai.ChatMessageRoleUser, Content: fmt.Sprintf("Summarize the following transcript in concise text (<500 words):\n%s", string(data))},
		},
	}
	ctx := context.Background()
	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, []byte(resp.Choices[0].Message.Content), 0644)
}
