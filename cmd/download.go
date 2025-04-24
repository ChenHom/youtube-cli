package cmd

import (
	"fmt"

	"github.com/ChenHom/ytcli/pkg/downloader"
	"github.com/spf13/cobra"
)

var downloadURL string
var downloadOutput string

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "下載影片檔與內嵌字幕（若有）",
	RunE: func(cmd *cobra.Command, args []string) error {
		video, subtitle, err := downloader.Download(downloadURL, downloadOutput)
		if err != nil {
			return err
		}
		fmt.Printf("Download complete: video=%s\n", video)
		if subtitle != "" {
			fmt.Printf("Subtitle file: %s\n", subtitle)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.Flags().StringVar(&downloadURL, "url", "", "YouTube 影片網址")
	downloadCmd.Flags().StringVar(&downloadOutput, "output", "", "輸出資料夾路徑")
	downloadCmd.MarkFlagRequired("url")
	downloadCmd.MarkFlagRequired("output")
}
