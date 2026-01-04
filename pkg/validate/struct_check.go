package validate

import (
	"fmt"
	"reflect"
	"strings"

	"git.oxl.at/go-validator/pkg/validate/regex"
)

// isFieldUnset checks if a reflect.Value is its zero value, or an empty collection.
func isFieldUnset(field reflect.Value) bool {
	if field.IsZero() {
		return true
	}

	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Array:
		return field.Len() == 0
	}

	return false
}

type StructValidator struct {
	ValidatorsCustom map[string]ValidatorCustom
	ValidatorsRegex  map[string]string
}

// Validate recursively validates a struct based on field tags.
func (v *StructValidator) Validate(s interface{}) []error {
	if v.ValidatorsCustom == nil {
		v.ValidatorsCustom = GetDefaultCustomValidators()
	}
	if v.ValidatorsRegex == nil {
		v.ValidatorsRegex = GetDefaultRegexValidators()
	}

	errors := []error{}

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	// Iterate over all fields in the struct
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := val.Type().Field(i)

		// handle Nested Structs/Pointers (Recursive)
		// If handled, we skip all other checks for this field.
		if nestedErrs, handled := v.validateNestedField(field, fieldType); handled {
			errors = append(errors, nestedErrs...)
			continue
		}

		// handle "required" tag
		if err := v.checkRequiredField(field, fieldType); err != nil {
			errors = append(errors, err)
			continue // Field is required but empty, skip other checks
		}

		// handle "required_if" tag
		if err := v.checkRequiredIfField(field, fieldType, val); err != nil {
			errors = append(errors, err)
			continue // Field is required_if but empty, skip other checks
		}

		// nNot required => skip empty
		// NOTE: fields with explicit default-values will not be skipped
		if isFieldUnset(field) {
			continue
		}

		// handle Slice/Array Element Validation
		// This validates elements *inside* the slice.
		// We don't 'continue' here, as the slice itself might be 'required'.
		if sliceErrs := v.validateSliceField(field, fieldType); sliceErrs != nil {
			errors = append(errors, sliceErrs...)
		}

		// handle "validate" tags (for NON-SLICE, non-struct fields)
		// Skips fields that are containers (struct, slice) or empty/non-required.
		if simpleErrs := v.validateSimpleField(field, fieldType); simpleErrs != nil {
			errors = append(errors, simpleErrs...)
		}
	}

	return errors
}

// validateNestedField handles recursive validation for nested structs and pointers-to-structs.
// It returns errors and a 'bool' indicating if the field was handled (and the loop should continue).
func (v *StructValidator) validateNestedField(field reflect.Value, fieldType reflect.StructField) (errors []error, handled bool) {
	switch field.Kind() {
	case reflect.Ptr:
		if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
			if err := v.Validate(field.Interface()); err != nil {
				for _, e := range err {
					errors = append(errors, fmt.Errorf("%s - %w", fieldType.Name, e))
				}
			}
		}
		return errors, true

	case reflect.Struct:
		if err := v.Validate(field.Interface()); err != nil {
			for _, e := range err {
				errors = append(errors, fmt.Errorf("%s - %w", fieldType.Name, e))
			}
		}
		return errors, true
	}

	return nil, false
}

// validateSliceField handles validation for elements within slices and arrays.
// It supports slices of structs, pointers-to-structs, and basic types.
func (v *StructValidator) validateSliceField(field reflect.Value, fieldType reflect.StructField) (errors []error) {

	fieldKind := field.Kind()
	if fieldKind != reflect.Slice && fieldKind != reflect.Array {
		return nil
	}

	elemType := field.Type().Elem()
	isStructElem := (elemType.Kind() == reflect.Struct) || (elemType.Kind() == reflect.Ptr && elemType.Elem().Kind() == reflect.Struct)

	tag := fieldType.Tag.Get("validate")
	regexTag := fieldType.Tag.Get("validate_regex")
	hasBasicValidation := (tag != "" || regexTag != "") && !isStructElem

	validatorCustom, existsCustom := v.ValidatorsCustom[tag]
	checkRegex, existsRegex := v.ValidatorsRegex[tag]

	// Pre-check for unknown tag on basic slice validation
	if hasBasicValidation && tag != "" && !existsCustom && !existsRegex {
		errors = append(errors, fmt.Errorf("no validator registered for tag: '%s' on field '%s'", tag, fieldType.Name))
	}

	for j := 0; j < field.Len(); j++ {
		elem := field.Index(j)

		if isStructElem {
			if (elem.Kind() == reflect.Ptr && !elem.IsNil() && elem.Elem().Kind() == reflect.Struct) || elem.Kind() == reflect.Struct {
				if errs := v.Validate(elem.Interface()); errs != nil {
					for _, e := range errs {
						errors = append(errors, fmt.Errorf("%s[%d] - %w", fieldType.Name, j, e))
					}
				}
			}

		} else if hasBasicValidation {
			elemVal := elem.Interface()
			if tag != "" {
				isValid := false
				if existsCustom {
					isValid = validatorCustom(elemVal)

				} else if existsRegex {
					isValid = regex.ValidateRegex(checkRegex, elemVal)
				}

				if !isValid {
					errors = append(errors, fmt.Errorf("field '%s[%d]' ('%v') failed validation '%s'", fieldType.Name, j, elemVal, tag))
				}

			} else if regexTag != "" {
				if !regex.ValidateRegex(regexTag, elemVal) {
					errors = append(errors, fmt.Errorf("field '%s[%d]' ('%v') failed validation '%s'", fieldType.Name, j, elemVal, regexTag))
				}
			}
		}
	}

	return errors
}

// checkRequiredField handles the 'required:"true"' tag.
func (v *StructValidator) checkRequiredField(field reflect.Value, fieldType reflect.StructField) error {
	requiredTag := fieldType.Tag.Get("required")
	if requiredTag == "true" && isFieldUnset(field) {
		return fmt.Errorf("field '%s' is required but is empty", fieldType.Name)
	}
	return nil
}

// checkRequiredIfField handles the 'required_if' tag.
func (v *StructValidator) checkRequiredIfField(field reflect.Value, fieldType reflect.StructField, structVal reflect.Value) error {
	requiredIfTag := fieldType.Tag.Get("required_if")
	if requiredIfTag == "" {
		return nil
	}

	parts := strings.SplitN(requiredIfTag, "=", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid required_if tag format on field '%s'", fieldType.Name)
	}

	otherFieldName := parts[0]
	regexPattern := parts[1]

	otherField := structVal.FieldByName(otherFieldName)
	if !otherField.IsValid() {
		return fmt.Errorf("invalid field name '%s' in required_if tag on field '%s'", otherFieldName, fieldType.Name)
	}

	if regex.ValidateRegex(regexPattern, otherField.Interface()) {
		if isFieldUnset(field) {
			return fmt.Errorf("field '%s' is required because '%s' is set", fieldType.Name, otherFieldName)
		}
	}
	return nil
}

// validateSimpleField handles 'validate' and 'validate_regex' for simple (non-container) fields.
func (v *StructValidator) validateSimpleField(field reflect.Value, fieldType reflect.StructField) (errors []error) {

	fieldKind := field.Kind()

	// If field is a container, its elements were already validated.
	if fieldKind == reflect.Slice || fieldKind == reflect.Array || fieldKind == reflect.Struct {
		return nil
	}

	tag := fieldType.Tag.Get("validate")
	regexTag := fieldType.Tag.Get("validate_regex")

	if tag == "" && regexTag == "" {
		return nil // No validation to perform
	}

	// 'validate' tag takes precedence
	if tag != "" {
		validatorCustom, existsCustom := v.ValidatorsCustom[tag]
		checkRegex, existsRegex := v.ValidatorsRegex[tag]

		isValid := false
		if existsCustom {
			isValid = validatorCustom(field.Interface())
		} else if existsRegex {
			isValid = regex.ValidateRegex(checkRegex, field.Interface())
		} else {
			errors = append(errors, fmt.Errorf("no validator registered for tag: '%s' on field '%s'", tag, fieldType.Name))
			return errors
		}

		if !isValid {
			errors = append(errors, fmt.Errorf("field '%s' failed validation '%s'", fieldType.Name, tag))
		}

	} else if regexTag != "" { // Only use 'validate_regex' if 'validate' is not set
		if !regex.ValidateRegex(regexTag, field.Interface()) {
			errors = append(errors, fmt.Errorf("field '%s' failed validation '%s'", fieldType.Name, regexTag))
		}
	}

	return errors
}
