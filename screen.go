package gui

import (
	"github.com/nsf/termbox-go"
)

type (
	Screen struct {
		initHandlers       []InitHandler
		backgroundHandlers []BackgroundHandler
		globalMiddlewares  []Handler
		globalPostwares    []Handler
		handlers           map[State][]Handler

		context *Context
	}
)

func NewScreen(screenConfig ...ScreenConfig) (*Screen, error) {
	var config ScreenConfig
	if len(screenConfig) != 0 {
		config = screenConfig[0]
	} else {
		config.DefaultCell = DefaultCell
	}

	handlers := make(map[State][]Handler)

	context, err := newContext(config.DefaultCell)
	if err != nil {
		return nil, err
	}

	screen := Screen{
		initHandlers:       nil,
		backgroundHandlers: nil,
		globalMiddlewares:  nil,
		globalPostwares:    nil,
		handlers:           handlers,
		context:            context,
	}

	return &screen, nil
}

func (s *Screen) BindInitHandlers(initHandlers ...InitHandler) {
	s.initHandlers = initHandlers
}

func (s *Screen) BindBackgroundHandlers(backgroundHandlers ...BackgroundHandler) {
	s.backgroundHandlers = backgroundHandlers
}

func (s *Screen) BindGlobalMiddlewares(globalMiddlewares ...Handler) {
	s.globalMiddlewares = globalMiddlewares
}

func (s *Screen) BindGlobalPostwares(globalPostwares ...Handler) {
	s.globalPostwares = globalPostwares
}

func (s *Screen) BindHandlers(state State, handlers ...Handler) {
	s.handlers[state] = handlers
}

func (s *Screen) Run() {
	for i := range s.initHandlers {
		s.initHandlers[i](s.context)
	}

	for i := range s.backgroundHandlers {
		go s.backgroundHandlers[i](s.context)
	}

	eventChannel := make(chan Event)

	go s.getEvents(eventChannel)

RunLoop:
	for {
		select {
		case <-s.context.killChannel:
			break RunLoop
		case event := <-eventChannel:
			go s.handleEvent(event)
		}
	}

	s.context.Cancel()
}

func (s *Screen) getEvents(eventChannel chan<- Event) {
	for {
		termboxEvent := termbox.PollEvent()
		event := termboxEventToEvent(termboxEvent)
		eventChannel <- event
	}
}

func (s *Screen) handleEvent(eventType Event) {
	switch event := eventType.(type) {
	case *EventResize:
		s.context.setViewSize(event.X, event.Y)
	}

	currentState := s.context.getCurrentState()

	handlers := s.getHandlers(currentState)
	childContext := s.context.newChildContext()

	childContext.resetData(childContext)

	for i := range s.globalMiddlewares {
		s.globalMiddlewares[i](childContext, eventType)
	}

	for handlerIndex := childContext.getHandlerIndex(); handlerIndex < len(handlers); handlerIndex = childContext.getHandlerIndex() {
		handlers[handlerIndex](childContext, eventType)
		childContext.addHandlerIndex()
	}

	for i := range s.globalPostwares {
		s.globalPostwares[i](childContext, eventType)
	}

	childContext.Cancel()
}

// init

func (s *Screen) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	termbox.SetOutputMode(termbox.OutputRGB)

	viewSizeX, viewSizeY := termbox.Size()
	s.context.setViewSize(viewSizeX, viewSizeY)

	return nil
}

func (s *Screen) Close() {
	termbox.Close()
}

// util

func (s *Screen) getHandlers(state State) []Handler {
	handlers := s.handlers[state]

	return handlers
}
