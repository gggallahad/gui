package gui

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

type (
	Screen struct {
		context *Context

		initHandler InitHandler
		handlers    map[State][]Handler

		handlersMutex sync.RWMutex
	}
)

func NewScreen() (*Screen, error) {
	context, err := newContext()
	if err != nil {
		return nil, err
	}

	screen := Screen{
		context:     context,
		initHandler: emptyInitHandler,
		handlers:    make(map[State][]Handler),
	}

	return &screen, nil
}

func (s *Screen) BindInitHandler(handler InitHandler) {
	s.handlersMutex.Lock()
	s.initHandler = handler
	s.handlersMutex.Unlock()
}

func (s *Screen) BindHandlers(state State, handlers ...Handler) {
	s.handlersMutex.Lock()
	s.handlers[state] = handlers
	s.handlersMutex.Unlock()
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
	}

	s.context.Cancel()
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
