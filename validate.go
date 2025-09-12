package vee

import (
	"github.com/go-playground/validator/v10"
)

// validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate validates a struct using go-playground/validator.
// It works alongside VEE tags - you can use both vee and validate tags on the same struct.
//
// Example:
//
//	type User struct {
//	    Name  string `vee:"required" validate:"required,min=2,max=50"`
//	    Email string `vee:"type:'email',required" validate:"required,email"`
//	    Age   int    `vee:"min:18,max:120" validate:"required,gte=18,lte=120"`
//	}
//
//	user := User{Name: "John", Email: "john@example.com", Age: 25}
//	if err := vee.Validate(user); err != nil {
//	    // Handle validation errors
//	}
func Validate(s any) error {
	return validate.Struct(s)
}

// ValidateVar validates a single variable using validation tags.
// This is useful for validating individual values outside of structs.
//
// Example:
//
//	email := "invalid-email"
//	if err := vee.ValidateVar(email, "required,email"); err != nil {
//	    // Handle validation error
//	}
func ValidateVar(field any, tag string) error {
	return validate.Var(field, tag)
}
