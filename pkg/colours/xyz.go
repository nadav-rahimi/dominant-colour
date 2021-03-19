package colours

// XYZ
type XYZ struct {
	X, Y, Z float64
}

// Create a new XYZ colour
func NewXYZ(x, y, z float64) *XYZ {
	return &XYZ{x, y, z}
}

// Converts an XYZ colour to the RGB colour space
func (xyz *XYZ) RGB() *RGB {
	r := Clamp((3.2404542*xyz.X+-1.5371385*xyz.Y+-0.4985314*xyz.Z)*255, 0, 255)
	g := Clamp((-0.9692660*xyz.X+1.8760108*xyz.Y+0.0415560*xyz.Z)*255, 0, 255)
	b := Clamp((0.0556434*xyz.X+-0.2040259*xyz.Y+1.0572252*xyz.Z)*255, 0, 255)

	return &RGB{r, g, b}
}

// Converts an XYZ colour to the LAB colour space
func (xyz *XYZ) LAB() *LAB {
	L := 116*f(xyz.Y/Yn) - 16
	a := 500 * (f(xyz.X/Xn) - f(xyz.Y/Yn))
	b := 200 * (f(xyz.Y/Yn) - f(xyz.Z/Zn))

	return &LAB{L, a, b}
}

// Standard Illuminati D65 used for XYZ conversion
const (
	Xn float64 = 95.0489
	Yn float64 = 100
	Zn float64 = 109.8840
)
