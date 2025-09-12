package vee

import (
	"strings"
	"testing"
	"time"
)

func TestBind(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any                     // Function to create fresh target struct
		check   func(t *testing.T, target any) // Function to verify result
		wantErr bool
	}{
		{
			name: "bind simple string fields",
			input: map[string][]string{
				"name":  {"John Doe"},
				"email": {"john@example.com"},
			},
			target: func() any {
				return &struct {
					Name  string
					Email string
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name  string
					Email string
				})
				if s.Name != "John Doe" || s.Email != "john@example.com" {
					t.Errorf("Bind() result = %+v, want Name='John Doe', Email='john@example.com'", s)
				}
			},
			wantErr: false,
		},
		{
			name: "bind with custom field names",
			input: map[string][]string{
				"first_name": {"John"},
				"last_name":  {"Doe"},
			},
			target: func() any {
				return &struct {
					FirstName string `vee:"$first_name"`
					LastName  string `vee:"$last_name"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					FirstName string `vee:"$first_name"`
					LastName  string `vee:"$last_name"`
				})
				if s.FirstName != "John" || s.LastName != "Doe" {
					t.Errorf("Bind() result = %+v, want FirstName='John', LastName='Doe'", s)
				}
			},
			wantErr: false,
		},
		{
			name: "missing form fields leave struct unchanged",
			input: map[string][]string{
				"name": {"John"},
			},
			target: func() any {
				return &struct {
					Name  string
					Email string
				}{
					Email: "original@example.com",
				}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name  string
					Email string
				})
				if s.Name != "John" || s.Email != "original@example.com" {
					t.Errorf("Bind() result = %+v, want Name='John', Email='original@example.com'", s)
				}
			},
			wantErr: false,
		},
		{
			name: "skip fields are ignored",
			input: map[string][]string{
				"name":     {"John"},
				"internal": {"secret"},
			},
			target: func() any {
				return &struct {
					Name     string
					Internal string `vee:"-"`
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name     string
					Internal string `vee:"-"`
				})
				if s.Name != "John" || s.Internal != "" {
					t.Errorf("Bind() result = %+v, want Name='John', Internal=''", s)
				}
			},
			wantErr: false,
		},
		{
			name: "numeric fields bound correctly",
			input: map[string][]string{
				"name": {"John"},
				"age":  {"30"},
			},
			target: func() any {
				return &struct {
					Name string
					Age  int
				}{}
			},
			check: func(t *testing.T, target any) {
				s := target.(*struct {
					Name string
					Age  int
				})
				if s.Name != "John" || s.Age != 30 {
					t.Errorf("Bind() result = %+v, want Name='John', Age=30", s)
				}
			},
			wantErr: false,
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

func TestBindErrors(t *testing.T) {
	tests := []struct {
		name   string
		input  any
		target any
	}{
		{
			name:   "non-pointer target returns error",
			input:  map[string][]string{},
			target: struct{}{},
		},
		{
			name:   "wrong input type returns error",
			input:  "wrong type",
			target: &struct{}{},
		},
		{
			name:   "pointer to non-struct returns error",
			input:  map[string][]string{},
			target: new(string),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Bind(tt.input, tt.target)
			if err == nil {
				t.Errorf("Bind() expected error but got none")
			}
		})
	}
}

func TestBindParseErrors(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string][]string
		target  func() any
		wantErr string
	}{
		{
			name: "invalid int value returns parse error",
			input: map[string][]string{
				"age": {"not-a-number"},
			},
			target: func() any {
				return &struct {
					Age int
				}{}
			},
			wantErr: "cannot parse 'not-a-number' as integer for field 'age'",
		},
		{
			name: "invalid int64 value returns parse error",
			input: map[string][]string{
				"user_id": {"invalid"},
			},
			target: func() any {
				return &struct {
					UserID int64
				}{}
			},
			wantErr: "cannot parse 'invalid' as integer for field 'user_id'",
		},
		{
			name: "invalid float64 value returns parse error",
			input: map[string][]string{
				"price": {"abc"},
			},
			target: func() any {
				return &struct {
					Price float64
				}{}
			},
			wantErr: "cannot parse 'abc' as float for field 'price'",
		},
		{
			name: "invalid time value returns parse error",
			input: map[string][]string{
				"created": {"not-a-date"},
			},
			target: func() any {
				return &struct {
					Created time.Time
				}{}
			},
			wantErr: "cannot parse 'not-a-date' as time for field 'created'",
		},
		{
			name: "invalid duration value returns parse error",
			input: map[string][]string{
				"timeout": {"not-a-number"},
			},
			target: func() any {
				return &struct {
					Timeout time.Duration
				}{}
			},
			wantErr: "cannot parse 'not-a-number' as duration for field 'timeout'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := tt.target()
			err := Bind(tt.input, target)

			if err == nil {
				t.Errorf("Bind() expected error but got none")
				return
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Errorf("Bind() error = %v, want error containing %v", err, tt.wantErr)
			}
		})
	}
}
