package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

// Return true if valid
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
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

// Return true if a value is in a list of permitted ints
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
