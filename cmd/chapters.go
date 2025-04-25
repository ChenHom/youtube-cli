package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/AlecAivazis/survey/v2"

	"github.com/ChenHom/ytcli/pkg/chapterizer"
	"github.com/spf13/cobra"
)

var (
	chaptersTranscript string
	chaptersSelected   []int
)

var chaptersCmd = &cobra.Command{
	Use:           "chapters",
	Short:         "偵測章節並輸出時間區間",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		chapters, err := chapterizer.DetectParagraphChapters(chaptersTranscript)
		if err != nil {
			return err
		}
		fmt.Printf("共偵測到 %d 個章節\n", len(chapters))
		if len(chaptersSelected) == 0 {
			labels := make([]string, len(chapters))
			for i, ch := range chapters {
				labels[i] = fmt.Sprintf("Chapter %d: %s (%s–%s)", i, ch.Title, ch.Start, ch.End)
			}
			var sel []int
			prompt := &survey.MultiSelect{Message: "選擇章節:", Options: labels}
			if err := survey.AskOne(prompt, &sel); err != nil {
				return err
			}
			chaptersSelected = sel
			fmt.Printf("已選取章節索引: %v\n", chaptersSelected)
		}
		// 顯示所選章節內容
		for _, idx := range chaptersSelected {
			if idx >= 0 && idx < len(chapters) {
				fmt.Printf("\n--- Chapter %d ---\n%s\n", idx, chapters[idx].Text)
			}
		}
		data, err := json.MarshalIndent(chapters, "", "  ")
		if err != nil {
			return err
		}
		outputPath := filepath.Join(filepath.Dir(chaptersTranscript), "chapters.json")
		if err := os.WriteFile(outputPath, data, 0644); err != nil {
			return err
		}
		fmt.Println("Chapters saved to", outputPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(chaptersCmd)
	chaptersCmd.Flags().StringVar(&chaptersTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	chaptersCmd.Flags().IntSliceVar(&chaptersSelected, "chapters-selected", nil, "要選擇的章節索引 (多選)")
	chaptersCmd.MarkFlagRequired("transcript")
}
