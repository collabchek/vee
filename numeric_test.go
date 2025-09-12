package vee

import (
	"testing"
)

func TestNumericRendering(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "int field renders as number input",
			input: struct {
				Age int
			}{Age: 25},
			want: `<form method="POST">
<label for="age">Age</label>
<input type="number" name="age" value="25" id="age">
</form>
`,
		},
		{
			name: "int64 field renders as number input",
			input: struct {
				ID int64
			}{ID: 1234567890},
			want: `<form method="POST">
<label for="id">I D</label>
<input type="number" name="id" value="1234567890" id="id">
</form>
`,
		},
		{
			name: "float64 field renders with step=any",
			input: struct {
				Price float64
			}{Price: 29.99},
			want: `<form method="POST">
<label for="price">Price</label>
<input type="number" name="price" value="29.99" step="any" id="price">
</form>
`,
		},
		{
			name: "mixed string and numeric fields",
			input: struct {
				Name  string
				Age   int
				Score float64
			}{
				Name:  "John",
				Age:   30,
				Score: 95.5,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="age">Age</label>
<input type="number" name="age" value="30" id="age">
<label for="score">Score</label>
<input type="number" name="score" value="95.5" step="any" id="score">
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

func TestNumericAttributes(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "int with min/max attributes",
			input: struct {
				Age int `vee:"min:18,max:65"`
			}{Age: 30},
			want: `<form method="POST">
<label for="age">Age</label>
<input type="number" name="age" value="30" min="18" max="65" id="age">
</form>
`,
		},
		{
			name: "int with step attribute",
			input: struct {
				Rating int `vee:"min:1,max:10,step:1"`
			}{Rating: 5},
			want: `<form method="POST">
<label for="rating">Rating</label>
<input type="number" name="rating" value="5" min="1" max="10" step="1" id="rating">
</form>
`,
		},
		{
			name: "float64 with custom step overriding default",
			input: struct {
				Price float64 `vee:"min:0,step:0.01"`
			}{Price: 19.99},
			want: `<form method="POST">
<label for="price">Price</label>
<input type="number" name="price" value="19.99" min="0" step="0.01" id="price">
</form>
`,
		},
		{
			name: "numeric field with custom name override",
			input: struct {
				UserAge int `vee:"$user_age,min:13,max:120"`
			}{UserAge: 25},
			want: `<form method="POST">
<label for="user_age">User Age</label>
<input type="number" name="user_age" value="25" min="13" max="120" id="user_age">
</form>
`,
		},
		{
			name: "numeric attributes with CSS classes",
			input: struct {
				Count int `vee:"min:0,max:100" css:"w-20 text-center"`
			}{Count: 42},
			want: `<form method="POST">
<label for="count">Count</label>
<input type="number" name="count" value="42" min="0" max="100" class="w-20 text-center" id="count">
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

func TestNumericBinding(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		check   func(t *testing.T, target any)
		wantErr bool
	}{
		{
			name: "bind int field",
			input: map[string][]string{
				"age": {"30"},
			},
			target: func() any {
				return &struct {
					Age int
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Age int })
				if s.Age != 30 {
					t.Errorf("Expected Age=30, got Age=%d", s.Age)
				}
			},
		},
		{
			name: "bind int64 field",
			input: map[string][]string{
				"id": {"1234567890"},
			},
			target: func() any {
				return &struct {
					ID int64
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ ID int64 })
				if s.ID != 1234567890 {
					t.Errorf("Expected ID=1234567890, got ID=%d", s.ID)
				}
			},
		},
		{
			name: "bind float64 field",
			input: map[string][]string{
				"price": {"29.99"},
			},
			target: func() any {
				return &struct {
					Price float64
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct{ Price float64 })
				if s.Price != 29.99 {
					t.Errorf("Expected Price=29.99, got Price=%f", s.Price)
				}
			},
		},
		{
			name: "bind mixed numeric and string fields",
			input: map[string][]string{
				"name":  {"John"},
				"age":   {"25"},
				"score": {"95.5"},
			},
			target: func() any {
				return &struct {
					Name  string
					Age   int
					Score float64
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name  string
					Age   int
					Score float64
				})
				if s.Name != "John" || s.Age != 25 || s.Score != 95.5 {
					t.Errorf("Expected Name='John' Age=25 Score=95.5, got Name='%s' Age=%d Score=%f", s.Name, s.Age, s.Score)
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
