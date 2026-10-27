package pdf417

import "testing"

func TestEncodeModuleWidth(t *testing.T) {
	const content = "PDF417 MODULE WIDTH TEST 0001"
	reference, err := Encode(content, 5, 8, 4, 1)
	if err != nil {
		t.Fatal(err)
	}
	if reference.Bounds().Dx() != 137 || reference.Bounds().Dy() != 168 {
		t.Fatalf("unexpected reference bounds: %v", reference.Bounds())
	}

	testCases := []struct {
		name  string
		width int
		scale int
	}{
		{name: "negative width uses minimum", width: -1, scale: 1},
		{name: "zero width uses minimum", width: 0, scale: 1},
		{name: "one dot", width: 1, scale: 1},
		{name: "two dots", width: 2, scale: 2},
		{name: "three dots", width: 3, scale: 3},
		{name: "ten dots", width: 10, scale: 10},
		{name: "eleven dots are not clamped", width: 11, scale: 11},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			img, err := Encode(content, 5, 8, 4, tC.width)
			if err != nil {
				t.Fatal(err)
			}
			if img.Bounds().Dx() != 137*tC.scale || img.Bounds().Dy() != 168 {
				t.Fatalf("unexpected scaled bounds: %v", img.Bounds())
			}
			for y := 0; y < img.Bounds().Dy(); y++ {
				for x := 0; x < img.Bounds().Dx(); x++ {
					if img.At(x, y) != reference.At(x/tC.scale, y) {
						t.Fatalf("encoded pattern differs at (%d,%d)", x, y)
					}
				}
			}
		})
	}
}
