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

	screen.Clear()

	// go Tick(*screen)

	cursor := Cursor{
		Symbol:      '?',
		SymbolStyle: tcell.Style.Foreground(tcell.StyleDefault, tcell.ColorRed),
		X:           5,
		Y:           5,
	}

	DrawCursorPosition(screen, cursor)
	screen.Flush()

	for {
		eventType := screen.PollEvent()
		switch event := eventType.(type) {
		case *tcell.EventKey:
			if event.Key() == tcell.KeyEscape || event.Rune() == 'q' {
				return
			}
			if event.Key() == tcell.KeyUp {
				ClearCursorPosition(screen, cursor)
				cursor = UpdateCursorPosition(screen, cursor, 0, -1)
				DrawCursorPosition(screen, cursor)
			}
			if event.Key() == tcell.KeyDown {
				ClearCursorPosition(screen, cursor)
				cursor = UpdateCursorPosition(screen, cursor, 0, 1)
				DrawCursorPosition(screen, cursor)
			}
			if event.Key() == tcell.KeyLeft {
				ClearCursorPosition(screen, cursor)
				cursor = UpdateCursorPosition(screen, cursor, -1, 0)
				DrawCursorPosition(screen, cursor)
			}
			if event.Key() == tcell.KeyRight {
				ClearCursorPosition(screen, cursor)
				cursor = UpdateCursorPosition(screen, cursor, 1, 0)
				DrawCursorPosition(screen, cursor)
			}

			screen.Flush()
		}
	}
}

// func Tick(screen gui.Screen) {
// 	ticker := time.Tick(10 * time.Millisecond)
// 	for range ticker {
// 		x, y := generatePosition()
// 		placeDot(screen, x, y)
// 		screen.Flush()
// 	}
// }

// func generatePosition() (int, int) {
// 	x := rand.IntN(190)
// 	y := rand.IntN(45)

// 	return x, y
// }

// func placeDot(screen gui.Screen, x, y int) {
// 	screen.SetContent(x, y, '.', nil, tcell.StyleDefault)
// }

func ClearCursorPosition(screen *gui.Screen, cursor Cursor) {
	screen.SetContent(cursor.X, cursor.Y, ' ', nil, tcell.StyleDefault)
}

func DrawCursorPosition(screen *gui.Screen, cursor Cursor) {
	screen.SetContent(cursor.X, cursor.Y, cursor.Symbol, nil, cursor.SymbolStyle)
}

func UpdateCursorPosition(screen *gui.Screen, cursor Cursor, xOffset, yOffset int) Cursor {
	screenX, screenY := screen.Size()

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
