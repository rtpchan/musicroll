package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// create button image for Normal,Hover,Press,Disable
// helper file to create button images
func ButtonImages(w, h int, color [4]color.Color) []*ebiten.Image {
	out := []*ebiten.Image{}

	for _, co := range color {
		c := gg.NewContext(w, h)
		c.DrawRoundedRectangle(0, 0, float64(w), float64(h), 5)
		c.SetColor(co)
		c.Fill()
		img := ebiten.NewImageFromImage(c.Image())
		out = append(out, img)
	}
	return out
}
