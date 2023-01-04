package baseenv

import (
	"context"
	"io"
)

type contextCloser func()

func (c contextCloser) Close() error {
	c()
	return nil
}

// CloserFromCancel context.CancelFunc is sealed as the io.Closer
func CloserFromCancel(cancel context.CancelFunc) io.Closer {
	return contextCloser(cancel)
}

type funcCloser struct {
	fn func() error
}

func (f funcCloser) Close() error {
	return f.fn()
}

// CloserFromFunc fn is sealed as the io.Closer
func CloserFromFunc(fn func() error) io.Closer {
	return funcCloser{fn}
}
