package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(packedString string) (string, error) {

	var result strings.Builder
	var prevRune rune

	for _, currRune := range packedString {

		isDigitCurrentChar := unicode.IsDigit(currRune)
		isDigitPrevChar := unicode.IsDigit(prevRune)
		currChar := string(currRune)

		if prevRune == 0 {
			if isDigitCurrentChar {
				return "", ErrInvalidString
			}
			prevRune = currRune
			continue
		}

		if isDigitCurrentChar && !isDigitPrevChar {

			repeat, _ := strconv.Atoi(currChar)
			result.WriteString(strings.Repeat(string(prevRune), repeat))
			prevRune = 0
			continue
		}

		if isDigitCurrentChar && isDigitPrevChar {
			return "", ErrInvalidString
		}

		if !isDigitCurrentChar {
			result.WriteRune(prevRune)
			prevRune = currRune
		}
	}

	if prevRune != 0 && unicode.IsDigit(prevRune) {
		return "", ErrInvalidString
	}

	if prevRune != 0 {
		result.WriteRune(prevRune)
	}

	return result.String(), nil
}
