package drawers

import (
	"image"
	"strings"
	"unicode"

	"github.com/ingridhq/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// The bundled font 0 has a much narrower NBSP glyph than its space glyph.
// Use the space metrics and ink without changing the rune used for wrapping.
type nbspFace struct{ font.Face }

func spaceGlyph(r rune) rune {
	if r == '\u00a0' {
		return ' '
	}
	return r
}

func (f nbspFace) Glyph(dot fixed.Point26_6, r rune) (image.Rectangle, image.Image, image.Point, fixed.Int26_6, bool) {
	return f.Face.Glyph(dot, spaceGlyph(r))
}

func (f nbspFace) GlyphBounds(r rune) (fixed.Rectangle26_6, fixed.Int26_6, bool) {
	return f.Face.GlyphBounds(spaceGlyph(r))
}

func (f nbspFace) GlyphAdvance(r rune) (fixed.Int26_6, bool) {
	return f.Face.GlyphAdvance(spaceGlyph(r))
}

func (f nbspFace) Kern(left, right rune) fixed.Int26_6 {
	return f.Face.Kern(spaceGlyph(left), spaceGlyph(right))
}

func isBreakableSpace(r rune) bool {
	return unicode.IsSpace(r) && r != '\u00a0'
}

func wrapFieldText(ctx *gg.Context, text string, width float64) []string {
	if !strings.ContainsRune(text, '\u00a0') {
		return ctx.WordWrap(text, width)
	}
	var lines []string
	for _, paragraph := range strings.Split(text, "\n") {
		// Keep the same wrapping as before for paragraphs without an NBSP.
		if !strings.ContainsRune(paragraph, '\u00a0') {
			lines = append(lines, ctx.WordWrap(paragraph, width)...)
			continue
		}
		rest := strings.TrimFunc(paragraph, isBreakableSpace)
		for rest != "" {
			end, lastBreak := len(rest), 0
			wasSpace := false
			for i, r := range rest {
				space := isBreakableSpace(r)
				if space && !wasSpace {
					if w, _ := ctx.MeasureString(rest[:i]); w > width {
						end = i
						if lastBreak > 0 {
							end = lastBreak
						}
						break
					}
					lastBreak = i
				}
				wasSpace = space
			}
			if end == len(rest) && lastBreak > 0 {
				if w, _ := ctx.MeasureString(rest); w > width {
					end = lastBreak
				}
			}
			lines = append(lines, rest[:end])
			rest = strings.TrimLeftFunc(rest[end:], isBreakableSpace)
		}
	}
	return lines
}
