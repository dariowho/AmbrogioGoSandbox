package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("No negative numbers bitch! (%f received)",
		e)
}

func Sqrt(x float64) (float64, error) {
	if x < 0 {
		var e ErrNegativeSqrt
		e = ErrNegativeSqrt(x)
		return 0, e
	}

	z      := float64(1)
	z_prev := float64(0)

	i := 1
	for math.Abs(z-z_prev)>= 0.00000000000002 {
		z_prev = z
		z = z - ((z*z - x)/(2*z))
		i++
	}
	
	return z,nil
}

func main() {
	fmt.Println(Sqrt(2))
	fmt.Println(Sqrt(-2))
}
