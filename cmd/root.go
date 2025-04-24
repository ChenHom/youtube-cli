package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "ytcli",
	Short: "YouTube AI CLI：下載影片、轉錄、章節偵測、摘要、深入與相關知識檢索",
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
