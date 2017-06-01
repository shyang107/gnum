package integrate

import (
	"fmt"
	"math"
)

// Romberg is integrate a function using Romberg method
func Romberg(xa, xb float64, f func(float64) float64, eps float64) (area float64) {
	//-----------------------------------------------------
	// area = Int_xa^xb f(x) dx
	//-----------------------------------------------------
	var A [20]float64
	var T [20][20]float64
	//-----------------------------------------------------
	h := xb - xa
	A[0] = 0.5 * (f(xa) + f(xb)) * h
	area = A[0]
	T[0][0] = A[0]
	fmt.Printf("%13.6e\n", T[0][0])
	//-----------------------------------------------------
	// compute T^{(1)}_N
	//-----------------------------------------------------
	jj := 1
	var aold float64
	for n := 1; n < 20; n++ {
		// fmt.Println("n = ", n)
		aold = area
		an := 0.0
		x := xa + h*0.5
		for j := 1; j < jj+1; j++ {
			// fmt.Println("\tj = ", j)
			an += f(x)
			x += h
		}
		A[n] = 0.5 * (A[n-1] + h*an)
		T[n][0] = A[n]
		//-----------------------------------------------------
		// compute T^{(m)}_N
		//-----------------------------------------------------
		r := 1.0
		for i := n - 1; i >= 0; i-- {
			r *= 4.
			A[i] = A[i+1] + (A[i+1]-A[i])/(r-1.0)
			T[n][n-i] = A[i]
			// fmt.Println("\t\ti = ", i, "n-i = ", n-i)
		}
		//-----------------------------------------------------
		for m := 0; m <= n; m++ {
			// fmt.Printf("T[%1d][%1d] = %13.6e ", n, m, T[n][m])
			fmt.Printf("%13.6e ", T[n][m])
		}
		fmt.Println()
		//-----------------------------------------------------
		area = A[0]
		if math.Abs(aold-area) < eps*math.Abs(area) {
			return area
		}
		h *= 0.5
		jj += jj
	}
	return area
}
