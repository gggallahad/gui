package gui

import "github.com/gdamore/tcell/v2"

type (
	Screen struct {
		tcellScreen tcell.Screen

		context *Context
	}
)

func NewScreen() (*Screen, error) {
	tcellScreen, err := tcell.NewScreen()
	if err != nil {
		return nil, err
	}

	screen := Screen{
		tcellScreen: tcellScreen,
		context:     newContext(),
	}

	return &screen, nil
}

// init

func (s *Screen) Init() error {
	err := s.tcellScreen.Init()
	if err != nil {
		return err
	}

	return nil
}

func (s *Screen) Close() {
	s.tcellScreen.Fini()
}

func (s *Screen) PollEvent() tcell.Event {
	event := s.tcellScreen.PollEvent()

	return event
}

// draw

func (s *Screen) SetContent(x, y int, symbol rune, combining []rune, style tcell.Style) {
	s.tcellScreen.SetContent(x, y, symbol, combining, style)
}

func (s *Screen) GetContent(x, y int) (rune, []rune, tcell.Style, int) {
	symbol, combining, style, width := s.tcellScreen.GetContent(x, y)

	return symbol, combining, style, width
}

func (s *Screen) Flush() {
	s.tcellScreen.Show()
}

func (s *Screen) Fill(symbol rune, style tcell.Style) {
	s.tcellScreen.Fill(symbol, style)
}

func (s *Screen) Clear() {
	s.tcellScreen.Clear()
}

// util

func (s *Screen) HideCursor() {
	s.tcellScreen.HideCursor()
}

func (s *Screen) ShowCursor(x, y int) {
	s.tcellScreen.ShowCursor(x, y)
}

func (s *Screen) Size() (int, int) {
	screenX, screenY := s.tcellScreen.Size()

	return screenX, screenY
}
