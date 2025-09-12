package vee

import (
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode"
)

// Render generates HTML form fields from a Go struct.
// Accepts optional RenderOptions to customize form rendering.
func Render(v any, opts ...RenderOption) (string, error) {
	options := ConsolidateOptions(opts...)
	// if len(opts) > 0 && opts[0] != nil {
	// 	options = opts[0]
	// } else {
	// 	options = &RenderOption{}
	// }
	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Handle pointer to struct
	if typ.Kind() == reflect.Ptr {
		val = val.Elem()
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return "", fmt.Errorf("vee: expected struct, got %v", typ.Kind())
	}

	// First pass: validate hidden field restrictions before other validations
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		veeTag := field.Tag.Get("vee")
		config := parseVeeTag(veeTag, field.Name)

		// Skip if requested
		if config.Skip {
			continue
		}

		// Validate hidden field restrictions
		if config.Hidden {
			// Check if this is a pointer type
			if field.Type.Kind() == reflect.Ptr {
				return "", fmt.Errorf("vee: hidden attribute not supported for pointer type '%s'", field.Name)
			}

			// Check if this is a multi-value field (Choices or Chosen)
			if strings.HasSuffix(field.Name, "Choices") || strings.HasSuffix(field.Name, "Chosen") {
				return "", fmt.Errorf("vee: hidden attribute not supported for multi-value field '%s'", field.Name)
			}

			// Check if field type is a slice/array
			if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
				return "", fmt.Errorf("vee: hidden attribute not supported for slice/array type '%s'", field.Name)
			}
		}
	}

	// Validate Choices/Chosen pairs
	choicesChosenPairs, err := validateChoicesChosen(typ, val)
	if err != nil {
		return "", err
	}

	var html strings.Builder

	// Always wrap in form tag
	html.WriteString("<form")
	if options.FormID != "" {
		html.WriteString(fmt.Sprintf(` id="%s"`, escapeHTML(options.FormID)))
	}
	if options.FormCSS != "" {
		html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(options.FormCSS)))
	}
	// Skip method and action if we're going to submit the form via Javascript
	if options.FormAction != "script" {
		method := options.FormMethod
		if method == "" {
			method = "POST"
		}
		html.WriteString(fmt.Sprintf(` method="%s"`, method))
		if options.FormAction != "" {
			html.WriteString(fmt.Sprintf(` action="%s"`, escapeHTML(options.FormAction)))
		}
	}
	html.WriteString(">\n")

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

		// Build CSS classes
		var cssClass string
		cssTag := field.Tag.Get("css")
		if cssTag != "" {
			cssClass = cssTag
		} else if options.DefaultInputCSS != "" {
			cssClass = options.DefaultInputCSS
		}

		// Skip Choices fields (they're not rendered, only used for Chosen fields)
		if strings.HasSuffix(field.Name, "Choices") {
			continue
		}

		// Handle Chosen fields specially
		if strings.HasSuffix(field.Name, "Chosen") {
			baseName := strings.TrimSuffix(field.Name, "Chosen")
			if pair, exists := choicesChosenPairs[baseName]; exists {
				err := renderMultiValueField(&html, pair, config, cssClass)
				if err != nil {
					return "", err
				}
				continue
			}
		}

		// Handle hidden fields early - they override normal rendering
		if config.Hidden {
			err := renderHiddenField(&html, field, fieldVal, config, cssClass)
			if err != nil {
				return "", err
			}
			continue
		}

		// Handle pointer types
		actualType := field.Type
		actualVal := fieldVal
		isPointer := false

		if actualType.Kind() == reflect.Ptr {
			isPointer = true
			actualType = actualType.Elem()

			// If pointer is nil, we'll use zero values for rendering
			if fieldVal.IsNil() {
				actualVal = reflect.Zero(actualType)
			} else {
				actualVal = fieldVal.Elem()
			}
		}

		// Check for specific types first (before generic kind matching)
		if actualType == reflect.TypeOf(time.Time{}) {
			timeVal := actualVal.Interface().(time.Time)

			// Render label first
			renderLabel(&html, config, field.Name)

			// Determine input type (default to datetime-local)
			inputType := "datetime-local"
			if typeAttr, ok := config.Attributes["type"]; ok {
				switch typeAttr {
				case "date", "datetime-local", "time":
					inputType = typeAttr
				}
			}

			html.WriteString(fmt.Sprintf(`<input type="%s"`, inputType))
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))

			// Format the value based on input type
			var value string
			if !isPointer || !fieldVal.IsNil() {
				if !timeVal.IsZero() {
					switch inputType {
					case "date":
						value = timeVal.Format("2006-01-02")
					case "time":
						value = timeVal.Format("15:04")
					case "datetime-local":
						value = timeVal.Format("2006-01-02T15:04")
					}
					html.WriteString(fmt.Sprintf(` value="%s"`, escapeHTML(value)))
				}
			}

			// Add min/max attributes
			if min, ok := config.Attributes["min"]; ok {
				html.WriteString(fmt.Sprintf(` min="%s"`, escapeHTML(min)))
			}
			if max, ok := config.Attributes["max"]; ok {
				html.WriteString(fmt.Sprintf(` max="%s"`, escapeHTML(max)))
			}

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")
			continue
		}

		if actualType == reflect.TypeOf(time.Duration(0)) {
			durationVal := actualVal.Interface().(time.Duration)

			// Render label first
			renderLabel(&html, config, field.Name)

			// Get units (default to seconds)
			units := "s"
			if unitsAttr, ok := config.Attributes["units"]; ok {
				switch unitsAttr {
				case "ms", "s", "m", "h":
					units = unitsAttr
				}
			}

			html.WriteString(`<input type="number"`)
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))

			// Convert duration to specified units and render value
			if (!isPointer || !fieldVal.IsNil()) && durationVal != 0 {
				var value float64
				switch units {
				case "ms":
					value = float64(durationVal / time.Millisecond)
				case "s":
					value = float64(durationVal / time.Second)
				case "m":
					value = float64(durationVal / time.Minute)
				case "h":
					value = float64(durationVal / time.Hour)
				}
				html.WriteString(fmt.Sprintf(` value="%g"`, value))
			}

			// Add numeric attributes
			if min, ok := config.Attributes["min"]; ok {
				html.WriteString(fmt.Sprintf(` min="%s"`, escapeHTML(min)))
			}
			if max, ok := config.Attributes["max"]; ok {
				html.WriteString(fmt.Sprintf(` max="%s"`, escapeHTML(max)))
			}
			if step, ok := config.Attributes["step"]; ok {
				html.WriteString(fmt.Sprintf(` step="%s"`, escapeHTML(step)))
			}

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")
			continue
		}

		// Render field based on type
		switch actualType.Kind() {
		case reflect.String:
			value := actualVal.String()

			// Render label first
			renderLabel(&html, config, field.Name)

			// Determine input type (default to text, but allow override)
			inputType := "text"
			if typeAttr, ok := config.Attributes["type"]; ok {
				switch typeAttr {
				case "email", "password", "tel", "url":
					inputType = typeAttr
				}
			}

			html.WriteString(fmt.Sprintf(`<input type="%s"`, inputType))
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
			html.WriteString(fmt.Sprintf(` value="%s"`, escapeHTML(value)))

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")

		case reflect.Int, reflect.Int64:
			value := actualVal.Int()

			// Render label first
			renderLabel(&html, config, field.Name)

			html.WriteString(`<input type="number"`)
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
			html.WriteString(fmt.Sprintf(` value="%d"`, value))

			// Add numeric attributes
			if min, ok := config.Attributes["min"]; ok {
				html.WriteString(fmt.Sprintf(` min="%s"`, escapeHTML(min)))
			}
			if max, ok := config.Attributes["max"]; ok {
				html.WriteString(fmt.Sprintf(` max="%s"`, escapeHTML(max)))
			}
			if step, ok := config.Attributes["step"]; ok {
				html.WriteString(fmt.Sprintf(` step="%s"`, escapeHTML(step)))
			}

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")

		case reflect.Float64:
			value := actualVal.Float()

			// Render label first
			renderLabel(&html, config, field.Name)

			html.WriteString(`<input type="number"`)
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
			html.WriteString(fmt.Sprintf(` value="%g"`, value))

			// Add numeric attributes (step defaults to "any" for floats if not specified)
			if min, ok := config.Attributes["min"]; ok {
				html.WriteString(fmt.Sprintf(` min="%s"`, escapeHTML(min)))
			}
			if max, ok := config.Attributes["max"]; ok {
				html.WriteString(fmt.Sprintf(` max="%s"`, escapeHTML(max)))
			}
			if step, ok := config.Attributes["step"]; ok {
				html.WriteString(fmt.Sprintf(` step="%s"`, escapeHTML(step)))
			} else {
				html.WriteString(` step="any"`) // Default for float64
			}

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")

		case reflect.Bool:
			isChecked := actualVal.Bool()

			// Render label first
			renderLabel(&html, config, field.Name)

			html.WriteString(`<input type="checkbox"`)
			html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
			html.WriteString(` value="true"`)
			if isChecked {
				html.WriteString(` checked`)
			}

			// Add CSS class
			if cssClass != "" {
				html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
			}

			// Add universal attributes
			addUniversalAttributes(&html, config)

			html.WriteString(">\n")

		default:
			// Skip unsupported types
			continue
		}
	}

	// Always close form tag
	html.WriteString("</form>\n")

	return html.String(), nil
}

// validateChoicesChosen validates Choices/Chosen field pairs and returns information about them
func validateChoicesChosen(typ reflect.Type, val reflect.Value) (map[string]ChoicesChosenPair, error) {
	pairs := make(map[string]ChoicesChosenPair)
	choicesFields := make(map[string]reflect.StructField)
	chosenFields := make(map[string]reflect.StructField)

	// First pass: identify Choices and Chosen fields
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		if strings.HasSuffix(fieldName, "Choices") {
			baseName := strings.TrimSuffix(fieldName, "Choices")
			choicesFields[baseName] = field
		} else if strings.HasSuffix(fieldName, "Chosen") {
			baseName := strings.TrimSuffix(fieldName, "Chosen")
			chosenFields[baseName] = field
		}
	}

	// Validate that Choices and Chosen fields come in pairs
	for baseName, choicesField := range choicesFields {
		chosenField, hasChosen := chosenFields[baseName]
		if !hasChosen {
			return nil, fmt.Errorf("vee: field '%s' requires corresponding '%sChosen' field", choicesField.Name, baseName)
		}

		// Validate Choices field type (must be slice)
		if choicesField.Type.Kind() != reflect.Slice {
			return nil, fmt.Errorf("vee: field '%s' must be a slice type, got %s", choicesField.Name, choicesField.Type.Kind())
		}

		// Validate Chosen field type (must be int or []int)
		chosenKind := chosenField.Type.Kind()
		isMultiSelect := false
		if chosenKind == reflect.Slice {
			if chosenField.Type.Elem().Kind() != reflect.Int {
				return nil, fmt.Errorf("vee: field '%s' must be int or []int, got %s", chosenField.Name, chosenField.Type)
			}
			isMultiSelect = true
		} else if chosenKind != reflect.Int {
			return nil, fmt.Errorf("vee: field '%s' must be int or []int, got %s", chosenField.Name, chosenField.Type)
		}

		// Get field values
		choicesFieldVal := val.FieldByName(choicesField.Name)
		chosenFieldVal := val.FieldByName(chosenField.Name)

		// Validate choices are not empty
		if choicesFieldVal.Len() == 0 {
			return nil, fmt.Errorf("vee: field '%s' cannot be empty", choicesField.Name)
		}

		// Validate chosen indices are in range
		if isMultiSelect {
			for i := 0; i < chosenFieldVal.Len(); i++ {
				index := int(chosenFieldVal.Index(i).Int())
				if index < 0 || index >= choicesFieldVal.Len() {
					return nil, fmt.Errorf("vee: field '%s' index %d out of range for %d choices", chosenField.Name, index, choicesFieldVal.Len())
				}
			}
		} else {
			index := int(chosenFieldVal.Int())
			if index < 0 || index >= choicesFieldVal.Len() {
				return nil, fmt.Errorf("vee: field '%s' index %d out of range for %d choices", chosenField.Name, index, choicesFieldVal.Len())
			}
		}

		pairs[baseName] = ChoicesChosenPair{
			ChoicesField:  choicesField,
			ChosenField:   chosenField,
			ChoicesValue:  choicesFieldVal,
			ChosenValue:   chosenFieldVal,
			IsMultiSelect: isMultiSelect,
		}
	}

	// Check for orphaned Chosen fields
	for baseName, chosenField := range chosenFields {
		_, hasChoices := choicesFields[baseName]
		if hasChoices {
			continue
		}
		return nil, fmt.Errorf("vee: field '%s' requires corresponding '%sChoices' field", chosenField.Name, baseName)
	}

	return pairs, nil
}

// ChoicesChosenPair represents a validated pair of Choices and Chosen fields
type ChoicesChosenPair struct {
	ChoicesField  reflect.StructField
	ChosenField   reflect.StructField
	ChoicesValue  reflect.Value
	ChosenValue   reflect.Value
	IsMultiSelect bool
}

// renderMultiValueField renders a Chosen field as select, radio, or checkbox group
func renderMultiValueField(html *strings.Builder, pair ChoicesChosenPair, config FieldConfig, cssClass string) error {
	// Determine the input type from attributes (defaults to select)
	inputType := "select"
	if typeAttr, ok := config.Attributes["type"]; ok {
		switch typeAttr {
		case "select", "radio", "checkbox":
			inputType = typeAttr
		}
	}

	// Get selected indices
	var selectedIndices []int
	if pair.IsMultiSelect {
		for i := 0; i < pair.ChosenValue.Len(); i++ {
			selectedIndices = append(selectedIndices, int(pair.ChosenValue.Index(i).Int()))
		}
	} else {
		selectedIndices = []int{int(pair.ChosenValue.Int())}
	}

	switch inputType {
	case "select":
		return renderSelectField(html, pair, config, cssClass, selectedIndices)
	case "radio":
		if pair.IsMultiSelect {
			return fmt.Errorf("vee: radio buttons cannot be used with multi-select field '%s'", pair.ChosenField.Name)
		}
		return renderRadioField(html, pair, config, cssClass, selectedIndices[0])
	case "checkbox":
		return renderCheckboxField(html, pair, config, cssClass, selectedIndices)
	}

	return nil
}

// renderSelectField renders a select element
func renderSelectField(html *strings.Builder, pair ChoicesChosenPair, config FieldConfig, cssClass string, selectedIndices []int) error {
	// Render label first
	renderLabel(html, config, pair.ChosenField.Name)

	html.WriteString("<select")
	html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))

	if pair.IsMultiSelect {
		html.WriteString(" multiple")
	}

	if cssClass != "" {
		html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
	}

	// Add universal attributes
	addUniversalAttributes(html, config)

	html.WriteString(">\n")

	// Add options
	for i := 0; i < pair.ChoicesValue.Len(); i++ {
		choice := pair.ChoicesValue.Index(i).String()
		html.WriteString(fmt.Sprintf(`<option value="%d"`, i))

		// Check if this option is selected
		for _, selectedIndex := range selectedIndices {
			if i == selectedIndex {
				html.WriteString(" selected")
				break
			}
		}

		html.WriteString(fmt.Sprintf(">%s</option>\n", escapeHTML(choice)))
	}

	html.WriteString("</select>\n")
	return nil
}

// renderRadioField renders a radio button group
func renderRadioField(html *strings.Builder, pair ChoicesChosenPair, config FieldConfig, cssClass string, selectedIndex int) error {
	// Render group label first (if not disabled)
	if !config.NoLabel {
		labelText := generateLabel(config, pair.ChosenField.Name)
		html.WriteString(fmt.Sprintf(`<fieldset><legend>%s</legend>`, escapeHTML(labelText)))
		html.WriteString("\n")
	}

	for i := 0; i < pair.ChoicesValue.Len(); i++ {
		choice := pair.ChoicesValue.Index(i).String()
		radioID := fmt.Sprintf("%s_%d", config.Name, i)

		html.WriteString(`<input type="radio"`)
		html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
		html.WriteString(fmt.Sprintf(` value="%d"`, i))

		if i == selectedIndex {
			html.WriteString(" checked")
		}

		if cssClass != "" {
			html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
		}

		html.WriteString(fmt.Sprintf(` id="%s"`, radioID))

		// Add other universal attributes (except id since we set it specifically)
		if placeholder, ok := config.Attributes["placeholder"]; ok {
			html.WriteString(fmt.Sprintf(` placeholder="%s"`, escapeHTML(placeholder)))
		}
		if _, ok := config.Attributes["required"]; ok {
			html.WriteString(` required`)
		}
		if _, ok := config.Attributes["readonly"]; ok {
			html.WriteString(` readonly`)
		}
		if _, ok := config.Attributes["disabled"]; ok {
			html.WriteString(` disabled`)
		}

		html.WriteString(fmt.Sprintf(`><label for="%s">%s</label>`, radioID, escapeHTML(choice)))
		html.WriteString("\n")
	}

	// Close fieldset if we opened one
	if !config.NoLabel {
		html.WriteString("</fieldset>\n")
	}

	return nil
}

// renderCheckboxField renders a checkbox group
func renderCheckboxField(html *strings.Builder, pair ChoicesChosenPair, config FieldConfig, cssClass string, selectedIndices []int) error {
	// Render group label first (if not disabled)
	if !config.NoLabel {
		labelText := generateLabel(config, pair.ChosenField.Name)
		html.WriteString(fmt.Sprintf(`<fieldset><legend>%s</legend>`, escapeHTML(labelText)))
		html.WriteString("\n")
	}

	for i := 0; i < pair.ChoicesValue.Len(); i++ {
		choice := pair.ChoicesValue.Index(i).String()
		checkboxID := fmt.Sprintf("%s_%d", config.Name, i)

		html.WriteString(`<input type="checkbox"`)
		html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))
		html.WriteString(fmt.Sprintf(` value="%d"`, i))

		// Check if this checkbox is selected
		for _, selectedIndex := range selectedIndices {
			if i == selectedIndex {
				html.WriteString(" checked")
				break
			}
		}

		if cssClass != "" {
			html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
		}

		html.WriteString(fmt.Sprintf(` id="%s"`, checkboxID))

		// Add other universal attributes (except id since we set it specifically)
		if placeholder, ok := config.Attributes["placeholder"]; ok {
			html.WriteString(fmt.Sprintf(` placeholder="%s"`, escapeHTML(placeholder)))
		}
		if _, ok := config.Attributes["required"]; ok {
			html.WriteString(` required`)
		}
		if _, ok := config.Attributes["readonly"]; ok {
			html.WriteString(` readonly`)
		}
		if _, ok := config.Attributes["disabled"]; ok {
			html.WriteString(` disabled`)
		}

		html.WriteString(fmt.Sprintf(`><label for="%s">%s</label>`, checkboxID, escapeHTML(choice)))
		html.WriteString("\n")
	}

	// Close fieldset if we opened one
	if !config.NoLabel {
		html.WriteString("</fieldset>\n")
	}

	return nil
}

// addUniversalAttributes adds universal HTML attributes (required, readonly, disabled, placeholder, id)
func addUniversalAttributes(html *strings.Builder, config FieldConfig) {
	// Add id attribute (custom or default to field name)
	if id, ok := config.Attributes["id"]; ok {
		html.WriteString(fmt.Sprintf(` id="%s"`, escapeHTML(id)))
	} else {
		html.WriteString(fmt.Sprintf(` id="%s"`, escapeHTML(config.Name)))
	}

	// Add placeholder attribute
	if placeholder, ok := config.Attributes["placeholder"]; ok {
		html.WriteString(fmt.Sprintf(` placeholder="%s"`, escapeHTML(placeholder)))
	}

	// Add boolean attributes (required, readonly, disabled)
	if _, ok := config.Attributes["required"]; ok {
		html.WriteString(` required`)
	}
	if _, ok := config.Attributes["readonly"]; ok {
		html.WriteString(` readonly`)
	}
	if _, ok := config.Attributes["disabled"]; ok {
		html.WriteString(` disabled`)
	}
}

// escapeHTML escapes HTML characters in attribute values
func escapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// generateLabel creates a human-readable label for a field
func generateLabel(config FieldConfig, fieldName string) string {
	// Check if custom label is provided
	if label, ok := config.Attributes["label"]; ok {
		return label
	}

	// Convert field name to human-readable format
	// AccountName -> Account Name, EmailAddress -> Email Address
	return fieldNameToLabel(fieldName)
}

// fieldNameToLabel converts a field name to a human-readable label
// This properly handles international characters (Ä, É, Α, А, etc.)
func fieldNameToLabel(fieldName string) string {
	var result strings.Builder
	for i, r := range fieldName {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune(' ')
		}
		result.WriteRune(r)
	}
	return result.String()
}

// renderLabel generates a <label> element for a field if not disabled
func renderLabel(html *strings.Builder, config FieldConfig, fieldName string) {
	if config.NoLabel {
		return
	}

	labelText := generateLabel(config, fieldName)
	fieldID := config.Name
	if customID, ok := config.Attributes["id"]; ok {
		fieldID = customID
	}

	html.WriteString(fmt.Sprintf(`<label for="%s">%s</label>`, escapeHTML(fieldID), escapeHTML(labelText)))
	html.WriteString("\n")
}

// renderHiddenField renders a hidden input field for any supported field type
func renderHiddenField(html *strings.Builder, field reflect.StructField, fieldVal reflect.Value, config FieldConfig, cssClass string) error {
	// Hidden fields never render labels
	html.WriteString(`<input type="hidden"`)
	html.WriteString(fmt.Sprintf(` name="%s"`, config.Name))

	// Handle different field types and extract their values
	actualType := field.Type
	actualVal := fieldVal

	// Check for specific types first (before generic kind matching)
	if actualType == reflect.TypeOf(time.Time{}) {
		timeVal := actualVal.Interface().(time.Time)
		if !timeVal.IsZero() {
			// Use ISO format for hidden time fields
			value := timeVal.Format("2006-01-02T15:04:05Z07:00")
			html.WriteString(fmt.Sprintf(` value="%s"`, escapeHTML(value)))
		}
	} else if actualType == reflect.TypeOf(time.Duration(0)) {
		durationVal := actualVal.Interface().(time.Duration)
		if durationVal != 0 {
			// Store duration as nanoseconds for hidden fields
			html.WriteString(fmt.Sprintf(` value="%d"`, int64(durationVal)))
		}
	} else {
		// Handle by kind for basic types
		switch actualType.Kind() {
		case reflect.String:
			value := actualVal.String()
			html.WriteString(fmt.Sprintf(` value="%s"`, escapeHTML(value)))

		case reflect.Int, reflect.Int64:
			value := actualVal.Int()
			html.WriteString(fmt.Sprintf(` value="%d"`, value))

		case reflect.Float64:
			value := actualVal.Float()
			html.WriteString(fmt.Sprintf(` value="%g"`, value))

		case reflect.Bool:
			isTrue := actualVal.Bool()
			if isTrue {
				html.WriteString(` value="true"`)
			} else {
				html.WriteString(` value="false"`)
			}

		default:
			return fmt.Errorf("vee: unsupported type for hidden field '%s': %s", field.Name, actualType.Kind())
		}
	}

	// Add CSS class if provided
	if cssClass != "" {
		html.WriteString(fmt.Sprintf(` class="%s"`, escapeHTML(cssClass)))
	}

	// Add universal attributes (id is still useful, others may not be but we'll include them)
	addUniversalAttributes(html, config)

	html.WriteString(">\n")
	return nil
}
