package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var lang string
var verbose bool

var rootCmd = &cobra.Command{
	Use:              "ytcli",
	Short:            "YouTube AI CLI：下載影片、轉錄、章節偵測、摘要、深入與相關知識檢索",
	TraverseChildren: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("請使用子命令，例如 download、transcribe 等")
	},
}

// Execute 執行根命令
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ytcli.yaml)")
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	rootCmd.PersistentFlags().StringVar(&lang, "lang", "zh-TW", "輸出文字說明語言，預設正體中文；可選 en、ja、…")
	viper.BindPFlag("lang", rootCmd.PersistentFlags().Lookup("lang"))
	// 新增 verbose flag
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "啟用詳細輸出")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME")
		viper.SetConfigName(".ytcli")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("使用 config: ", viper.ConfigFileUsed())
	}
}
