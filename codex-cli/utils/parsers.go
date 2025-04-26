package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
	"codex-cli/config"
)

type Config struct {
	Model                string `json:"model,omitempty" yaml:"model,omitempty"`
	Provider             string `json:"provider,omitempty" yaml:"provider,omitempty"`
	Notify               bool   `json:"notify,omitempty" yaml:"notify,omitempty"`
	FlexMode             bool   `json:"flexMode,omitempty" yaml:"flexMode,omitempty"`
	DisableResponseStorage bool `json:"disableResponseStorage,omitempty" yaml:"disableResponseStorage,omitempty"`
}

func ParseConfigFile(configPath string) (Config, error) {
	var config Config

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, fmt.Errorf("config file not found: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, err
	}

	switch filepath.Ext(configPath) {
	case ".json":
		err = json.Unmarshal(data, &config)
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &config)
	default:
		return config, fmt.Errorf("unsupported config file format: %s", configPath)
	}

	if err != nil {
		return config, err
	}

	return config, nil
}

func ParseCommandLineArgs(args []string) (map[string]string, error) {
	parsedArgs := make(map[string]string)

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			key := strings.TrimPrefix(arg, "--")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "--") {
				parsedArgs[key] = args[i+1]
				i++
			} else {
				parsedArgs[key] = ""
			}
		} else if strings.HasPrefix(arg, "-") {
			key := strings.TrimPrefix(arg, "-")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				parsedArgs[key] = args[i+1]
				i++
			} else {
				parsedArgs[key] = ""
			}
		} else {
			return nil, fmt.Errorf("invalid argument: %s", arg)
		}
	}

	return parsedArgs, nil
}
