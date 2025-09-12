package vee

import (
	"strings"
	"testing"
	"time"
)

func TestHiddenFieldRendering(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected []string
		notExpected []string
	}{
		{
			name: "hidden string field",
			input: struct {
				Name   string `vee:"hidden"`
				Visible string
			}{
				Name:   "secret",
				Visible: "public",
			},
			expected: []string{
				`<input type="hidden" name="name" value="secret"`,
				`<label for="visible">Visible</label>`,
				`<input type="text" name="visible" value="public"`,
			},
			notExpected: []string{
				`<label for="name">Name</label>`, // No label for hidden field
			},
		},
		{
			name: "hidden vs nolabel distinction",
			input: struct {
				Hidden   string `vee:"hidden"`
				NoLabel  string `vee:"nolabel"`
				Normal   string
			}{
				Hidden:  "secret",
				NoLabel: "unlabeled",
				Normal:  "labeled",
			},
			expected: []string{
				`<input type="hidden" name="hidden" value="secret"`,      // Hidden field
				`<input type="text" name="no_label" value="unlabeled"`,  // NoLabel field (text input)
				`<label for="normal">Normal</label>`,                    // Normal field has label
				`<input type="text" name="normal" value="labeled"`,      // Normal field
			},
			notExpected: []string{
				`<label for="hidden">Hidden</label>`,     // No label for hidden
				`<label for="no_label">No Label</label>`, // No label for nolabel
			},
		},
		{
			name: "hidden field with other attributes",
			input: struct {
				Phase int `vee:"hidden,required,id:'phase_id'" css:"hidden-field"`
			}{
				Phase: 3,
			},
			expected: []string{
				`<input type="hidden" name="phase" value="3"`,
				`id="phase_id"`,
				`class="hidden-field"`,
				`required`, // Required still works for validation
			},
			notExpected: []string{
				`<label for="phase_id">Phase</label>`, // No label even with custom id
			},
		},
		{
			name: "hidden fields with different data types",
			input: struct {
				ID       int           `vee:"hidden"`
				Score    float64       `vee:"hidden"`
				Active   bool          `vee:"hidden"`
				Created  time.Time     `vee:"hidden"`
				Timeout  time.Duration `vee:"hidden"`
			}{
				ID:      42,
				Score:   3.14,
				Active:  true,
				Created: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
				Timeout: 5 * time.Second,
			},
			expected: []string{
				`<input type="hidden" name="id" value="42"`,
				`<input type="hidden" name="score" value="3.14"`,
				`<input type="hidden" name="active" value="true"`,
				`<input type="hidden" name="created" value="2023-12-25T10:30:00Z"`,
				`<input type="hidden" name="timeout" value="5000000000"`, // Duration as nanoseconds
			},
			notExpected: []string{
				`<label`, // No labels for any hidden field
			},
		},
		{
			name: "hidden boolean false value",
			input: struct {
				Inactive bool `vee:"hidden"`
			}{
				Inactive: false,
			},
			expected: []string{
				`<input type="hidden" name="inactive" value="false"`,
			},
		},
		{
			name: "hidden field with name override",
			input: struct {
				InternalID int `vee:"$session_id,hidden"`
			}{
				InternalID: 123,
			},
			expected: []string{
				`<input type="hidden" name="session_id" value="123"`,
			},
			notExpected: []string{
				`<label for="session_id">Internal ID</label>`,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Render(tt.input)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			// Check expected strings
			for _, expected := range tt.expected {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected to find %q in output:\n%s", expected, result)
				}
			}

			// Check strings that should not be present
			for _, notExpected := range tt.notExpected {
				if strings.Contains(result, notExpected) {
					t.Errorf("Expected NOT to find %q in output:\n%s", notExpected, result)
				}
			}
		})
	}
}

func TestHiddenFieldErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    any
		expected string
	}{
		{
			name: "pointer type not supported",
			input: struct {
				Name *string `vee:"hidden"`
			}{},
			expected: "hidden attribute not supported for pointer type 'Name'",
		},
		{
			name: "slice type not supported",
			input: struct {
				Items []string `vee:"hidden"`
			}{},
			expected: "hidden attribute not supported for slice/array type 'Items'",
		},
		{
			name: "array type not supported",
			input: struct {
				Items [3]string `vee:"hidden"`
			}{},
			expected: "hidden attribute not supported for slice/array type 'Items'",
		},
		{
			name: "choices field not supported",
			input: struct {
				ColorChoices []string `vee:"hidden"`
				ColorChosen  int
			}{},
			expected: "hidden attribute not supported for multi-value field 'ColorChoices'",
		},
		{
			name: "chosen field not supported",
			input: struct {
				ColorChoices []string
				ColorChosen  int `vee:"hidden"`
			}{},
			expected: "hidden attribute not supported for multi-value field 'ColorChosen'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Render(tt.input)
			if err == nil {
				t.Fatal("Expected error but got none")
			}
			if !strings.Contains(err.Error(), tt.expected) {
				t.Errorf("Expected error to contain %q, got %v", tt.expected, err)
			}
		})
	}
}

func TestHiddenFieldBinding(t *testing.T) {
	// Test that hidden fields bind properly like normal fields
	tests := []struct {
		name     string
		formData map[string][]string
		target   func() any                     // Function to create fresh target struct
		check    func(t *testing.T, target any) // Function to verify result
		wantErr  bool
	}{
		{
			name: "bind hidden string field",
			formData: map[string][]string{
				"session_id": {"abc123"},
				"visible":    {"test"},
			},
			target: func() any {
				return &struct {
					SessionID string `vee:"$session_id,hidden"`
					Visible   string
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					SessionID string `vee:"$session_id,hidden"`
					Visible   string
				})
				if s.SessionID != "abc123" || s.Visible != "test" {
					t.Errorf("Bind() result = %+v, want SessionID='abc123', Visible='test'", s)
				}
			},
			wantErr: false,
		},
		{
			name: "bind hidden numeric fields",
			formData: map[string][]string{
				"id":    {"42"},
				"score": {"3.14"},
			},
			target: func() any {
				return &struct {
					ID    int     `vee:"hidden"`
					Score float64 `vee:"hidden"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					ID    int     `vee:"hidden"`
					Score float64 `vee:"hidden"`
				})
				if s.ID != 42 || s.Score != 3.14 {
					t.Errorf("Bind() result = %+v, want ID=42, Score=3.14", s)
				}
			},
			wantErr: false,
		},
		{
			name: "bind hidden boolean field",
			formData: map[string][]string{
				"active": {"true"},
			},
			target: func() any {
				return &struct {
					Active bool `vee:"hidden"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Active bool `vee:"hidden"`
				})
				if s.Active != true {
					t.Errorf("Bind() result = %+v, want Active=true", s)
				}
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.target()
			err := Bind(tt.formData, target)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Bind() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				tt.check(t, target)
			}
		})
	}
}