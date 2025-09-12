package vee

import (
	"strings"

	"github.com/iancoleman/strcase"
)

// FieldConfig holds the configuration for a struct field parsed from tags.
type FieldConfig struct {
	Name       string            // HTML form field name
	Skip       bool              // Whether to skip this field
	NoLabel    bool              // Whether to skip label generation
	Hidden     bool              // Whether to render as hidden input
	Attributes map[string]string // HTML attributes (min, max, step, etc.)
}

// parseVeeTag parses a "vee" struct tag and extracts the field name and attributes.
// Fields are processed by default. Supports:
//   - vee:"-" to skip field
//   - vee:"$override_name" to override field name
//   - vee:"" to use auto-derived field name
//   - vee:"min:10,max:100,step:5" for numeric attributes
func parseVeeTag(tag, fieldName string) FieldConfig {
	config := FieldConfig{
		Attributes: make(map[string]string),
	}

	if tag == "" {
		// Default behavior: process all fields with auto-derived name
		config.Name = strcase.ToSnake(fieldName)
		return config
	}

	if tag == "-" {
		config.Skip = true
		return config
	}

	// Default behavior: process all fields with auto-derived name
	config.Name = strcase.ToSnake(fieldName)

	// Split by comma
	parts := strings.Split(tag, ",")

	// Check if first part is name override
	if strings.HasPrefix(parts[0], "$") {
		config.Name = parts[0][1:] // Remove $ prefix
		parts = parts[1:]          // Process remaining parts as attributes
	}
	// Otherwise keep the auto-derived name

	// Process remaining parts as attributes
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if strings.Contains(part, ":") {
			// Key-value attribute (e.g., min:10, type:'email')
			kv := strings.SplitN(part, ":", 2)
			if len(kv) == 2 {
				key := strings.TrimSpace(kv[0])
				value := strings.TrimSpace(kv[1])
				// Strip surrounding single quotes from value
				if len(value) >= 2 && value[0] == '\'' && value[len(value)-1] == '\'' {
					value = value[1 : len(value)-1]
				}
				config.Attributes[key] = value
			}
		} else {
			// Check for special boolean attributes
			if part == "nolabel" {
				config.NoLabel = true
			} else if part == "hidden" {
				config.Hidden = true
			} else {
				// Boolean attribute (e.g., required, readonly, disabled)
				config.Attributes[part] = ""
			}
		}
	}

	return config
}
