package drawers

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/Gobbo89/zebrash/drawers"
	"github.com/Gobbo89/zebrash/internal/elements"
	"github.com/Gobbo89/zebrash/internal/images"
	"github.com/ingridhq/gg"
)

func NewGraphicFieldDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			field, ok := element.(*elements.GraphicField)
			if !ok || field.RowBytes <= 0 {
				return nil
			}

			dataLen := len(field.Data)
			if field.TotalBytes > 0 {
				dataLen = min(field.TotalBytes, dataLen)
			}

			height := dataLen / field.RowBytes
			sx, sy := field.MagnificationX, field.MagnificationY
			if height == 0 || sx <= 0 || sy <= 0 {
				return nil
			}
			if field.RowBytes > math.MaxInt/8/sx || height > math.MaxInt/sy {
				return fmt.Errorf("embedded graphic dimensions overflow")
			}

			// Keep the graphic packed and scale it on demand. Drawing a
			// magnified field should not allocate an intermediate bitmap.
			img := packedGraphic{
				data:     field.Data,
				rowBytes: field.RowBytes,
				bounds:   image.Rect(0, 0, field.RowBytes*8*sx, height*sy),
				sx:       sx,
				sy:       sy,
			}
			gCtx.DrawImage(img, field.Position.X, field.Position.Y)

			return nil
		},
	}
}

type packedGraphic struct {
	data     []byte
	rowBytes int
	bounds   image.Rectangle
	sx, sy   int
}

func (g packedGraphic) ColorModel() color.Model { return color.RGBAModel }
func (g packedGraphic) Bounds() image.Rectangle { return g.bounds }
func (g packedGraphic) At(x, y int) color.Color {
	if image.Pt(x, y).In(g.bounds) {
		x /= g.sx
		y /= g.sy
		if g.data[y*g.rowBytes+x/8]&(0x80>>uint(x%8)) != 0 {
			return images.ColorBlack
		}
	}
	return color.Transparent
}
