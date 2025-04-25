package config

import (
	"fmt"
	"os"
)

// Config 儲存應用程式設定
type Config struct {
	OpenAIKey string
	Model     string // 新增: OpenAI model 名稱
}

// Load 從環境變數讀取設定
func Load() (*Config, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("環境變數 OPENAI_API_KEY 未設定")
	}
	model := os.Getenv("OPENAI_MODEL")
	if model == "" {
		model = "gpt-4o-mini" // 預設值，可依需求調整
	}
	return &Config{OpenAIKey: key, Model: model}, nil
}
