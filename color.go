package gui

import "github.com/nsf/termbox-go"

type (
	Color struct {
		R int
		G int
		B int
	}
)

var (
	DefaultColor Color = Color{
		R: -1,
		G: -1,
		B: -1,
	}
)

func (c *Color) toAttribute() termbox.Attribute {
	if *c == DefaultColor {
		return termbox.ColorDefault
	}

	attribute := termbox.RGBToAttribute(uint8(c.R), uint8(c.G), uint8(c.B))

	return attribute
}

// func (c *Color) fromAttribute(attribute termbox.Attribute) Color {
// 	if attribute == termbox.ColorDefault {
// 		color := DefaultColor

// 		return color
// 	}

// 	r, g, b := termbox.AttributeToRGB(attribute)
// 	color := Color{
// 		R: int(r),
// 		G: int(g),
// 		B: int(b),
// 	}

// 	return color
// }
