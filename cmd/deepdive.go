package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/knowledge"
	"github.com/spf13/cobra"
)

var deepdiveTranscript string
var deepdiveSection int

var deepdiveCmd = &cobra.Command{
	Use:   "deep-dive",
	Short: "對指定章節進行深入探討與分析",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		outputPath := filepath.Join(filepath.Dir(deepdiveTranscript), fmt.Sprintf("deep_dive_%d.txt", deepdiveSection))
		if err := knowledge.DeepDive(cfg.OpenAIKey, deepdiveTranscript, deepdiveSection, outputPath); err != nil {
			return err
		}
		fmt.Println("Deep dive saved to", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deepdiveCmd)
	deepdiveCmd.Flags().StringVar(&deepdiveTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	deepdiveCmd.Flags().IntVar(&deepdiveSection, "chapter", 0, "要深入探討的章節索引 (從 0 開始)")
	deepdiveCmd.MarkFlagRequired("transcript")
	deepdiveCmd.MarkFlagRequired("chapter")
}
