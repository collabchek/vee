package vee

// RenderOptions configures form rendering behavior.
type RenderOptions struct {
	// DefaultInputCSS sets default CSS classes for all input elements
	DefaultInputCSS string

	// FormID sets the HTML id attribute for the form wrapper
	FormID string

	// FormCSS sets CSS classes for the form wrapper
	FormCSS string

	// FormMethod sets the HTTP method for the form (defaults to "POST")
	FormMethod string

	// FormAction sets the action URL for the form
	FormAction string
}
