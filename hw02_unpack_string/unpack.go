package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var (
	ErrWrongEscape          = errors.New("invalid escape sequence")
	ErrInvalidDigitPosition = errors.New("invalid digit position in string")
)

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
				return "", ErrInvalidDigitPosition
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
				return "", ErrWrongEscape
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
			return "", ErrInvalidDigitPosition
		}
		result.WriteRune(previousRune)
	}

	return result.String(), nil
}

func isEscapeChar(r rune) bool {
	return string(r) == "\\"
}
