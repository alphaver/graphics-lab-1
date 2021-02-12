package graphicslab1

type CMYK struct {
	c, m, y, k float64
}

func NewCMYK(c, m, y, k float64) *CMYK {
	return &CMYK{c, m, y, k}
}

func (cmyk *CMYK) C() float64 {
	return cmyk.c
}

func (cmyk *CMYK) M() float64 {
	return cmyk.m
}

func (cmyk *CMYK) Y() float64 {
	return cmyk.y
}

func (cmyk *CMYK) K() float64 {
	return cmyk.k
}

func (cmyk *CMYK) EqualTo(other *CMYK) bool {
	return cmyk.c == other.c && cmyk.m == other.m && cmyk.y == other.y && cmyk.k == other.k
}
