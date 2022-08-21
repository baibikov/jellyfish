package timeoutgroup

import (
	"context"

	"github.com/pkg/errors"
)

type Group struct {
	err    chan error
	closed bool
}

func New(ctx context.Context) *Group {
	g := &Group{
		err: make(chan error),
	}
	go func() {
		for range ctx.Done() {
			g.err <- errors.New("group ctx has done")
		}
	}()

	return &Group{
		err: make(chan error),
	}
}

func (g *Group) Go(f func() error) {
	go func() {
		if err := f(); err != nil {
			g.errWrite(err)
			return
		}

		g.errWrite(nil)
	}()
}

func (g *Group) errWrite(err error) {
	if g.closed {
		return
	}

	g.err <- err
}

func (g *Group) Wait() error {
	defer func() {
		if g.closed {
			return
		}
		close(g.err)
		g.closed = true
	}()
	return <-g.err
}
