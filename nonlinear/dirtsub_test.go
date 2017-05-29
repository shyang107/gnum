package nonlinear

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

func Test_dirtsub(t *testing.T) {
	type args struct {
		x0     float64
		eps    float64
		nloops int
		g      func(float64) float64
	}
	tests := []struct {
		name     string
		args     args
		wantXres float64
		wantEpsf float64
		wantErr  bool
	}{
		{
			"case 1: f(x)= x^2-3 = 0",
			args{2.0, 1e-5, -1, g1},
			1.73205,
			1e-5,
			false,
		},
		{
			"case 2: Solve f(x)= tan x - x = 0 using Newton Raphson method",
			args{0.43, 1e-5, 10, g2},
			4.493471e-1,
			1e-5,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotXres, gotEpsf, err := dirtsub(tt.args.x0, tt.args.eps, tt.args.nloops, tt.args.g)
			fmt.Printf("gotXres = %.6e, gotEpsf = %.6e, err = %v\n", gotXres, gotEpsf, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("dirtsub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if math.Abs(gotXres-tt.wantXres) > tt.args.eps {
				t.Errorf("dirtsub() gotXres = %v, want %v", gotXres, tt.wantXres)
			}
			if gotEpsf > tt.wantEpsf {
				t.Errorf("dirtsub() gotEpsf = %v, want %v", gotEpsf, tt.wantEpsf)
			}
		})
	}
}

func g1(x float64) float64 {
	return 0.5 * (x + 3./x)
}
func g2(x float64) float64 {
	// g(x) = x - f(x)/f'(x)
	L := 1.0
	xL := x * L
	cs := math.Cos(xL)
	ss := math.Sin(xL)
	f := x/(ss*ss) - cs/(L*ss)
	return f
}
func g3(x float64) float64 {
	L := 1.0
	xL := x * L
	f := math.Tan(xL) - xL
	if math.Abs(f) <= 1e-6 {
		return 0.0
	}
	return f
}

func Test_searchHI(t *testing.T) {
	type args struct {
		f    func(float64) float64
		xmin float64
		xmax float64
		dx   float64
		icut int
		flmt float64
	}
	tests := []struct {
		name    string
		args    args
		wantX   float64
		wantFx  float64
		wantErr bool
	}{
		{
			"Case 1 : Solve f(x)= tan xL - xL = 0",
			args{g3, 1., 100., 0.1, 20, 1000.},
			1.5703125000000004,
			2065.2848772485004,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotFx, err := searchHI(tt.args.f, tt.args.xmin, tt.args.xmax, tt.args.dx, tt.args.icut, tt.args.flmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotX != tt.wantX {
				t.Errorf("search() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotFx != tt.wantFx {
				t.Errorf("search() gotFx = %v, want %v", gotFx, tt.wantFx)
			}
		})
	}
	fmt.Println(strings.Repeat("=", 60))
	xmax, xmin, dx, icut, flmt := 100., 1., 0.1, 20, 1000.
	for i := 1; i < 11; i++ {
		x, fx, err := searchHI(g3, xmin, xmax, dx, icut, flmt)
		if err != nil {
			// fmt.Println(err)
		}
		if math.Abs(fx) >= flmt {
			fmt.Println("# Following may be a discontinuous point")
		}
		fmt.Printf("x[%2d] = %13.6e, f(x) = %13.6e\n", i, x, fx)
		xmin = x + dx/1
		if x > xmax {
			break
		}
	}
	fmt.Println(strings.Repeat("=", 60))
}
func g4(x float64) float64 {
	f := x*x*x - 2.*x - 5.
	if math.Abs(f) <= 1.0e-6 {
		return 0.0
	}
	return f
}
func Test_searchFS(t *testing.T) {
	type args struct {
		f    func(float64) float64
		xmin float64
		xmax float64
		dx   float64
		icut int
		flmt float64
	}
	tests := []struct {
		name    string
		args    args
		wantX   float64
		wantFx  float64
		wantErr bool
	}{
		{
			"Case 1 : Solve f(x)= tan xL - xL = 0",
			args{g3, 1., 100., 0.1, 20, 1000.},
			1.5706533187057712,
			6991.040738794766,
			true,
		},
		{
			"Case 2 : Solve f(x)= x^3 - 2x -5 = 0",
			args{g4, 2., 2.5, 0.1, 20, 100000.},
			2.0945514816982445,
			0.,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotX, gotFx, err := searchFS(tt.args.f, tt.args.xmin, tt.args.xmax, tt.args.dx, tt.args.icut, tt.args.flmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchFS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotX != tt.wantX {
				t.Errorf("searchFS() gotX = %v, want %v", gotX, tt.wantX)
			}
			if gotFx != tt.wantFx {
				t.Errorf("searchFS() gotFx = %v, want %v", gotFx, tt.wantFx)
			}
		})
	}
	fmt.Println(strings.Repeat("=", 60))
	xmax, xmin, dx, icut, flmt := 100., 1., 0.1, 20, 1000.
	for i := 1; i < 11; i++ {
		x, fx, err := searchFS(g3, xmin, xmax, dx, icut, flmt)
		if err != nil {
			// fmt.Println(err)
		}
		if math.Abs(fx) >= flmt {
			fmt.Println("# Following may be a discontinuous point")
		}
		fmt.Printf("x[%2d] = %13.6e, f(x) = %13.6e\n", i, x, fx)
		xmin = x + dx/1
		if x > xmax {
			break
		}
	}
	fmt.Println(strings.Repeat("=", 60))
}

func TestFroot(t *testing.T) {
	type args struct {
		f     func(float64) float64
		xini  float64
		dx    float64
		eps   float64
		itmax int
		flmt  float64
	}
	tests := []struct {
		name    string
		args    args
		wantXx  float64
		wantFx  float64
		wantErr bool
	}{
		{
			"Case 1 : Solve f(x)= tan xL - xL = 0",
			args{g3, 0.1, 0.1, 1e-5, 20, 100.},
			1.5706533187057712,
			6991.040738794766,
			false,
		},
		// {
		// 	"Case 2 : Solve f(x)= x^3 - 2x -5 = 0",
		// 	args{g4, 0.0, 0.1, 1e-5, 20, 0.},
		// 	2.0945514816982445,
		// 	0,
		// 	false,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotXx, gotFx, err := Froot(tt.args.f, tt.args.xini, tt.args.dx, tt.args.eps, tt.args.itmax, tt.args.flmt)
			if (err != nil) != tt.wantErr {
				t.Errorf("Froot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotXx != tt.wantXx {
				t.Errorf("Froot() gotXx = %v, want %v", gotXx, tt.wantXx)
			}
			if gotFx != tt.wantFx {
				t.Errorf("Froot() gotFx = %v, want %v", gotFx, tt.wantFx)
			}
		})
	}
	fmt.Println(strings.Repeat("=", 60))
	xini, dx, eps, itmax, flmt := 1.5, 0.1, 1.0e-5, 20, 1000.
	for i := 1; i < 11; i++ {
		x, fx, err := Froot(g3, xini, dx, eps, itmax, flmt)
		if err != nil {
			fmt.Println(err)
		}
		if math.Abs(fx) >= flmt {
			fmt.Println("# Following may be a discontinuous point")
		}
		fmt.Printf("x%02d = %13.6e, f(x) = %13.6e\n", i, x, fx)
		xini = x + dx*2.0
		// if x > xmax {
		// 	break
		// }
	}
	fmt.Println(strings.Repeat("=", 60))
}
