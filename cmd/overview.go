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
	overviewTranscript string
	overviewChapter    int
)

var overviewCmd = &cobra.Command{
	Use:           "overview",
	Short:         "列出章節標題、開始時間及重點句",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 載入 OpenAI API 金鑰
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// 互動式選擇章節
		if overviewChapter < 0 {
			chapFile := filepath.Join(filepath.Dir(overviewTranscript), "chapters.json")
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
			overviewChapter = sel
			fmt.Printf("已選擇章節 %d\n", overviewChapter)
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
	overviewCmd.Flags().IntVar(&overviewChapter, "chapter", -1, "要顯示概要的章節索引 (從 0 開始)")
	overviewCmd.MarkFlagRequired("transcript")
}
