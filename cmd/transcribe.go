package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/ChenHom/ytcli/pkg/audio"
	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/transcript"
	"github.com/spf13/cobra"
)

var transcribeInput string
var transcribeModel string

var transcribeCmd = &cobra.Command{
	Use:   "transcribe",
	Short: "抽取音訊並用 Whisper 轉文字，產生字幕檔",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 載入 OpenAI API 金鑰
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// 音訊檔路徑：同名 wav
		audioPath := strings.TrimSuffix(transcribeInput, filepath.Ext(transcribeInput)) + ".wav"
		if err := audio.Extract(transcribeInput, audioPath); err != nil {
			return err
		}
		// 產生 VTT 檔路徑
		outputVTT := strings.TrimSuffix(transcribeInput, filepath.Ext(transcribeInput)) + ".vtt"
		if err := transcript.WhisperTranscribe(cfg.OpenAIKey, audioPath, transcribeModel, outputVTT); err != nil {
			return err
		}
		fmt.Println("Transcript:", outputVTT)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(transcribeCmd)
	transcribeCmd.Flags().StringVar(&transcribeInput, "input", "", "輸入影片檔案 (*.mp4)")
	transcribeCmd.Flags().StringVar(&transcribeModel, "model", "whisper", "Whisper 模型名稱")
	transcribeCmd.MarkFlagRequired("input")
}
