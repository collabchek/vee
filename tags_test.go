package vee

import "testing"

func TestParseVeeTag(t *testing.T) {
	tests := []struct {
		name      string
		tag       string
		fieldName string
		want      FieldConfig
	}{
		{
			name:      "empty tag uses auto-derived name",
			tag:       "",
			fieldName: "FirstName",
			want:      FieldConfig{Name: "first_name"},
		},
		{
			name:      "skip field with dash",
			tag:       "-",
			fieldName: "Internal",
			want:      FieldConfig{Skip: true},
		},
		{
			name:      "override name with dollar prefix",
			tag:       "$customName",
			fieldName: "FirstName",
			want:      FieldConfig{Name: "customName"},
		},
		{
			name:      "override with comma (ignores rest for now)",
			tag:       "$userName,required",
			fieldName: "Name",
			want:      FieldConfig{Name: "userName"},
		},
		{
			name:      "no override auto-derives name",
			tag:       "required",
			fieldName: "EmailAddress",
			want:      FieldConfig{Name: "email_address"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseVeeTag(tt.tag, tt.fieldName)
			if got.Name != tt.want.Name || got.Skip != tt.want.Skip {
				t.Errorf("parseVeeTag() = %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestStrCaseConversion(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"FirstName", "first_name"},
		{"EmailAddress", "email_address"},
		{"UserID", "user_id"},         // Fixed: should be user_id not user_i_d
		{"HTMLParser", "html_parser"}, // Fixed: should be html_parser not h_t_m_l_parser
		{"name", "name"},
		{"Name", "name"},
		{"", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			config := parseVeeTag("", tt.input)
			if config.Name != tt.want {
				t.Errorf("parseVeeTag(\"\", %q).Name = %q, want %q", tt.input, config.Name, tt.want)
			}
		})
	}
}
