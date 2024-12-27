package gui

import "github.com/gdamore/tcell/v2"

type (
	InitHandler func(*Context)
	Handler     func(*Context, tcell.Event)
)

var (
	emptyInitHandler InitHandler = func(ctx *Context) {}
	emptyHandler     Handler     = func(ctx *Context, e tcell.Event) {}
)
