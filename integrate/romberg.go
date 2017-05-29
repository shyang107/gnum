package integrate

import (
	"fmt"
	"math"
)

// romberg is integrate a function using Romberg method
func romberg(xa, xb float64, f func(float64) float64, eps float64) (area float64) {
	//-----------------------------------------------------
	// area = Int_xa^xb f(x) dx
	//-----------------------------------------------------
	var A [21]float64
	var T [21][21]float64
	//-----------------------------------------------------
	h := xb - xa
	A[1] = 0.5 * (f(xa) + f(xb)) * h
	area = A[1]
	T[1][1] = A[1]
	fmt.Printf("%13.6e\n", T[1][1])
	//-----------------------------------------------------
	// compute T^{(1)}_N
	//-----------------------------------------------------
	jj := 1
	var aold float64
	for n := 2; n < 21; n++ {
		aold = area
		an := 0.0
		x := xa + h*0.5
		for j := 1; j < jj+1; j++ {
			an += f(x)
			x += h
		}
		A[n] = 0.5 * (A[n-1] + h*an)
		T[n][1] = A[n]
		//-----------------------------------------------------
		// compute T^{(m)}_N
		//-----------------------------------------------------
		r := 1.0
		for i := n - 1; i > 0; i-- {
			r *= 4.
			A[i] = A[i+1] + (A[i+1]-A[i])/(r-1.0)
			T[n][n+1-i] = A[i]
		}
		//-----------------------------------------------------
		for m := 1; m < n+1; m++ {
			fmt.Printf("%13.6e ", T[n][m])
		}
		fmt.Println()
		//-----------------------------------------------------
		area = A[1]
		if math.Abs(aold-area) < eps*math.Abs(area) {
			return area
		}
		h *= 0.5
		jj += jj
	}
	return area
}
