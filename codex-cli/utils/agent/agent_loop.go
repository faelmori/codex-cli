package agent

import (
	"context"
	"fmt"
	"time"
)

// AgentLoop represents the main loop for the agent.
type AgentLoop struct {
	ctx    context.Context
	cancel context.CancelFunc
	ticker *time.Ticker
}

// NewAgentLoop creates a new AgentLoop instance.
func NewAgentLoop(interval time.Duration) *AgentLoop {
	ctx, cancel := context.WithCancel(context.Background())
	return &AgentLoop{
		ctx:    ctx,
		cancel: cancel,
		ticker: time.NewTicker(interval),
	}
}

// Start begins the agent loop.
func (a *AgentLoop) Start() {
	for {
		select {
		case <-a.ctx.Done():
			fmt.Println("Agent loop stopped")
			return
		case <-a.ticker.C:
			a.executeTask()
		}
	}
}

// Stop stops the agent loop.
func (a *AgentLoop) Stop() {
	a.cancel()
	a.ticker.Stop()
}

// executeTask represents the task to be executed in each loop iteration.
func (a *AgentLoop) executeTask() {
	// Implement the task logic here
	fmt.Println("Executing task")
}
