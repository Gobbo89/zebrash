package code128

import "testing"

func TestSubsetCSwitch(t *testing.T) {
	cases := []struct {
		input, text string
		symbols     []byte
	}{
		{">:ABC>512345678", "ABC12345678", []byte{104, 33, 34, 35, 99, 12, 34, 56, 78}},
		{">51234", "1234", []byte{104, 99, 12, 34}},
		{">912>53456", ",3456", []byte{103, 12, 99, 34, 56}},
		{">;12>534", "1234", []byte{105, 12, 34}},
		{">:AB>51234>6CD>55678", "AB1234CD5678", []byte{104, 33, 34, 99, 12, 34, 100, 35, 36, 99, 56, 78}},
		{">:AB>512>834", "AB1234", []byte{104, 33, 34, 99, 12, 102, 34}},
	}
	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			actual, text, err := EncodeNoMode(tc.input, 20, 3)
			if err != nil {
				t.Fatal(err)
			}
			if text != tc.text {
				t.Fatalf("text = %q, want %q", text, tc.text)
			}
			expected, err := encode(tc.symbols, 20, 3)
			if err != nil {
				t.Fatal(err)
			}
			if actual.Bounds() != expected.Bounds() {
				t.Fatalf("bounds = %v, want %v", actual.Bounds(), expected.Bounds())
			}
			for x := 0; x < actual.Bounds().Dx(); x++ {
				if actual.At(x, 0) != expected.At(x, 0) {
					t.Fatalf("bar pattern differs at x=%d", x)
				}
			}
		})
	}
}
