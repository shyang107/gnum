package gnum

import (
	"math"
	"strings"
)

const (
	Space = ' '
)

type Alignment byte

const (
	AlignLeft   Alignment = itoa // 0
	AlignCenter                  // 1
	AlignRight                   //2
)

// Spaces returns the number of spaces
func Spaces(num int) string {
	return strings.Repeat(string(Space), num)
}

// var align = map[string]int{
// 	"left":   1,
// 	"center": 0,
// 	"right":  -1,
// }

// Hstring return the formatted string what you wants.
// parameters
// 	str 		: the original string
//	explen 		: the length of 'rstr' string
//	fillchar 	: fill 'fillchar' in the spaces of 'rstr' string
//	align		= 1 	: align at the left (default)
//				= 0		: align at the center
//				= -1	: align at the right
// outputs
//	rstr		: the return string
//	err			: the error if therer are errors
func Hstring(str string, explen int, fillchar byte, alg Alignment) (rstr string) {
	lenstr := len(str)
	if lenstr > explen {
		return str[:explen]
	}
	// if fillchar == '' {
	// 	fillchar = space
	// }
	switch alg {
	case AlignCenter: // align at the center
		nleft := (explen - lenstr) / 2
		rstr = strings.Repeat(string(fillchar), nleft) + str + strings.Repeat(string(fillchar), explen-lenstr-nleft)
	case AlignRight: // align at the right
		rstr = strings.Repeat(string(fillchar), explen-lenstr) + str
	default: // align at the left (default)
		rstr = str + strings.Repeat(string(fillchar), explen-lenstr)
	}
	return rstr
}

// sign : sign(A,B) returns the value of A with the sign of B.
func Sign(a, b float64) float64 {
	return math.Copysign(a, b)
}
