package integrate

import (
	"math"
	"testing"
)

func Test_romberg(t *testing.T) {
	type args struct {
		xa  float64
		xb  float64
		f   func(float64) float64
		eps float64
	}
	tests := []struct {
		name     string
		args     args
		wantArea float64
	}{
		{
			"Case 1 : f(x) = 1.0 / math.Sqrt(1.0+x*x)",
			args{0.0, 1.0, Fa, 1.0e-6},
			0.8813735883780485,
		},
		{
			"Case 2 : f(x) = x * math.Sin(x)",
			args{-1.0, 1.0, Fb, 1.0e-6},
			0.602337357879467,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotArea := romberg(tt.args.xa, tt.args.xb, tt.args.f, tt.args.eps); gotArea != tt.wantArea {
				t.Errorf("romberg() = %v, want %v", gotArea, tt.wantArea)
			}
		})
	}
}

func Fa(x float64) float64 {
	return 1.0 / math.Sqrt(1.0+x*x)
}
func Fb(x float64) float64 {
	return x * math.Sin(x)
}
