package integrate

import "math"

// Gauss integrate a function using Gauss integration method
var R, S, R1, S1, RS, SR float64
func F0(x float64) (f0 float64) {
	if R != 0.0 {
		z := R1 * math.Pow(x, SR)
		if S == 1.0 {
			f0 = G0(z) * math.Pow(x, 1-R)
		} else {
			f0 = G0(z) * R1 * SR * math.Pow(x, SR-1)
		}
	} else {
		z := math.Exp(x)
		f0 = G0(z) * z
	}
	return f0
}

func F1(x float64) (f1 float64) {
	if R != 0.0 {
		z := R1 * math.Pow(x, SR)
		if S == 1.0 {
			f1 = G1(z) * math.Pow(x, 1-R)
		} else {
			f1 = G1(z) * R1 * SR * math.Pow(x, SR-1)
		}
	} else {
		z := math.Exp(x)
		f1 = G1(z) * z
	}
	return f1
}
func G0(z float64) (g float64) {
	if z <= 0.0 {
		g = 0.0
	} else if z >= 1.0 {
		g = 1.0e7
	} else {
		g = 0.5 / math.Sqrt(-math.Log(z))
	}
	return g
}
func G1(z float64) (g float64) {
	// g = G0(1.0e-z)
	if z <= 0.0 {
		g = 1.0e30
	} else if z >= 1.0e0 {
		g = 0.0e0
	} else if z <= 1.0e-3 {
		g = 0.5 / math.Sqrt(z+z*z/2.0+math.Pow(z, 3.0)/3.0+math.Pow(z, 4.0)/4.0+math.Pow(z, 5.0)/5.0)
	} else {
		g = 0.5 / math.Sqrt(-math.Log(math.Pow(10.0, -z)))
	}
	return g
}
