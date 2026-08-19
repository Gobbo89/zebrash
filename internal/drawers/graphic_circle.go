package drawers

import (
	"github.com/Gobbo89/zebrash/drawers"
	"github.com/Gobbo89/zebrash/internal/elements"
	"github.com/ingridhq/gg"
)

func NewGraphicCircleDrawer() *ElementDrawer {
	return &ElementDrawer{
		Draw: func(gCtx *gg.Context, element any, _ drawers.DrawerOptions, _ *DrawerState) error {
			circle, ok := element.(*elements.GraphicCircle)
			if !ok {
				return nil
			}

			setLineColor(gCtx, circle.LineColor)
			gCtx.SetLineWidth(float64(circle.BorderThickness))

			radius := float64(circle.CircleDiameter) / 2.0
			gCtx.DrawCircle(float64(circle.Position.X)+radius, float64(circle.Position.Y)+radius, radius)

			gCtx.Stroke()

			return nil
		},
	}
}
