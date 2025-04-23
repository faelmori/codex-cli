package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/blang/semver"
	"github.com/fatih/color"
)

const (
	updateCheckURL      = "https://api.example.com/latest-version"
	updateCheckFileName = "update-check.json"
	updateCheckInterval = 24 * time.Hour
)

type UpdateCheckState struct {
	LastUpdateCheck time.Time `json:"last_update_check"`
}

type UpdateInfo struct {
	CurrentVersion string `json:"current_version"`
	LatestVersion  string `json:"latest_version"`
}

func CheckForUpdates(currentVersion string) {
	stateFilePath := getStateFilePath()
	state := loadUpdateCheckState(stateFilePath)

	if time.Since(state.LastUpdateCheck) < updateCheckInterval {
		return
	}

	latestVersion, err := fetchLatestVersion()
	if err != nil {
		fmt.Println("Error checking for updates:", err)
		return
	}

	if isNewerVersion(currentVersion, latestVersion) {
		printUpdateMessage(currentVersion, latestVersion)
	}

	state.LastUpdateCheck = time.Now()
	saveUpdateCheckState(stateFilePath, state)
}

func getStateFilePath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory:", err)
		return ""
	}
	return filepath.Join(configDir, updateCheckFileName)
}

func loadUpdateCheckState(filePath string) UpdateCheckState {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return UpdateCheckState{}
	}

	var state UpdateCheckState
	if err := json.Unmarshal(data, &state); err != nil {
		return UpdateCheckState{}
	}

	return state
}

func saveUpdateCheckState(filePath string, state UpdateCheckState) {
	data, err := json.Marshal(state)
	if err != nil {
		fmt.Println("Error saving update check state:", err)
		return
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Println("Error writing update check state file:", err)
	}
}

func fetchLatestVersion() (string, error) {
	resp, err := http.Get(updateCheckURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var updateInfo UpdateInfo
	if err := json.NewDecoder(resp.Body).Decode(&updateInfo); err != nil {
		return "", err
	}

	return updateInfo.LatestVersion, nil
}

func isNewerVersion(currentVersion, latestVersion string) bool {
	current, err := semver.Parse(currentVersion)
	if err != nil {
		fmt.Println("Error parsing current version:", err)
		return false
	}

	latest, err := semver.Parse(latestVersion)
	if err != nil {
		fmt.Println("Error parsing latest version:", err)
		return false
	}

	return latest.GT(current)
}

func printUpdateMessage(currentVersion, latestVersion string) {
	color.Yellow("A new version of the application is available!")
	color.Yellow("Current version: %s", currentVersion)
	color.Yellow("Latest version: %s", latestVersion)
	color.Yellow("Please update to the latest version.")
}
