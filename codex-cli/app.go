package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type AppConfig struct {
	Model                string
	Provider             string
	Notify               bool
	FlexMode             bool
	DisableResponseStorage bool
}

type AppRollout struct {
	Session TerminalChatSession
	Items   []ResponseItem
}

type TerminalChatSession struct {
	Version   string
	ID        string
	Model     string
	Timestamp string
	User      string
}

type ResponseItem struct {
	Type    string
	Role    string
	Content []ResponseContent
}

type ResponseContent struct {
	Type     string
	Text     string
	Filename string
	Refusal  string
}

var (
	configFile string
	config     AppConfig
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
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"prompt": prompt,
		})
	})

	r.POST("/submit", func(c *gin.Context) {
		var input struct {
			Text string `form:"text"`
		}
		if err := c.ShouldBind(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		response := processInput(input.Text)
		c.JSON(http.StatusOK, gin.H{"response": response})
	})

	r.Run()
}

func processInput(input string) string {
	// Implement the logic to process the input and generate a response
	return "Processed: " + input
}
