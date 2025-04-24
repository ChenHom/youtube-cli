package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ChenHom/ytcli/pkg/chapterizer"
	"github.com/spf13/cobra"
)

var chaptersTranscript string

var chaptersCmd = &cobra.Command{
	Use:   "chapters",
	Short: "偵測章節並輸出時間區間",
	RunE: func(cmd *cobra.Command, args []string) error {
		chapters, err := chapterizer.DetectChapters(chaptersTranscript)
		if err != nil {
			return err
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
	chaptersCmd.MarkFlagRequired("transcript")
}
