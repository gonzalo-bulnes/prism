package ciexyz

import (
	"github.com/mandykoh/prism/cielab"
	"math"
)

const constantE = 216.0 / 24389.0
const constantK = 24389.0 / 27.0

// Color represents a linear normalised colour in CIE XYZ space.
type Color struct {
	X float32
	Y float32
	Z float32
	A float32
}

// ToLAB converts this colour to a CIE Lab colour given a reference white point.
func (c Color) ToLAB(whitePoint Color) cielab.Color {
	fx := componentToLAB(c.X, whitePoint.X)
	fy := componentToLAB(c.Y, whitePoint.Y)
	fz := componentToLAB(c.Z, whitePoint.Z)

	return cielab.Color{
		L:     float32(116*fy - 16),
		A:     float32(500 * (fx - fy)),
		B:     float32(200 * (fy - fz)),
		Alpha: c.A,
	}
}

// ColorFromLAB creates a CIE XYZ Color instance from a CIE LAB representation
// given a reference white point.
func ColorFromLAB(lab cielab.Color, whitePoint Color) Color {
	fy := (float64(lab.L) + 16) / 116
	fx := float64(lab.A)/500 + fy
	fz := fy - float64(lab.B)/200

	xr := componentFromLAB(fx)
	zr := componentFromLAB(fz)

	var yr float64
	if lab.L > constantK*constantE {
		yr = math.Pow((float64(lab.L)+16)/116, 3)
	} else {
		yr = float64(lab.L) / constantK
	}

	return Color{
		X: float32(xr * float64(whitePoint.X)),
		Y: float32(yr * float64(whitePoint.Y)),
		Z: float32(zr * float64(whitePoint.Z)),
		A: lab.Alpha,
	}
}

func componentFromLAB(f float64) float64 {
	if f3 := math.Pow(f, 3); f3 > constantE {
		return f3
	}
	return (116*f - 16) / constantK
}

func componentToLAB(v float32, wp float32) float64 {
	r := float64(v) / float64(wp)
	if r > constantE {
		return math.Pow(r, 1.0/3.0)
	}
	return (constantK*r + 16) / 116.0
}