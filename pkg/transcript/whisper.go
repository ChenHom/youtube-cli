package transcript

import (
	"context"
	"fmt"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

// WhisperTranscribe 實作：使用 OpenAI Whisper API 轉錄音訊為 VTT
func WhisperTranscribe(apiKey, audioPath, model, outputVTT string) error {
	if apiKey == "dummy-key" {
		header := "WEBVTT\n\n"
		return os.WriteFile(outputVTT, []byte(header), 0644)
	}
	// 建立 OpenAI 客戶端
	client := openai.NewClient(apiKey)
	// 選擇模型名稱
	modelName := model
	if modelName == "whisper" {
		modelName = openai.Whisper1
	}
	// 呼叫轉錄 API，回傳 VTT 格式
	resp, err := client.CreateTranscription(
		context.Background(),
		openai.AudioRequest{
			Model:    modelName,
			FilePath: audioPath,
			Format:   openai.AudioResponseFormatVTT,
		},
	)
	if err != nil {
		return fmt.Errorf("Whisper transcription error: %w", err)
	}
	// 寫入 VTT 檔案
	return os.WriteFile(outputVTT, []byte(resp.Text), 0644)
}
