package gui

import (
	"sync"

	"github.com/nsf/termbox-go"
)

type (
	Screen struct {
		initHandler InitHandler
		handlers    map[State][]Handler

		context *Context

		mutex sync.RWMutex
	}
)

func NewScreen(screenConfig ...ScreenConfig) (*Screen, error) {
	var config ScreenConfig
	if len(screenConfig) != 0 {
		config = screenConfig[0]
	} else {
		config.DefaultCell = DefaultCell
	}

	initHandler := emptyInitHandler
	handlers := make(map[State][]Handler)

	context, err := newContext(config.DefaultCell)
	if err != nil {
		return nil, err
	}

	screen := Screen{
		initHandler: initHandler,
		handlers:    handlers,
		context:     context,
	}

	return &screen, nil
}

func (s *Screen) BindInitHandler(handler InitHandler) {
	s.mutex.Lock()
	s.initHandler = handler
	s.mutex.Unlock()
}

func (s *Screen) BindHandlers(state State, handlers ...Handler) {
	s.mutex.Lock()
	s.handlers[state] = handlers
	s.mutex.Unlock()
}

func (s *Screen) Run() {
	s.initHandler(s.context)

	eventChannel := make(chan termbox.Event)

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

func (s *Screen) getEvents(eventChannel chan<- termbox.Event) {
	for {
		event := termbox.PollEvent()
		eventChannel <- event
	}
}

func (s *Screen) handleEvent(event termbox.Event) {
	currentState := s.context.getCurrentState()

	handlers := s.getHandlers(currentState)
	childContext := s.context.newChildContext()

	childContext.resetData(childContext)

	for handlerIndex := childContext.getHandlerIndex(); handlerIndex < len(handlers); handlerIndex = childContext.getHandlerIndex() {
		handlers[handlerIndex](childContext, event)
		childContext.addHandlerIndex()
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
