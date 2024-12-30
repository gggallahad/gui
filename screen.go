package gui

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

type (
	Screen struct {
		initHandler InitHandler
		handlers    map[State][]Handler

		context *Context

		mutex sync.RWMutex
	}
)

func NewScreen() (*Screen, error) {
	initHandler := emptyInitHandler
	handlers := make(map[State][]Handler)

	context, err := newContext()
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

	eventChannel := make(chan tcell.Event)
	quitChannel := make(chan struct{})

	go s.context.tcellScreen.ChannelEvents(eventChannel, quitChannel)

RunLoop:
	for {
		select {
		case <-s.context.killChannel:
			break RunLoop
		case <-quitChannel:
			break RunLoop
		case event := <-eventChannel:
			go s.handleEvent(event)
		}
	}

	s.context.Cancel()
}

func (s *Screen) handleEvent(event tcell.Event) {
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
	err := s.context.tcellScreen.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Screen) Close() {
	s.context.tcellScreen.Fini()
}

// util

func (s *Screen) getHandlers(state State) []Handler {
	handlers := s.handlers[state]

	return handlers
}
