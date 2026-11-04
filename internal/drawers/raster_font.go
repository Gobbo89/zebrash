package drawers

import (
	"fmt"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	maxRasterFontSize = 4096
	maxGlyphMaskBytes = 16 * 1024 * 1024
	defaultGlyphCache = 512
)

// Freetype allocates one mask of the font-wide bounding-box size per
// cache entry. Reduce the cache for large fonts instead of resizing text
// to the label, and reject sizes that cannot safely be rasterized.
func newRasterFontFace(ttf *truetype.Font, size float64) (font.Face, error) {
	if !(size > 0 && size <= maxRasterFontSize) {
		return nil, fmt.Errorf("font size %g outside supported range (0, %d] dots", size, maxRasterFontSize)
	}
	units := int64(ttf.FUnitsPerEm())
	if units <= 0 {
		return nil, fmt.Errorf("font has invalid units per em")
	}
	scale := int64(0.5 + size*64)
	raw := ttf.Bounds(fixed.Int26_6(units))
	for _, coord := range []fixed.Int26_6{raw.Min.X, raw.Min.Y, raw.Max.X, raw.Max.Y} {
		product := int64(coord) * scale
		if product > math.MaxInt32-units/2 || product < math.MinInt32+units/2 {
			return nil, fmt.Errorf("font size %g exceeds rasterizer coordinate range", size)
		}
	}
	box := ttf.Bounds(fixed.Int26_6(scale))
	// Match the rasterizer's rounding when estimating the allocation.
	w := ((int64(box.Max.X) + 63) >> 6) - (int64(box.Min.X) >> 6)
	h := ((-int64(box.Min.Y) + 63) >> 6) - (-int64(box.Max.Y) >> 6)
	if w < 0 || h < 0 || w*h > maxGlyphMaskBytes {
		return nil, fmt.Errorf("font size %g exceeds glyph mask memory limit", size)
	}
	entries := defaultGlyphCache
	for w*h*int64(entries) > maxGlyphMaskBytes {
		entries /= 2
	}
	return truetype.NewFace(ttf, &truetype.Options{Size: size, GlyphCacheEntries: entries}), nil
}
