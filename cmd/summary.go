package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ChenHom/ytcli/pkg/chapterizer"
	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/summarizer"
	"github.com/spf13/cobra"
)

var (
	summaryTranscript string
	summaryChapter    int
)

var summaryCmd = &cobra.Command{
	Use:           "summary",
	Short:         "生成影片摘要",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 互動式選擇章節
		if summaryChapter < 0 {
			chapFile := filepath.Join(filepath.Dir(summaryTranscript), "chapters.json")
			data, err := os.ReadFile(chapFile)
			if err != nil {
				return err
			}
			var chapters []chapterizer.Chapter
			if err := json.Unmarshal(data, &chapters); err != nil {
				return err
			}
			labels := make([]string, len(chapters))
			for i, ch := range chapters {
				labels[i] = fmt.Sprintf("Chapter %d: %s (%s–%s)", i, ch.Title, ch.Start, ch.End)
			}
			var sel int
			prompt := &survey.Select{Message: "選擇章節:", Options: labels}
			if err := survey.AskOne(prompt, &sel); err != nil {
				return err
			}
			summaryChapter = sel
			fmt.Printf("已選擇章節 %d\n", summaryChapter)
		}
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
	summaryCmd.Flags().IntVar(&summaryChapter, "chapter", -1, "要摘要的章節索引 (從 0 開始)")
	summaryCmd.MarkFlagRequired("transcript")
}
