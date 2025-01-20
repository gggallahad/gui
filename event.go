package gui

import "github.com/nsf/termbox-go"

type (
	Event interface {
		IsEvent()
	}

	EventKey struct {
		Symbol   rune
		Key      Key
		Modifier Modifier
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
			Key:      Key(termboxEvent.Key),
			Modifier: Modifier(termboxEvent.Mod),
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
