package revalid

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

const defaultTagName = "regex"

var (
	ErrUnexpected = errors.New("unexpected error occurred")
	ErrValidating = errors.New("validation error")
)

// ValidationError is what it is
// we can catch it type in logistics level.
type ValidationError struct {
	entity   string
	property string
	isZero   bool
	code     string
	pattern  string
}

// Error implements error interface
// and can distinct if non-zero value was provided.
func (vErr *ValidationError) Error() string {
	return strings.ToLower(fmt.Sprintf("%s.%s is required according to pattern: `%s`",
		vErr.entity,
		vErr.property,
		vErr.pattern,
	))
}

// ValidateStruct validates struct fields
// according to given regex tag
func ValidateStruct(src any) (err error) {
	// check if src is a struct
	srcValue, err := isStructValue(src)
	if err != nil {
		return err
	}

	// top level struct name (in case we are using nested structs)
	structName := srcValue.Type().Name()

	// iterate  all over struct fields
	for i := 0; i < srcValue.NumField(); i++ {
		fieldValue := srcValue.Field(i)
		fieldType := srcValue.Type().Field(i)
		// check presence of regex tag (.Tag.Lookup() would not work here)
		if pattern, ok := getTagValue(fieldType.Tag, defaultTagName); ok {
			// field validation according to pattern
			if !regexp.MustCompile(pattern).MatchString(fmt.Sprint(fieldValue)) {
				return fmt.Errorf("%s: %w", ErrValidating,
					&ValidationError{
						entity:   structName,
						property: fieldType.Name,
						pattern:  pattern},
				)
			}
		}
		// recursive call for nested structs
		if fieldValue.Type().Kind() != reflect.Struct {
			continue
		}

		if err := ValidateStruct(fieldValue.Interface()); err != nil {
			return fmt.Errorf("error validating nested struct: %w", err)
		}
	}
	// in case of panic we will return ErrUnexpected, but we won't panic.
	return err
}

// getTagValue is designed because luck of functionality in reflect.Tag.Lookup()
// and help retrieve <value> in given <key> from struct fields
func getTagValue(tag reflect.StructTag, key string) (string, bool) {
	structTag := fmt.Sprintf("%v", tag)
	tagValue := fmt.Sprintf(`(?s)(?i)\s*(?P<key>%s):\"(?P<value>[^\"]+)\"`, key)

	if match := regexp.MustCompile(tagValue).
		FindStringSubmatch(structTag); match != nil {
		return match[2], true
	}

	return "", false
}

// isStructValue check if validation object struct or not.
func isStructValue(src any) (*reflect.Value, error) {
	var err error
	defer func() {
		if recover() != nil {
			err = ErrUnexpected
		}
	}()

	srcValue := reflect.Indirect(reflect.ValueOf(src))

	if srcType := srcValue.Kind(); srcType != reflect.Struct {
		return nil, fmt.Errorf("input value must be struct, got: %v", srcType)
	}

	return &srcValue, err
}
