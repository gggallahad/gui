package gui

import "github.com/gdamore/tcell/v2"

type (
	cell struct {
		symbol    rune
		combining []rune
		style     tcell.Style
	}
)

var defaultCell cell = cell{
	symbol:    ' ',
	combining: nil,
	style:     tcell.StyleDefault,
}
