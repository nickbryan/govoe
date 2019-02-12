package noise

import simplex "github.com/ojrac/opensimplex-go"

type OctaveNoise struct {
	baseNoise []simplex.Noise
}

func NewOctave(octaves int) *OctaveNoise {
	o := &OctaveNoise{
		baseNoise: make([]simplex.Noise, octaves),
	}

	for i := 0; i < octaves; i++ {
		o.baseNoise[i] = simplex.New(420)
	}

	return o
}

func (o *OctaveNoise) Compute(x, y float64) float64 {
	amp, freq, sum := 1.0, 1.0, 0.0

	for i := 0; i < len(o.baseNoise); i++ {
		sum += o.baseNoise[i].Eval2(x*freq, y*freq) * amp
		amp *= 2.0
		freq *= 0.5
	}

	return sum
}
