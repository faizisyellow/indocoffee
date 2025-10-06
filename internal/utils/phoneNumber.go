package utils

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

// Exported, prefixed error variables for phone number validation
var (
	ErrPhoneValidationEmpty       = errors.New("phone validation: phone number is empty")
	ErrPhoneValidationInvalidPlus = errors.New("phone validation: '+' is only allowed at the beginning")
	ErrPhoneValidationInvalidChar = errors.New("phone validation: contains invalid character")
	ErrPhoneValidationTooShort    = errors.New("phone validation: must have at least 10 digits")
	ErrPhoneValidationTooLong     = errors.New("phone validation: must have at most 15 digits")
)

// ValidateAndFormatPhoneNumber performs basic validation and returns a cleaned phone number.
// It strips common formatting characters and ensures the number is valid for storage.
func ValidateAndFormatPhoneNumber(raw string) (string, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", ErrPhoneValidationEmpty
	}

	var cleaned strings.Builder
	plusSeen := false

	for i, r := range raw {
		switch {
		case r == '+':
			if i == 0 {
				plusSeen = true
				cleaned.WriteRune('+')
			} else {
				return "", ErrPhoneValidationInvalidPlus
			}
		case unicode.IsDigit(r):
			cleaned.WriteRune(r)
		case unicode.IsSpace(r) || r == '-' || r == '(' || r == ')':
			continue // skip common formatting
		default:
			// Wrap the base error to include which character caused the issue
			return "", fmt.Errorf("%w: '%c'", ErrPhoneValidationInvalidChar, r)
		}
	}

	cleanedStr := cleaned.String()

	// Basic digit count check
	digitCount := len(cleanedStr)
	if plusSeen {
		digitCount--
	}

	if digitCount < 10 {
		return "", ErrPhoneValidationTooShort
	}
	if digitCount > 15 {
		return "", ErrPhoneValidationTooLong
	}

	return cleanedStr, nil
}
