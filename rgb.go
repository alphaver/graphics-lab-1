package graphicslab1

import "math"

type RGB uint32

func NewRGB(r, g, b byte) RGB {
	return RGB((int32(b) << 16) | (int32(g) << 8) | (int32(r)))
}

func (rgb RGB) B() byte {
	return byte((rgb & 0xff0000) >> 16)
}

func (rgb RGB) G() byte {
	return byte((rgb & 0xff00) >> 8)
}

func (rgb RGB) R() byte {
	return byte(rgb & 0xff)
}

func (rgb RGB) Normalize() (r, g, b float64) {
	return float64(rgb.R()) / 255.0, float64(rgb.G()) / 255.0, float64(rgb.B()) / 255.0
}

func (rgb RGB) EqualTo(other RGB) bool {
	return rgb == other
}

func DenormalizeRGB(r, g, b float64) RGB {
	return NewRGB(byte(math.Round(r*255)), byte(math.Round(g*255)), byte(math.Round(b*255)))
}
