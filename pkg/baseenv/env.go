package baseenv

import (
	"context"
	"io"
	"log"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Environment struct {
	context    context.Context
	cancel     context.CancelFunc
	closeWG    *sync.WaitGroup
	cancelWG   *sync.WaitGroup
	logger     *zap.Logger
	destroyers []io.Closer
}

func NewBaseEnv(ctx context.Context, logger *zap.Logger) *Environment {
	ctx, cancel := context.WithCancel(ctx)
	env := &Environment{
		context:  ctx,
		cancel:   cancel,
		closeWG:  new(sync.WaitGroup),
		cancelWG: new(sync.WaitGroup),
		logger:   logger,
		destroyers: []io.Closer{
			CloserFromCancel(cancel),
			CloserFromFunc(logger.Sync),
		},
	}

	return env
}

func (env *Environment) Logger() *zap.Logger {
	return env.logger
}

// CancelWG is used to block the waitGroup of cancel
func (env *Environment) CancelWG() *sync.WaitGroup {
	return env.cancelWG
}

// CloseWG is used to block the waitGroup of close。
func (env *Environment) CloseWG() *sync.WaitGroup {
	return env.closeWG
}

func (env *Environment) Wg() *sync.WaitGroup {
	return env.closeWG
}

func (env *Environment) Context() context.Context {
	return env.context
}

func (env *Environment) CancelRootContext() {
	timeoutErr := WaitWithTimeout(env.cancelWG, 30*time.Second)
	log.Println(timeoutErr)
	env.cancel()
}

func (env *Environment) WithContext(ctx context.Context) *Environment {
	next := &Environment{}
	*next = *env
	next.context = ctx

	return next
}

func (env *Environment) AddCloser(closer io.Closer) {
	env.destroyers = append(env.destroyers, closer)
}

func (env *Environment) Close() (err error) {
	env.CancelRootContext()
	wgErr := WaitWithTimeout(env.closeWG, 30*time.Second)
	log.Println(wgErr.Error())
	for _, destroyer := range env.destroyers {
		e := destroyer.Close()
		if e != nil {
			err = e
			log.Printf("%+v\n", e)
		}
	}

	return
}

func WaitWithTimeout(wg *sync.WaitGroup, d time.Duration) error {
	waitCh := GoChan(func() error {
		wg.Wait()
		return nil
	})

	select {
	case <-waitCh:
		return nil
	case <-time.After(d):
		return ErrWGTimeout
	}
}

const ErrWGTimeout TimeoutError = "waitGroup timeout"

type TimeoutError string

func (s TimeoutError) Error() string {
	return string(s)
}

type Job func() error

func GoChan(job Job) <-chan error {
	ch := make(chan error, 1)
	go func() {
		ch <- job()
	}()

	return ch
}
