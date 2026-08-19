package ean8

import (
	"fmt"
	"image"
	"strings"
	"unicode"

	"github.com/Gobbo89/zebrash/internal/barcodes/utils"
)

type encodedNumber struct {
	Left  []bool
	Right []bool
}

var encoderTable = map[rune]encodedNumber{
	'0': {[]bool{false, false, false, true, true, false, true}, []bool{true, true, true, false, false, true, false}},
	'1': {[]bool{false, false, true, true, false, false, true}, []bool{true, true, false, false, true, true, false}},
	'2': {[]bool{false, false, true, false, false, true, true}, []bool{true, true, false, true, true, false, false}},
	'3': {[]bool{false, true, true, true, true, false, true}, []bool{true, false, false, false, false, true, false}},
	'4': {[]bool{false, true, false, false, false, true, true}, []bool{true, false, true, true, true, false, false}},
	'5': {[]bool{false, true, true, false, false, false, true}, []bool{true, false, false, true, true, true, false}},
	'6': {[]bool{false, true, false, true, true, true, true}, []bool{true, false, true, false, false, false, false}},
	'7': {[]bool{false, true, true, true, false, true, true}, []bool{true, false, false, false, true, false, false}},
	'8': {[]bool{false, true, true, false, true, true, true}, []bool{true, false, false, true, false, false, false}},
	'9': {[]bool{false, false, false, true, false, true, true}, []bool{true, true, true, false, true, false, false}},
}

func calcCheckNum(code string) rune {
	x3 := len(code) == 7
	sum := 0
	for _, r := range code {
		curNum := utils.RuneToInt(r)
		if curNum < 0 || curNum > 9 {
			return 'B'
		}
		if x3 {
			curNum = curNum * 3
		}
		x3 = !x3
		sum += curNum
	}

	return utils.IntToRune((10 - (sum % 10)) % 10)
}

func sanitizeContent(content string) string {
	content = strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return '0'
	}, content)

	content = fmt.Sprintf("%07s", content)
	if len(content) > 8 {
		content = content[0:7]
	}

	return withCheckDigit(content[0:7])
}

func withCheckDigit(content string) string {
	return content + string(calcCheckNum(content))
}

func Encode(content string, height, barWidth int) (image.Image, string, error) {
	content = sanitizeContent(content)
	code := encodeEAN8(content)
	return newEan8(code, height, barWidth), content, nil
}

func encodeEAN8(code string) *utils.BitList {
	result := new(utils.BitList)
	result.AddBit(true, false, true) // Start Guard

	for cpos, r := range code {
		num, ok := encoderTable[r]
		if !ok {
			return nil
		}

		if cpos == 4 { // Middle Guard
			result.AddBit(false, true, false, true, false)
		}

		if cpos < 4 {
			result.AddBit(num.Left...)
		} else {
			result.AddBit(num.Right...)
		}
	}
	result.AddBit(true, false, true) // End Guard
	return result
}
