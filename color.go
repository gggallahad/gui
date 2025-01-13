package gui

import "github.com/nsf/termbox-go"

type (
	Color struct {
		R int
		G int
		B int
	}
)

func (c *Color) toAttribute() termbox.Attribute {
	if c.R == -1 && c.G == -1 && c.B == -1 {
		return termbox.ColorDefault
	}

	attribute := termbox.RGBToAttribute(uint8(c.R), uint8(c.G), uint8(c.B))

	return attribute
}

func (c *Color) fromAttribute(attribute termbox.Attribute) {
	r, g, b := termbox.AttributeToRGB(attribute)
	c.R = int(r)
	c.G = int(g)
	c.B = int(b)
}
