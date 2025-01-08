package gui

import "github.com/nsf/termbox-go"

type (
	InitHandler func(*Context)
	Handler     func(*Context, termbox.Event)
)

var (
	emptyInitHandler InitHandler = func(*Context) {}
	emptyHandler     Handler     = func(*Context, termbox.Event) {}
)
