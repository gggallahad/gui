package gui

type (
	InitHandler func(*Context)
	Handler     func(*Context, Event)
)

var (
	emptyInitHandler InitHandler = func(*Context) {}
	emptyHandler     Handler     = func(*Context, Event) {}
)
