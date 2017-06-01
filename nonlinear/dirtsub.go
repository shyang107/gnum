package nonlinear

import (
	"fmt"
	"math"
)

// var targetfunc func(float64) float64

// dirtsub : direct substitution
// usage: x = dirtsub(x0, eps, nloop, g(x))
//	1. f(x) = 0 改寫成 x = g(x)
//  2.	x0 		: 	起始值
//		eps		:	誤差控制；i.e. abs(x-x0) < eps
// 		N		: 	迭代次數。
//					nloop > 0 	: return x=a and epsf at loop=nloops or abs(x-a) < eps
//					nloop < 0	: return x=a epsf when abs(xi-a) < eps (default) or loop > nloops=20
// 3. 收斂條件： d{g(x)}/dx < 1 at x=a
func dirtsub(x0, eps float64, N int, g func(float64) float64) (x, epsf float64, err error) {
	// fmt.Printf("x0 = %f; eps = %f, N = %d\n", x0, eps, N)
	if N < 0 {
		N = 20
	}
	// fmt.Printf("N = %d\n", N)
	x = x0
	var xold float64
	for i := 0; i < N; i++ {
		xold = x
		x = g(x)
		epsf = math.Abs(x - xold)
		// fmt.Printf(num.Spaces(3)+"loop=%d, xold=%f, x = %f; epsf = %f\n", i, xold, x, epsf)
		if epsf <= eps {
			msg := fmt.Sprintf("The sequence convergence in %4d iterations within %10.3e \nThe last two successive x values are %13.6e and %13.6e", i+1, eps, xold, x)
			fmt.Println(msg)
			return x, epsf, nil
		}
	}
	msg := fmt.Sprintf("Not convergence in %4d iterations within %10.3e \nThe last two successive x values are %13.6e and %13.6e", N, eps, xold, x)
	fmt.Println(msg)
	return x, epsf, fmt.Errorf(msg)
}

// NewtonRaphson method to solve nonlinear equation
// This method is equivalent to select
// g(x) = x - f(x)/f'(x) in direct substitution (dirtsub)
// func NewtonRaphson() float64 {
// }

// searchHI find the roots in [xmin,xmax] using half-interval method
//
// search the root from range [xmin, xmin+dx] to the range
// where the roots is included; then dx <- dx/2 (half-interval), find the
// more precise root, the error < dx * 2**(-m).
// if the root does not exist to [..., xmax], return to the
// calling procedure and push back error message.
//
// f 			: the target function f(x) = ... = 0
// xmin, xmax 	: suppose the root exist in the range [xmin, xmax]
// dx			: the increment of x,
//				  x = xmin, xmin+dx, xmin+2*dx, ... , xmin+i*dx, ... , xmax
// icut 		: The number of interations
// flmt			: the upper limit of f(x)
//
// output
//	xx			: root of function f(x)=0
//  fx			: the corresponding value of f(x) of xx
func searchHI(f func(float64) float64, xmin, xmax, dx float64, icut int, flmt float64) (xx, fx float64, err error) {
	//=====================================================
	// err = nil
	fa, fb, xa, xb := 1.0, -1., 0., 0.
	xx = xmin
	ir, ie := 0, 0
	//-----------------------------------------------------
	// compute function value
	//-----------------------------------------------------
L30:
	fx = f(xx)
	ir++
	if fx == 0. {
		return xx, fx, nil
	}
	// fmt.Printf("ir = %5d, fx = %15.6e, xx = %15.6e\n", ir, fx, xx)
	//-----------------------------------------------------
	// return on discontinuous point
	//-----------------------------------------------------
	if math.Abs(fx) > flmt {
		msg := fmt.Sprintf("f(%15.6e) =  %15.6e > flmt = %15.6e", xx, fx, flmt)
		// fmt.Println(msg)
		return xx, fx, fmt.Errorf(msg)
	}
	//-----------------------------------------------------
	// set limits for next new root
	//-----------------------------------------------------
	if fx < 0. {
		xa, fa = xx, fx
	} else {
		xb, fb = xx, fx
	}
	//-----------------------------------------------------
	// search interval for changing sign
	//-----------------------------------------------------
	if fa*fb > 0. {
		if xx > xmax {
			msg := fmt.Sprintf("No root between %15.6e and %15.6e", xmin, xmax)
			// fmt.Println(msg)
			return xx, fx, fmt.Errorf(msg)
		}
		xx += dx
		ie = ir + icut
		goto L30
	}
	//-----------------------------------------------------
	// get new root by half interval method
	//-----------------------------------------------------
	xx = 0.5 * (xa + xb)
	//-----------------------------------------------------
	// get new root by false position method
	//-----------------------------------------------------
	// L75:
	// 	xx = xa - fa*(xb-xa)/(fb-fa)
	//-----------------------------------------------------
	// checl the iteration numbers
	//-----------------------------------------------------
	if ir <= ie {
		goto L30
	}
	//=====================================================
	msg := fmt.Sprintf("The number of interations exceeds, icut = %4d", icut)
	// fmt.Println(msg)
	return xx, fx, fmt.Errorf(msg)
}

// searchFS find the roots in [xmin,xmax] using the combination of false-position and secant methods
//
// search the root from range [xmin, xmin+dx] to the range
// where the roots is included; then dx <- dx/2 (half-interval), find the
// more precise root, the error < dx * 2**(-m).
// if the root does not exist to [..., xmax], return to the
// calling procedure and push back error message.
//
// f 				: function to calculate, f(x) = ... = 0
// xmin, xmax, dx	: search limits and increment for root
//				  	  x = xmin, xmin+dx, xmin+2*dx, ... , xmin+i*dx, ... , xmax
// icut 		: maximum number of interations
// flmt			: limited function value of f(xx)
//
// output
//	xx			: argument of function
//  fx			: function value of xx, f(xx)
func searchFS(f func(float64) float64, xmin, xmax, dx float64, icut int, flmt float64) (xx, fx float64, err error) {
	//=====================================================
	// err = nil
	var x1, x2, f1, f2, xn, xp float64
	fn, fp := 1.e30, -1.e30
	xx = xmin
	ir, ie := 0, 0
	//-----------------------------------------------------
	// compute function value
	//-----------------------------------------------------
L30:
	fx = f(xx)
	ir++
	if fx == 0. {
		return xx, fx, nil
	}
	// fmt.Printf("ir = %5d, fx = %15.6e, xx = %15.6e\n", ir, fx, xx)
	//-----------------------------------------------------
	// return on discontinuous point
	//-----------------------------------------------------
	if math.Abs(fx) > flmt {
		msg := fmt.Sprintf("f(%15.6e) =  %15.6e > flmt = %15.6e", xx, fx, flmt)
		// fmt.Println(msg)
		return xx, fx, fmt.Errorf(msg)
	}
	//-----------------------------------------------------
	// push down old values, and put the new one on top
	//-----------------------------------------------------
	x1 = x2
	x2 = xx
	f1 = f2
	f2 = fx
	//-----------------------------------------------------
	// set limits for next new root
	//-----------------------------------------------------
	if fx < 0. {
		xn, fn = xx, fx
		if fx < fn {
			ie--
		}
	} else {
		xp, fp = xx, fx
		if fx > fp {
			ie--
		}
	}
	//-----------------------------------------------------
	// search interval for changing sign
	//-----------------------------------------------------
	if fn*fp > 0. {
		if xx > xmax {
			msg := fmt.Sprintf("No root between %15.6e and %15.6e", xmin, xmax)
			// fmt.Println(msg)
			return xx, fx, fmt.Errorf(msg)
		}
		xx += dx
		ie = ir + icut
		goto L30
	}
	//-----------------------------------------------------
	// get new root by secant method
	//-----------------------------------------------------
	if f1 != f2 {
		xx = x1 - f1*(x1-x2)/(f1-f2)
		//-----------------------------------------------------
		//	check if the new root goes outside the limits
		//-----------------------------------------------------
		if (xx-xp)*(xx-xn) > 0 {
			//-----------------------------------------------------
			// true : get new root by false position method
			//-----------------------------------------------------
			xx = xp - fp*(xp-xn)/(fp-fn)
		}
	}
	//-----------------------------------------------------
	// check the iteration numbers
	//-----------------------------------------------------
	if ir <= ie {
		goto L30
	}
	//=====================================================
	msg := fmt.Sprintf("The number of interations exceeds, icut = %4d", icut)
	// fmt.Println(msg)
	return xx, fx, fmt.Errorf(msg)
}

// Xzero find root
//-----------------------------------------------------
// x 		= Argument of function
// f		= function value of x, f(x)
// xzero	= a better guess for f(xzero) = 0
// dx		= in case no better formula to guess xzero
// 			  this function will return:
// 			  either xzero = x + dx, if dx != 0
// 			  or	 xzero = x + f,  if dx == 0
//					 this is simulate xzero = g(x)
// 					 as in the direct substitution method
// 					 if set f(x) = g(x) - x
// istep	= 1 : for the first call of a new function
//			> 1 : otherwise
// for istep = 1 :
// 	xzero = x + dx or xzero = x + f
// for istep > 1 :
// 	xzero by secant method or
// 		  by false-position method (if applicable)
// 		  if xzero by secant mehtod goes outside (xp, xn)
//-----------------------------------------------------
func Xzero(x, f, dx float64, istep int, b *xzeroParameters) (xzero float64) {
	var a xzeroParameters
	//-----------------------------------------------------
	if istep == 1 {
		a.fp, a.fn = -1.0, 1.0
	} else {
		a = *b
	}
	a.st++
	//-----------------------------------------------------
	// save the most newly 3 points
	//-----------------------------------------------------
	a.x0 = a.x1
	a.x1 = a.x2
	a.x2 = x
	a.f0 = a.f1
	a.f2 = a.f1
	a.f2 = f
	//-----------------------------------------------------
	// find new points (xp,xn) that bracket the root
	//-----------------------------------------------------
	if f < 0.0 {
		a.xn, a.fn = x, f
	} else {
		a.xp, a.fp = x, f
	}
	iget, xx := 0, 0.0
	//-----------------------------------------------------
	// for slow convergence due to multiple roots
	// then use secant method on u(x) = f(x) / f'(x) = 0
	// after 13 iterations (13 in statement below, >= 3)
	//-----------------------------------------------------
	if a.st >= 13 {
		if (a.f2-a.f1) != 0.0 && (a.f1-a.f0) != 0.0 {
			u1 := a.f1 * (a.x1 - a.x0) / (a.f1 - a.f0)
			u2 := a.f2 * (a.x2 - a.x1) / (a.f2 - a.f1)
			if (u2 - u1) != 0.0 {
				xx = a.x2 - u2*(a.x2-a.x1)/(u2-u1)
				iget = 4
			}
		}
	}
	//-----------------------------------------------------
	// if we have at least 2 functions and no xx get
	// then use secant method on f(x) = 0
	//-----------------------------------------------------
	if (a.st >= 2) && (iget == 0) {
		if (a.f2 - a.f1) != 0.0 {
			xx = a.x2 - a.f2*(a.x2-a.x1)/(a.f2-a.f1)
			iget = 3
		}
	}
	//-----------------------------------------------------
	// if root is bracketed by (xp,xn)
	// and if no xx get or xx goes outside (xp,xn)
	// then use false-position method by (xp,xn)
	//-----------------------------------------------------
	if (a.fp >= 0.0) && (a.fn < 0.0) {
		if (iget == 0) || ((xx-a.xp)*(xx-a.xn) > 0.0) {
			xx = a.xp - a.fp*(a.xp-a.xn)/(a.fp-a.fn)
			iget = 2
		}
	}
	//-----------------------------------------------------
	// if secant method & false-position method cannot be used:
	// if dx != 0.0	: arbitrary set xx = x + dx
	// else			: set xx = x + f(x)
	// 				  to simulate direct substitution method
	// 				  if set f(x) = g(x) - x in calling program
	//-----------------------------------------------------
	if iget == 0 {
		if dx != 0.0 {
			xx = x + dx
		} else {
			xx = x + f
		}
		iget = 1
	}
	//-----------------------------------------------------
	//  for debug
	//-----------------------------------------------------
	// a.x0 = float64(iget)
	// a.f0 = xx
	//-----------------------------------------------------
	*b = a
	xzero = xx
	// fmt.Printf("Xzero : %+v\n", b)
	return xzero
}

type xzeroParameters struct {
	x0, x1, x2, f0, f1, f2, xp, xn, fp, fn float64
	st                                     int
}

// Froot is used to call Xzero to find roots
func Froot(f func(float64) float64, xini, dx, eps float64, itmax int, flmt float64) (xx, fx float64, err error) {
	//-----------------------------------------------------
	// find root of f(xx) = 0 from the initial guess xini
	//-----------------------------------------------------
	nstp := itmax
	if nstp == 0 {
		nstp = 12
	}
	errx := eps
	if errx == 0.0 {
		errx = 1.0e-6
	}
	if xini != 0.0 {
		errx *= xini
	}
	flmy := flmt
	if flmy == 0.0 {
		flmy = 1.0e30
	}
	//-----------------------------------------------------
	var buf xzeroParameters
	xn := xini
	for istep := 1; istep < nstp+1; istep++ {
		xx, fx = xn, f(xn)
		xn = Xzero(xx, fx, dx, istep, &buf)
		// fmt.Printf("Froot :  %+v\n", buf)
		if math.Abs(xn-xx) <= errx {
			return xx, fx, nil
		}
		if math.Abs(fx) > flmy {
			return xx, fx, fmt.Errorf("f(xx) > limited function value of f(xx)")
		}
	}
	//-----------------------------------------------------
	return xx, fx, nil
}

// Groot is used to call Xzero to find roots
func Groot(g func(float64) float64, xini, eps float64, itmax int) (xn, xx float64, err error) {
	//-----------------------------------------------------
	// find root of xx = g(xx) from the initial guess xini
	//-----------------------------------------------------
	nstp := itmax
	if nstp == 0 {
		nstp = 12
	}
	errx := eps
	if errx == 0.0 {
		errx = 1.0e-6
	}
	if xini != 0.0 {
		errx *= xini
	}
	//-----------------------------------------------------
	var buf xzeroParameters
	xn = xini
	zero := 0.0
	for istep := 1; istep < nstp+1; istep++ {
		xx, xn = xn, g(xx)
		xn = Xzero(xx, xn-xx, zero, istep, &buf)
		fmt.Printf("Dirsub : %v\n", buf)
		if math.Abs(xn-xx) <= errx {
			return xx, xn, nil
		}
	}
	//-----------------------------------------------------
	return xx, xn, nil
}

// Sroot is used to call Xzero to find roots
func Sroot(f func(float64) float64, xini, xfin, dx, eps float64, itmax int, flmt float64) (xx, fx float64, err error) {
	//-----------------------------------------------------
	// find root of f(xx) = 0 from the initial guess xini
	//-----------------------------------------------------
	nstp := itmax
	if nstp == 0 {
		nstp = 12
	}
	errx := eps
	if errx == 0.0 {
		errx = 1.0e-6
	}
	if xini != 0.0 {
		errx *= xini
	}
	flmy := flmt
	if flmy == 0.0 {
		flmy = 1.0e30
	}
	dxx := (xfin - xini) / 16.0
	if dxx != 0.0 {
		dxx = math.Copysign(math.Abs(dx), xfin-xini)
	}
	//-----------------------------------------------------
	var buf xzeroParameters
	xn := xini
	fx = 1.0
	var f0 float64
	for istep := 1; istep < nstp+1; istep++ {
	L10:
		xx, f0, fx = xn, fx, f(xx)
		//-----------------------------------------------------
		//	search interval with root
		//-----------------------------------------------------
		if istep == 2 {
			if (xfin-xx)*(xini-xx) > 0.0 {
				return xx, fx, fmt.Errorf("No root between %15.6e and %15.6e", xini, xfin)
			}
			if fx*f0 > 0.0 {
				xn = Xzero(xx, fx, dxx, 1, &buf)
				// fmt.Printf("Sroot : %v\n", buf)
				goto L10
			}
		}
		//-----------------------------------------------------
		xn = Xzero(xx, fx, dxx, istep, &buf)
		// fmt.Printf("Sroot : %v\n", buf)
		if math.Abs(xn-xx) <= errx {
			return xx, fx, nil
		}
		if math.Abs(fx) > flmy {
			return xx, fx, fmt.Errorf("f(xx) > limited function value of f(xx)")
		}
	}
	//-----------------------------------------------------
	return xx, fx, nil
}

// sign : sign(A,B) returns the value of A with the sign of B.
func sign(a, b float64) float64 {
	if b < 0 {
		return -math.Abs(a)
	}
	return math.Abs(a)
}
