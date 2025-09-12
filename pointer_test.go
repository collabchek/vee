package vee

import (
	"testing"
	"time"
)

func TestPointerRendering(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "pointer to string with value renders with value",
			input: struct {
				Name *string
			}{Name: stringPtr("John")},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
</form>
`,
		},
		{
			name: "nil pointer to string renders without value",
			input: struct {
				Name *string
			}{Name: nil},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="" id="name">
</form>
`,
		},
		{
			name: "pointer to int with value",
			input: struct {
				Age *int
			}{Age: intPtr(25)},
			want: `<form method="POST">
<label for="age">Age</label>
<input type="number" name="age" value="25" id="age">
</form>
`,
		},
		{
			name: "nil pointer to int renders without value",
			input: struct {
				Age *int
			}{Age: nil},
			want: `<form method="POST">
<label for="age">Age</label>
<input type="number" name="age" value="0" id="age">
</form>
`,
		},
		{
			name: "pointer to bool with true value",
			input: struct {
				Active *bool
			}{Active: boolPtr(true)},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
</form>
`,
		},
		{
			name: "pointer to bool with false value",
			input: struct {
				Active *bool
			}{Active: boolPtr(false)},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" id="active">
</form>
`,
		},
		{
			name: "nil pointer to bool renders unchecked",
			input: struct {
				Active *bool
			}{Active: nil},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" id="active">
</form>
`,
		},
		{
			name: "pointer to float64 with value",
			input: struct {
				Price *float64
			}{Price: float64Ptr(29.99)},
			want: `<form method="POST">
<label for="price">Price</label>
<input type="number" name="price" value="29.99" step="any" id="price">
</form>
`,
		},
		{
			name: "nil pointer to float64 renders without value",
			input: struct {
				Price *float64
			}{Price: nil},
			want: `<form method="POST">
<label for="price">Price</label>
<input type="number" name="price" value="0" step="any" id="price">
</form>
`,
		},
		{
			name: "pointer to time.Time with value",
			input: struct {
				Created *time.Time
			}{Created: timePtr(time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC))},
			want: `<form method="POST">
<label for="created">Created</label>
<input type="datetime-local" name="created" value="2023-12-25T15:30" id="created">
</form>
`,
		},
		{
			name: "nil pointer to time.Time renders without value",
			input: struct {
				Created *time.Time
			}{Created: nil},
			want: `<form method="POST">
<label for="created">Created</label>
<input type="datetime-local" name="created" id="created">
</form>
`,
		},
		{
			name: "pointer to time.Duration with value",
			input: struct {
				Timeout *time.Duration
			}{Timeout: durationPtr(5 * time.Second)},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="5" id="timeout">
</form>
`,
		},
		{
			name: "nil pointer to time.Duration renders without value",
			input: struct {
				Timeout *time.Duration
			}{Timeout: nil},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" id="timeout">
</form>
`,
		},
		{
			name: "mixed pointer and non-pointer fields",
			input: struct {
				Name     string
				Email    *string
				Age      *int
				IsActive bool
			}{
				Name:     "John",
				Email:    stringPtr("john@example.com"),
				Age:      nil,
				IsActive: true,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" id="email">
<label for="age">Age</label>
<input type="number" name="age" value="0" id="age">
<label for="is_active">Is Active</label>
<input type="checkbox" name="is_active" value="true" checked id="is_active">
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

func TestPointerBinding(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		check   func(t *testing.T, target any)
		wantErr bool
	}{
		{
			name: "bind string to pointer field",
			input: map[string][]string{
				"name": {"John Doe"},
			},
			target: func() any {
				return &struct {
					Name *string
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Name *string })
				if s.Name == nil {
					t.Errorf("Expected Name to be non-nil")
					return
				}
				if *s.Name != "John Doe" {
					t.Errorf("Expected Name='John Doe', got Name='%s'", *s.Name)
				}
			},
		},
		{
			name: "bind int to pointer field",
			input: map[string][]string{
				"age": {"30"},
			},
			target: func() any {
				return &struct {
					Age *int
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Age *int })
				if s.Age == nil {
					t.Errorf("Expected Age to be non-nil")
					return
				}
				if *s.Age != 30 {
					t.Errorf("Expected Age=30, got Age=%d", *s.Age)
				}
			},
		},
		{
			name: "bind bool to pointer field",
			input: map[string][]string{
				"active": {"true"},
			},
			target: func() any {
				return &struct {
					Active *bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Active *bool })
				if s.Active == nil {
					t.Errorf("Expected Active to be non-nil")
					return
				}
				if *s.Active != true {
					t.Errorf("Expected Active=true, got Active=%t", *s.Active)
				}
			},
		},
		{
			name: "bind float to pointer field",
			input: map[string][]string{
				"price": {"29.99"},
			},
			target: func() any {
				return &struct {
					Price *float64
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Price *float64 })
				if s.Price == nil {
					t.Errorf("Expected Price to be non-nil")
					return
				}
				if *s.Price != 29.99 {
					t.Errorf("Expected Price=29.99, got Price=%f", *s.Price)
				}
			},
		},
		{
			name: "bind time to pointer field",
			input: map[string][]string{
				"created": {"2023-12-25T15:30"},
			},
			target: func() any {
				return &struct {
					Created *time.Time
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Created *time.Time })
				if s.Created == nil {
					t.Errorf("Expected Created to be non-nil")
					return
				}
				expected := time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC)
				if !s.Created.Equal(expected) {
					t.Errorf("Expected Created=%v, got Created=%v", expected, *s.Created)
				}
			},
		},
		{
			name: "bind duration to pointer field",
			input: map[string][]string{
				"timeout": {"5"},
			},
			target: func() any {
				return &struct {
					Timeout *time.Duration
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Timeout *time.Duration })
				if s.Timeout == nil {
					t.Errorf("Expected Timeout to be non-nil")
					return
				}
				expected := 5 * time.Second
				if *s.Timeout != expected {
					t.Errorf("Expected Timeout=%v, got Timeout=%v", expected, *s.Timeout)
				}
			},
		},
		{
			name: "missing form data leaves pointer nil",
			input: map[string][]string{
				"other_field": {"value"},
			},
			target: func() any {
				return &struct {
					Name *string
					Age  *int
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name *string
					Age  *int
				})
				if s.Name != nil {
					t.Errorf("Expected Name to be nil, got %v", *s.Name)
				}
				if s.Age != nil {
					t.Errorf("Expected Age to be nil, got %v", *s.Age)
				}
			},
		},
		{
			name: "mixed pointer and non-pointer binding",
			input: map[string][]string{
				"name":      {"John"},
				"email":     {"john@example.com"},
				"is_active": {"true"},
			},
			target: func() any {
				return &struct {
					Name     string
					Email    *string
					Age      *int // Not provided in form
					IsActive bool
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name     string
					Email    *string
					Age      *int
					IsActive bool
				})
				if s.Name != "John" {
					t.Errorf("Expected Name='John', got Name='%s'", s.Name)
				}
				if s.Email == nil || *s.Email != "john@example.com" {
					t.Errorf("Expected Email='john@example.com', got Email=%v", s.Email)
				}
				if s.Age != nil {
					t.Errorf("Expected Age to be nil, got %v", *s.Age)
				}
				if !s.IsActive {
					t.Errorf("Expected IsActive=true, got IsActive=%t", s.IsActive)
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

// Helper functions to create pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func durationPtr(d time.Duration) *time.Duration {
	return &d
}
