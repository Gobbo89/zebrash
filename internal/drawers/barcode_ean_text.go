package drawers

import (
	"image/color"
	"math"
	"strings"

	"github.com/Gobbo89/zebrash/internal/assets"
	"github.com/Gobbo89/zebrash/internal/barcodes/utils"
	"github.com/Gobbo89/zebrash/internal/elements"
	"github.com/ingridhq/gg"
)

var fontEANDigits = mustLoadFont(assets.FontOCRB)

func drawEANText(ctx *gg.Context, text string, pos elements.LabelPosition, above bool, barHeight, moduleWidth int) error {
	ctx.SetColor(color.Black)
	w := float64(moduleWidth)
	x, y := float64(pos.X), float64(pos.Y)
	if above {
		// Preserve the separate above-bar interpretation layout.
		divisor, modules := 11.0, 67.0
		if len(text) == 13 {
			divisor, modules = 13, 95
		}
		face, err := newRasterFontFace(font0, math.Round(modules*w/divisor))
		if err != nil {
			return err
		}
		defer face.Close()
		ctx.SetFontFace(face)
		var spaced strings.Builder
		for _, digit := range text {
			spaced.WriteRune(digit)
			spaced.WriteByte(' ')
		}
		drawStringJustified(ctx, spaced.String(), x+8*w, y-float64(utils.EANGuardExtension)/2, 0, 0, (modules-16)*w)
		return nil
	}

	// The interpretation digits scale through three-dot modules and then
	// retain their height as the bars grow wider.
	face, err := newRasterFontFace(fontEANDigits, float64(10*min(moduleWidth, 3)))
	if err != nil {
		return err
	}
	defer face.Close()
	ctx.SetFontFace(face)
	// Place the bottom of the digit ink independently of the guard bars.
	// The 1- and 2-dot variants use a smaller interpretation line.
	bottom := [...]float64{0, 11, 18, 27}[min(moduleWidth, 3)]
	var descent float64
	for digit := '0'; digit <= '9'; digit++ {
		bounds, _, _ := face.GlyphBounds(digit)
		descent = max(descent, float64(bounds.Max.Y)/64)
	}
	baseline := y + float64(barHeight) + bottom - descent
	drawDigit := func(digit byte, center float64) {
		bounds, _, _ := face.GlyphBounds(rune(digit))
		inkCenter := (float64(bounds.Min.X) + float64(bounds.Max.X)) / 128
		ctx.DrawString(string(digit), x+center-inkCenter, baseline)
	}
	rightStart := 36
	if len(text) == 13 {
		drawDigit(text[0], -7.5*w)
		text = text[1:]
		rightStart = 50
	}
	half := len(text) / 2
	// Keep the groups next to the center guard. At wide modules, the
	// interpretation digits stop spreading out with the bars.
	pitch := min(7*w, 20)
	leftEdge := float64(rightStart-5) * w
	rightEdge := float64(rightStart) * w
	for i := 0; i < half; i++ {
		drawDigit(text[i], leftEdge-(float64(half-i)-0.5)*pitch)
		drawDigit(text[half+i], rightEdge+(float64(i)+0.5)*pitch)
	}
	return nil
}
