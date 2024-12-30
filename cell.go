package gui

import "github.com/gdamore/tcell/v2"

type (
	cell struct {
		symbol    rune
		combining []rune
		style     tcell.Style
	}
)
