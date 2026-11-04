package hex

import (
	"bufio"
	"compress/zlib"
	"encoding/base64"
	hx "encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"
)

const (
	maxEmbeddedImageBytes = 3 * 1024 * 1024
	maxEmbeddedFontBytes  = 3 * 1024 * 1024
)

func DecodeEscapedString(value string, escapeChar byte) (string, error) {
	replaceChar := string(escapeChar)

	hexEscapeRegex, err := regexp.Compile(regexp.QuoteMeta(replaceChar) + `([0-9A-Fa-f]{2})`)
	if err != nil {
		return "", fmt.Errorf("failed to compile regex: %w", err)
	}

	res := hexEscapeRegex.ReplaceAllStringFunc(value, func(s string) string {
		v, err := hx.DecodeString(strings.TrimPrefix(s, replaceChar))
		if err == nil {
			return string(v)
		}

		return s
	})

	return res, nil
}

func DecodeFontData(data string, totalBytes int) ([]byte, error) {
	return decodeFileData(data, totalBytes, totalBytes, maxEmbeddedFontBytes)
}

func DecodeGraphicFieldData(data string, rowBytes, totalBytes int) ([]byte, error) {
	if rowBytes <= 0 {
		return nil, fmt.Errorf("embedded graphic row width must be positive")
	}
	return decodeFileData(data, rowBytes, totalBytes, maxEmbeddedImageBytes)
}

var compressCounts = map[byte]int{
	'G': 1,
	'H': 2,
	'I': 3,
	'J': 4,
	'K': 5,
	'L': 6,
	'M': 7,
	'N': 8,
	'O': 9,
	'P': 10,
	'Q': 11,
	'R': 12,
	'S': 13,
	'T': 14,
	'U': 15,
	'V': 16,
	'W': 17,
	'X': 18,
	'Y': 19,
	'g': 20,
	'h': 40,
	'i': 60,
	'j': 80,
	'k': 100,
	'l': 120,
	'm': 140,
	'n': 160,
	'o': 180,
	'p': 200,
	'q': 220,
	'r': 240,
	's': 260,
	't': 280,
	'u': 300,
	'v': 320,
	'w': 340,
	'x': 360,
	'y': 380,
	'z': 400,
}

func decodeFileData(data string, rowBytes, totalBytes, maxBytes int) ([]byte, error) {
	if rowBytes < 0 || rowBytes > maxBytes {
		return nil, fmt.Errorf("invalid embedded file row size %d", rowBytes)
	}
	if z64Encoded(data) {
		limit := maxBytes
		if totalBytes > 0 {
			limit = min(totalBytes, limit)
		}
		return decodeZ64(data, limit)
	}

	var (
		result strings.Builder
		line   strings.Builder
	)

	prevLine := ""
	compressCount := 0
	rowHex := rowBytes * 2

	for i := 0; i < len(data); i++ {
		char := data[i]

		if line.Len() >= rowHex {
			if err := validateEmbeddedFileSize(&result, &line, 0, maxBytes); err != nil {
				return nil, err
			}

			prevLine = line.String()
			line.Reset()
			result.WriteString(prevLine)
		}

		if c, ok := compressCounts[char]; ok {
			if c > 2*maxBytes-compressCount {
				return nil, embeddedSizeError(maxBytes)
			}
			compressCount += c
			continue
		}

		switch char {
		case ' ', '\t', '\r', '\n':
			continue
		case ',':
			l := rowHex - line.Len()
			if err := validateEmbeddedFileSize(&result, &line, l, maxBytes); err != nil {
				return nil, err
			}

			if rowHex > line.Len() {
				line.WriteString(strings.Repeat("0", l))
			}

			continue
		case '!':
			l := rowHex - line.Len()
			if err := validateEmbeddedFileSize(&result, &line, l, maxBytes); err != nil {
				return nil, err
			}

			if rowHex > line.Len() {
				line.WriteString(strings.Repeat("1", l))
			}

			continue
		case ':':
			if err := validateEmbeddedFileSize(&result, &line, len(prevLine), maxBytes); err != nil {
				return nil, err
			}

			line.WriteString(prevLine)
			continue
		}

		l := max(compressCount, 1)
		if err := validateEmbeddedFileSize(&result, &line, l, maxBytes); err != nil {
			return nil, err
		}

		line.WriteString(strings.Repeat(string(char), l))
		compressCount = 0
	}

	if line.Len() > 0 {
		if err := validateEmbeddedFileSize(&result, &line, 0, maxBytes); err != nil {
			return nil, err
		}

		result.WriteString(line.String())
	}

	return hx.DecodeString(result.String())
}

func validateEmbeddedFileSize(result, line *strings.Builder, l, maxBytes int) error {
	totalHexL := result.Len() + line.Len() + l
	if totalHexL > 2*maxBytes {
		return embeddedSizeError(maxBytes)
	}

	return nil
}

func embeddedSizeError(maxBytes int) error {
	return fmt.Errorf("embedded file exceeds %d-byte limit", maxBytes)
}

const z64Prefix = ":Z64:"

func z64Encoded(value string) bool {
	return strings.HasPrefix(value, z64Prefix)
}

func decodeZ64(value string, maxBytes int) ([]byte, error) {
	value = strings.TrimPrefix(value, z64Prefix)

	idx := strings.LastIndex(value, ":")
	if idx >= 0 {
		value = value[:idx]
	}

	b := bufio.NewReader(base64.NewDecoder(base64.StdEncoding, strings.NewReader(value)))
	z, err := zlib.NewReader(b)
	if err != nil {
		return nil, err
	}

	defer z.Close()

	data, err := io.ReadAll(io.LimitReader(z, int64(maxBytes)+1))
	if err != nil {
		return nil, err
	}
	if len(data) > maxBytes {
		return nil, embeddedSizeError(maxBytes)
	}
	// zlib may finish before the encoded input. Validate any remaining
	// base64 without allocating a buffer for the entire compressed stream.
	if _, err := io.Copy(io.Discard, b); err != nil {
		return nil, err
	}
	return data, nil
}
