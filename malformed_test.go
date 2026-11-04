package zebrash

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"strings"
	"testing"

	"github.com/Gobbo89/zebrash/drawers"
)

func TestDrawMalformedFields(t *testing.T) {
	const box = "^FO10,10^GB40,40,2^FS"
	options := drawers.DrawerOptions{LabelWidthMm: 30, LabelHeightMm: 20, Dpmm: 8}
	render := func(t *testing.T, fields string) ([]byte, error) {
		t.Helper()
		labels, err := NewParser().Parse([]byte("^XA" + fields + "^XZ"))
		if err != nil {
			t.Fatal(err)
		}
		if len(labels) != 1 {
			t.Fatalf("got %d labels, want 1", len(labels))
		}
		var b bytes.Buffer
		err = NewDrawer().DrawLabelAsPng(labels[0], &b, options)
		return b.Bytes(), err
	}
	want, err := render(t, box)
	if err != nil {
		t.Fatal(err)
	}
	for _, field := range []string{"^GF", "^GFA,0,0,0,^FS", "^GFA,1,1,-1,FF^FS", "^GFA,1,1,0,:Z64:invalid:0000^FS"} {
		got, err := render(t, box+field)
		if err != nil || !bytes.Equal(got, want) {
			t.Errorf("%q should be skipped without affecting the box: %v", field, err)
		}
	}
	for _, tc := range []struct{ field, wantError string }{
		{"^FO10,10^BCN,100^FDJOÃO^FS", "Code 128 cannot encode"},
		{"^FO0,0^A0N,9000,9000^FDW^FS", "font size 9000"},
	} {
		got, err := render(t, box+tc.field)
		if err == nil || !strings.Contains(err.Error(), tc.wantError) || len(got) != 0 {
			t.Errorf("%q yielded %d PNG bytes and error %v", tc.field, len(got), err)
		}
	}
	if _, err := render(t, "^FO0,0^A0N,500,500^FDW^FS"); err != nil {
		t.Fatalf("a valid font larger than the label should still render: %v", err)
	}
}

func TestParseZ64RejectsDataBeyondDeclaredSize(t *testing.T) {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	if _, err := w.Write([]byte{0xff, 0xff}); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	data := ":Z64:" + base64.StdEncoding.EncodeToString(compressed.Bytes()) + ":0000"
	for _, prefix := range []string{"^GFA,1,1,1,", "~DGR:TEST.GRF,1,1,", "~DUR:TEST.TTF,1,"} {
		_, err := NewParser().Parse([]byte("^XA" + prefix + data + "^FS^XZ"))
		if err == nil || !strings.Contains(err.Error(), "exceeds 1-byte limit") {
			t.Errorf("%s returned %v; expected an expansion limit error", prefix, err)
		}
	}
}

func TestDrawMagnifiedStoredGraphic(t *testing.T) {
	const zpl = "~DGR:TEST.GRF,2,1,8040^XA^FO10,10^XGR:TEST.GRF,2,3^FS^XZ"
	labels, err := NewParser().Parse([]byte(zpl))
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Fatalf("got %d labels, want 1", len(labels))
	}
	var b bytes.Buffer
	options := drawers.DrawerOptions{LabelWidthMm: 4, LabelHeightMm: 3, Dpmm: 8}
	if err := NewDrawer().DrawLabelAsPng(labels[0], &b, options); err != nil {
		t.Fatal(err)
	}
	img := mustDecodePng(b.Bytes(), t)
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			wantBlack := (x >= 10 && x < 12 && y >= 10 && y < 13) ||
				(x >= 12 && x < 14 && y >= 13 && y < 16)
			gray, _, _, _ := img.At(x, y).RGBA()
			if (gray < 0x8000) != wantBlack {
				t.Fatalf("unexpected pixel at (%d, %d): %d", x, y, gray)
			}
		}
	}
}
