package closer

import (
	"context"
	"fmt"
	"strings"
	"sync"
)

type Closer struct {
	mu    sync.Mutex
	funcs []Func
}

type Func func(ctx context.Context) error

func (c *Closer) Add(f Func) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, f)
}

func (c *Closer) Close(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var (
		wg         = &sync.WaitGroup{}
		messagesCh = make(chan string, len(c.funcs))
	)

	defer close(messagesCh)

	wg.Add(len(c.funcs))

	go func() {
		for _, f := range c.funcs {
			go func(f Func) {
				defer wg.Done()

				if err := f(ctx); err != nil {
					messagesCh <- fmt.Sprintf("[!] %v", err)
				}
			}(f)
		}
	}()

	select {
	case <-ctx.Done():
		return fmt.Errorf("shutdown cancelled: %v", ctx.Err())
	default:
	}

	wg.Wait()

	msgs := make([]string, 0, len(c.funcs))

	for msg := range messagesCh {
		msgs = append(msgs, msg)
	}

	if len(msgs) > 0 {
		return fmt.Errorf(
			"shutdown finished with error(s): \n%s",
			strings.Join(msgs, "\n"),
		)
	}

	return nil
}
