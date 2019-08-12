package shutdown

// Hook handles cleanup function within.
type Hook interface {
	Cleanup() error
}

// Cleanup represents the pre-action before the program shuts down.
type Cleanup func() error

// hook is a pure function that implements Hook.
type hook func() error

func (h hook) Cleanup() error {
	return h()
}

// NewHook offers a convinient way to generate a Hook from a pure function.
func NewHook(fn func() error) Hook {
	return hook(fn)
}
