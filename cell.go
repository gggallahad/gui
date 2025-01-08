package gui

type (
	Cell struct {
		Symbol     rune
		Foreground Color
		Background Color
	}
)

var DefaultCell Cell = Cell{
	Symbol: ' ',
	Foreground: Color{
		R: 0,
		G: 0,
		B: 0,
	},
	Background: Color{
		R: -1,
		G: -1,
		B: -1,
	},
}
