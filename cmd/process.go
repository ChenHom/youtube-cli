package cmd

import (
	"fmt"

	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/pipeline"
	"github.com/spf13/cobra"
)

var processURL string
var processOutput string
var processModel string
var doChapters bool
var doOverview bool
var doSummary bool
var deepDiveChapter int
var processRelatedChapter int

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "執行完整處理流程：下載、轉錄、章節、摘要等",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 載入設定檔
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// 執行 pipeline
		fmt.Printf("Processing %s -> %s\n", processURL, processOutput)
		return pipeline.Process(cfg.OpenAIKey, processURL, processOutput, processModel, doChapters, doOverview, doSummary, deepDiveChapter, processRelatedChapter)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)
	processCmd.Flags().StringVar(&processURL, "url", "", "YouTube 影片網址")
	processCmd.Flags().StringVar(&processOutput, "output", "", "輸出資料夾路徑")
	processCmd.Flags().StringVar(&processModel, "model", "whisper", "Whisper 模型名稱")
	processCmd.Flags().BoolVar(&doChapters, "chapters", false, "是否偵測章節")
	processCmd.Flags().BoolVar(&doOverview, "overview", false, "是否列出章節概要")
	processCmd.Flags().BoolVar(&doSummary, "summary", false, "是否生成影片摘要")
	processCmd.Flags().IntVar(&deepDiveChapter, "deep-dive", -1, "要深入探討的章節索引 (從 0 開始)")
	processCmd.Flags().IntVar(&processRelatedChapter, "related", -1, "要檢索相關知識的章節索引 (從 0 開始)")
	processCmd.MarkFlagRequired("url")
	processCmd.MarkFlagRequired("output")
}
