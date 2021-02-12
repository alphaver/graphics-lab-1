package graphicslab1

import "math"

func normRGBToHLS(r, g, b float64) *HLS {
	maxC := math.Max(math.Max(r, g), b)
	minC := math.Min(math.Min(r, g), b)
	delta := maxC - minC

	var h, l, s float64

	if delta == 0 {
		h = 0
	} else if maxC == r {
		h = (g - b) / delta
	} else if maxC == g {
		h = (b-r)/delta + 2
	} else if maxC == b {
		h = (r-g)/delta + 4
	}
	h *= 60
	for h < 360 {
		h += 360
	}
	for h >= 360 {
		h -= 360
	}

	l = (maxC + minC) / 2
	if l == 0 || l == 1 {
		s = 0
	} else {
		s = (maxC - l) / math.Min(l, 1-l)
	}

	return NewHLS(h, l, s)
}

func hueToRGB(p, q, t float64) float64 {
	if t < 0 {
		t++
	}
	if t > 1 {
		t--
	}
	switch {
	case t < 1.0/6:
		return p + (q-p)*6*t
	case t >= 1.0/6 && t < 1.0/2:
		return q
	case t >= 1.0/2 && t < 2.0/3:
		return p + (q-p)*(2.0/3-t)*6
	default:
		return p
	}
}

func hlsToNormRGB(hls *HLS) (r, g, b float64) {
	if hls.S() == 0 {
		return hls.L(), hls.L(), hls.L()
	}
	hN := hls.H() / 360
	var p, q float64
	if hls.L() < 0.5 {
		q = hls.L() * (1 + hls.S())
	} else {
		q = hls.L() + hls.S() - hls.L()*hls.S()
	}
	p = 2*hls.L() - q
	return hueToRGB(p, q, hN+1.0/3), hueToRGB(p, q, hN), hueToRGB(p, q, hN-1.0/3)
}

func cmykToNormRGB(cmyk *CMYK) (r, g, b float64) {
	return (1 - cmyk.C()) * (1 - cmyk.K()),
		(1 - cmyk.M()) * (1 - cmyk.K()),
		(1 - cmyk.Y()) * (1 - cmyk.K())
}

func normRGBToCMYK(r, g, b float64) *CMYK {
	key := 1 - math.Max(math.Max(r, g), b)
	if key == 1.0 {
		return NewCMYK(0, 0, 0, 1)
	} else {
		return NewCMYK((1-r-key)/(1-key),
			(1-g-key)/(1-key),
			(1-b-key)/(1-key),
			key)
	}
}

func ToRGB(color interface{}) RGB {
	switch color.(type) {
	case RGB:
		return color.(RGB)

	case *CMYK:
		return DenormalizeRGB(cmykToNormRGB(color.(*CMYK)))

	case *HLS:
		return DenormalizeRGB(hlsToNormRGB(color.(*HLS)))

	default:
		panic("unknown color type")
	}
}

func ToCMYK(color interface{}) *CMYK {
	switch color.(type) {
	case RGB:
		return normRGBToCMYK(color.(RGB).Normalize())

	case *CMYK:
		return color.(*CMYK)

	case *HLS:
		return normRGBToCMYK(hlsToNormRGB(color.(*HLS)))

	default:
		panic("unknown color type")
	}
}

func ToHLS(color interface{}) *HLS {
	switch color.(type) {
	case RGB:
		return normRGBToHLS(color.(RGB).Normalize())

	case *CMYK:
		return normRGBToHLS(cmykToNormRGB(color.(*CMYK)))

	case *HLS:
		return color.(*HLS)

	default:
		panic("unknown color type")
	}
}
