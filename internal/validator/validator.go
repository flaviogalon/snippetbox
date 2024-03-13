package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
}

var EmailRX = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

// Return true if valid
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// Add an error message to the errors mapping
func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// Add an error message to the errors mapping only if a validation is nok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// Return true if a value if not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// Return true if a value contains no more than N chars
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// Return true if a value is in a list of permitted values
func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

// Return true if value has at least n chars
func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

// Return true if value matches the provided compiled regex
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
