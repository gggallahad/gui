package main

import (
	"log"

	"github.com/gggallahad/gui"
	"github.com/nsf/termbox-go"
)

type (
	Cursor struct {
		Cell gui.Cell
		X    int
		Y    int
	}
)

var (
	cursor Cursor = Cursor{
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
		X: 5,
		Y: 5,
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
	err := ctx.Clear()
	if err != nil {
		return
	}

	DrawCursorPosition(ctx, cursor)
	ctx.Flush()
}

func KillMiddleware(ctx *gui.Context, event termbox.Event) {
	switch event.Type {
	case termbox.EventKey:
		if event.Key == termbox.KeyEsc || event.Ch == 'q' {
			ctx.Abort()
			ctx.Kill()
		}
	}
}

func NoStateHandler(ctx *gui.Context, event termbox.Event) {
	switch event.Type {
	case termbox.EventKey:
		if event.Ch == 'w' {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 0, -1)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Ch == 's' {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 0, 1)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Ch == 'a' {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, -1, 0)
			DrawCursorPosition(ctx, cursor)
		}
		if event.Ch == 'd' {
			ClearCursorPosition(ctx, cursor)
			cursor = UpdateCursorPosition(ctx, cursor, 1, 0)
			DrawCursorPosition(ctx, cursor)
		}

		ctx.Flush()
	}
}

func ClearCursorPosition(ctx *gui.Context, cursor Cursor) {
	ctx.SetCell(cursor.X, cursor.Y, gui.DefaultCell)
}

func DrawCursorPosition(ctx *gui.Context, cursor Cursor) {
	ctx.SetCell(cursor.X, cursor.Y, cursor.Cell)
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
