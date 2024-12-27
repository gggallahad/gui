package gui

import (
	"context"
	"time"
)

type (
	Context struct {
		states []State

		context    context.Context
		cancelFunc context.CancelFunc
	}
)

func newContext() *Context {
	context, cancelFunc := context.WithCancel(context.Background())
	return &Context{
		context:    context,
		cancelFunc: cancelFunc,
	}
}

func (ctx *Context) newChildContext() *Context {
	context, cancelFunc := context.WithCancel(ctx)
	return &Context{
		context:    context,
		cancelFunc: cancelFunc,
	}
}

// имплементация интерфейса context.Context

func (ctx *Context) Deadline() (deadline time.Time, ok bool) {
	return ctx.context.Deadline()
}

func (ctx *Context) Done() <-chan struct{} {
	return ctx.context.Done()
}

func (ctx *Context) Err() error {
	return ctx.context.Err()
}

func (ctx *Context) Value(key any) any {
	return ctx.context.Value(key)
}

// отмена контекста

func (ctx *Context) Cancel() {
	ctx.cancelFunc()
}
