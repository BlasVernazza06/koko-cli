package ui

import (
	"fmt"
	"time"
)

type Spinner struct {
	message string
	stop    chan struct{}
}

func NewSpinner(message string) *Spinner {
	return &Spinner{
		message: message,
		stop:    make(chan struct{}),
	}
}

func (s *Spinner) Start() {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	go func() {
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		i := 0
		for {
			select {
			case <-s.stop:
				return
			case <-ticker.C:
				fmt.Printf("\r\033[36m%s\033[0m %s", frames[i%len(frames)], s.message)
				i++
			}
		}
	}()
}

func (s *Spinner) Stop(success bool) {
	close(s.stop)
	// Clear line first
	fmt.Print("\r\033[K")
	if success {
		fmt.Printf("\033[32m✓\033[0m %s\n", s.message)
	} else {
		fmt.Printf("\033[31m✗\033[0m %s\n", s.message)
	}
}
