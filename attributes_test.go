package vee

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"testing"
	"time"
)

func TestUniversalAttributes(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "required attribute on string field",
			input: struct {
				Name string `vee:"required"`
			}{Name: "John"},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name" required>
</form>
`,
		},
		{
			name: "readonly attribute on string field",
			input: struct {
				Name string `vee:"readonly"`
			}{Name: "John"},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name" readonly>
</form>
`,
		},
		{
			name: "disabled attribute on string field",
			input: struct {
				Name string `vee:"disabled"`
			}{Name: "John"},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name" disabled>
</form>
`,
		},
		{
			name: "placeholder attribute on string field",
			input: struct {
				Name string `vee:"placeholder:'Enter your name'"`
			}{Name: ""},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="" id="name" placeholder="Enter your name">
</form>
`,
		},
		{
			name: "custom id attribute on string field",
			input: struct {
				Name string `vee:"id:'user_name'"`
			}{Name: "John"},
			want: `<form method="POST">
<label for="user_name">Name</label>
<input type="text" name="name" value="John" id="user_name">
</form>
`,
		},
		{
			name: "multiple universal attributes on string field",
			input: struct {
				Email string `vee:"type:'email',required,placeholder:'Enter email',id:'user_email'"`
			}{Email: "john@example.com"},
			want: `<form method="POST">
<label for="user_email">Email</label>
<input type="email" name="email" value="john@example.com" id="user_email" placeholder="Enter email" required>
</form>
`,
		},
		{
			name: "universal attributes on numeric field",
			input: struct {
				Age int `vee:"required,min:18,max:120,placeholder:'Age'"`
			}{Age: 25},
			want: `<form method="POST">
<label for="age">Age</label>
<input type="number" name="age" value="25" min="18" max="120" id="age" placeholder="Age" required>
</form>
`,
		},
		{
			name: "universal attributes on float field",
			input: struct {
				Price float64 `vee:"required,min:0,step:0.01,readonly"`
			}{Price: 19.99},
			want: `<form method="POST">
<label for="price">Price</label>
<input type="number" name="price" value="19.99" min="0" step="0.01" id="price" required readonly>
</form>
`,
		},
		{
			name: "universal attributes on boolean field",
			input: struct {
				Active bool `vee:"required,disabled"`
			}{Active: true},
			want: `<form method="POST">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active" required disabled>
</form>
`,
		},
		{
			name: "universal attributes on time field",
			input: struct {
				Birthday time.Time `vee:"type:'date',required,min:'1900-01-01',max:'2023-12-31'"`
			}{Birthday: time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC)},
			want: `<form method="POST">
<label for="birthday">Birthday</label>
<input type="date" name="birthday" value="1990-06-15" min="1900-01-01" max="2023-12-31" id="birthday" required>
</form>
`,
		},
		{
			name: "universal attributes on duration field",
			input: struct {
				Timeout time.Duration `vee:"units:'s',required,min:1,max:3600"`
			}{Timeout: 30 * time.Second},
			want: `<form method="POST">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="30" min="1" max="3600" id="timeout" required>
</form>
`,
		},
		{
			name: "default id attribute (field name to snake_case)",
			input: struct {
				FirstName string
			}{FirstName: "John"},
			want: `<form method="POST">
<label for="first_name">First Name</label>
<input type="text" name="first_name" value="John" id="first_name">
</form>
`,
		},
		{
			name: "custom name with default id",
			input: struct {
				FirstName string `vee:"$user_first_name"`
			}{FirstName: "John"},
			want: `<form method="POST">
<label for="user_first_name">First Name</label>
<input type="text" name="user_first_name" value="John" id="user_first_name">
</form>
`,
		},
		{
			name: "custom name with custom id",
			input: struct {
				FirstName string `vee:"$user_first_name,id:'fname'"`
			}{FirstName: "John"},
			want: `<form method="POST">
<label for="fname">First Name</label>
<input type="text" name="user_first_name" value="John" id="fname">
</form>
`,
		},
		{
			name: "mixed fields with various universal attributes",
			input: struct {
				Name   string `vee:"required,placeholder:'Full name'"`
				Email  string `vee:"type:'email',required"`
				Age    int    `vee:"min:18,readonly"`
				Active bool   `vee:"disabled"`
			}{
				Name:   "John Doe",
				Email:  "john@example.com",
				Age:    25,
				Active: true,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John Doe" id="name" placeholder="Full name" required>
<label for="email">Email</label>
<input type="email" name="email" value="john@example.com" id="email" required>
<label for="age">Age</label>
<input type="number" name="age" value="25" min="18" id="age" readonly>
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active" disabled>
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

func TestStringTypeOverrides(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "email type override",
			input: struct {
				Email string `vee:"type:'email'"`
			}{Email: "john@example.com"},
			want: `<form method="POST">
<label for="email">Email</label>
<input type="email" name="email" value="john@example.com" id="email">
</form>
`,
		},
		{
			name: "password type override",
			input: struct {
				Password string `vee:"type:'password'"`
			}{Password: "secret"},
			want: `<form method="POST">
<label for="password">Password</label>
<input type="password" name="password" value="secret" id="password">
</form>
`,
		},
		{
			name: "tel type override",
			input: struct {
				Phone string `vee:"type:'tel'"`
			}{Phone: "555-1234"},
			want: `<form method="POST">
<label for="phone">Phone</label>
<input type="tel" name="phone" value="555-1234" id="phone">
</form>
`,
		},
		{
			name: "url type override",
			input: struct {
				Website string `vee:"type:'url'"`
			}{Website: "https://example.com"},
			want: `<form method="POST">
<label for="website">Website</label>
<input type="url" name="website" value="https://example.com" id="website">
</form>
`,
		},
		{
			name: "email type with universal attributes",
			input: struct {
				Email string `vee:"type:'email',required,placeholder:'Enter your email'"`
			}{Email: ""},
			want: `<form method="POST">
<label for="email">Email</label>
<input type="email" name="email" value="" id="email" placeholder="Enter your email" required>
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

const iterations = 1000

var allOptions = [][]string{
	{"placeholder", "testing"},
	{"required", ""},
	{"id", "abc"},
	{"readonly", ""},
	{"disabled", ""},
}

var optionCount = len(allOptions)

func selectRandomOption(previous []int) (string, string, []int) {
	found := false
	var index int
	for !found {
		index = rand.Intn(optionCount - 1)
		found = true
		if slices.Contains(previous, index) {
			found = false
			break
		}
	}
	return allOptions[index][0], allOptions[index][1], append(previous, index)
}

func buildExpected(config FieldConfig) string {
	buf := &strings.Builder{}
	if id, ok := config.Attributes["id"]; ok {
		buf.WriteString(fmt.Sprintf(` id="%s"`, escapeHTML(id)))
	} else {
		buf.WriteString(fmt.Sprintf(` id="%s"`, escapeHTML(config.Name)))
	}

	// Add placeholder attribute
	if placeholder, ok := config.Attributes["placeholder"]; ok {
		buf.WriteString(fmt.Sprintf(` placeholder="%s"`, escapeHTML(placeholder)))
	}

	// Add boolean attributes (required, readonly, disabled)
	if _, ok := config.Attributes["required"]; ok {
		buf.WriteString(` required`)
	}
	if _, ok := config.Attributes["readonly"]; ok {
		buf.WriteString(` readonly`)
	}
	if _, ok := config.Attributes["disabled"]; ok {
		buf.WriteString(` disabled`)
	}
	return buf.String()
}

func TestUniversalAttributePermutations(t *testing.T) {
	for i := range iterations {
		config := FieldConfig{
			Name:       "foo",
			Attributes: make(map[string]string, 0),
		}
		var key, value string
		previous := make([]int, 0)
		for j := 0; j < 3; j++ {
			key, value, previous = selectRandomOption(previous)
			config.Attributes[key] = value
		}
		output := &strings.Builder{}
		addUniversalAttributes(output, config)
		expected := buildExpected(config)
		if output.String() != expected {
			t.Errorf("iteration %d failed. expected '%s', got '%s'", i, expected, output.String())
		}
	}
}
