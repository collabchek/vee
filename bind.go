package vee

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Bind parses HTTP form data and populates the provided struct.
// The struct pointer v will be populated with form data.
func Bind(r any, v any) error {
	// For now, expect r to be url.Values (we'll enhance this later for http.Request)
	values, ok := r.(map[string][]string)
	if !ok {
		return fmt.Errorf("vee: expected url.Values or map[string][]string, got %T", r)
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Must be pointer to struct
	if typ.Kind() != reflect.Ptr {
		return fmt.Errorf("vee: expected pointer to struct, got %v", typ.Kind())
	}

	val = val.Elem()
	typ = typ.Elem()

	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("vee: expected pointer to struct, got pointer to %v", typ.Kind())
	}

	// Validate Choices/Chosen pairs
	choicesChosenPairs, err := validateChoicesChosen(typ, val)
	if err != nil {
		return err
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldVal := val.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Parse vee tag
		veeTag := field.Tag.Get("vee")
		config := parseVeeTag(veeTag, field.Name)

		// Skip if requested
		if config.Skip {
			continue
		}

		// Skip Choices fields (they're not bound from form data)
		if strings.HasSuffix(field.Name, "Choices") {
			continue
		}

		// Handle Chosen fields specially
		if strings.HasSuffix(field.Name, "Chosen") {
			baseName := strings.TrimSuffix(field.Name, "Chosen")
			if pair, exists := choicesChosenPairs[baseName]; exists {
				err := bindMultiValueField(values, fieldVal, pair, config)
				if err != nil {
					return err
				}
				continue
			}
		}

		// Bind based on field type
		if !fieldVal.CanSet() {
			continue
		}

		// Handle pointer types
		actualType := field.Type
		isPointer := false

		if actualType.Kind() == reflect.Ptr {
			isPointer = true
			actualType = actualType.Elem()
		}

		// Check for specific types first (before generic kind matching)
		if actualType == reflect.TypeOf(time.Time{}) {
			// For time fields, skip if no form data
			formValues, exists := values[config.Name]
			if !exists || len(formValues) == 0 {
				continue
			}

			formValue := formValues[0]

			// Determine expected format based on type attribute
			inputType := "datetime-local" // default
			if typeAttr, ok := config.Attributes["type"]; ok {
				switch typeAttr {
				case "date", "datetime-local", "time":
					inputType = typeAttr
				}
			}

			// Parse based on input type
			var timeVal time.Time
			var err error
			switch inputType {
			case "date":
				timeVal, err = time.Parse("2006-01-02", formValue)
			case "time":
				timeVal, err = time.Parse("15:04", formValue)
			case "datetime-local":
				timeVal, err = time.Parse("2006-01-02T15:04", formValue)
			}

			if err != nil {
				return fmt.Errorf("vee: cannot parse '%s' as time for field '%s': %w", formValue, config.Name, err)
			}

			if isPointer {
				fieldVal.Set(reflect.ValueOf(&timeVal))
			} else {
				fieldVal.Set(reflect.ValueOf(timeVal))
			}
			continue
		}

		if actualType == reflect.TypeOf(time.Duration(0)) {
			// For duration fields, skip if no form data
			formValues, exists := values[config.Name]
			if !exists || len(formValues) == 0 {
				continue
			}

			formValue := formValues[0]

			// Get units (default to seconds)
			units := "s"
			if unitsAttr, ok := config.Attributes["units"]; ok {
				switch unitsAttr {
				case "ms", "s", "m", "h":
					units = unitsAttr
				}
			}

			// Parse the numeric value and multiply by unit constant
			floatVal, err := strconv.ParseFloat(formValue, 64)
			if err != nil {
				return fmt.Errorf("vee: cannot parse '%s' as duration for field '%s': %w", formValue, config.Name, err)
			}

			var duration time.Duration
			switch units {
			case "ms":
				duration = time.Duration(floatVal) * time.Millisecond
			case "s":
				duration = time.Duration(floatVal) * time.Second
			case "m":
				duration = time.Duration(floatVal) * time.Minute
			case "h":
				duration = time.Duration(floatVal) * time.Hour
			}

			if isPointer {
				fieldVal.Set(reflect.ValueOf(&duration))
			} else {
				fieldVal.Set(reflect.ValueOf(duration))
			}
			continue
		}

		switch actualType.Kind() {
		case reflect.Bool:
			// For checkboxes: present in form data = true, absent = false
			formValues, exists := values[config.Name]
			boolVal := exists && len(formValues) > 0

			if isPointer {
				fieldVal.Set(reflect.ValueOf(&boolVal))
			} else {
				fieldVal.SetBool(boolVal)
			}

		default:
			// For non-boolean fields, skip if no form data
			formValues, exists := values[config.Name]
			if !exists || len(formValues) == 0 {
				continue
			}

			formValue := formValues[0]

			switch actualType.Kind() {
			case reflect.String:
				if isPointer {
					fieldVal.Set(reflect.ValueOf(&formValue))
				} else {
					fieldVal.SetString(formValue)
				}

			case reflect.Int, reflect.Int64:
				intVal, err := strconv.ParseInt(formValue, 10, 64)
				if err != nil {
					return fmt.Errorf("vee: cannot parse '%s' as integer for field '%s': %w", formValue, config.Name, err)
				}

				if isPointer {
					if actualType.Kind() == reflect.Int {
						intPtr := int(intVal)
						fieldVal.Set(reflect.ValueOf(&intPtr))
					} else {
						fieldVal.Set(reflect.ValueOf(&intVal))
					}
				} else {
					fieldVal.SetInt(intVal)
				}

			case reflect.Float64:
				floatVal, err := strconv.ParseFloat(formValue, 64)
				if err != nil {
					return fmt.Errorf("vee: cannot parse '%s' as float for field '%s': %w", formValue, config.Name, err)
				}

				if isPointer {
					fieldVal.Set(reflect.ValueOf(&floatVal))
				} else {
					fieldVal.SetFloat(floatVal)
				}

			}
		}
	}

	return nil
}

// bindMultiValueField binds form data to a Chosen field
func bindMultiValueField(values map[string][]string, fieldVal reflect.Value, pair ChoicesChosenPair, config FieldConfig) error {
	formValues, exists := values[config.Name]
	if !exists || len(formValues) == 0 {
		return nil // No form data, leave field unchanged
	}

	if pair.IsMultiSelect {
		// Multi-select: bind []int
		var indices []int
		for _, formValue := range formValues {
			index, err := strconv.Atoi(formValue)
			if err != nil {
				return fmt.Errorf("vee: invalid index '%s' for multi-select field '%s'", formValue, config.Name)
			}
			// Validate index is in range
			if index < 0 || index >= pair.ChoicesValue.Len() {
				return fmt.Errorf("vee: index %d out of range for %d choices in field '%s'", index, pair.ChoicesValue.Len(), config.Name)
			}
			indices = append(indices, index)
		}

		// Set the slice
		sliceVal := reflect.MakeSlice(fieldVal.Type(), len(indices), len(indices))
		for i, index := range indices {
			sliceVal.Index(i).SetInt(int64(index))
		}
		fieldVal.Set(sliceVal)
	} else {
		// Single select: bind int
		index, err := strconv.Atoi(formValues[0])
		if err != nil {
			return fmt.Errorf("vee: invalid index '%s' for single-select field '%s'", formValues[0], config.Name)
		}
		// Validate index is in range
		if index < 0 || index >= pair.ChoicesValue.Len() {
			return fmt.Errorf("vee: index %d out of range for %d choices in field '%s'", index, pair.ChoicesValue.Len(), config.Name)
		}
		fieldVal.SetInt(int64(index))
	}

	return nil
}
