package cmd

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/ChenHom/ytcli/pkg/audio"
	"github.com/ChenHom/ytcli/pkg/config"
	"github.com/ChenHom/ytcli/pkg/transcript"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var (
	transcribeInput string
	transcribeModel string
)

var transcribeCmd = &cobra.Command{
	Use:           "transcribe",
	Short:         "抽取音訊並用 Whisper 轉文字，產生字幕檔",
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}
		// 顯示音訊抽取進度
		sp := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		audioPath := strings.TrimSuffix(transcribeInput, filepath.Ext(transcribeInput)) + ".wav"
		sp.Prefix = "Extracting audio... "
		sp.Start()
		if err := audio.Extract(transcribeInput, audioPath); err != nil {
			sp.Stop()
			return err
		}
		sp.Stop()
		// 顯示轉錄進度
		sp2 := spinner.New(spinner.CharSets[14], 150*time.Millisecond)
		outputVTT := strings.TrimSuffix(transcribeInput, filepath.Ext(transcribeInput)) + ".vtt"
		sp2.Prefix = "Transcribing... "
		sp2.Start()
		if err := transcript.WhisperTranscribe(cfg.OpenAIKey, audioPath, transcribeModel, outputVTT); err != nil {
			sp2.Stop()
			return err
		}
		sp2.Stop()
		fmt.Println("Transcript:", outputVTT)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(transcribeCmd)
	transcribeCmd.Flags().StringVar(&transcribeInput, "input", "", "輸入影片檔案 (*.mp4)")
	transcribeCmd.Flags().StringVar(&transcribeModel, "model", "whisper", "Whisper 模型名稱")
	transcribeCmd.MarkFlagRequired("input")
}
