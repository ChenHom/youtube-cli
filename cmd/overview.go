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
	"github.com/ChenHom/ytcli/pkg/summarizer"
	"github.com/ChenHom/ytcli/pkg/translator"
	"github.com/spf13/cobra"
)

var (
	overviewTranscript string
	overviewChapter    int
	overviewLang       string
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
		// 讀取章節資料
		chapFile := filepath.Join(filepath.Dir(overviewTranscript), "chapters.json")
		data, err := os.ReadFile(chapFile)
		var chapters []chapterizer.Chapter
		if err != nil || json.Unmarshal(data, &chapters) != nil {
			// 若讀取或解析失敗，重新偵測章節並寫入章節檔
			chapters, err = chapterizer.DetectParagraphChapters(overviewTranscript)
			if err != nil {
				return err
			}
			data, err = json.MarshalIndent(chapters, "", "  ")
			if err != nil {
				return err
			}
			if err := os.WriteFile(chapFile, data, 0644); err != nil {
				return err
			}
		}
		// 主迴圈: 可重選章節或執行相關操作
		for {
			// 互動式選章
			if overviewChapter < 0 {
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
				// 立即顯示該章節內容
				if overviewChapter >= 0 && overviewChapter < len(chapters) {
					fmt.Println("\n▶ 章節全文：")
					fmt.Println(chapters[overviewChapter].Text)
				}
			}
			// 產生並存檔 overview
			outputPath := filepath.Join(filepath.Dir(overviewTranscript), "overview.txt")
			if err := summarizer.Overview(cfg.OpenAIKey, overviewTranscript, outputPath); err != nil {
				return err
			}
			// 翻譯至指定語言
			if overviewLang != "" {
				if err := translator.TranslateFile(cfg.OpenAIKey, outputPath, overviewLang); err != nil {
					return err
				}
			}
			// 讀取並印出內容
			content, err := os.ReadFile(outputPath)
			if err != nil {
				return err
			}
			fmt.Println("\n▶ 章節內容（立即顯示）：")
			fmt.Println(string(content))
			fmt.Println("Overview saved to", outputPath)
			// 下一步操作選單
			actions := []string{"重選章節", "檢索相關知識", "深入探討此章節", "結束"}
			var act string
			actPrompt := &survey.Select{Message: "下一步:", Options: actions}
			if err := survey.AskOne(actPrompt, &act); err != nil {
				return err
			}
			switch act {
			case "重選章節":
				overviewChapter = -1
				continue
			case "檢索相關知識":
				// 呼叫 knowledge.Related
				relPath := filepath.Join(filepath.Dir(overviewTranscript), fmt.Sprintf("related_%d.txt", overviewChapter))
				if err := knowledge.Related(cfg.OpenAIKey, overviewTranscript, overviewChapter, relPath); err != nil {
					return err
				}
				if overviewLang != "" {
					if err := translator.TranslateFile(cfg.OpenAIKey, relPath, overviewLang); err != nil {
						return err
					}
				}
				relContent, err := os.ReadFile(relPath)
				if err != nil {
					return err
				}
				fmt.Println("\n▶ 相關知識：")
				fmt.Println(string(relContent))
				fmt.Println("Related saved to", relPath)
			case "深入探討此章節":
				ddPath := filepath.Join(filepath.Dir(overviewTranscript), fmt.Sprintf("deep_dive_%d.txt", overviewChapter))
				if err := knowledge.DeepDive(cfg.OpenAIKey, overviewTranscript, overviewChapter, ddPath); err != nil {
					return err
				}
				if overviewLang != "" {
					if err := translator.TranslateFile(cfg.OpenAIKey, ddPath, overviewLang); err != nil {
						return err
					}
				}
				ddContent, err := os.ReadFile(ddPath)
				if err != nil {
					return err
				}
				fmt.Println("\n▶ 深入探討：")
				fmt.Println(string(ddContent))
				fmt.Println("Deep dive saved to", ddPath)
			case "結束":
				return nil
			}
			// 完成相關操作後，可選擇是否繼續
			overviewChapter = -1
		}
	},
}

func init() {
	rootCmd.AddCommand(overviewCmd)
	overviewCmd.Flags().StringVar(&overviewTranscript, "transcript", "", "字幕檔案 (*.vtt)")
	overviewCmd.Flags().IntVar(&overviewChapter, "chapter", -1, "要顯示概要的章節索引 (從 0 開始)")
	overviewCmd.Flags().StringVar(&overviewLang, "lang", "zh-TW", "輸出文字說明語言，預設正體中文；可選 en、ja、…")
	overviewCmd.MarkFlagRequired("transcript")
}
