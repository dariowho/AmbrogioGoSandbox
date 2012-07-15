package main

import (
	"image"
	"image/color"
	"code.google.com/p/go-tour/pic"
)

type MyImage struct{
	width int
	height int
}

func (self MyImage) ColorModel() color.Model {
	return color.RGBAModel
}

func (self MyImage) Bounds() image.Rectangle {
	return image.Rect(0,0,self.width,self.height)
}

func(self MyImage) At(x,y int) color.Color {
	v := uint8(y^x)
	return color.RGBA{v,v,128,255}
}

func main() {
	m := MyImage{256,256}
	pic.ShowImage(m)
}
