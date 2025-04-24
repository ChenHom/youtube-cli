package config

import (
	"fmt"
	"os"
)

// Config 儲存應用程式設定
type Config struct {
	OpenAIKey string
}

// Load 從環境變數讀取設定
func Load() (*Config, error) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("環境變數 OPENAI_API_KEY 未設定")
	}
	return &Config{OpenAIKey: key}, nil
}
