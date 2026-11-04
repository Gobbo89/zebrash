package zebrash

import (
	"bytes"
	"fmt"
	"strings"
	"testing"

	"github.com/Gobbo89/zebrash/drawers"
)

func renderFieldBlock(t *testing.T, fields string) []byte {
	t.Helper()
	labels, err := NewParser().Parse([]byte("^XA^CI28" + fields + "^XZ"))
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Fatalf("got %d labels, want 1", len(labels))
	}
	var output bytes.Buffer
	options := drawers.DrawerOptions{LabelWidthMm: 60, LabelHeightMm: 36, Dpmm: 8, GrayscaleOutput: true}
	if err := NewDrawer().DrawLabelAsPng(labels[0], &output, options); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

func TestDrawFieldBlockOverprintsOverflow(t *testing.T) {
	const lineBreak = `\&`
	testCases := []struct {
		name, origin, font, alignment, separator string
		maxLines, spacing, width                 int
		lines                                    []string
	}{
		{"printer comparison", "^FO25,70", "^A0N,30,30", "L", lineBreak, 2, 0, 150, []string{"FIRST", "MMMM", "I", "X"}},
		{"single line", "^FO25,70", "^A0N,30,30", "L", lineBreak, 1, 0, 150, []string{"MMMM", "I", "X"}},
		{"bottom anchored", "^FT25,190", "^A0N,30,30", "C", lineBreak, 2, 6, 150, []string{"FIRST", "MMMM", "I", "X"}},
		{"rotated right", "^FO100,40", "^A0R,30,30", "R", lineBreak, 2, 6, 150, []string{"FIRST", "MMMM", "I", "X"}},
		{"rotated inverted", "^FO100,40", "^A0I,30,30", "L", lineBreak, 2, -6, 150, []string{"FIRST", "MMMM", "I", "X"}},
		{"automatically wrapped", "^FO25,70", "^A0N,30,30", "L", " ", 2, 0, 115, []string{"WWWW", "MMMM", "HHHH", "NNNN"}},
		{"nonbreaking text", "^FO25,70", "^A0N,30,30", "L", " ", 1, 0, 115, []string{"MMMM", "10,21\u00a0€", "NNNN"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			field := func(content string) string {
				return fmt.Sprintf("%s%s^FB%d,%d,%d,%s,0^FD%s^FS",
					tc.origin, tc.font, tc.width, tc.maxLines, tc.spacing, tc.alignment, content)
			}
			got := renderFieldBlock(t, field(strings.Join(tc.lines, tc.separator)))
			firstBlock := field(strings.Join(tc.lines[:tc.maxLines], lineBreak))
			// Independently draw each extra line at the last slot, using
			// blank NBSP lines to position it without redrawing earlier ink.
			wantFields := firstBlock
			for _, line := range tc.lines[tc.maxLines:] {
				wantFields += field(strings.Repeat("\u00a0"+lineBreak, tc.maxLines-1) + line)
			}
			want := renderFieldBlock(t, wantFields)
			if bytes.Equal(want, renderFieldBlock(t, firstBlock)) {
				t.Fatal("test data does not distinguish overprinting from truncation")
			}
			if !bytes.Equal(got, want) {
				t.Fatal("overflow render differs from independently overprinted fields")
			}
		})
	}
}

func TestDrawJustifiedFieldBlockOverprintsOverflow(t *testing.T) {
	field := func(maxLines int, content string) string {
		return fmt.Sprintf("^FO25,70^A0N,30,30^FB150,%d,0,J,0^FD%s^FS", maxLines, content)
	}
	got := renderFieldBlock(t, field(2, `P Q\&A B\&C D\&E F`))
	// Non-final logical lines remain justified; only the final line is
	// left-aligned, even when all three share the same baseline.
	want := renderFieldBlock(t, field(3, `P Q\&A B\&`+"\u00a0")+
		field(3, "\u00a0"+`\&C D\&`+"\u00a0")+
		field(2, "\u00a0"+`\&E F`))
	if !bytes.Equal(got, want) {
		t.Fatal("justified overflow differs from independently overprinted fields")
	}
}
