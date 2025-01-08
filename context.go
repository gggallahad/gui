package gui

import (
	"context"
	"math"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

type (
	Context struct {
		cells       *[][]Cell
		defaultCell *Cell

		viewPositionX *int
		viewPositionY *int

		viewSizeX *int
		viewSizeY *int

		stateIndex *int
		states     *[]State

		killChannel chan struct{}

		handlerIndex int
		context      context.Context
		cancelFunc   context.CancelFunc

		mutex sync.RWMutex
	}
)

func newContext() (*Context, error) {
	cells := make([][]Cell, 0)
	defaultCell := DefaultCell

	viewPositionX := 0
	viewPositionY := 0

	viewSizeX := 0
	viewSizeY := 0

	stateIndex := 0
	states := []State{NoState}

	killChannel := make(chan struct{})

	handlerIndex := 0
	context, cancelFunc := context.WithCancel(context.Background())

	ctx := Context{
		cells:         &cells,
		defaultCell:   &defaultCell,
		viewPositionX: &viewPositionX,
		viewPositionY: &viewPositionY,
		viewSizeX:     &viewSizeX,
		viewSizeY:     &viewSizeY,
		stateIndex:    &stateIndex,
		states:        &states,
		killChannel:   killChannel,
		handlerIndex:  handlerIndex,
		context:       context,
		cancelFunc:    cancelFunc,
	}

	return &ctx, nil
}

func (ctx *Context) newChildContext() *Context {
	handlerIndex := 0
	context, cancelFunc := context.WithCancel(ctx)

	childContext := Context{
		cells:         ctx.cells,
		defaultCell:   ctx.defaultCell,
		viewPositionX: ctx.viewPositionX,
		viewPositionY: ctx.viewPositionY,
		viewSizeX:     ctx.viewSizeX,
		viewSizeY:     ctx.viewSizeY,
		stateIndex:    ctx.stateIndex,
		states:        ctx.states,
		killChannel:   ctx.killChannel,
		handlerIndex:  handlerIndex,
		context:       context,
		cancelFunc:    cancelFunc,
	}

	return &childContext
}

// draw

func (ctx *Context) UpdateView(viewPositionOffsetX, viewPositionOffsetY int) {
	*ctx.viewPositionX += viewPositionOffsetX
	*ctx.viewPositionY += viewPositionOffsetY

	if *ctx.viewPositionX < 0 {
		*ctx.viewPositionX = 0
	}
	if *ctx.viewPositionY < 0 {
		*ctx.viewPositionY = 0
	}
}

func (ctx *Context) SetView(viewPositionX, viewPositionY int) {
	*ctx.viewPositionX = viewPositionX
	*ctx.viewPositionY = viewPositionY

	if *ctx.viewPositionX < 0 {
		*ctx.viewPositionX = 0
	}
	if *ctx.viewPositionY < 0 {
		*ctx.viewPositionY = 0
	}
}

// func (ctx *Context) UpdateViewContent() {
// 	if ctx.viewPositionY < len(*ctx.cells) {
// 		for row := range *ctx.cells {
// 			if ctx.viewPositionX < len(*ctx.cells[row]) {

// 			}
// 		}
// 	}
// }

func (ctx *Context) SetCell(x, y int, cell Cell) {
	ctx.setLocalCell(x, y, cell)

	foregroundAttribute := cell.Foreground.toAttribute()
	backgroundAttribute := cell.Background.toAttribute()

	termbox.SetCell(x-*ctx.viewPositionX, y-*ctx.viewPositionY, cell.Symbol, foregroundAttribute, backgroundAttribute)
}

func (ctx *Context) setLocalCell(x, y int, cell Cell) {
	rowsMaxIndex := len(*ctx.cells) - 1
	if y > rowsMaxIndex {
		*ctx.cells = growSlice(*ctx.cells, y+1)
		newRows := y - rowsMaxIndex
		for range newRows {
			*ctx.cells = append(*ctx.cells, []Cell{})
		}
	}

	columnsMaxIndex := len((*ctx.cells)[y]) - 1
	if x > columnsMaxIndex {
		(*ctx.cells)[y] = growSlice((*ctx.cells)[y], x+1)
		newColumns := x - columnsMaxIndex
		for range newColumns {
			(*ctx.cells)[y] = append((*ctx.cells)[y], *ctx.defaultCell)
		}
	}

	(*ctx.cells)[y][x] = cell
}

func (ctx *Context) GetContent(x, y int) Cell {
	if y < len(*ctx.cells)-1 || x < len((*ctx.cells)[y]) {
		return *ctx.defaultCell
	}

	return (*ctx.cells)[y][x]
}

func (ctx *Context) Flush() error {
	err := termbox.Flush()
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) ClearWithColor(foreground, background Color) error {
	foregroundAttribute := foreground.toAttribute()
	backgroundAttribute := foreground.toAttribute()

	err := termbox.Clear(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) Clear() error {
	err := termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if err != nil {
		return err
	}

	return nil
}

// user util

func (ctx *Context) HideCursor() {
	termbox.HideCursor()
}

func (ctx *Context) ShowCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (ctx *Context) Size() (int, int) {
	x, y := termbox.Size()

	return x, y
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

	state := (*ctx.states)[*ctx.stateIndex]

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
