package main

import (
	"log"

	"github.com/gdamore/tcell/v2"
	"github.com/gggallahad/gui"
)

type (
	Cursor struct {
		Symbol      rune
		SymbolStyle tcell.Style
		X           int
		Y           int
	}
)

var (
	cursor Cursor = Cursor{
		Symbol:      '?',
		SymbolStyle: tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorRed),
		X:           5,
		Y:           5,
	}
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

	screen.BindInitHandler(InitHandler)

	screen.BindHandlers(gui.NoState, KillMiddleware, NoStateHandler)

	screen.Run()
}

func InitHandler(ctx *gui.Context) {
	ctx.Clear()
	DrawCursorPosition(ctx, cursor)
	ctx.Flush()
}

func KillMiddleware(ctx *gui.Context, eventType tcell.Event) {
	switch event := eventType.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
			ctx.Abort()
			ctx.Kill()
		}
	}
}

func NoStateHandler(ctx *gui.Context, eventType tcell.Event) {
	switch event := eventType.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyUp {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 0, -1)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Key() == tcell.KeyDown {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 0, 1)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Key() == tcell.KeyLeft {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, -1, 0)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Key() == tcell.KeyRight {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 1, 0)
			DrawCursorPosition(ctx, cursor)
		}

		ctx.Flush()
	}
}

func ClearCursorPosition(ctx *gui.Context, cursor Cursor) {
	ctx.SetContent(cursor.X, cursor.Y, ' ', nil, tcell.StyleDefault)
}

func DrawCursorPosition(ctx *gui.Context, cursor Cursor) {
	ctx.SetContent(cursor.X, cursor.Y, cursor.Symbol, nil, cursor.SymbolStyle)
}

func UpdateCursorPosition(ctx *gui.Context, cursor Cursor, xOffset, yOffset int) Cursor {
	screenX, screenY := ctx.Size()

	cursor.X += xOffset
	cursor.Y += yOffset

	if cursor.X >= screenX {
		cursor.X = 0
	}

	if cursor.Y >= screenY {
		cursor.Y = 0
	}

	if cursor.X < 0 {
		cursor.X = screenX - 1
	}

	if cursor.Y < 0 {
		cursor.Y = screenY - 1
	}

	return cursor
}
