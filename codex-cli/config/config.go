package config

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

type AppConfig struct {
	Model                string
	Provider             string
	Notify               bool
	FlexMode             bool
	DisableResponseStorage bool
}

var (
	configFile string
	config     AppConfig
)

func InitConfig() {
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
