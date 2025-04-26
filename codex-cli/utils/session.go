package session

import (
	"sync"
)

const (
	CLIVersion = "0.1.2504221401"
	Origin     = "codex_cli_ts"
)

type TerminalChatSession struct {
	ID          string
	User        string
	Version     string
	Model       string
	Timestamp   string
	Instructions string
}

var (
	sessionID    string
	sessionMutex sync.RWMutex
	currentModel string
	modelMutex   sync.RWMutex
)

func SetSessionID(id string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	sessionID = id
}

func GetSessionID() string {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	return sessionID
}

func SetCurrentModel(model string) {
	modelMutex.Lock()
	defer modelMutex.Unlock()
	currentModel = model
}

func GetCurrentModel() string {
	modelMutex.RLock()
	defer modelMutex.RUnlock()
	return currentModel
}
