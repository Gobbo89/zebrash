package drawers

import (
	"fmt"
	"image/color"
	"math"

	"github.com/Gobbo89/zebrash/drawers"
	"github.com/Gobbo89/zebrash/internal/barcodes/ean8"
	"github.com/Gobbo89/zebrash/internal/elements"
	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/gg"
)

func NewBarcodeEan8Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeEan8WithData)
			if !ok {
				return nil
			}

			// We will create this ean8 package next!
			img, text, err := ean8.Encode(barcode.Data, barcode.Height, barcode.Width)
			if err != nil {
				return fmt.Errorf("failed to encode EAN-8 barcode: %w", err)
			}

			width := float64(img.Bounds().Dx())
			height := float64(img.Bounds().Dy())
			guardExtension := ean8.CalculateGuardExtension(barcode.Width)

			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)
			pos = adjustEan8Position(pos, barcode.Orientation, guardExtension)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				applyEan8TextToCtx(gCtx, text, pos, barcode.LineAbove, width, height, float64(barcode.Width), float64(guardExtension))
			}

			return nil
		},
	}
}

func adjustEan8Position(pos elements.LabelPosition, ori elements.FieldOrientation, guardExtension int) elements.LabelPosition {
	if pos.CalculateFromBottom {
		return pos
	}

	x := pos.X
	y := pos.Y

	switch ori {
	case elements.FieldOrientation90:
		x -= guardExtension
	case elements.FieldOrientation180:
		y -= guardExtension
	}

	return elements.LabelPosition{
		X: x,
		Y: y,
	}
}

func applyEan8TextToCtx(gCtx *gg.Context, text string, pos elements.LabelPosition, lineAbove bool, width, height, barWidth, guardExtension float64) {
	gCtx.SetColor(color.Black)

	fontSize := math.Round(width / 11)
	face := truetype.NewFace(font0, &truetype.Options{Size: fontSize})
	gCtx.SetFontFace(face)

	if len(text) == 8 && !lineAbove {
		formattedText := fmt.Sprintf("%s||%s", text[0:4], text[4:8])
		formattedText = wrapWithSpaces(formattedText)

		// Add a visual margin (2 modules) so the text doesn't touch the guard bars
		padding := float64(barWidth * 2)

		// Shift X to the right by the padding amount
		x := float64(pos.X) + float64(barWidth*3) + padding
		y := float64(pos.Y) + height + fontSize - guardExtension

		// Reduce the total width by padding on BOTH sides
		w := width - float64(barWidth*6) - (padding * 2)

		drawStringJustified(gCtx, formattedText, x, y, 0, 0, w, []string{"|"})
	} else {
		x := float64(pos.X) + float64(barWidth*8)
		y := float64(pos.Y) - float64(guardExtension/2)
		w := width - float64(barWidth*16)
		formattedText := wrapWithSpaces(text)
		drawStringJustified(gCtx, formattedText, x, y, 0, 0, w, nil)
	}
}
