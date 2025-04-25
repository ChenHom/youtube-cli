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
	deepdiveTranscript string
	deepdiveSection    int
)

var deepdiveCmd = &cobra.Command{
	Use:           "deep-dive",
	Short:         "對指定章節進行深入探討與分析",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		// 互動式選擇章節
		if deepdiveSection < 0 {
			chapFile := filepath.Join(filepath.Dir(deepdiveTranscript), "chapters.json")
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
			deepdiveSection = sel
			fmt.Printf("已選擇章節 %d\n", deepdiveSection)
		}
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
	deepdiveCmd.Flags().IntVar(&deepdiveSection, "chapter", -1, "要深入探討的章節索引 (從 0 開始)")
	deepdiveCmd.MarkFlagRequired("transcript")
	deepdiveCmd.MarkFlagRequired("chapter")
}
