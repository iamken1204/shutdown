package shutdown

import (
	"os"
	"syscall"
	"time"
)

const defaultTimeout = 5 * time.Second

func defaultSignals() []os.Signal {
	return []os.Signal{
		os.Kill,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGKILL,
	}
}
