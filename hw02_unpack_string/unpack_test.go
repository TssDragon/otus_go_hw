package hw02unpackstring

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "a4bc2d5e", expected: "aaaabccddddde"},
		{input: "abccd", expected: "abccd"},
		{input: "", expected: ""},
		{input: "aaa0b", expected: "aab"},
		{input: "qwe.rty", expected: "qwe.rty"},
		{input: "as\n5qwe", expected: "as\n\n\n\n\nqwe"},
		{input: `qwe\4\55`, expected: `qwe455555`},
		{input: `qwe\4\5`, expected: `qwe45`},
		{input: `qwe\45`, expected: `qwe44444`},
		{input: `qwe\\5`, expected: `qwe\\\\\`},
		{input: `qwe\\\3`, expected: `qwe\3`},
		{input: "абываыва", expected: "абываыва"},
		{input: "аа3пп4я", expected: "аааапппппя"},
		{input: "      ", expected: "      "},
		{input: "str-3str0", expected: "str---st"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}

func TestUnpackStringHasInvalidDigitPosition(t *testing.T) {
	invalidStrings := []string{"3abc", "45", "aaa10b", "-42", "aaa10b", "abc33"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrInvalidDigitPosition), "actual error %q", err)
		})
	}
}

func TestUnpackStringWrongEscapeSequence(t *testing.T) {
	invalidStrings := []string{"as\\`asd", "\\asd"}
	for _, tc := range invalidStrings {
		tc := tc
		t.Run(tc, func(t *testing.T) {
			_, err := Unpack(tc)
			require.Truef(t, errors.Is(err, ErrWrongEscape), "actual error %q", err)
		})
	}
}

func TestUnpackSpecialCharacters(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{input: "-+)(*&^^%$@", expected: "-+)(*&^^%$@"},
		{input: "_2!4'2;8", expected: "__!!!!'';;;;;;;;"},
		{input: "`", expected: "`"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			result, err := Unpack(tc.input)
			require.NoError(t, err)
			require.Equal(t, tc.expected, result)
		})
	}
}
