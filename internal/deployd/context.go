package deployd

import (
	"time"

	"github.com/iceber/iouring-go"
)

type (
	Context struct {
		ring *iouring.IOURing
		done chan struct{}
	}
)

func NewContext() (ctx *Context, err error) {
	ring, err := iouring.New(64)
	if err != nil {
		return nil, err
	}

	ctx = &Context{
		ring: ring,
	}

	return ctx, nil
}

func (ctx *Context) Close() error {
	err := ctx.ring.Close()
	close(ctx.done)
	return err
}

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return time.Time{}, false
}

func (ctx *Context) Done() chan struct{} {
	return ctx.done
}

func (ctx *Context) Err() error {
	return nil
}

func (ctx *Context) Value(key any) any {
	return nil
}
