package vee

import (
	"strings"
	"testing"
)

func TestMultiValueConventions(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		wantError bool
		errorMsg  string
	}{
		{
			name: "valid single select convention",
			input: struct {
				ColorChoices []string
				ColorChosen  int
			}{
				ColorChoices: []string{"Red", "Blue", "Green"},
				ColorChosen:  1,
			},
			wantError: false,
		},
		{
			name: "valid multi-select convention",
			input: struct {
				SkillChoices []string
				SkillChosen  []int
			}{
				SkillChoices: []string{"Go", "JavaScript", "Python"},
				SkillChosen:  []int{0, 2},
			},
			wantError: false,
		},
		{
			name: "valid mixed fields with multi-value",
			input: struct {
				Name         string
				ColorChoices []string
				ColorChosen  int
				Age          int
			}{
				Name:         "John",
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  0,
				Age:          30,
			},
			wantError: false,
		},
		{
			name: "error: Chosen field without corresponding Choices",
			input: struct {
				ColorChosen int
			}{
				ColorChosen: 1,
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' requires corresponding 'ColorChoices' field",
		},
		{
			name: "error: Choices field without corresponding Chosen",
			input: struct {
				ColorChoices []string
			}{
				ColorChoices: []string{"Red", "Blue"},
			},
			wantError: true,
			errorMsg:  "field 'ColorChoices' requires corresponding 'ColorChosen' field",
		},
		{
			name: "error: invalid Choices type (not slice)",
			input: struct {
				ColorChoices string
				ColorChosen  int
			}{
				ColorChoices: "Red,Blue",
				ColorChosen:  1,
			},
			wantError: true,
			errorMsg:  "field 'ColorChoices' must be a slice type, got string",
		},
		{
			name: "error: invalid Chosen type for single select (not int)",
			input: struct {
				ColorChoices []string
				ColorChosen  string
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  "Red",
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' must be int or []int, got string",
		},
		{
			name: "error: invalid Chosen type for multi-select (not []int)",
			input: struct {
				ColorChoices []string
				ColorChosen  []string
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  []string{"Red"},
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' must be int or []int, got []string",
		},
		{
			name: "error: empty Choices slice",
			input: struct {
				ColorChoices []string
				ColorChosen  int
			}{
				ColorChoices: []string{},
				ColorChosen:  0,
			},
			wantError: true,
			errorMsg:  "field 'ColorChoices' cannot be empty",
		},
		{
			name: "error: Chosen index out of range for single select",
			input: struct {
				ColorChoices []string
				ColorChosen  int
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  5, // Invalid index
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' index 5 out of range for 2 choices",
		},
		{
			name: "error: Chosen indices out of range for multi-select",
			input: struct {
				ColorChoices []string
				ColorChosen  []int
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  []int{0, 5}, // Index 5 invalid
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' index 5 out of range for 2 choices",
		},
		{
			name: "error: negative Chosen index",
			input: struct {
				ColorChoices []string
				ColorChosen  int
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  -1,
			},
			wantError: true,
			errorMsg:  "field 'ColorChosen' index -1 out of range for 2 choices",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Render(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error, but got none")
					return
				}
				if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}
			}
		})
	}
}

func TestMultiValueRendering(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "single select dropdown",
			input: struct {
				ColorChoices []string
				ColorChosen  int `vee:"type:'select'"`
			}{
				ColorChoices: []string{"Red", "Blue", "Green"},
				ColorChosen:  1, // "Blue" selected
			},
			want: `<form method="POST">
<label for="color_chosen">Color Chosen</label>
<select name="color_chosen" id="color_chosen">
<option value="0">Red</option>
<option value="1" selected>Blue</option>
<option value="2">Green</option>
</select>
</form>
`,
		},
		{
			name: "multi-select dropdown",
			input: struct {
				SkillChoices []string
				SkillChosen  []int `vee:"type:'select',multiple"`
			}{
				SkillChoices: []string{"Go", "JavaScript", "Python"},
				SkillChosen:  []int{0, 2}, // "Go" and "Python" selected
			},
			want: `<form method="POST">
<label for="skill_chosen">Skill Chosen</label>
<select name="skill_chosen" multiple id="skill_chosen">
<option value="0" selected>Go</option>
<option value="1">JavaScript</option>
<option value="2" selected>Python</option>
</select>
</form>
`,
		},
		{
			name: "radio button group",
			input: struct {
				SizeChoices []string
				SizeChosen  int `vee:"type:'radio'"`
			}{
				SizeChoices: []string{"Small", "Medium", "Large"},
				SizeChosen:  1, // "Medium" selected
			},
			want: `<form method="POST">
<fieldset><legend>Size Chosen</legend>
<input type="radio" name="size_chosen" value="0" id="size_chosen_0"><label for="size_chosen_0">Small</label>
<input type="radio" name="size_chosen" value="1" checked id="size_chosen_1"><label for="size_chosen_1">Medium</label>
<input type="radio" name="size_chosen" value="2" id="size_chosen_2"><label for="size_chosen_2">Large</label>
</fieldset>
</form>
`,
		},
		{
			name: "checkbox group",
			input: struct {
				FeatureChoices []string
				FeatureChosen  []int `vee:"type:'checkbox'"`
			}{
				FeatureChoices: []string{"WiFi", "Bluetooth", "GPS"},
				FeatureChosen:  []int{0, 2}, // "WiFi" and "GPS" selected
			},
			want: `<form method="POST">
<fieldset><legend>Feature Chosen</legend>
<input type="checkbox" name="feature_chosen" value="0" checked id="feature_chosen_0"><label for="feature_chosen_0">WiFi</label>
<input type="checkbox" name="feature_chosen" value="1" id="feature_chosen_1"><label for="feature_chosen_1">Bluetooth</label>
<input type="checkbox" name="feature_chosen" value="2" checked id="feature_chosen_2"><label for="feature_chosen_2">GPS</label>
</fieldset>
</form>
`,
		},
		{
			name: "mixed regular and multi-value fields",
			input: struct {
				Name         string
				Email        string
				ColorChoices []string
				ColorChosen  int `vee:"type:'select'"`
			}{
				Name:         "John",
				Email:        "john@example.com",
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  0,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="email">Email</label>
<input type="text" name="email" value="john@example.com" id="email">
<label for="color_chosen">Color Chosen</label>
<select name="color_chosen" id="color_chosen">
<option value="0" selected>Red</option>
<option value="1">Blue</option>
</select>
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

func TestMultiValueBinding(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		check   func(t *testing.T, target any)
		wantErr bool
	}{
		{
			name: "bind single select",
			input: map[string][]string{
				"color_chosen": {"2"},
			},
			target: func() any {
				return &struct {
					ColorChoices []string
					ColorChosen  int
				}{
					ColorChoices: []string{"Red", "Blue", "Green"},
					ColorChosen:  0,
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					ColorChoices []string
					ColorChosen  int
				})
				if s.ColorChosen != 2 {
					t.Errorf("Expected ColorChosen=2, got %d", s.ColorChosen)
				}
			},
		},
		{
			name: "bind multi-select",
			input: map[string][]string{
				"skill_chosen": {"0", "2"},
			},
			target: func() any {
				return &struct {
					SkillChoices []string
					SkillChosen  []int
				}{
					SkillChoices: []string{"Go", "JS", "Python"},
					SkillChosen:  []int{},
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					SkillChoices []string
					SkillChosen  []int
				})
				expected := []int{0, 2}
				if len(s.SkillChosen) != len(expected) {
					t.Errorf("Expected SkillChosen length %d, got %d", len(expected), len(s.SkillChosen))
					return
				}
				for i, v := range expected {
					if s.SkillChosen[i] != v {
						t.Errorf("Expected SkillChosen[%d]=%d, got %d", i, v, s.SkillChosen[i])
					}
				}
			},
		},
		{
			name: "bind with validation - out of range index",
			input: map[string][]string{
				"color_chosen": {"5"}, // Invalid index
			},
			target: func() any {
				return &struct {
					ColorChoices []string
					ColorChosen  int
				}{
					ColorChoices: []string{"Red", "Blue"},
					ColorChosen:  0,
				}
			},
			wantErr: true,
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

			if !tt.wantErr && tt.check != nil {
				tt.check(t, target)
			}
		})
	}
}
