package vee

import (
	"testing"
)

func TestRender(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    string
		wantErr bool
	}{
		{
			name: "simple struct with string fields",
			input: struct {
				Name  string
				Email string
			}{
				Name:  "John Doe",
				Email: "john@example.com",
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John Doe" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" id="email">
</form>
`,
			wantErr: false,
		},
		{
			name: "struct with custom field names",
			input: struct {
				FirstName string `vee:"$first_name"`
				LastName  string `vee:"$last_name"`
			}{
				FirstName: "John",
				LastName:  "Doe",
			},
			want: `<form method="POST">
<label for="first_name">First Name</label>
<input type="text" name="first_name" value="John" id="first_name">
<label for="last_name">Last Name</label>
<input type="text" name="last_name" value="Doe" id="last_name">
</form>
`,
			wantErr: false,
		},
		{
			name: "struct with skipped field",
			input: struct {
				Name     string
				Internal string `vee:"-"`
				Email    string
			}{
				Name:     "John",
				Internal: "secret",
				Email:    "john@example.com",
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" id="email">
</form>
`,
			wantErr: false,
		},
		{
			name: "struct with numeric fields rendered",
			input: struct {
				Name string
				Age  int
			}{
				Name: "John",
				Age:  30,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="age">Age</label>
<input type="number" name="age" value="30" id="age">
</form>
`,
			wantErr: false,
		},
		{
			name: "struct with HTML characters escaped",
			input: struct {
				Name string
			}{
				Name: `John "The Great" <smith@example.com>`,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John &quot;The Great&quot; &lt;smith@example.com&gt;" id="name">
</form>
`,
			wantErr: false,
		},
		{
			name: "pointer to struct",
			input: &struct {
				Name string
			}{
				Name: "John",
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
</form>
`,
			wantErr: false,
		},
		{
			name:    "non-struct input returns error",
			input:   "not a struct",
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Render(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Render() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestEscapeHTML(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"normal text", "normal text"},
		{"<script>", "&lt;script&gt;"},
		{`"quoted"`, "&quot;quoted&quot;"},
		{"& ampersand", "&amp; ampersand"},
		{`<>"&`, "&lt;&gt;&quot;&amp;"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := escapeHTML(tt.input)
			if got != tt.want {
				t.Errorf("escapeHTML(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
