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
		R: -1,
		G: -1,
		B: -1,
	},
	Background: Color{
		R: -1,
		G: -1,
		B: -1,
	},
}
