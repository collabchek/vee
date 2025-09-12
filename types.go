package vee

// RenderOption configures form rendering behavior.
type RenderOption struct {
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

func InputCSSOption(css string) RenderOption {
	return RenderOption{
		DefaultInputCSS: css,
	}
}

func FormIDOption(id string) RenderOption {
	return RenderOption{
		FormID: id,
	}
}

func FormCSSOption(css string) RenderOption {
	return RenderOption{
		FormCSS: css,
	}
}

func FormMethodOption(method string) RenderOption {
	return RenderOption{
		FormMethod: method,
	}
}

func FormActionOption(action string) RenderOption {
	return RenderOption{
		FormAction: action,
	}
}

func (option RenderOption) IsEqual(other RenderOption) bool {
	return option.DefaultInputCSS == other.DefaultInputCSS &&
		option.FormAction == other.FormAction &&
		option.FormCSS == other.FormCSS &&
		option.FormID == other.FormID &&
		option.FormMethod == other.FormMethod
}

func (option *RenderOption) apply(other RenderOption) {
	if other.DefaultInputCSS != "" {
		option.DefaultInputCSS = other.DefaultInputCSS
	}
	if other.FormID != "" {
		option.FormID = other.FormID
	}
	if other.FormCSS != "" {
		option.FormCSS = other.FormCSS
	}
	if other.FormMethod != "" {
		option.FormMethod = other.FormMethod
	}
	if other.FormAction != "" {
		option.FormAction = other.FormAction
	}
}

func ConsolidateOptions(opts ...RenderOption) *RenderOption {
	target := &RenderOption{}
	for _, opt := range opts {
		target.apply(opt)
	}
	return target
}
