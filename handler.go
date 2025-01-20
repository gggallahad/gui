package gui

type (
	InitHandler       func(*Context)
	BackgroundHandler func(*Context)
	Handler           func(*Context, Event)
)

var (
	emptyInitHandler       InitHandler       = func(*Context) {}
	emptyBackgroundHandler BackgroundHandler = func(*Context) {}
	emptyHandler           Handler           = func(*Context, Event) {}
)
