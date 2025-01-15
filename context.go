package gui

import (
	"context"
	"math"
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

func (ctx *Context) SetViewPosition(viewPositionX, viewPositionY int) {
	*ctx.viewPositionX = viewPositionX
	*ctx.viewPositionY = viewPositionY
}

func (ctx *Context) UpdateViewContent() error {
	foregroundAttribute := ctx.defaultCell.Foreground.toAttribute()
	backgroundAttribute := ctx.defaultCell.Background.toAttribute()

	// очистка видимой области
	err := ctx.clearTermboxScreen(foregroundAttribute, backgroundAttribute)
	if err != nil {
		return err
	}

	// перерисовка клеток которые попадают в видимую область
	for y := *ctx.viewPositionY; y < len(*ctx.cells); y++ {
		for x := *ctx.viewPositionX; x < len((*ctx.cells)[y]); x++ {
			ctx.setTermboxCell(x, y, (*ctx.cells)[y][x])
		}
	}

	return nil
}

func (ctx *Context) SetCell(x, y int, cell Cell) {
	ctx.setlocalCell(x, y, cell)

	ctx.setTermboxCell(x, y, cell)
}

func (ctx *Context) setlocalCell(x, y int, cell Cell) {
	// добавление недостающих строк
	maxY := len(*ctx.cells) - 1
	if y > maxY {
		*ctx.cells = growSlice(*ctx.cells, y+1)
		newRows := y - maxY
		for range newRows {
			*ctx.cells = append(*ctx.cells, []Cell{})
		}
	}

	// добавление недостающего столбца в нужную строку
	maxX := len((*ctx.cells)[y]) - 1
	if x > maxX {
		(*ctx.cells)[y] = growSlice((*ctx.cells)[y], x+1)
		newColumns := x - maxX
		for range newColumns {
			(*ctx.cells)[y] = append((*ctx.cells)[y], *ctx.defaultCell)
		}
	}

	(*ctx.cells)[y][x] = cell
}

func (ctx *Context) setTermboxCell(x, y int, cell Cell) {
	// отрисовка клетки со смещением по видимой области

	x -= *ctx.viewPositionX
	y -= *ctx.viewPositionY

	foregroundAttribute := cell.Foreground.toAttribute()
	backgroundAttribute := cell.Background.toAttribute()

	termbox.SetCell(x, y, cell.Symbol, foregroundAttribute, backgroundAttribute)
}

func (ctx *Context) SetRow(y int, cells []Cell) {
	ctx.ClearRow(y)

	ctx.setLocalRow(y, cells)

	ctx.setTermboxRow(y, cells)
}

func (ctx *Context) setLocalRow(y int, cells []Cell) {
	// добавление недостающих строк
	maxY := len(*ctx.cells) - 1
	if y > maxY {
		*ctx.cells = growSlice(*ctx.cells, y+1)
		newRows := y - maxY
		for range newRows {
			*ctx.cells = append(*ctx.cells, []Cell{})
		}
	}

	(*ctx.cells)[y] = cells
}

func (ctx *Context) setTermboxRow(y int, cells []Cell) {
	// отрисовка всех клеток новой строки
	for x := *ctx.viewPositionX; x < len(cells); x++ {
		ctx.setTermboxCell(x, y, cells[x])
	}
}

func (ctx *Context) SetColumn(x int, cells []Cell) {
	ctx.ClearColumn(x)

	ctx.setLocalColumn(x, cells)

	ctx.setTermboxColumn(x, cells)
}

func (ctx *Context) setLocalColumn(x int, cells []Cell) {
	// добавление недостающих строк
	maxY := len(*ctx.cells) - 1
	newMaxY := len(cells) - 1
	if newMaxY > maxY {
		*ctx.cells = growSlice(*ctx.cells, newMaxY+1)
		newRows := newMaxY - maxY
		for range newRows {
			*ctx.cells = append(*ctx.cells, make([]Cell, 0, len(cells)))
		}
	}

	// добавление недостающих столбцов в нужные строки
	for y := range len(cells) {
		maxX := len((*ctx.cells)[y]) - 1
		if x > maxX {
			(*ctx.cells)[y] = growSlice((*ctx.cells)[y], x+1)
			newColumns := x - maxX
			for range newColumns {
				(*ctx.cells)[y] = append((*ctx.cells)[y], *ctx.defaultCell)
			}
		}
	}

	// устанавливаем значения в столбец. Можно это было делать и в предыдущем цикле, но я не хочу засорять логику. Исправлю когда буду релизить
	for y := range len(cells) {
		(*ctx.cells)[y][x] = cells[y]
	}
}

func (ctx *Context) setTermboxColumn(x int, cells []Cell) {
	// отрисовка всех клеток нового столбца
	for y := *ctx.viewPositionY; y < len(cells); y++ {
		ctx.setTermboxCell(x, y, cells[y])
	}
}

func (ctx *Context) GetCell(x, y int) Cell {
	if y >= len(*ctx.cells) || x >= len((*ctx.cells)[y]) {
		return *ctx.defaultCell
	}

	localCell := ctx.getLocalCell(x, y)

	return localCell
}

func (ctx *Context) getLocalCell(x, y int) Cell {
	cell := (*ctx.cells)[y][x]

	return cell
}

// func (ctx *Context) getTermboxCell(x, y int) Cell {
// 	termboxCell := termbox.GetCell(x, y)

// 	!!!var foregroundColor Color
// 	!!!var backgroundColor Color

// 	!!!foregroundColor.fromAttribute(termboxCell.Fg)
// 	!!!backgroundColor.fromAttribute(termboxCell.Bg)

// 	!!!cell := Cell{
// 	!!!	Symbol:     termboxCell.Ch,
// 	!!!	Foreground: foregroundColor,
// 	!!!	Background: backgroundColor,
// 	!!!}

// 	return cell
// }

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

func (ctx *Context) ClearRow(y int) {
	ctx.clearTermboxRow(y)

	ctx.clearLocalRow(y)
}

func (ctx *Context) clearLocalRow(y int) {
	maxY := len((*ctx.cells)[y]) - 1
	if y > maxY {
		return
	}

	(*ctx.cells)[y] = nil
}

func (ctx *Context) clearTermboxRow(y int) {
	maxY := len((*ctx.cells)[y]) - 1
	if y > maxY {
		return
	}

	maxX := len((*ctx.cells)[y])
	for x := range maxX {
		ctx.setTermboxCell(x, y, *ctx.defaultCell)
	}
}

func (ctx *Context) ClearColumn(x int) {
	ctx.clearTermboxColumn(x)

	ctx.clearLocalColumn(x)
}

func (ctx *Context) clearLocalColumn(x int) {
	for y := range len((*ctx.cells)) {
		maxX := len((*ctx.cells)[y]) - 1
		if x > maxX {
			continue
		}

		(*ctx.cells)[y][x] = *ctx.defaultCell
	}
}

func (ctx *Context) clearTermboxColumn(x int) {
	for y := range len((*ctx.cells)) {
		maxX := len((*ctx.cells)[y]) - 1
		if x > maxX {
			continue
		}

		ctx.setTermboxCell(x, y, *ctx.defaultCell)
	}
}

// state

func (ctx *Context) Abort() {
	ctx.handlerIndex = math.MaxInt - 1
}

// user util

func (ctx *Context) Kill() {
	ctx.killChannel <- struct{}{}
}

// util

func (ctx *Context) getCurrentState() State {
	state := (*ctx.states)[*ctx.stateIndex]

	return state
}

func (ctx *Context) resetData(context *Context) {
	if ctx.cancelFunc != nil {
		ctx.cancelFunc()
	}

	ctx.cancelFunc = context.cancelFunc
	ctx.handlerIndex = 0
}

func (ctx *Context) addHandlerIndex() {
	ctx.handlerIndex++
}

func (ctx *Context) getHandlerIndex() int {
	handlerIndex := ctx.handlerIndex

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
