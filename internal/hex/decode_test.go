package hex

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"testing"
)

func TestHexWhitespace(t *testing.T) {
	for _, plain := range []string{"FFFF,FFFFFF", "GFGF", "HF", "FF,:00", "F,00"} {
		want, err := DecodeGraphicFieldData(plain, 1, 0)
		if err != nil {
			t.Fatal(err)
		}
		for _, whitespace := range []string{" ", "\t", "\r", "\n", " \t\r\n"} {
			input := whitespace
			for _, c := range plain {
				input += string(c) + whitespace
			}
			got, err := DecodeGraphicFieldData(input, 1, 0)
			if err != nil || !bytes.Equal(got, want) {
				t.Errorf("%q: got %x, %v; want %x", input, got, err, want)
			}
		}
	}
	got, err := DecodeGraphicFieldData("FF FF,FFFFFF", 1, 0)
	if err != nil || !bytes.Equal(got, []byte{255, 255, 0, 255, 255, 255}) {
		t.Fatalf("PR example: %x, %v", got, err)
	}
	// The PR's b/c counts say 8 bytes, but its actual data decodes to 6.
	got, err = DecodeFontData("F F", 1)
	if err != nil || !bytes.Equal(got, []byte{255}) {
		t.Fatalf("font: %x, %v", got, err)
	}
	if _, err := DecodeGraphicFieldData("F?", 1, 0); err == nil {
		t.Error("invalid hex accepted")
	}
	if _, err := decodeFileData("H F", 1, 0, 1); err != nil {
		t.Fatal(err)
	}
	if _, err := decodeFileData("J F", 1, 0, 1); err == nil {
		t.Error("size limit bypassed")
	}
}

func TestZ64StillDecodes(t *testing.T) {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write([]byte{255, 0, 170})
	w.Close()
	got, err := DecodeGraphicFieldData(":Z64:"+base64.StdEncoding.EncodeToString(compressed.Bytes())+":0000", 1, 3)
	if err != nil || !bytes.Equal(got, []byte{255, 0, 170}) {
		t.Fatalf("Z64: %x, %v", got, err)
	}
}

func TestZ64ExpansionLimits(t *testing.T) {
	makeZ64 := func(size int) string {
		t.Helper()
		var compressed bytes.Buffer
		w := zlib.NewWriter(&compressed)
		if _, err := w.Write(make([]byte, size)); err != nil {
			t.Fatal(err)
		}
		if err := w.Close(); err != nil {
			t.Fatal(err)
		}
		return z64Prefix + base64.StdEncoding.EncodeToString(compressed.Bytes()) + ":0000"
	}

	testCases := []struct {
		name     string
		size     int
		declared int
		wantErr  bool
	}{
		{name: "at declared size", size: 16, declared: 16},
		{name: "over declared size", size: 17, declared: 16, wantErr: true},
		{name: "at hard limit", size: maxEmbeddedImageBytes, declared: maxEmbeddedImageBytes},
		{name: "declaration cannot raise hard limit", size: maxEmbeddedImageBytes + 1, declared: maxEmbeddedImageBytes + 1, wantErr: true},
		{name: "omission cannot bypass hard limit", size: maxEmbeddedImageBytes + 1, wantErr: true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := DecodeGraphicFieldData(makeZ64(tc.size), 1, tc.declared)
			if tc.wantErr {
				if err == nil || data != nil {
					t.Fatalf("got %d bytes and error %v; expected rejection without partial data", len(data), err)
				}
				return
			}
			if err != nil || len(data) != tc.size {
				t.Fatalf("got %d bytes and error %v; want %d bytes", len(data), err, tc.size)
			}
		})
	}
	if data, err := DecodeFontData(makeZ64(maxEmbeddedFontBytes+1), 0); err == nil || data != nil {
		t.Fatalf("oversized Z64 font yielded %d bytes and error %v", len(data), err)
	}
}
