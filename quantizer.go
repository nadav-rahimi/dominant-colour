package dominant_colour

import "image/color"

var OtsuQuantizer Quantizer = OtsuStruct{}

type Quantizer interface {
	DominantGreyscaleValue(input string) (color.Gray, error)
	DominantGreyscaleValues(input string, m int) ([]color.Gray, error)
	DominantColourValue(input string) (color.RGBA, error)
	DominantColourValues(input string, m int) ([]color.RGBA, error)
}
