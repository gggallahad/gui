package gui

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

type (
	Context struct {
		tcellScreen tcell.Screen

		stateIndex int
		states     []State

		handlerIndex int

		context    context.Context
		cancelFunc context.CancelFunc

		killChannel chan struct{}

		mutex sync.RWMutex
	}
)

func newContext() (*Context, error) {
	tcellScreen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	context, cancelFunc := context.WithCancel(context.Background())

	initStates := []State{NoState}

	ctx := Context{
		tcellScreen: tcellScreen,
		states:      initStates,
		context:     context,
		cancelFunc:  cancelFunc,
		killChannel: make(chan struct{}, 1),
	}

	return &ctx, nil
}

func (ctx *Context) newChildContext() *Context {
	context, cancelFunc := context.WithCancel(ctx)
	return &Context{
		tcellScreen: ctx.tcellScreen,
		states:      ctx.states,
		context:     context,
		cancelFunc:  cancelFunc,
		killChannel: ctx.killChannel,
	}
}

// draw

func (ctx *Context) SetContent(x, y int, symbol rune, combining []rune, style tcell.Style) {
	ctx.tcellScreen.SetContent(x, y, symbol, combining, style)
}

func (ctx *Context) GetContent(x, y int) (rune, []rune, tcell.Style, int) {
	symbol, combining, style, width := ctx.tcellScreen.GetContent(x, y)

	return symbol, combining, style, width
}

func (ctx *Context) Flush() {
	ctx.tcellScreen.Show()
}

func (ctx *Context) Fill(symbol rune, style tcell.Style) {
	ctx.tcellScreen.Fill(symbol, style)
}

func (ctx *Context) Clear() {
	ctx.tcellScreen.Clear()
}

// user util

func (ctx *Context) HideCursor() {
	ctx.tcellScreen.HideCursor()
}

func (ctx *Context) ShowCursor(x, y int) {
	ctx.tcellScreen.ShowCursor(x, y)
}

func (ctx *Context) Size() (int, int) {
	screenX, screenY := ctx.tcellScreen.Size()

	return screenX, screenY
}

func (ctx *Context) Kill() {
	ctx.killChannel <- struct{}{}
}

// state

func (ctx *Context) Abort() {
	ctx.mutex.Lock()

	ctx.handlerIndex = math.MaxInt - 1

	ctx.mutex.Unlock()
}

// util

func (ctx *Context) getCurrentState() State {
	ctx.mutex.RLock()

	state := ctx.states[ctx.stateIndex]

	ctx.mutex.RUnlock()

	return state
}

func (ctx *Context) resetData(context *Context) {
	ctx.mutex.Lock()

	if ctx.cancelFunc != nil {
		ctx.cancelFunc()
	}

	ctx.cancelFunc = context.cancelFunc
	ctx.handlerIndex = 0

	ctx.mutex.Unlock()
}

func (ctx *Context) addHandlerIndex() {
	ctx.mutex.Lock()

	ctx.handlerIndex++

	ctx.mutex.Unlock()
}

func (ctx *Context) getHandlerIndex() int {
	ctx.mutex.RLock()

	handlerIndex := ctx.handlerIndex

	ctx.mutex.RUnlock()

	return handlerIndex
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
