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
	var previousRune rune
	isEscapingSequence := false
	lastSymbolApproved := false

	for _, currentRune := range packedString {
		isDigitCurrentChar := unicode.IsDigit(currentRune)
		currentChar := string(currentRune)
		previousChar := string(previousRune)
		lastSymbolApproved = false

		// Для первого символа в подпоследовательности
		if previousRune == 0 {
			if isDigitCurrentChar {
				return "", ErrInvalidString
			}

			previousRune = currentRune
			if isEscapeChar(currentRune) {
				isEscapingSequence = true
			}
			continue
		}

		if isEscapingSequence {
			// Экранировать можно только обратный слэш и цифры
			if !isDigitCurrentChar && !isEscapeChar(currentRune) {
				return "", ErrInvalidString
			}

			previousRune = currentRune
			isEscapingSequence = false
			lastSymbolApproved = true
			continue
		}

		if isDigitCurrentChar {
			charactersCount, _ := strconv.Atoi(currentChar)
			result.WriteString(strings.Repeat(previousChar, charactersCount))
			previousRune = 0
			continue
		}

		result.WriteRune(previousRune)
		previousRune = currentRune

		if isEscapeChar(currentRune) {
			isEscapingSequence = true
			continue
		}
	}

	if previousRune != 0 {
		if !lastSymbolApproved && unicode.IsDigit(previousRune) {
			return "", ErrInvalidString
		}
		result.WriteRune(previousRune)
	}

	return result.String(), nil
}

func isEscapeChar(r rune) bool {
	return string(r) == "\\"
}
