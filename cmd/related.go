package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"
	"github.com/ChenHom/ytcli/pkg/chapterizer"
	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/knowledge"
	"github.com/spf13/cobra"
)

var (
	relatedTranscript string
	relatedChapter    int
)

var relatedCmd = &cobra.Command{
	Use:           "related",
	Short:         "根據影片主題檢索相關知識",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 互動式選擇章節
		if relatedChapter < 0 {
			chapFile := filepath.Join(filepath.Dir(relatedTranscript), "chapters.json")
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
			relatedChapter = sel
			fmt.Printf("已選擇章節 %d\n", relatedChapter)
		}
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
	relatedCmd.Flags().IntVar(&relatedChapter, "chapter", -1, "要檢索相關知識的章節索引 (從 0 開始)")
	relatedCmd.MarkFlagRequired("transcript")
	relatedCmd.MarkFlagRequired("chapter")
}
