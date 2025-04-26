package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"codex-cli/config"
)

var (
	configFile string
	config     config.AppConfig
)

func main() {
	cobra.OnInitialize(initConfig)

	rootCmd := &cobra.Command{
		Use:   "codex",
		Short: "Codex CLI",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) > 0 {
				prompt := strings.Join(args, " ")
				runApp(prompt)
			} else {
				cmd.Help()
			}
		},
	}

	rootCmd.PersistentFlags().StringVar(&configFile, "config", "", "config file (default is $HOME/.codex.yaml)")
	rootCmd.PersistentFlags().StringVarP(&config.Model, "model", "m", "o4-mini", "Model to use for completions")
	rootCmd.PersistentFlags().StringVarP(&config.Provider, "provider", "p", "openai", "Provider to use for completions")
	rootCmd.PersistentFlags().BoolVarP(&config.Notify, "notify", "", false, "Enable desktop notifications for responses")
	rootCmd.PersistentFlags().BoolVarP(&config.FlexMode, "flex-mode", "", false, "Use flex-mode processing mode for the request")
	rootCmd.PersistentFlags().BoolVarP(&config.DisableResponseStorage, "disable-response-storage", "", false, "Disable server-side response storage")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".codex")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}
}

func runApp(prompt string) {
	fmt.Println("Running app with prompt:", prompt)
}
