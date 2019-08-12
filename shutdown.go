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
	ctx     context.Context
	cancel  context.CancelFunc
	signals []os.Signal
	hooks   []Hook
}

// New initializes a Shutdown instance.
func New() *Shutdown {
	ctx, cancel := context.WithCancel(context.Background())
	return &Shutdown{
		Timeout: defaultTimeout,
		ctx:     ctx,
		cancel:  cancel,
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

	go func() {
		time.Sleep(s.Timeout)
		s.cancel()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-done:
			count++
			if count >= total {
				s.cancel()
			}
		}
	}
}

// Context returns Shutdown's underlying Context.
func (s *Shutdown) Context() context.Context {
	return s.ctx
}
