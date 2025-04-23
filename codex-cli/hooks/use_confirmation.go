package hooks

import (
	"sync"
)

type ConfirmationResult struct {
	Decision         string
	CustomDenyMessage string
}

type ConfirmationItem struct {
	Prompt      string
	Resolve     func(result ConfirmationResult)
	Explanation string
}

type UseConfirmation struct {
	current            *ConfirmationItem
	queue              []*ConfirmationItem
	mutex              sync.Mutex
}

func NewUseConfirmation() *UseConfirmation {
	return &UseConfirmation{
		queue: make([]*ConfirmationItem, 0),
	}
}

func (uc *UseConfirmation) SubmitConfirmation(result ConfirmationResult) {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.current != nil {
		uc.current.Resolve(result)
		uc.advanceQueue()
	}
}

func (uc *UseConfirmation) RequestConfirmation(prompt string, explanation string) chan ConfirmationResult {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	resultChan := make(chan ConfirmationResult)
	item := &ConfirmationItem{
		Prompt:      prompt,
		Resolve:     func(result ConfirmationResult) { resultChan <- result },
		Explanation: explanation,
	}

	wasEmpty := len(uc.queue) == 0
	uc.queue = append(uc.queue, item)

	if wasEmpty {
		uc.advanceQueue()
	}

	return resultChan
}

func (uc *UseConfirmation) advanceQueue() {
	if len(uc.queue) > 0 {
		uc.current = uc.queue[0]
		uc.queue = uc.queue[1:]
	} else {
		uc.current = nil
	}
}

func (uc *UseConfirmation) GetCurrentPrompt() string {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.current != nil {
		return uc.current.Prompt
	}
	return ""
}

func (uc *UseConfirmation) GetCurrentExplanation() string {
	uc.mutex.Lock()
	defer uc.mutex.Unlock()

	if uc.current != nil {
		return uc.current.Explanation
	}
	return ""
}
