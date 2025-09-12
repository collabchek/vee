package vee

import (
	"strings"
	"testing"
	"time"
)

// Helper function to simulate form submission by parsing HTML form and extracting field values
func parseFormHTML(html string) (map[string][]string, error) {
	formData := make(map[string][]string)

	// Simple HTML parsing to extract input field values
	lines := strings.Split(html, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "<input") {
			// Extract name and value attributes
			var name, value string
			var isChecked bool

			// Parse name attribute
			if nameStart := strings.Index(line, `name="`); nameStart != -1 {
				nameStart += 6 // len(`name="`)
				nameEnd := strings.Index(line[nameStart:], `"`)
				if nameEnd != -1 {
					name = line[nameStart : nameStart+nameEnd]
				}
			}

			// Parse value attribute
			if valueStart := strings.Index(line, `value="`); valueStart != -1 {
				valueStart += 7 // len(`value="`)
				valueEnd := strings.Index(line[valueStart:], `"`)
				if valueEnd != -1 {
					value = line[valueStart : valueStart+valueEnd]
				}
			}

			// Check if checkbox is checked
			isChecked = strings.Contains(line, " checked")

			// Add to form data if we have a name
			if name != "" {
				// For checkboxes, only add if checked
				if strings.Contains(line, `type="checkbox"`) {
					if isChecked {
						formData[name] = []string{value}
					}
					// If not checked, don't add to form data (simulates unchecked checkbox behavior)
				} else if value != "" {
					// For other input types, add the value
					formData[name] = []string{value}
				}
			}
		}
	}

	return formData, nil
}

func TestRoundTripBasicTypes(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		verify func(t *testing.T, original, result any)
	}{
		{
			name: "string field",
			input: struct {
				Name string
			}{Name: "John Doe"},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Name string })
				res := result.(*struct{ Name string })
				if orig.Name != res.Name {
					t.Errorf("Expected Name=%q, got Name=%q", orig.Name, res.Name)
				}
			},
		},
		{
			name: "int field",
			input: struct {
				Age int
			}{Age: 25},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Age int })
				res := result.(*struct{ Age int })
				if orig.Age != res.Age {
					t.Errorf("Expected Age=%d, got Age=%d", orig.Age, res.Age)
				}
			},
		},
		{
			name: "int64 field",
			input: struct {
				ID int64
			}{ID: 12345678901234},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ ID int64 })
				res := result.(*struct{ ID int64 })
				if orig.ID != res.ID {
					t.Errorf("Expected ID=%d, got ID=%d", orig.ID, res.ID)
				}
			},
		},
		{
			name: "float64 field",
			input: struct {
				Price float64
			}{Price: 19.99},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Price float64 })
				res := result.(*struct{ Price float64 })
				if orig.Price != res.Price {
					t.Errorf("Expected Price=%f, got Price=%f", orig.Price, res.Price)
				}
			},
		},
		{
			name: "bool field (true)",
			input: struct {
				Active bool
			}{Active: true},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Active bool })
				res := result.(*struct{ Active bool })
				if orig.Active != res.Active {
					t.Errorf("Expected Active=%t, got Active=%t", orig.Active, res.Active)
				}
			},
		},
		{
			name: "bool field (false)",
			input: struct {
				Active bool
			}{Active: false},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Active bool })
				res := result.(*struct{ Active bool })
				if orig.Active != res.Active {
					t.Errorf("Expected Active=%t, got Active=%t", orig.Active, res.Active)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Render struct to HTML
			html, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			// Step 2: Parse HTML form to simulate form submission
			formData, err := parseFormHTML(html)
			if err != nil {
				t.Fatalf("parseFormHTML() error = %v", err)
			}

			// Step 3: Create new struct instance and bind form data
			result := tt.input // Get the same type
			// Use reflection to create a pointer to a new zero instance
			switch result.(type) {
			case struct{ Name string }:
				newStruct := &struct{ Name string }{}
				err = Bind(formData, newStruct)
				if err != nil {
					t.Fatalf("Bind() error = %v", err)
				}
				tt.verify(t, tt.input, newStruct)
			case struct{ Age int }:
				newStruct := &struct{ Age int }{}
				err = Bind(formData, newStruct)
				if err != nil {
					t.Fatalf("Bind() error = %v", err)
				}
				tt.verify(t, tt.input, newStruct)
			case struct{ ID int64 }:
				newStruct := &struct{ ID int64 }{}
				err = Bind(formData, newStruct)
				if err != nil {
					t.Fatalf("Bind() error = %v", err)
				}
				tt.verify(t, tt.input, newStruct)
			case struct{ Price float64 }:
				newStruct := &struct{ Price float64 }{}
				err = Bind(formData, newStruct)
				if err != nil {
					t.Fatalf("Bind() error = %v", err)
				}
				tt.verify(t, tt.input, newStruct)
			case struct{ Active bool }:
				newStruct := &struct{ Active bool }{}
				err = Bind(formData, newStruct)
				if err != nil {
					t.Fatalf("Bind() error = %v", err)
				}
				tt.verify(t, tt.input, newStruct)
			}
		})
	}
}

func TestRoundTripTimeTypes(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		target func() any
		verify func(t *testing.T, original, result any)
	}{
		{
			name: "time.Time datetime-local",
			input: struct {
				Created time.Time
			}{Created: time.Date(2023, 12, 25, 14, 30, 0, 0, time.UTC)},
			target: func() any { return &struct{ Created time.Time }{} },
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Created time.Time })
				res := result.(*struct{ Created time.Time })
				if !orig.Created.Equal(res.Created) {
					t.Errorf("Expected Created=%v, got Created=%v", orig.Created, res.Created)
				}
			},
		},
		{
			name: "time.Time date",
			input: struct {
				Birthday time.Time `vee:"type:'date'"`
			}{Birthday: time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)},
			target: func() any {
				return &struct {
					Birthday time.Time `vee:"type:'date'"`
				}{}
			},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct {
					Birthday time.Time `vee:"type:'date'"`
				})
				res := result.(*struct {
					Birthday time.Time `vee:"type:'date'"`
				})
				if !orig.Birthday.Equal(res.Birthday) {
					t.Errorf("Expected Birthday=%v, got Birthday=%v", orig.Birthday, res.Birthday)
				}
			},
		},
		{
			name: "time.Duration seconds",
			input: struct {
				Timeout time.Duration
			}{Timeout: 30 * time.Second},
			target: func() any { return &struct{ Timeout time.Duration }{} },
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct{ Timeout time.Duration })
				res := result.(*struct{ Timeout time.Duration })
				if orig.Timeout != res.Timeout {
					t.Errorf("Expected Timeout=%v, got Timeout=%v", orig.Timeout, res.Timeout)
				}
			},
		},
		{
			name: "time.Duration minutes",
			input: struct {
				Duration time.Duration `vee:"units:'m'"`
			}{Duration: 2*time.Hour + 30*time.Minute},
			target: func() any {
				return &struct {
					Duration time.Duration `vee:"units:'m'"`
				}{}
			},
			verify: func(t *testing.T, original, result any) {
				orig := original.(struct {
					Duration time.Duration `vee:"units:'m'"`
				})
				res := result.(*struct {
					Duration time.Duration `vee:"units:'m'"`
				})
				if orig.Duration != res.Duration {
					t.Errorf("Expected Duration=%v, got Duration=%v", orig.Duration, res.Duration)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Render struct to HTML
			html, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			// Step 2: Parse HTML form to simulate form submission
			formData, err := parseFormHTML(html)
			if err != nil {
				t.Fatalf("parseFormHTML() error = %v", err)
			}

			// Step 3: Create new struct instance and bind form data
			result := tt.target()
			err = Bind(formData, result)
			if err != nil {
				t.Fatalf("Bind() error = %v", err)
			}

			// Step 4: Verify round-trip fidelity
			tt.verify(t, tt.input, result)
		})
	}
}

func TestRoundTripMixedTypes(t *testing.T) {
	type User struct {
		Name      string        `vee:"required"`
		Age       int           `vee:"min:18,max:120"`
		Email     string        `vee:"$user_email,type:'email'"`
		Active    bool          `vee:"label:'Account Active'"`
		CreatedAt time.Time     `vee:"type:'datetime-local'"`
		Timeout   time.Duration `vee:"units:'s'"`
	}

	original := User{
		Name:      "John Doe",
		Age:       25,
		Email:     "john@example.com",
		Active:    true,
		CreatedAt: time.Date(2023, 12, 1, 10, 0, 0, 0, time.UTC),
		Timeout:   60 * time.Second,
	}

	// Step 1: Render struct to HTML
	html, err := Render(original)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Step 2: Parse HTML form to simulate form submission
	formData, err := parseFormHTML(html)
	if err != nil {
		t.Fatalf("parseFormHTML() error = %v", err)
	}

	// Step 3: Create new struct instance and bind form data
	var result User
	err = Bind(formData, &result)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	// Step 4: Verify round-trip fidelity for all fields
	if original.Name != result.Name {
		t.Errorf("Expected Name=%q, got Name=%q", original.Name, result.Name)
	}
	if original.Age != result.Age {
		t.Errorf("Expected Age=%d, got Age=%d", original.Age, result.Age)
	}
	if original.Email != result.Email {
		t.Errorf("Expected Email=%q, got Email=%q", original.Email, result.Email)
	}
	if original.Active != result.Active {
		t.Errorf("Expected Active=%t, got Active=%t", original.Active, result.Active)
	}
	if !original.CreatedAt.Equal(result.CreatedAt) {
		t.Errorf("Expected CreatedAt=%v, got CreatedAt=%v", original.CreatedAt, result.CreatedAt)
	}
	if original.Timeout != result.Timeout {
		t.Errorf("Expected Timeout=%v, got Timeout=%v", original.Timeout, result.Timeout)
	}
}

func TestRoundTripZeroValues(t *testing.T) {
	type ZeroStruct struct {
		Name      string
		Age       int
		Price     float64
		Active    bool
		CreatedAt time.Time
		Timeout   time.Duration
	}

	// Test with all zero values
	original := ZeroStruct{}

	// Step 1: Render struct to HTML
	html, err := Render(original)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	// Step 2: Parse HTML form to simulate form submission
	formData, err := parseFormHTML(html)
	if err != nil {
		t.Fatalf("parseFormHTML() error = %v", err)
	}

	// Step 3: Create new struct instance and bind form data
	var result ZeroStruct
	err = Bind(formData, &result)
	if err != nil {
		t.Fatalf("Bind() error = %v", err)
	}

	// Step 4: Verify zero values are preserved
	if result.Name != "" {
		t.Errorf("Expected empty Name, got Name=%q", result.Name)
	}
	if result.Age != 0 {
		t.Errorf("Expected Age=0, got Age=%d", result.Age)
	}
	if result.Price != 0 {
		t.Errorf("Expected Price=0, got Price=%f", result.Price)
	}
	if result.Active != false {
		t.Errorf("Expected Active=false, got Active=%t", result.Active)
	}
	if !result.CreatedAt.IsZero() {
		t.Errorf("Expected zero time, got CreatedAt=%v", result.CreatedAt)
	}
	if result.Timeout != 0 {
		t.Errorf("Expected Timeout=0, got Timeout=%v", result.Timeout)
	}
}
