package hex

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"
)

func TestHexWhitespace(t *testing.T) {
	for _, plain := range []string{"FFFF,FFFFFF", "GFGF", "HF", "FF,:00", "F,00"} {
		want, err := DecodeGraphicFieldData(plain, 1)
		if err != nil {
			t.Fatal(err)
		}
		for _, whitespace := range []string{" ", "\t", "\r", "\n", " \t\r\n"} {
			input := whitespace
			for _, c := range plain {
				input += string(c) + whitespace
			}
			got, err := DecodeGraphicFieldData(input, 1)
			if err != nil || !bytes.Equal(got, want) {
				t.Errorf("%q: got %x, %v; want %x", input, got, err, want)
			}
		}
	}
	got, err := DecodeGraphicFieldData("FF FF,FFFFFF", 1)
	if err != nil || !bytes.Equal(got, []byte{255, 255, 0, 255, 255, 255}) {
		t.Fatalf("PR example: %x, %v", got, err)
	}
	// The PR's b/c counts say 8 bytes, but its actual data decodes to 6.
	got, err = DecodeFontData("F F", 1)
	if err != nil || !bytes.Equal(got, []byte{255}) {
		t.Fatalf("font: %x, %v", got, err)
	}
	if _, err := DecodeGraphicFieldData("F?", 1); err == nil {
		t.Error("invalid hex accepted")
	}
	if _, err := decodeFileData("H F", 1, 0); err == nil {
		t.Error("size limit bypassed")
	}
}

func TestZ64StillDecodes(t *testing.T) {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write([]byte{255, 0, 170})
	w.Close()
	got, err := DecodeGraphicFieldData(":Z64:"+base64.StdEncoding.EncodeToString(compressed.Bytes())+":0000", 1)
	if err != nil || !bytes.Equal(got, []byte{255, 0, 170}) {
		t.Fatalf("Z64: %x, %v", got, err)
	}
}
