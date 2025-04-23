package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	DefaultAgenticModel       = "o4-mini"
	DefaultFullContextModel   = "gpt-4.1"
	DefaultApprovalMode       = "suggest"
	DefaultInstructions       = ""
	ConfigDir                 = ".codex"
	ConfigJSONFilePath        = "config.json"
	ConfigYAMLFilePath        = "config.yaml"
	ConfigYMLFilePath         = "config.yml"
	InstructionsFilePath      = "instructions.md"
	OpenAITimeoutEnvVar       = "OPENAI_TIMEOUT_MS"
	OpenAIBaseURLEnvVar       = "OPENAI_BASE_URL"
	OpenAIAPIKeyEnvVar        = "OPENAI_API_KEY"
	PrettyPrintEnvVar         = "PRETTY_PRINT"
	ProjectDocMaxBytes        = 32 * 1024 // 32 kB
	ProjectDocFilenames       = "codex.md,.codex.md,CODEX.md"
	CodexDisableProjectDocEnv = "CODEX_DISABLE_PROJECT_DOC"
)

type StoredConfig struct {
	Model                string `json:"model,omitempty" yaml:"model,omitempty"`
	Provider             string `json:"provider,omitempty" yaml:"provider,omitempty"`
	ApprovalMode         string `json:"approvalMode,omitempty" yaml:"approvalMode,omitempty"`
	FullAutoErrorMode    string `json:"fullAutoErrorMode,omitempty" yaml:"fullAutoErrorMode,omitempty"`
	Notify               bool   `json:"notify,omitempty" yaml:"notify,omitempty"`
	DisableResponseStorage bool `json:"disableResponseStorage,omitempty" yaml:"disableResponseStorage,omitempty"`
	History              struct {
		MaxSize          int      `json:"maxSize,omitempty" yaml:"maxSize,omitempty"`
		SaveHistory      bool     `json:"saveHistory,omitempty" yaml:"saveHistory,omitempty"`
		SensitivePatterns []string `json:"sensitivePatterns,omitempty" yaml:"sensitivePatterns,omitempty"`
	} `json:"history,omitempty" yaml:"history,omitempty"`
}

type AppConfig struct {
	APIKey               string
	Model                string
	Provider             string
	Instructions         string
	ApprovalMode         string
	FullAutoErrorMode    string
	Notify               bool
	DisableResponseStorage bool
	FlexMode             bool
	History              struct {
		MaxSize          int
		SaveHistory      bool
		SensitivePatterns []string
	}
}

func LoadConfig(configPath string, instructionsPath string, options map[string]interface{}) (AppConfig, error) {
	var config AppConfig
	var storedConfig StoredConfig

	if configPath == "" {
		configPath = ConfigJSONFilePath
	}

	if instructionsPath == "" {
		instructionsPath = InstructionsFilePath
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		if strings.HasSuffix(configPath, ".json") {
			configPath = ConfigYAMLFilePath
		} else if strings.HasSuffix(configPath, ".yaml") {
			configPath = ConfigYMLFilePath
		}
	}

	if _, err := os.Stat(configPath); err == nil {
		data, err := os.ReadFile(configPath)
		if err != nil {
			return config, err
		}

		if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
			err = yaml.Unmarshal(data, &storedConfig)
		} else {
			err = json.Unmarshal(data, &storedConfig)
		}

		if err != nil {
			return config, err
		}
	}

	config.Model = storedConfig.Model
	config.Provider = storedConfig.Provider
	config.ApprovalMode = storedConfig.ApprovalMode
	config.FullAutoErrorMode = storedConfig.FullAutoErrorMode
	config.Notify = storedConfig.Notify
	config.DisableResponseStorage = storedConfig.DisableResponseStorage
	config.History.MaxSize = storedConfig.History.MaxSize
	config.History.SaveHistory = storedConfig.History.SaveHistory
	config.History.SensitivePatterns = storedConfig.History.SensitivePatterns

	if _, err := os.Stat(instructionsPath); err == nil {
		data, err := os.ReadFile(instructionsPath)
		if err != nil {
			return config, err
		}
		config.Instructions = string(data)
	} else {
		config.Instructions = DefaultInstructions
	}

	if apiKey := os.Getenv(OpenAIAPIKeyEnvVar); apiKey != "" {
		config.APIKey = apiKey
	}

	if flexMode, ok := options["flexMode"].(bool); ok {
		config.FlexMode = flexMode
	}

	return config, nil
}

func SaveConfig(config AppConfig, configPath string, instructionsPath string) error {
	if configPath == "" {
		configPath = ConfigJSONFilePath
	}

	if instructionsPath == "" {
		instructionsPath = InstructionsFilePath
	}

	dir := filepath.Dir(configPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	storedConfig := StoredConfig{
		Model:                config.Model,
		Provider:             config.Provider,
		ApprovalMode:         config.ApprovalMode,
		FullAutoErrorMode:    config.FullAutoErrorMode,
		Notify:               config.Notify,
		DisableResponseStorage: config.DisableResponseStorage,
	}
	storedConfig.History.MaxSize = config.History.MaxSize
	storedConfig.History.SaveHistory = config.History.SaveHistory
	storedConfig.History.SensitivePatterns = config.History.SensitivePatterns

	var data []byte
	var err error

	if strings.HasSuffix(configPath, ".yaml") || strings.HasSuffix(configPath, ".yml") {
		data, err = yaml.Marshal(storedConfig)
	} else {
		data, err = json.MarshalIndent(storedConfig, "", "  ")
	}

	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return err
	}

	err = os.WriteFile(instructionsPath, []byte(config.Instructions), 0644)
	if err != nil {
		return err
	}

	return nil
}

func GetBaseURL(provider string) string {
	if baseURL := os.Getenv(OpenAIBaseURLEnvVar); baseURL != "" {
		return baseURL
	}

	if provider == "openai" {
		return "https://api.openai.com"
	}

	return ""
}

func GetAPIKey(provider string) string {
	if apiKey := os.Getenv(OpenAIAPIKeyEnvVar); apiKey != "" {
		return apiKey
	}

	return ""
}

func LoadProjectDoc(cwd string, explicitPath string) (string, error) {
	if explicitPath != "" {
		if _, err := os.Stat(explicitPath); os.IsNotExist(err) {
			return "", fmt.Errorf("project doc not found at %s", explicitPath)
		}
		return os.ReadFile(explicitPath)
	}

	for _, name := range strings.Split(ProjectDocFilenames, ",") {
		path := filepath.Join(cwd, name)
		if _, err := os.Stat(path); err == nil {
			return os.ReadFile(path)
		}
	}

	return "", nil
}

func DiscoverProjectDocPath(startDir string) (string, error) {
	cwd, err := filepath.Abs(startDir)
	if err != nil {
		return "", err
	}

	for {
		for _, name := range strings.Split(ProjectDocFilenames, ",") {
			path := filepath.Join(cwd, name)
			if _, err := os.Stat(path); err == nil {
				return path, nil
			}
		}

		parent := filepath.Dir(cwd)
		if parent == cwd {
			break
		}
		cwd = parent
	}

	return "", nil
}

func InitConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.codex")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Error reading config file:", err)
	}
}
