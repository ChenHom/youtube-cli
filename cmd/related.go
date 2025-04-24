package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/knowledge"
	"github.com/spf13/cobra"
)

var relatedTranscript string
var relatedChapter int

var relatedCmd = &cobra.Command{
	Use:   "related",
	Short: "根據影片主題檢索相關知識",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		outputPath := filepath.Join(filepath.Dir(relatedTranscript), fmt.Sprintf("related_%d.txt", relatedChapter))
		if err := knowledge.Related(cfg.OpenAIKey, relatedTranscript, relatedChapter, outputPath); err != nil {
			return err
		}
		fmt.Println("Related saved to", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(relatedCmd)
	relatedCmd.Flags().StringVar(&relatedTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	relatedCmd.Flags().IntVar(&relatedChapter, "chapter", 0, "要檢索相關知識的章節索引 (從 0 開始)")
	relatedCmd.MarkFlagRequired("transcript")
	relatedCmd.MarkFlagRequired("chapter")
}
