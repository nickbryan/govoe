package noise

type CombinedNoise struct {
	n1, n2 *OctaveNoise
}

func NewCombined(n1, n2 *OctaveNoise) *CombinedNoise {
	return &CombinedNoise{n1, n2}
}

func (c *CombinedNoise) Compute(x, y float64) float64 {
	offset := c.n2.Compute(x, y)
	return c.n1.Compute(x+offset, y)
}
