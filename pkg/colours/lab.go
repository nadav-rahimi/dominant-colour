package colours

import "math"

// LAB Colour
type LAB struct {
	L, A, B float64
}

// Converts a LAB colour to XYZ
func (lab *LAB) XYZ() *XYZ {
	x := Xn * finv(((lab.L+16)/116)+(lab.A/500))
	y := Yn * finv((lab.L+16)/116)
	z := Zn * finv(((lab.L+16)/116)-(lab.B/200))

	return &XYZ{x, y, z}
}

// Converts a LAB colour to RGB
func (lab *LAB) RGB() *RGB {
	return lab.XYZ().RGB()
}

// CIEDE2000 - http://www2.ece.rochester.edu/~gsharma/ciede2000/ciede2000noteCRNA.pdf
func LABDistance(lab1, lab2 *LAB) float64 {
	// Stage 1
	C1 := math.Sqrt(Sqr(lab1.A) + Sqr(lab1.B))
	C2 := math.Sqrt(Sqr(lab2.A) + Sqr(lab2.B))
	Cmean := (C1 + C2) / float64(2)

	G := 0.5 * (1 - math.Sqrt(math.Pow(Cmean, 7)/(math.Pow(Cmean, 7)+math.Pow(25, 7))))

	a1Prime := (1 + G) * lab1.A
	a2Prime := (1 + G) * lab2.A

	C1Prime := math.Sqrt(Sqr(a1Prime) + Sqr(lab1.B))
	C2Prime := math.Sqrt(Sqr(a2Prime) + Sqr(lab2.B))

	var h1Prime, h2Prime float64
	if a1Prime == lab1.B && lab1.B == 0 {
		h1Prime = 0
	} else {
		h1Prime = HueAtan2(lab1.B, a1Prime)
	}
	if a2Prime == lab2.B && lab2.B == 0 {
		h2Prime = 0
	} else {
		h2Prime = HueAtan2(lab2.B, a2Prime)
	}

	// Stage 2
	deltaLPrime := lab2.L - lab1.L
	deltaCPrime := C2Prime - C1Prime

	var deltahPrime float64
	if C1Prime*C2Prime == 0 {
		deltahPrime = 0
	} else if math.Abs(h2Prime-h1Prime) <= 180 {
		deltahPrime = h2Prime - h1Prime
	} else if (h2Prime - h1Prime) > 180 {
		deltahPrime = (h2Prime - h1Prime) - 360
	} else if (h2Prime - h1Prime) < -180 {
		deltahPrime = (h2Prime - h1Prime) + 360
	}

	deltaHPrime := 2 * math.Sqrt(C1Prime*C2Prime) * Sin(deltahPrime/2)

	// Stage 3
	meanLPrime := (lab1.L + lab2.L) / 2
	meanCPrime := (C1Prime + C2Prime) / 2

	var deltahMean float64
	if C1Prime*C2Prime == 0 {
		deltahMean = h1Prime + h2Prime
	} else if math.Abs(h1Prime-h2Prime) <= 180 {
		deltahMean = (h1Prime + h2Prime) / 2
	} else if math.Abs(h1Prime-h2Prime) > 180 && (h1Prime+h2Prime) < 360 {
		deltahMean = (h1Prime + h2Prime + 360) / 2
	} else if math.Abs(h1Prime-h2Prime) > 180 && (h1Prime+h2Prime) >= 360 {
		deltahMean = (h1Prime + h2Prime - 360) / 2
	}

	T := 1 - 0.17*Cos(deltahMean-30) + 0.24*Cos(2*deltahMean) + 0.32*Cos(3*deltahMean+6) - 0.20*Cos(4*deltahMean-63)

	deltaTheta := 30 * math.Exp(-1*Sqr((deltahMean-275)/25))

	Rc := 2 * math.Sqrt(math.Pow(meanCPrime, 7)/(math.Pow(meanCPrime, 7)+math.Pow(25, 7)))

	Sl := 1 + (0.015*Sqr(meanLPrime-50))/(math.Sqrt(20+Sqr(meanLPrime-50)))
	Sc := 1 + 0.045*meanCPrime
	Sh := 1 + 0.015*meanCPrime*T

	Rt := -Sin(2*deltaTheta) * Rc

	// Weighting factors
	var kL, kC, kH float64
	kL, kC, kH = 1, 1, 1

	// Calculating the final difference
	p1 := Sqr(deltaLPrime / (kL * Sl))
	p2 := Sqr(deltaCPrime / (kC * Sc))
	p3 := Sqr(deltaHPrime / (kH * Sh))
	p4 := Rt * (deltaCPrime / (kC * Sc)) * (deltaHPrime / (kH * Sh))

	return math.Sqrt(p1 + p2 + p3 + p4)
}

// Functions used in XYZ --> LAB and LAB --> XYZ conversion
func f(t float64) float64 {
	if t > math.Pow(d, 3) {
		return math.Cbrt(t)
	} else {
		lhs := t / (3 * math.Pow(d, 2))
		rhs := float64(4) / float64(29)
		return lhs + rhs
	}
}

func finv(t float64) float64 {
	if t > d {
		return math.Pow(t, 3)
	} else {
		lhs := (3 * math.Pow(d, 2))
		rhs := t - (float64(4) / float64(29))
		return lhs * rhs
	}
}

// Constants for LAB conversion
const (
	d float64 = 6 / 29
)
