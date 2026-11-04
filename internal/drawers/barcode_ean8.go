package drawers

import (
	"fmt"

	"github.com/Gobbo89/zebrash/drawers"
	"github.com/Gobbo89/zebrash/internal/barcodes/ean8"
	"github.com/Gobbo89/zebrash/internal/barcodes/utils"
	"github.com/Gobbo89/zebrash/internal/elements"
	"github.com/ingridhq/gg"
)

func NewBarcodeEan8Drawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			barcode, ok := element.(*elements.BarcodeEan8WithData)
			if !ok {
				return nil
			}

			img, text, err := ean8.Encode(barcode.Data, barcode.Height, barcode.Width)
			if err != nil {
				return fmt.Errorf("failed to encode EAN-8 barcode: %w", err)
			}

			guardExtension := utils.EANGuardExtension

			pos := adjustImageTypeSetPosition(img, barcode.Position, barcode.Orientation)
			pos = adjustEan8Position(pos, barcode.Orientation, guardExtension)

			rotateImage(gCtx, img, pos, barcode.Orientation)

			defer gCtx.Identity()

			gCtx.DrawImage(img, pos.X, pos.Y)
			if barcode.Line {
				return drawEANText(gCtx, text, pos, barcode.LineAbove, img.Bounds().Dy()-guardExtension, max(1, barcode.Width))
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
