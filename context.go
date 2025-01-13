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

func newContext(defaultCell Cell) (*Context, error) {
	cells := make([][]Cell, 0)

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

func (ctx *Context) UpdateViewPosition(viewPositionOffsetX, viewPositionOffsetY int) {
	*ctx.viewPositionX += viewPositionOffsetX
	*ctx.viewPositionY += viewPositionOffsetY
}

func (ctx *Context) SetViewPosition(viewPositionX, viewPositionY int) {
	*ctx.viewPositionX = viewPositionX
	*ctx.viewPositionY = viewPositionY
}

func (ctx *Context) UpdateViewContent() error {
	foregroundAttribute := ctx.defaultCell.Foreground.toAttribute()
	backgroundAttribute := ctx.defaultCell.Background.toAttribute()

	err := ctx.clearTermboxScreen(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	for rowIndex := *ctx.viewPositionY; rowIndex < len(*ctx.cells); rowIndex++ {
		for columnIndex := *ctx.viewPositionX; columnIndex < len((*ctx.cells)[rowIndex]); columnIndex++ {
			cell := (*ctx.cells)[rowIndex][columnIndex]
			positionWithViewOffsetX := columnIndex - *ctx.viewPositionX
			positionWithViewOffsetY := rowIndex - *ctx.viewPositionY
			ctx.setTermboxCell(positionWithViewOffsetX, positionWithViewOffsetY, cell)
		}
	}

	return nil
}

func (ctx *Context) SetCell(x, y int, cell Cell) {
	ctx.setlocalCell(x, y, cell)

	ctx.setTermboxCell(x, y, cell)
}

func (ctx *Context) setlocalCell(x, y int, cell Cell) {
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

func (ctx *Context) setTermboxCell(x, y int, cell Cell) {
	foregroundAttribute := cell.Foreground.toAttribute()
	backgroundAttribute := cell.Background.toAttribute()

	termbox.SetCell(x, y, cell.Symbol, foregroundAttribute, backgroundAttribute)
}

func (ctx *Context) GetCell(x, y int) Cell {
	// ?
	if y >= len(*ctx.cells) || x >= len((*ctx.cells)[y]) {
		termboxCell := ctx.getTermboxCell(x, y)

		return termboxCell
	}

	localCell := ctx.getLocalCell(x, y)

	return localCell
}

func (ctx *Context) getLocalCell(x, y int) Cell {
	cell := (*ctx.cells)[y][x]

	return cell
}

func (ctx *Context) getTermboxCell(x, y int) Cell {
	termboxCell := termbox.GetCell(x, y)

	var foregroundColor Color
	var backgroundColor Color

	foregroundColor.fromAttribute(termboxCell.Fg)
	backgroundColor.fromAttribute(termboxCell.Bg)

	cell := Cell{
		Symbol:     termboxCell.Ch,
		Foreground: foregroundColor,
		Background: backgroundColor,
	}

	return cell
}

func (ctx *Context) Flush() error {
	err := termbox.Flush()
	if err != nil {
		return err
	}

	return nil
}

// clear

func (ctx *Context) Clear() error {
	foregroundAttribute := ctx.defaultCell.Foreground.toAttribute()
	backgroundAttribute := ctx.defaultCell.Background.toAttribute()

	err := ctx.clear(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) ClearWithColor(foreground, background Color) error {
	foregroundAttribute := foreground.toAttribute()
	backgroundAttribute := foreground.toAttribute()

	err := ctx.clear(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) clear(foregroundAttribute, backgroundAttribute termbox.Attribute) error {
	ctx.clearLocalScreen()

	err := ctx.clearTermboxScreen(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	return nil
}

func (ctx *Context) clearLocalScreen() {
	*ctx.cells = nil
}

func (ctx *Context) clearTermboxScreen(foregroundAttribute, backgroundAttribute termbox.Attribute) error {
	err := termbox.Clear(foregroundAttribute, backgroundAttribute)
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
