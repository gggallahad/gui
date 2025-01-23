package gui

type (
	Cell struct {
		Symbol     rune
		Foreground Color
		Background Color
	}
)

var (
	DefaultSymbol rune = ' '

	DefaultCell Cell = Cell{
		Symbol:     DefaultSymbol,
		Foreground: DefaultColor,
		Background: DefaultColor,
	}
)
