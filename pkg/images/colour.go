package images

type colour struct {
	R, G, B uint8
}

func NewColour(r, g, b uint8) *colour {
	return &colour{r, g, b}
}

func (c *colour) Add(c1, c2 colour) {
	c.R = c1.R + c2.R
	c.G = c1.G + c2.G
	c.B = c1.B + c2.B
}

func (c *colour) Sub(c1, c2 colour) {
	c.R = c1.R - c2.R
	c.G = c1.G - c2.G
	c.B = c1.B - c2.B
}

func (c *colour) DivScalar(c1 colour, s uint8) {
	c.R = c1.R / s
	c.G = c1.G / s
	c.B = c1.B / s
}
