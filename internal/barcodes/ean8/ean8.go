package ean8

import (
	"image"
	"image/color"

	"github.com/Gobbo89/zebrash/internal/barcodes/utils"
	"github.com/Gobbo89/zebrash/internal/images"
)

type ean8 struct {
	code           *utils.BitList
	width          int
	height         int
	barWidth       int
	guardExtension int
}

func newEan8(code *utils.BitList, height, barWidth int) *ean8 {
	barWidth = max(1, barWidth)
	height = max(1, height)

	return &ean8{
		code:           code,
		width:          code.Len() * barWidth,
		height:         height,
		barWidth:       barWidth,
		guardExtension: CalculateGuardExtension(barWidth),
	}
}

func (c *ean8) ColorModel() color.Model {
	return color.RGBAModel
}

func (c *ean8) Bounds() image.Rectangle {
	return image.Rect(0, 0, c.width, c.height+c.guardExtension)
}

func (c *ean8) At(x, y int) color.Color {
	x /= c.barWidth

	if x < 0 || x >= c.code.Len() || !c.code.GetBit(x) {
		return images.ColorTransparent
	}

	if isGuardBar(x) {
		return images.ColorBlack
	}

	if y < c.height {
		return images.ColorBlack
	}
	return images.ColorTransparent
}

func isGuardBar(x int) bool {
	// Start guard (first 3 modules: 0-2)
	if x >= 0 && x <= 2 {
		return true
	}
	// Middle guard (5 modules: 31-35)
	if x >= 31 && x <= 35 {
		return true
	}
	// End guard (last 3 modules: 64-66)
	if x >= 64 && x <= 66 {
		return true
	}
	return false
}

func CalculateGuardExtension(barWidth int) int {
	return min(barWidth*5, 20)
}
