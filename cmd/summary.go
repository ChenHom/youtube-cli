package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/summarizer"
	"github.com/spf13/cobra"
)

var summaryTranscript string

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "生成整體摘要",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		outputPath := filepath.Join(filepath.Dir(summaryTranscript), "summary.txt")
		if err := summarizer.Summarize(cfg.OpenAIKey, summaryTranscript, outputPath); err != nil {
			return err
		}
		fmt.Println("Summary saved to", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
	summaryCmd.Flags().StringVar(&summaryTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	summaryCmd.MarkFlagRequired("transcript")
}
