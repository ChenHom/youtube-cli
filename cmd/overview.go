package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/summarizer"
	"github.com/spf13/cobra"
)

var overviewTranscript string

var overviewCmd = &cobra.Command{
	Use:   "overview",
	Short: "列出章節標題、開始時間及重點句",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 載入 OpenAI API 金鑰
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// 產生輸出檔案路徑
		outputPath := filepath.Join(filepath.Dir(overviewTranscript), "overview.txt")
		if err := summarizer.Overview(cfg.OpenAIKey, overviewTranscript, outputPath); err != nil {
			return err
		}
		fmt.Println("Overview saved to", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(overviewCmd)
	overviewCmd.Flags().StringVar(&overviewTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	overviewCmd.MarkFlagRequired("transcript")
}
