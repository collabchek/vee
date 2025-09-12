package vee

import (
	"testing"
	"time"
)

func TestPointerRoundTrip(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		formData map[string][]string
		checkFn  func(t *testing.T, input any)
	}{
		{
			name: "all pointer fields with values",
			input: struct {
				Name    *string
				Age     *int
				Score   *float64
				Active  *bool
				Created *time.Time
				Timeout *time.Duration
			}{
				Name:    stringPtr("John"),
				Age:     intPtr(30),
				Score:   float64Ptr(95.5),
				Active:  boolPtr(true),
				Created: timePtr(time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC)),
				Timeout: durationPtr(5 * time.Second),
			},
			formData: map[string][]string{
				"name":    {"Alice"},
				"age":     {"25"},
				"score":   {"87.3"},
				"active":  {"true"},
				"created": {"2024-01-15T10:45"},
				"timeout": {"3"},
			},
			checkFn: func(t *testing.T, input any) {
				s := input.(*struct {
					Name    *string
					Age     *int
					Score   *float64
					Active  *bool
					Created *time.Time
					Timeout *time.Duration
				})

				if s.Name == nil || *s.Name != "Alice" {
					t.Errorf("Expected Name='Alice', got %v", s.Name)
				}
				if s.Age == nil || *s.Age != 25 {
					t.Errorf("Expected Age=25, got %v", s.Age)
				}
				if s.Score == nil || *s.Score != 87.3 {
					t.Errorf("Expected Score=87.3, got %v", s.Score)
				}
				if s.Active == nil || *s.Active != true {
					t.Errorf("Expected Active=true, got %v", s.Active)
				}
				if s.Created == nil {
					t.Errorf("Expected Created to be non-nil")
				} else {
					expected := time.Date(2024, 1, 15, 10, 45, 0, 0, time.UTC)
					if !s.Created.Equal(expected) {
						t.Errorf("Expected Created=%v, got %v", expected, *s.Created)
					}
				}
				if s.Timeout == nil || *s.Timeout != 3*time.Second {
					t.Errorf("Expected Timeout=3s, got %v", s.Timeout)
				}
			},
		},
		{
			name: "nil pointer fields remain nil when no form data",
			input: struct {
				Name    *string
				Age     *int
				Score   *float64
				Active  *bool
				Created *time.Time
				Timeout *time.Duration
			}{},
			formData: map[string][]string{
				"other_field": {"value"},
			},
			checkFn: func(t *testing.T, input any) {
				s := input.(*struct {
					Name    *string
					Age     *int
					Score   *float64
					Active  *bool
					Created *time.Time
					Timeout *time.Duration
				})

				if s.Name != nil {
					t.Errorf("Expected Name to be nil, got %v", *s.Name)
				}
				if s.Age != nil {
					t.Errorf("Expected Age to be nil, got %v", *s.Age)
				}
				if s.Score != nil {
					t.Errorf("Expected Score to be nil, got %v", *s.Score)
				}
				if s.Active == nil || *s.Active != false {
					t.Errorf("Expected Active to be false (checkbox unchecked), got %v", s.Active)
				}
				if s.Created != nil {
					t.Errorf("Expected Created to be nil, got %v", *s.Created)
				}
				if s.Timeout != nil {
					t.Errorf("Expected Timeout to be nil, got %v", *s.Timeout)
				}
			},
		},
		{
			name: "mixed pointer and non-pointer with partial form data",
			input: struct {
				Name     string
				Email    *string
				Age      int
				Score    *float64
				IsActive bool
				IsPro    *bool
			}{
				Name:     "Original",
				Email:    stringPtr("original@example.com"),
				Age:      20,
				Score:    float64Ptr(50.0),
				IsActive: false,
				IsPro:    boolPtr(false),
			},
			formData: map[string][]string{
				"name": {"Updated"},
				// email not provided - should remain as original
				"age":       {"35"},
				"score":     {"75.5"},
				"is_active": {"true"},
				// is_pro not provided - checkbox should become nil? or false?
			},
			checkFn: func(t *testing.T, input any) {
				s := input.(*struct {
					Name     string
					Email    *string
					Age      int
					Score    *float64
					IsActive bool
					IsPro    *bool
				})

				if s.Name != "Updated" {
					t.Errorf("Expected Name='Updated', got %s", s.Name)
				}
				if s.Email == nil || *s.Email != "original@example.com" {
					t.Errorf("Expected Email to remain original, got %v", s.Email)
				}
				if s.Age != 35 {
					t.Errorf("Expected Age=35, got %d", s.Age)
				}
				if s.Score == nil || *s.Score != 75.5 {
					t.Errorf("Expected Score=75.5, got %v", s.Score)
				}
				if !s.IsActive {
					t.Errorf("Expected IsActive=true, got %t", s.IsActive)
				}
				if s.IsPro == nil || *s.IsPro != false {
					t.Errorf("Expected IsPro=false (checkbox unchecked), got %v", s.IsPro)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test render -> bind round trip
			rendered, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			// Verify rendering produces HTML (basic check)
			if rendered == "" {
				t.Fatalf("Render() produced empty HTML")
			}

			// Create a target for binding, starting with the original values
			var targetType any
			switch v := tt.input.(type) {
			case struct {
				Name    *string
				Age     *int
				Score   *float64
				Active  *bool
				Created *time.Time
				Timeout *time.Duration
			}:
				targetType = &struct {
					Name    *string
					Age     *int
					Score   *float64
					Active  *bool
					Created *time.Time
					Timeout *time.Duration
				}{
					Name:    v.Name,
					Age:     v.Age,
					Score:   v.Score,
					Active:  v.Active,
					Created: v.Created,
					Timeout: v.Timeout,
				}
			case struct {
				Name     string
				Email    *string
				Age      int
				Score    *float64
				IsActive bool
				IsPro    *bool
			}:
				targetType = &struct {
					Name     string
					Email    *string
					Age      int
					Score    *float64
					IsActive bool
					IsPro    *bool
				}{
					Name:     v.Name,
					Email:    v.Email,
					Age:      v.Age,
					Score:    v.Score,
					IsActive: v.IsActive,
					IsPro:    v.IsPro,
				}
			default:
				t.Fatalf("Unknown struct type: %T", v)
			}

			// Bind form data
			err = Bind(tt.formData, targetType)
			if err != nil {
				t.Fatalf("Bind() error = %v", err)
			}

			// Check the results
			tt.checkFn(t, targetType)
		})
	}
}
