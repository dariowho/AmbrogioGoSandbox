package main

import "code.google.com/p/go-tour/pic"

func Pic(dx, dy int) [][]uint8 {
	var r = make([][]uint8,dy)
	
	var y,x int
	for y=0; y<dy; y++ {
		r[y] = make([]uint8,dx)
		for x=0; x<dx; x++ {
			r[y][x] = uint8((y*x))
		}
	}
	
	return r

}

func main() {
	pic.Show(Pic)
}
