package accounts

import (
	"testing"
	"reflect"
)

func TestNormalizeEmail(t *testing.T) {
	tcs := map[string]struct{
		input, expected string
	} {
		"all_upper": {input:"MABEL@EXAMPLE.COM", expected:"mabel@example.com"},
		"suf_upper": {input:"maBEL@EXAMPLE.COM", expected:"mabel@example.com"},
		"pre_upper": {input:"MABel@EXAMPLE.COM", expected:"mabel@example.com"},
	}

	for name, tc := range tcs {
		t.Run(name, func(t *testing.T) {
			actual := NormalizeEmail(tc.input)
			if ! reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("actual: %v vs expected: %v",
					actual, tc.expected)
			}
		})
	}
}
