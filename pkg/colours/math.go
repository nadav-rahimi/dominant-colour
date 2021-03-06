package colours

import "math"

// General functions
func Clamp(n, a, b float64) float64 {
	if n < a {
		return a
	}
	if n > b {
		return b
	}
	return n
}

func ClampFloatToUint8(n float64) uint8 {
	if n < 0 {
		return 0
	}
	if n > 255 {
		return 255
	}
	return uint8(n)
}

func Sqr(a float64) float64 {
	return a * a
}

// Returns value is in degrees
func HueAtan2(x, y float64) float64 {
	return (math.Atan2(x, y) + 2*math.Pi) * (180 / math.Pi)
}

// Takes degrees as input
func Sin(x float64) float64 {
	return math.Sin(x * (math.Pi / 180))
}

// Takes degrees as input
func Cos(x float64) float64 {
	return math.Cos(x * (math.Pi / 180))
}
