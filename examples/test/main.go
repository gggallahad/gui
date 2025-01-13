package main

import (
	"log"

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
		X int
		Y int
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
		X: 0,
		Y: 0,
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

	drawCursorPosition(ctx)

	err = ctx.Flush()
	if err != nil {
		return
	}
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
			MoveCursor(ctx, 0, -1)
		}
		if event.Ch == 's' {
			MoveCursor(ctx, 0, 1)
		}
		if event.Ch == 'a' {
			MoveCursor(ctx, -1, 0)
		}
		if event.Ch == 'd' {
			MoveCursor(ctx, 1, 0)
		}

		// if event.Key == termbox.KeyArrowUp {
		// 	err := MoveCamera(ctx, 0, -1)
		// 	if err != nil {
		// 		return
		// 	}
		// }
		// if event.Key == termbox.KeyArrowDown {
		// 	err := MoveCamera(ctx, 0, 1)
		// 	if err != nil {
		// 		return
		// 	}
		// }
		// if event.Key == termbox.KeyArrowLeft {
		// 	err := MoveCamera(ctx, -1, 0)
		// 	if err != nil {
		// 		return
		// 	}
		// }
		// if event.Key == termbox.KeyArrowRight {
		// 	err := MoveCamera(ctx, 1, 0)
		// 	if err != nil {
		// 		return
		// 	}
		// }

		ctx.Flush()
	}
}

func MoveCursor(ctx *gui.Context, cursorPositionOffsetX, cursorPositionOffsetY int) {
	clearCursorPosition(ctx)
	updateCursorPosition(cursorPositionOffsetX, cursorPositionOffsetY)
	drawCursorPosition(ctx)
}

func clearCursorPosition(ctx *gui.Context) {
	ctx.SetCell(cursor.X-view.X, cursor.Y-view.Y, gui.DefaultCell)
}

func drawCursorPosition(ctx *gui.Context) {
	ctx.SetCell(cursor.X-view.X, cursor.Y-view.Y, cursor.Cell)
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

// func MoveCamera(ctx *gui.Context, viewPositionOffsetX, viewPositionOffsetY int) error {
// 	updateViewPosition(viewPositionOffsetX, viewPositionOffsetY)
// 	err := updateViewContent(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func updateViewPosition(viewPositionOffsetX, viewPositionOffsetY int) {
// 	view.X += viewPositionOffsetX
// 	view.Y += viewPositionOffsetY

// 	if view.X < 0 {
// 		view.X = 0
// 	}
// 	if view.Y < 0 {
// 		view.Y = 0
// 	}
// }

// func updateViewContent(ctx *gui.Context) error {
// 	ctx.SetViewPosition(view.X, view.Y)
// 	err := ctx.UpdateViewContent()
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
