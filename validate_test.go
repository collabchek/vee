package vee

import (
	"testing"
)

func TestValidateStruct(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantErr bool
	}{
		{
			name: "valid struct passes validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
		{
			name: "missing required field fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "", // Missing required field
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "invalid email fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "John Doe",
				Email: "invalid-email", // Invalid email
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "age below minimum fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   16, // Below minimum age
			},
			wantErr: true,
		},
		{
			name: "age above maximum fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   150, // Above maximum age
			},
			wantErr: true,
		},
		{
			name: "name too short fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "J", // Too short
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "name too long fails validation",
			input: struct {
				Name  string `validate:"required,min=2,max=50"`
				Email string `validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "This is a very long name that exceeds the maximum length allowed for this field", // Too long
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: true,
		},
		{
			name: "struct with vee and validate tags works together",
			input: struct {
				Name  string `vee:"required" validate:"required,min=2,max=50"`
				Email string `vee:"type:'email',required" validate:"required,email"`
				Age   int    `validate:"required,gte=18,lte=120"`
			}{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   25,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateVar(t *testing.T) {
	tests := []struct {
		name    string
		field   any
		tag     string
		wantErr bool
	}{
		{
			name:    "valid email passes validation",
			field:   "john@example.com",
			tag:     "email",
			wantErr: false,
		},
		{
			name:    "invalid email fails validation",
			field:   "invalid-email",
			tag:     "email",
			wantErr: true,
		},
		{
			name:    "required field with value passes validation",
			field:   "John",
			tag:     "required",
			wantErr: false,
		},
		{
			name:    "required field without value fails validation",
			field:   "",
			tag:     "required",
			wantErr: true,
		},
		{
			name:    "number within range passes validation",
			field:   25,
			tag:     "gte=18,lte=120",
			wantErr: false,
		},
		{
			name:    "number below range fails validation",
			field:   16,
			tag:     "gte=18,lte=120",
			wantErr: true,
		},
		{
			name:    "number above range fails validation",
			field:   150,
			tag:     "gte=18,lte=120",
			wantErr: true,
		},
		{
			name:    "string within length range passes validation",
			field:   "John Doe",
			tag:     "min=2,max=50",
			wantErr: false,
		},
		{
			name:    "string too short fails validation",
			field:   "J",
			tag:     "min=2,max=50",
			wantErr: true,
		},
		{
			name:    "string too long fails validation",
			field:   "This is a very long string that exceeds the maximum length allowed",
			tag:     "min=2,max=50",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateVar(tt.field, tt.tag)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateVar() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidationIntegrationWithVEE(t *testing.T) {
	// Test that VEE rendering/binding works alongside validation
	type User struct {
		Name  string `vee:"required" validate:"required,min=2,max=50"`
		Email string `vee:"type:'email',required" validate:"required,email"`
		Age   int    `validate:"required,gte=18,lte=120"`
	}

	// Test 1: Valid data should render, bind, and validate correctly
	t.Run("valid data works end-to-end", func(t *testing.T) {
		original := User{
			Name:  "John Doe",
			Email: "john@example.com",
			Age:   25,
		}

		// Test validation
		if err := Validate(original); err != nil {
			t.Errorf("Validation failed for valid data: %v", err)
		}

		// Test rendering
		html, err := Render(original)
		if err != nil {
			t.Errorf("Render failed: %v", err)
		}
		if html == "" {
			t.Error("Render returned empty HTML")
		}

		// Test binding (simulate form data)
		formData := map[string][]string{
			"name":  {"John Doe"},
			"email": {"john@example.com"},
			"age":   {"25"},
		}

		var bound User
		if err := Bind(formData, &bound); err != nil {
			t.Errorf("Bind failed: %v", err)
		}

		// Validate bound data
		if err := Validate(bound); err != nil {
			t.Errorf("Validation failed for bound data: %v", err)
		}

		// Verify bound data matches original
		if bound.Name != original.Name || bound.Email != original.Email || bound.Age != original.Age {
			t.Errorf("Bound data doesn't match original. Got %+v, want %+v", bound, original)
		}
	})

	// Test 2: Invalid data should fail validation
	t.Run("invalid data fails validation", func(t *testing.T) {
		invalid := User{
			Name:  "J", // Too short
			Email: "invalid-email",
			Age:   16, // Too young
		}

		if err := Validate(invalid); err == nil {
			t.Error("Expected validation to fail for invalid data, but it passed")
		}
	})
}
