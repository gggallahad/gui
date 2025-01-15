package gui

import "github.com/nsf/termbox-go"

type (
	Event interface {
		IsEvent()
	}

	EventKey struct {
		Symbol   rune
		Key      termbox.Key
		Modifier termbox.Modifier
	}

	EventMouse struct {
		X   int
		Y   int
		Key termbox.Key
	}
)

func (e EventKey) IsEvent() {
}

func (e EventMouse) IsEvent() {
}

func termboxEventToEvent(termboxEvent termbox.Event) Event {
	var event Event

	switch termboxEvent.Type {
	case termbox.EventKey:
		eventKey := EventKey{
			Symbol:   termboxEvent.Ch,
			Key:      termboxEvent.Key,
			Modifier: termboxEvent.Mod,
		}
		event = eventKey
	case termbox.EventMouse:
		eventMouse := EventMouse{
			X:   termboxEvent.MouseX,
			Y:   termboxEvent.MouseY,
			Key: termboxEvent.Key,
		}
		event = eventMouse
	}

	return event

}
