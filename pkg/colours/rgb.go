package colours

// RGB Data
type RGB struct {
	R, G, B float64
}

func NewRGB(r, g, b float64) *RGB {
	return &RGB{r, g, b}
}

func (rgb *RGB) scale() (r, g, b float64) {
	// Scaling down uint8 values to be in the range [0, 1]
	r = rgb.R / 255
	g = rgb.G / 255
	b = rgb.B / 255

	return r, g, b
}

func (rgb *RGB) XYZ() *XYZ {
	r, g, b := rgb.scale()

	x := 0.4124564*r + 0.3575761*g + 0.1804375*b
	y := 0.2126729*r + 0.7151522*g + 0.0721750*b
	z := 0.0193339*r + 0.1191920*g + 0.9503041*b

	return &XYZ{x, y, z}
}

func (rgb *RGB) LAB() *LAB {
	return rgb.XYZ().LAB()
}
