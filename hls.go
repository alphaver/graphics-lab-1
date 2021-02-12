package graphicslab1

type HLS struct {
	h, l, s float64
}

func NewHLS(h, l, s float64) *HLS {
	return &HLS{h, l, s}
}

func (hls *HLS) H() float64 {
	return hls.h
}

func (hls *HLS) L() float64 {
	return hls.l
}

func (hls *HLS) S() float64 {
	return hls.s
}

func (hls *HLS) EqualTo(other *HLS) bool {
	return hls.h == other.h && hls.l == other.l && hls.s == other.s
}
