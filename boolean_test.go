package vee

import (
	"testing"
)

func TestBooleanRendering(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "true bool field renders as checked checkbox",
			input: struct {
				Active bool
			}{Active: true},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
</form>
`,
		},
		{
			name: "false bool field renders as unchecked checkbox",
			input: struct {
				Active bool
			}{Active: false},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" id="active">
</form>
`,
		},
		{
			name: "mixed types with boolean",
			input: struct {
				Name   string
				Age    int
				Active bool
			}{
				Name:   "John",
				Age:    25,
				Active: true,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="age">Age</label>
<input type="number" name="age" value="25" id="age">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
</form>
`,
		},
		{
			name: "boolean with custom name override",
			input: struct {
				IsEnabled bool `vee:"$is_enabled"`
			}{IsEnabled: false},
			want: `<form method="POST">
<label for="is_enabled">Is Enabled</label>
<input type="checkbox" name="is_enabled" value="true" id="is_enabled">
</form>
`,
		},
		{
			name: "boolean with CSS classes",
			input: struct {
				Active bool `css:"form-check-input"`
			}{Active: true},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked class="form-check-input" id="active">
</form>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if err != nil {
				t.Errorf("Render() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBooleanBinding(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		check   func(t *testing.T, target any)
		wantErr bool
	}{
		{
			name: "checkbox present in form data sets bool to true",
			input: map[string][]string{
				"active": {"true"},
			},
			target: func() any {
				return &struct {
					Active bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Active bool })
				if !s.Active {
					t.Errorf("Expected Active=true, got Active=%t", s.Active)
				}
			},
		},
		{
			name: "checkbox absent from form data sets bool to false",
			input: map[string][]string{
				"name": {"John"}, // Other fields present, but not 'active'
			},
			target: func() any {
				return &struct {
					Name   string
					Active bool
				}{
					Active: true, // Initial value should be overridden
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name   string
					Active bool
				})
				if s.Name != "John" || s.Active {
					t.Errorf("Expected Name='John' Active=false, got Name='%s' Active=%t", s.Name, s.Active)
				}
			},
		},
		{
			name: "checkbox with custom name override",
			input: map[string][]string{
				"is_enabled": {"on"},
			},
			target: func() any {
				return &struct {
					IsEnabled bool `vee:"$is_enabled"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					IsEnabled bool `vee:"$is_enabled"`
				})
				if !s.IsEnabled {
					t.Errorf("Expected IsEnabled=true, got IsEnabled=%t", s.IsEnabled)
				}
			},
		},
		{
			name: "mixed field types with booleans",
			input: map[string][]string{
				"name":   {"John"},
				"age":    {"25"},
				"active": {"true"},
				// 'premium' checkbox absent - should be false
			},
			target: func() any {
				return &struct {
					Name    string
					Age     int
					Active  bool
					Premium bool
				}{
					Premium: true, // Initial value should be overridden to false
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name    string
					Age     int
					Active  bool
					Premium bool
				})
				if s.Name != "John" || s.Age != 25 || !s.Active || s.Premium {
					t.Errorf("Expected Name='John' Age=25 Active=true Premium=false, got Name='%s' Age=%d Active=%t Premium=%t",
						s.Name, s.Age, s.Active, s.Premium)
				}
			},
		},
		{
			name: "checkbox value doesn't matter - presence determines truth",
			input: map[string][]string{
				"active": {"false"}, // Even "false" string means checkbox was checked
			},
			target: func() any {
				return &struct {
					Active bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Active bool })
				if !s.Active {
					t.Errorf("Expected Active=true (checkbox present), got Active=%t", s.Active)
				}
			},
		},
		{
			name: "empty checkbox value still means true",
			input: map[string][]string{
				"active": {""}, // Empty value but key exists
			},
			target: func() any {
				return &struct {
					Active bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Active bool })
				if !s.Active {
					t.Errorf("Expected Active=true (checkbox key present), got Active=%t", s.Active)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.target()
			err := Bind(tt.input, target)

			if (err != nil) != tt.wantErr {
				t.Errorf("Bind() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				tt.check(t, target)
			}
		})
	}
}
