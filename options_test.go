package vee

import (
	"testing"
)

func TestRenderWithOptions(t *testing.T) {
	type TestStruct struct {
		Name  string
		Email string `css:"custom-email-class"`
	}

	tests := []struct {
		name    string
		input   TestStruct
		options RenderOption
		want    string
	}{
		{
			name:    "default behavior without options",
			input:   TestStruct{Name: "John", Email: "john@example.com"},
			options: RenderOption{},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" class="custom-email-class" id="email">
</form>
`,
		},
		{
			name:  "default CSS class for all inputs",
			input: TestStruct{Name: "John", Email: "john@example.com"},
			options: RenderOption{
				DefaultInputCSS: "form-control",
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" class="form-control" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" class="custom-email-class" id="email">
</form>
`,
		},
		{
			name:  "form with basic attributes",
			input: TestStruct{Name: "John", Email: "john@example.com"},
			options: RenderOption{
				FormID:  "user-form",
				FormCSS: "my-form",
			},
			want: `<form id="user-form" class="my-form" method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" class="custom-email-class" id="email">
</form>
`,
		},
		{
			name:  "form with custom method and action",
			input: TestStruct{Name: "John", Email: "john@example.com"},
			options: RenderOption{
				FormMethod: "GET",
				FormAction: "/submit",
			},
			want: `<form method="GET" action="/submit">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" class="custom-email-class" id="email">
</form>
`,
		},
		{
			name:  "combined options - form attributes, default CSS, and field CSS",
			input: TestStruct{Name: "John", Email: "john@example.com"},
			options: RenderOption{
				FormID:          "contact-form",
				FormCSS:         "p-4 border rounded",
				FormMethod:      "POST",
				FormAction:      "/contact",
				DefaultInputCSS: "px-3 py-2 border rounded",
			},
			want: `<form id="contact-form" class="p-4 border rounded" method="POST" action="/contact">
<label for="name">Name</label>
<input type="text" name="name" value="John" class="px-3 py-2 border rounded" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" class="custom-email-class" id="email">
</form>
`,
		},
		{
			name:  "HTML escaping in form attributes",
			input: TestStruct{Name: "John", Email: "test@example.com"},
			options: RenderOption{
				FormID:     `form"with"quotes`,
				FormCSS:    `class<with>brackets`,
				FormAction: `/path?param="value"`,
			},
			want: `<form id="form&quot;with&quot;quotes" class="class&lt;with&gt;brackets" method="POST" action="/path?param=&quot;value&quot;">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="test@example.com" class="custom-email-class" id="email">
</form>
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input, tt.options)
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

func TestCSSTagPrecedence(t *testing.T) {
	type TestStruct struct {
		WithCSS    string `css:"field-specific-class"`
		WithoutCSS string
	}

	input := TestStruct{WithCSS: "value1", WithoutCSS: "value2"}
	options := RenderOption{DefaultInputCSS: "default-class"}

	got, err := Render(input, options)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	expected := `<form method="POST">
<label for="with_css">With C S S</label>
<input type="text" name="with_css" value="value1" class="field-specific-class" id="with_css">
<label for="without_css">Without C S S</label>
<input type="text" name="without_css" value="value2" class="default-class" id="without_css">
</form>
`

	if got != expected {
		t.Errorf("CSS precedence test failed.\nGot: %q\nWant: %q", got, expected)
	}
}

func TestOptionConsolidation(t *testing.T) {
	type TestStruct struct {
		options []RenderOption
	}
	tests := []struct {
		name  string
		input TestStruct
		want  RenderOption
	}{
		{
			name: "Input CSS styling option properly applied",
			input: TestStruct{
				options: []RenderOption{
					InputCSSOption("a"),
				},
			},
			want: InputCSSOption("a"),
		},
		{
			name: "Form id option properly applied",
			input: TestStruct{
				options: []RenderOption{
					FormIDOption("user-reg"),
				},
			},
			want: FormIDOption("user-reg"),
		},
		{
			name: "Form CSS option properly applied",
			input: TestStruct{
				options: []RenderOption{
					FormCSSOption("a"),
				},
			},
			want: FormCSSOption("a"),
		},
		{
			name: "Form method option properly applied",
			input: TestStruct{
				options: []RenderOption{
					FormMethodOption("PUT"),
				},
			},
			want: FormMethodOption("PUT"),
		},
		{
			name: "Form action option properly applied",
			input: TestStruct{
				options: []RenderOption{
					FormActionOption("/user-registration"),
				},
			},
			want: FormActionOption("/user-registration"),
		},
		{
			name: "Multiple options properly applied",
			input: TestStruct{
				options: []RenderOption{
					FormIDOption("user-reg"),
					FormMethodOption("POST"),
				},
			},
			want: RenderOption{FormID: "user-reg", FormMethod: "POST"},
		},
		{
			name: "Last value of competing options wins",
			input: TestStruct{
				options: []RenderOption{
					FormMethodOption("POST"),
					FormMethodOption("PUT"),
					FormMethodOption("GET"),
				},
			},
			want: RenderOption{FormMethod: "GET"},
		},
	}
	for _, tt := range tests {
		result := ConsolidateOptions(tt.input.options...)
		if !(*result).IsEqual(tt.want) {
			t.Errorf("Option consolidation test '%s' failed. Got %q, wanted %q\n", tt.name, *result, tt.want)
		}
	}
}
