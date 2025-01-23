package gui

import "github.com/nsf/termbox-go"

type (
	Event interface {
		IsEvent()
	}

	EventKey struct {
		Symbol   rune
		Key      KeyboardKey
		Modifier Modifier
	}

	EventMouse struct {
		X   int
		Y   int
		Key MouseKey
	}

	EventResize struct {
		X int
		Y int
	}
)

func (e *EventKey) IsEvent() {
}

func (e *EventMouse) IsEvent() {
}

func (e *EventResize) IsEvent() {
}

func termboxEventToEvent(termboxEvent termbox.Event) Event {
	var event Event

	switch termboxEvent.Type {
	case termbox.EventKey:
		eventKey := &EventKey{
			Symbol:   termboxEvent.Ch,
			Key:      KeyboardKey(termboxEvent.Key),
			Modifier: Modifier(termboxEvent.Mod),
		}
		event = eventKey
	case termbox.EventMouse:
		eventMouse := &EventMouse{
			X:   termboxEvent.MouseX,
			Y:   termboxEvent.MouseY,
			Key: MouseKey(termboxEvent.Key),
		}
		event = eventMouse
	case termbox.EventResize:
		eventResize := &EventResize{
			X: termboxEvent.Width,
			Y: termboxEvent.Height,
		}
		event = eventResize
	}

	return event
}
