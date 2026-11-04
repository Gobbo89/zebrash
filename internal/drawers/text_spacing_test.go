package drawers

import (
	"slices"
	"testing"

	"github.com/golang/freetype/truetype"
	"github.com/ingridhq/gg"
)

func TestFieldBlockNonbreakingSpace(t *testing.T) {
	ctx := gg.NewContext(100, 50)
	face := truetype.NewFace(font0, &truetype.Options{Size: 20})
	defer face.Close()
	ctx.SetFontFace(nbspFace{face})
	const amount = "10,21\u00a0€"
	width, _ := ctx.MeasureString(amount)
	spaceWidth, _ := ctx.MeasureString("10,21 €")
	if width != spaceWidth {
		t.Fatalf("NBSP width %v differs from space width %v", width, spaceWidth)
	}
	for _, limit := range []float64{width, width - 1} {
		got := wrapFieldText(ctx, "X "+amount+" Y", limit)
		want := []string{"X", amount, "Y"}
		if !slices.Equal(got, want) {
			t.Errorf("wrap at width %v = %q, want %q", limit, got, want)
		}
	}
}
