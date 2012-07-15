package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) (float64, int) {
	z      := float64(1)
	z_prev := float64(0)

	i := 1
	for math.Abs(z-z_prev)>= 0.00000000000002 {
		z_prev = z
		z = z - ((z*z - x)/(2*z))
		i++
	}
	
	return z,i
}

func main() {
	v, n := Sqrt(2)
	fmt.Println(v)
	fmt.Println(n,"iterations.")
}
