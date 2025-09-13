package vee

import (
	"testing"
	"time"
)

func TestLabelGeneration(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "default labels generated for all field types",
			input: struct {
				Name      string
				Age       int
				Price     float64
				Active    bool
				CreatedAt time.Time
				Timeout   time.Duration
			}{
				Name:      "John",
				Age:       25,
				Price:     19.99,
				Active:    true,
				CreatedAt: time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC),
				Timeout:   30 * time.Second,
			},
			want: `<form method="POST">
<label for="name">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="age">Age</label>
<input type="number" name="age" value="25" id="age">
<label for="price">Price</label>
<input type="number" name="price" value="19.99" step="any" id="price">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
<label for="created_at">Created At</label>
<input type="datetime-local" name="created_at" value="2023-12-25T15:30" id="created_at">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="30" id="timeout">
</form>
`,
		},
		{
			name: "default labels generated for all field types using custom css",
			input: struct {
				Name      string `labelCss:"bg-gray-200 dark:bg-gray-800"`
				Age       int
				Price     float64
				Active    bool
				CreatedAt time.Time
				Timeout   time.Duration
			}{
				Name:      "John",
				Age:       25,
				Price:     19.99,
				Active:    true,
				CreatedAt: time.Date(2023, 12, 25, 15, 30, 0, 0, time.UTC),
				Timeout:   30 * time.Second,
			},
			want: `<form method="POST">
<label for="name" class="bg-gray-200 dark:bg-gray-800">Name</label>
<input type="text" name="name" value="John" id="name">
<label for="age">Age</label>
<input type="number" name="age" value="25" id="age">
<label for="price">Price</label>
<input type="number" name="price" value="19.99" step="any" id="price">
<label for="active">Active</label>
<input type="checkbox" name="active" value="true" checked id="active">
<label for="created_at">Created At</label>
<input type="datetime-local" name="created_at" value="2023-12-25T15:30" id="created_at">
<label for="timeout">Timeout</label>
<input type="number" name="timeout" value="30" id="timeout">
</form>
`,
		},
		{
			name: "custom labels override default field names",
			input: struct {
				FirstName string `vee:"label:'Full Name'"`
				Email     string `vee:"type:'email',label:'Email Address'"`
			}{
				FirstName: "John",
				Email:     "john@example.com",
			},
			want: `<form method="POST">
<label for="first_name">Full Name</label>
<input type="text" name="first_name" value="John" id="first_name">
<label for="email">Email Address</label>
<input type="email" name="email" value="john@example.com" id="email">
</form>
`,
		},
		{
			name: "custom labels override default field names and use custom css",
			input: struct {
				FirstName string `vee:"label:'Full Name'" labelCss:"bg-gray-200 dark:bg-gray-800"`
				Email     string `vee:"type:'email',label:'Email Address'"`
			}{
				FirstName: "John",
				Email:     "john@example.com",
			},
			want: `<form method="POST">
<label for="first_name" class="bg-gray-200 dark:bg-gray-800">Full Name</label>
<input type="text" name="first_name" value="John" id="first_name">
<label for="email">Email Address</label>
<input type="email" name="email" value="john@example.com" id="email">
</form>
`,
		},
		{
			name: "nolabel attribute skips label generation",
			input: struct {
				Name     string `vee:"label:'User Name'"`
				Password string `vee:"type:'password',nolabel"`
				Email    string `vee:"nolabel"`
			}{
				Name:     "John",
				Password: "secret",
				Email:    "john@example.com",
			},
			want: `<form method="POST">
<label for="name">User Name</label>
<input type="text" name="name" value="John" id="name">
<input type="password" name="password" value="secret" id="password">
<input type="text" name="email" value="john@example.com" id="email">
</form>
`,
		},
		{
			name: "nolabel attribute skips label generation including custom css",
			input: struct {
				Name     string `vee:"label:'User Name'"`
				Password string `vee:"type:'password',nolabel"`
				Email    string `vee:"nolabel" labelCss:"bg-gray-200"`
			}{
				Name:     "John",
				Password: "secret",
				Email:    "john@example.com",
			},
			want: `<form method="POST">
<label for="name">User Name</label>
<input type="text" name="name" value="John" id="name">
<input type="password" name="password" value="secret" id="password">
<input type="text" name="email" value="john@example.com" id="email">
</form>
`,
		},
		{
			name: "labels work with custom field names and ids",
			input: struct {
				UserName string `vee:"$user_name,id:'username_field',label:'Username'"`
			}{
				UserName: "johndoe",
			},
			want: `<form method="POST">
<label for="username_field">Username</label>
<input type="text" name="user_name" value="johndoe" id="username_field">
</form>
`,
		},
		{
			name: "camelCase field names convert to readable labels",
			input: struct {
				FirstName       string
				EmailAddress    string
				PhoneNumber     string
				DateOfBirth     time.Time `vee:"type:'date'"`
				IsAccountActive bool
			}{
				FirstName:       "John",
				EmailAddress:    "john@example.com",
				PhoneNumber:     "555-1234",
				DateOfBirth:     time.Date(1990, 6, 15, 0, 0, 0, 0, time.UTC),
				IsAccountActive: true,
			},
			want: `<form method="POST">
<label for="first_name">First Name</label>
<input type="text" name="first_name" value="John" id="first_name">
<label for="email_address">Email Address</label>
<input type="text" name="email_address" value="john@example.com" id="email_address">
<label for="phone_number">Phone Number</label>
<input type="text" name="phone_number" value="555-1234" id="phone_number">
<label for="date_of_birth">Date Of Birth</label>
<input type="date" name="date_of_birth" value="1990-06-15" id="date_of_birth">
<label for="is_account_active">Is Account Active</label>
<input type="checkbox" name="is_account_active" value="true" checked id="is_account_active">
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

func TestLabelWithMultiValueFields(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "select dropdown gets label",
			input: struct {
				ColorChoices []string
				ColorChosen  int `vee:"type:'select',label:'Favorite Color'"`
			}{
				ColorChoices: []string{"Red", "Blue", "Green"},
				ColorChosen:  1,
			},
			want: `<form method="POST">
<label for="color_chosen">Favorite Color</label>
<select name="color_chosen" id="color_chosen">
<option value="0">Red</option>
<option value="1" selected>Blue</option>
<option value="2">Green</option>
</select>
</form>
`,
		},
		{
			name: "radio group gets fieldset with legend",
			input: struct {
				SizeChoices []string
				SizeChosen  int `vee:"type:'radio',label:'Size'"`
			}{
				SizeChoices: []string{"Small", "Medium", "Large"},
				SizeChosen:  1,
			},
			want: `<form method="POST">
<fieldset><legend>Size</legend>
<input type="radio" name="size_chosen" value="0" id="size_chosen_0"><label for="size_chosen_0">Small</label>
<input type="radio" name="size_chosen" value="1" checked id="size_chosen_1"><label for="size_chosen_1">Medium</label>
<input type="radio" name="size_chosen" value="2" id="size_chosen_2"><label for="size_chosen_2">Large</label>
</fieldset>
</form>
`,
		},
		{
			name: "checkbox group gets fieldset with legend",
			input: struct {
				FeatureChoices []string
				FeatureChosen  []int `vee:"type:'checkbox',label:'Features'"`
			}{
				FeatureChoices: []string{"WiFi", "Bluetooth"},
				FeatureChosen:  []int{0},
			},
			want: `<form method="POST">
<fieldset><legend>Features</legend>
<input type="checkbox" name="feature_chosen" value="0" checked id="feature_chosen_0"><label for="feature_chosen_0">WiFi</label>
<input type="checkbox" name="feature_chosen" value="1" id="feature_chosen_1"><label for="feature_chosen_1">Bluetooth</label>
</fieldset>
</form>
`,
		},
		{
			name: "nolabel works with multi-value fields",
			input: struct {
				ColorChoices []string
				ColorChosen  int `vee:"type:'select',nolabel"`
			}{
				ColorChoices: []string{"Red", "Blue"},
				ColorChosen:  0,
			},
			want: `<form method="POST">
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
func TestInternationalCharacterSupport(t *testing.T) {
	tests := []struct {
		name  string
		input any
		want  string
	}{
		{
			name: "German field names with umlauts in values",
			input: struct {
				NachName string
				VorName  string
			}{
				NachName: "Müller",
				VorName:  "Jürgen",
			},
			want: `<form method="POST">
<label for="nach_name">Nach Name</label>
<input type="text" name="nach_name" value="Müller" id="nach_name">
<label for="vor_name">Vor Name</label>
<input type="text" name="vor_name" value="Jürgen" id="vor_name">
</form>
`,
		},
		{
			name: "Mixed case with various international patterns",
			input: struct {
				UserID       int
				EmailAddress string
				PhoneNumber  string
				IsActiveUser bool
			}{
				UserID:       12345,
				EmailAddress: "andré@example.fr",
				PhoneNumber:  "+33-1-23-45-67-89",
				IsActiveUser: true,
			},
			want: `<form method="POST">
<label for="user_id">User I D</label>
<input type="number" name="user_id" value="12345" id="user_id">
<label for="email_address">Email Address</label>
<input type="text" name="email_address" value="andré@example.fr" id="email_address">
<label for="phone_number">Phone Number</label>
<input type="text" name="phone_number" value="+33-1-23-45-67-89" id="phone_number">
<label for="is_active_user">Is Active User</label>
<input type="checkbox" name="is_active_user" value="true" checked id="is_active_user">
</form>
`,
		},
		{
			name: "Single letter and consecutive uppercase handling",
			input: struct {
				A              string
				AB             string
				ABC            string
				XMLHttpRequest string
				HTMLParser     string
				URLPath        string
			}{
				A:              "single",
				AB:             "double",
				ABC:            "triple",
				XMLHttpRequest: "xml-http",
				HTMLParser:     "html-parse",
				URLPath:        "/api/v1",
			},
			want: `<form method="POST">
<label for="a">A</label>
<input type="text" name="a" value="single" id="a">
<label for="ab">A B</label>
<input type="text" name="ab" value="double" id="ab">
<label for="abc">A B C</label>
<input type="text" name="abc" value="triple" id="abc">
<label for="xml_http_request">X M L Http Request</label>
<input type="text" name="xml_http_request" value="xml-http" id="xml_http_request">
<label for="html_parser">H T M L Parser</label>
<input type="text" name="html_parser" value="html-parse" id="html_parser">
<label for="url_path">U R L Path</label>
<input type="text" name="url_path" value="/api/v1" id="url_path">
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
