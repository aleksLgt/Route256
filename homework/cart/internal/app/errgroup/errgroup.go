package errgroup

import (
	"context"
	"sync"
)

// Package errgroup provides synchronization, error propagation, and Context
// cancelation for groups of goroutines working on subtasks of a common task.
//
// [errgroup.Group] is related to [sync.WaitGroup] but adds handling of tasks
// returning errors.

type token struct{}

// A Group is a collection of goroutines working on subtasks that are part of
// the same overall task.
//
// A zero Group is valid, has no limit on the number of active goroutines,
// and does not cancel on error.
type Group struct {
	cancel func(error)

	wg sync.WaitGroup

	sem chan token

	errOnce sync.Once
	err     error
}

func (g *Group) done() {
	if g.sem != nil {
		<-g.sem
	}

	g.wg.Done()
}

// WithContext returns a new Group and an associated Context derived from ctx.
//
// The derived Context is canceled the first time a function passed to Go
// returns a non-nil error or the first time Wait returns, whichever occurs
// first.
func WithContext(ctx context.Context) (*Group, context.Context) {
	ctx, cancel := context.WithCancelCause(ctx)
	return &Group{cancel: cancel}, ctx
}

// Wait blocks until all function calls from the Go method have returned, then
// returns the first non-nil error (if any) from them.
func (g *Group) Wait() error {
	g.wg.Wait()

	if g.cancel != nil {
		g.cancel(g.err)
	}

	return g.err
}

// Go calls the given function in a new goroutine.
// It blocks until the new goroutine can be added without the number of
// active goroutines in the group exceeding the configured limit.
//
// The first call to return a non-nil error cancels the group's context, if the
// group was created by calling WithContext. The error will be returned by Wait.
func (g *Group) Go(f func() error) {
	if g.sem != nil {
		g.sem <- token{}
	}

	g.wg.Add(1)

	goFunc := func() {
		go func() {
			defer g.done()

			if err := f(); err != nil {
				g.errOnce.Do(func() {
					g.err = err
					if g.cancel != nil {
						g.cancel(g.err)
					}
				})
			}
		}()
	}

	goFunc()
}
