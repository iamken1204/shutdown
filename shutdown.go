package shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"
)

// Shutdown contains cleanup hooks.
type Shutdown struct {
	Timeout time.Duration
	signals []os.Signal
	hooks   []Hook
}

// New initializes a Shutdown instance.
func New() *Shutdown {
	return &Shutdown{
		Timeout: defaultTimeout,
		signals: defaultSignals(),
	}
}

// AddHook appends given hooks.
func (s *Shutdown) AddHook(hooks ...Hook) {
	s.hooks = append(s.hooks, hooks...)
}

// Listen blocks the program until received terminating signals, then trigger
// all hook functions within.
func (s *Shutdown) Listen() {
	quit := make(chan os.Signal)
	signal.Notify(quit, s.signals...)
	<-quit

	done := make(chan struct{})
	count := 0
	total := len(s.hooks)

	var ctx context.Context
	var cancel context.CancelFunc
	if s.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), s.Timeout)
	} else {
		ctx, cancel = context.WithCancel(context.Background())
	}

	go func() {
		if total == 0 {
			done <- struct{}{}
		}
		for _, hook := range s.hooks {
			err := hook.Cleanup()
			if err != nil {
				log.Println("shutdown err:", err)
			}
			done <- struct{}{}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			count++
			if count >= total {
				cancel()
			}
		}
	}
}
