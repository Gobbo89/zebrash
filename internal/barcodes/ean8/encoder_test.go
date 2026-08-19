package ean8

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitizeContent(t *testing.T) {
	testCases := []struct {
		name    string
		content string
		result  string
	}{
		{
			name:    "less than 8 digits",
			content: "123",
			result:  "00001236", // check digit for 0000123 is 6
		},
		{
			name:    "exactly 8 digits with valid check digit",
			content: "96385074", // Standard EAN-8 test case
			result:  "96385074",
		},
		{
			name:    "exactly 8 digits with invalid check digit",
			content: "96385071", // will be fixed by generator
			result:  "96385074",
		},
		{
			name:    "more than 8 digits",
			content: "123456789", // truncates to 1234567, check digit 0
			result:  "12345670",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.name, func(t *testing.T) {
			res := sanitizeContent(tC.content)

			if diff := cmp.Diff(res, tC.result); diff != "" {
				t.Errorf("mismatched result (-got,+want):\n %s", diff)
			}
		})
	}
}
