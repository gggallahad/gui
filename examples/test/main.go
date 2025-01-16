package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gggallahad/gui"
	"github.com/nsf/termbox-go"
)

type (
	Cursor struct {
		X    int
		Y    int
		Cell gui.Cell
	}

	View struct {
		CurrentX  int
		CurrentY  int
		PreviousX int
		PreviousY int
	}
)

var (
	cursor Cursor = Cursor{
		X: 5,
		Y: 5,
		Cell: gui.Cell{
			Symbol: '?',
			Foreground: gui.Color{
				R: 255,
				G: 0,
				B: 0,
			},
			Background: gui.Color{
				R: -1,
				G: -1,
				B: -1,
			},
		},
	}

	view View = View{
		CurrentX:  0,
		CurrentY:  0,
		PreviousX: 0,
		PreviousY: 0,
	}

	setRowCell gui.Cell = gui.Cell{
		Symbol: ' ',
		Background: gui.Color{
			R: 0,
			G: 255,
			B: 0,
		},
	}

	setColumnCell gui.Cell = gui.Cell{
		Symbol: ' ',
		Background: gui.Color{
			R: 0,
			G: 0,
			B: 255,
		},
	}

	statusLineForeground gui.Color = gui.Color{
		R: 1,
		G: 229,
		B: 210,
	}

	statusLineBackground gui.Color = gui.Color{
		R: 21,
		G: 21,
		B: 21,
	}

	statusLineOffsetX int = 3
	statusLineOffsetY int = 40
)

func main() {
	screen, err := gui.NewScreen()
	if err != nil {
		log.Println(err)
		return
	}

	err = screen.Init()
	if err != nil {
		log.Println(err)
		return
	}
	defer screen.Close()

	screen.BindInitHandlers(InitHandler)

	screen.BindGlobalMiddlewares(KillMiddleware)

	screen.BindGlobalPostwares(DrawStatusLine, SetVariables)

	screen.BindHandlers(gui.NoState, NoStateHandler)

	screen.Run()
}

func InitHandler(ctx *gui.Context) {
	err := ctx.Clear()
	if err != nil {
		return
	}

	drawCursorPosition(ctx)
	DrawStatusLine(ctx, nil)

	err = ctx.Flush()
	if err != nil {
		return
	}
}

func KillMiddleware(ctx *gui.Context, eventType gui.Event) {
	switch event := eventType.(type) {
	case gui.EventKey:
		if event.Key == termbox.KeyEsc || event.Symbol == 'q' {
			ctx.Abort()
			ctx.Kill()
		}
	}
}

func NoStateHandler(ctx *gui.Context, eventType gui.Event) {
	switch event := eventType.(type) {
	case gui.EventKey:
		if event.Symbol == 'w' {
			MoveCursor(ctx, 0, -1)
		}
		if event.Symbol == 's' {
			MoveCursor(ctx, 0, 1)
		}
		if event.Symbol == 'a' {
			MoveCursor(ctx, -1, 0)
		}
		if event.Symbol == 'd' {
			MoveCursor(ctx, 1, 0)
		}

		if event.Key == termbox.KeyArrowUp {
			err := MoveCamera(ctx, 0, -1)
			if err != nil {
				return
			}
		}
		if event.Key == termbox.KeyArrowDown {
			err := MoveCamera(ctx, 0, 1)
			if err != nil {
				return
			}
		}
		if event.Key == termbox.KeyArrowLeft {
			err := MoveCamera(ctx, -1, 0)
			if err != nil {
				return
			}
		}
		if event.Key == termbox.KeyArrowRight {
			err := MoveCamera(ctx, 1, 0)
			if err != nil {
				return
			}
		}

		if event.Symbol == 'r' {
			SetRow(ctx)
		}

		if event.Symbol == 'c' {
			SetColumn(ctx)
		}

		if event.Symbol == 't' {
			SetText(ctx)
		}

		// ctx.Flush()
	}
}

func MoveCursor(ctx *gui.Context, cursorPositionOffsetX, cursorPositionOffsetY int) {
	clearCursorPosition(ctx)
	updateCursorPosition(cursorPositionOffsetX, cursorPositionOffsetY)
	drawCursorPosition(ctx)
}

func clearCursorPosition(ctx *gui.Context) {
	ctx.SetCell(cursor.X, cursor.Y, gui.DefaultCell)
}

func drawCursorPosition(ctx *gui.Context) {
	ctx.SetCell(cursor.X, cursor.Y, cursor.Cell)
}

func updateCursorPosition(cursorPositionOffsetX, cursorPositionOffsetY int) {
	cursor.X += cursorPositionOffsetX
	cursor.Y += cursorPositionOffsetY

	if cursor.X < 0 {
		cursor.X = 0
	}

	if cursor.Y < 0 {
		cursor.Y = 0
	}
}

func MoveCamera(ctx *gui.Context, viewPositionOffsetX, viewPositionOffsetY int) error {
	updateViewPosition(viewPositionOffsetX, viewPositionOffsetY)
	err := updateViewContent(ctx)
	if err != nil {
		return err
	}

	return nil
}

func updateViewPosition(viewPositionOffsetX, viewPositionOffsetY int) {
	view.CurrentX += viewPositionOffsetX
	view.CurrentY += viewPositionOffsetY

	if view.CurrentX < 0 {
		view.CurrentX = 0
	}
	if view.CurrentY < 0 {
		view.CurrentY = 0
	}
}

func updateViewContent(ctx *gui.Context) error {
	ctx.SetViewPosition(view.CurrentX, view.CurrentY)
	err := ctx.UpdateViewContent()
	if err != nil {
		return err
	}

	return nil
}

func SetRow(ctx *gui.Context) {
	rowCells := make([]gui.Cell, 0, cursor.X)
	for range cursor.X {
		rowCells = append(rowCells, setRowCell)
	}

	rowCells = append(rowCells, cursor.Cell)

	ctx.SetRow(cursor.Y, rowCells)
}

func SetColumn(ctx *gui.Context) {
	columnCells := make([]gui.Cell, 0, cursor.Y)
	for range cursor.Y {
		columnCells = append(columnCells, setColumnCell)
	}

	columnCells = append(columnCells, cursor.Cell)

	ctx.SetColumn(cursor.X, columnCells)
}

func SetText(ctx *gui.Context) {
	ctx.SetText(cursor.X+1, cursor.Y, "text", gui.DefaultCell.Foreground, gui.DefaultCell.Background)
}

func DrawStatusLine(ctx *gui.Context, eventType gui.Event) {
	ctx.ClearRow(view.PreviousY + statusLineOffsetY)
	spaceBetweenTypesCount := 5
	spaceBetweenElementsCount := 3
	spaceBetweenTypesString := strings.Repeat(" ", spaceBetweenTypesCount)
	spaceBetweenElementsString := strings.Repeat(" ", spaceBetweenElementsCount)
	text := fmt.Sprintf("cursorX: %d%scursorY: %d%scameraX: %d%scameraY: %d", cursor.X, spaceBetweenElementsString, cursor.Y, spaceBetweenTypesString, view.CurrentX, spaceBetweenElementsString, view.CurrentY)
	ctx.SetText(view.CurrentX+statusLineOffsetX, view.CurrentY+statusLineOffsetY, text, statusLineForeground, statusLineBackground)
	ctx.Flush()
}

func SetVariables(ctx *gui.Context, eventType gui.Event) {
	view.PreviousX = view.CurrentX
	view.PreviousY = view.CurrentY
}
