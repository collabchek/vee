# vee Library Specification

## Overview

vee is a bi-directional binding and rendering library for Go structs as HTML forms. It generates HTML forms from Go struct definitions and parses form data back into structs with validation.

## Design Principles

- **Framework Agnostic**: Works with standard `net/http` and any web framework
- **Simple**: Flat structs only, no nested structures
- **Convention-based**: Uses naming patterns for complex behaviors
- **Tag-driven**: Configuration through struct tags

## Supported Data Types

### Basic Types
- `string` → `<input type="text">`
- `int`, `int64` → `<input type="number">`
- `float64` → `<input type="number" step="any">`
- `bool` → `<input type="checkbox">`
- `time.Time` → `<input type="datetime-local">`
- `time.Duration` → `<input type="number">` (with units field)

### Pointer Types
Pointer types indicate optional fields and support all base types:
- `*string` → `<input type="text">` (empty value for nil)
- `*int`, `*int64` → `<input type="number">` (zero value for nil)
- `*float64` → `<input type="number" step="any">` (zero value for nil)
- `*bool` → `<input type="checkbox">` (unchecked for nil)
- `*time.Time` → `<input type="datetime-local">` (no value attribute for nil)
- `*time.Duration` → `<input type="number">` (no value attribute for nil)

**Rendering Behavior:**
- **Nil pointer**: Field rendered without value (or zero value for numeric types)
- **Non-nil pointer**: Field rendered with pointer's value
- **All pointer fields are always rendered** - nil vs non-nil only affects the value

## Pointer Type Behavior Details

Pointer types have subtle but important behavioral differences from their non-pointer equivalents, especially during form binding.

### String, Numeric, and Time Pointer Types

For `*string`, `*int`, `*int64`, `*float64`, `*time.Time`, and `*time.Duration`:

```go
type User struct {
    Name    string   // Regular field
    Email   *string  // Pointer field
    Age     int      // Regular field  
    Score   *float64 // Pointer field
}
```

**Form Binding Behavior:**
- **Field present in form data**: Creates new pointer and sets value
- **Field absent from form data**: Leaves pointer unchanged (preserves existing value)

```go
// Starting values
user := User{
    Name:  "John",
    Email: stringPtr("john@example.com"), 
    Age:   30,
    Score: nil,
}

// Form data: {"name": ["Jane"], "age": ["25"]}
// Email and Score are missing from form

vee.Bind(formData, &user)

// Result:
// user.Name = "Jane"        (updated from form)
// user.Email = "john@example.com" (unchanged - preserved)
// user.Age = 25             (updated from form)
// user.Score = nil          (unchanged - still nil)
```

### Boolean Pointer Types

Boolean pointer types (`*bool`) follow HTML checkbox semantics, which are different:

```go
type Settings struct {
    IsActive   bool   // Regular checkbox
    IsOptional *bool  // Pointer checkbox
}
```

**Form Binding Behavior:**
- **Checkbox checked (field present)**: Sets pointer to `&true`
- **Checkbox unchecked (field absent)**: Sets pointer to `&false`

**Important:** Unlike other pointer types, boolean pointers do NOT preserve existing values when absent from form data. This follows standard HTML checkbox behavior where unchecked boxes don't send any data.

```go
// Starting values
settings := Settings{
    IsActive:   true,
    IsOptional: boolPtr(true), // Previously checked
}

// Form data: {"is_active": ["true"]}
// is_optional is missing (checkbox was unchecked)

vee.Bind(formData, &settings)

// Result:
// settings.IsActive = true     (checked - present in form)
// settings.IsOptional = &false (unchecked - absent from form)
```

### Why Boolean Pointers Behave Differently

HTML checkboxes have unique behavior:
- **Checked checkbox**: Browser sends field in form data
- **Unchecked checkbox**: Browser sends NO data for that field

For regular fields, "no data" means "don't change the value". But for checkboxes, "no data" explicitly means "unchecked" (false). This creates the behavioral difference between boolean pointers and other pointer types.

### Practical Use Cases

**Optional String Fields:**
```go
type Profile struct {
    Name     string  `vee:"required"`           // Always required
    Bio      *string `vee:"placeholder:'Optional bio'"` // Optional
    Website  *string `vee:"type:'url'"`         // Optional URL
}
```

**Optional Numeric Fields:**
```go
type Product struct {
    Name  string   `vee:"required"`
    Price float64  `vee:"required,min:0"`
    Sale  *float64 `vee:"min:0"`  // Optional sale price
}
```

**Settings with Optional Toggles:**
```go
type UserSettings struct {
    EmailNotifications bool   // Default behavior
    SmsNotifications   *bool  // Optional setting (nil = not configured)
}
```

## Tag Syntax

### vee Tag Format
```go
`vee:"[${override_name},]param1,param2:value,param3:'quoted string'"`
```

### Name Override
Override the HTML form field name using `$` prefix as the first parameter:
```go
FirstName string `vee:"$firstName,required,label:'First Name'"`
```

### String Values
All string values must be wrapped in single quotes:
```go
Name string `vee:"label:'Full Name',placeholder:'Enter your name',help:'This is required'"`
```

### CSS Tag
CSS classes are passed directly to the HTML `class` attribute:
```go
Name string `vee:"required" css:"border-2 border-gray-300 rounded px-3 py-2"`
```

## Universal Attributes

Available for all field types:
- `required` - Adds HTML `required` attribute for client-side validation (see Validation section for server-side validation)
- `readonly` - Field is read-only
- `disabled` - Field is disabled  
- `hidden` - Renders as `<input type="hidden">` without label (not supported for pointer types or multi-value fields)
- `label:'Text'` - Custom label text (defaults to human-readable field name)
- `nolabel` - Skip automatic label generation
- `placeholder:'Text'` - Placeholder text (forces rendering for pointer types)
- `help:'Text'` - Help/description text
- `id:'custom_id'` - Custom HTML id (always defaults to field name if not specified)

## Type-Specific Attributes

### String Fields
```go
Name string `vee:"type:'email'"`
```
- `type:'email|password|tel|url'` - HTML input type override

### Numeric Fields
```go
Age int `vee:"step:1"`
Price float64 `vee:"step:0.01"`
```
- `step:N` - Step increment for HTML input

### Boolean Fields
```go
Active bool `vee:"label:'Is Active'"`
```
**Note**: Boolean fields are rendered as checkboxes with `value="true"`. The `checked` attribute is set based on the struct field value - no tag override is needed.

### Time Fields
```go
Birthday time.Time `vee:"type:'date'"`
```
- `type:'date|datetime-local|time'` - HTML input type (defaults to datetime-local)

### Duration Fields
```go
Timeout time.Duration `vee:"units:'s',label:'Timeout'"`
```
- `units:'ms|s|m|h'` - Duration units (milliseconds, seconds, minutes, hours, defaults to seconds)

**Rendering:** Creates a number input with the value converted to the specified units.
**Binding:** Converts the number back to `time.Duration` using the units.

## Validation

vee integrates with [go-playground/validator](https://github.com/go-playground/validator) for validation. Use standard `validate` tags alongside `vee` tags:

```go
type User struct {
    Name  string `vee:"required" validate:"required,min=2,max=50"`
    Email string `vee:"type:'email',required" validate:"required,email"`
    Age   int    `validate:"required,gte=18,lte=120"`
    Phase int    `vee:"hidden" validate:"required"`
}

// Validate the struct
user := User{Name: "John", Email: "john@example.com", Age: 25, Phase: 1}
if err := vee.Validate(user); err != nil {
    // Handle validation errors
}
```

**Available Functions:**
- `vee.Validate(struct)` - Validates a struct using validator tags
- `vee.ValidateVar(value, tag)` - Validates a single value

**Important Distinction:**
- **vee's `required` attribute**: Only affects HTML form generation by adding the `required` attribute to input elements for client-side validation
- **Validator's `required` tag**: Handles actual server-side validation logic

```go
type Examples struct {
    // Client + server validation
    Name string `vee:"required" validate:"required"`
    
    // Only client-side (HTML required attribute)
    Email string `vee:"required"`
    
    // Only server-side validation
    Age int `validate:"required"`
    
    // Hidden field with server validation but no HTML required
    Phase int `vee:"hidden" validate:"required"`
}
```

**Note:** vee handles form rendering and binding, while validator handles validation logic. This separation keeps each library focused on its strengths.

## Multi-Value Fields (Dropdowns/Selects)

Use convention-based paired fields: `{Name}Choices` + `{Name}Chosen`

### Single Selection
```go
type User struct {
    ColorChoices []string  // ["Red", "Blue", "Green"] - not rendered
    ColorChosen  int       `vee:"type:'select',label:'Favorite Color'"` // renders as <select>
}
```

### Multiple Selection
```go
type User struct {
    SkillChoices []string  // ["Go", "JavaScript", "Python"]
    SkillChosen  []int     `vee:"type:'select',multiple,label:'Skills'"` // multi-select
    
    InterestChoices []string
    InterestChosen  []int  `vee:"type:'checkbox',label:'Interests'"` // checkbox group
}
```

### Input Type Options

**Select Dropdown (default):**
```go
ColorChosen int `vee:"type:'select'"` // Single select dropdown
SkillChosen []int `vee:"type:'select',multiple"` // Multi-select dropdown
```

**Radio Button Group:**
```go
SizeChosen int `vee:"type:'radio'"` // Radio buttons (single-select only)
```

**Checkbox Group:**
```go
FeatureChosen []int `vee:"type:'checkbox'"` // Checkbox group (multi-select only)
```

### Convention Validation

vee enforces strict conventions for multi-value fields:

- **Paired fields required**: Every `{Name}Choices` must have a corresponding `{Name}Chosen`
- **Choices field type**: Must be `[]string` or slice of any type implementing `String()`
- **Chosen field type**: Must be `int` (single-select) or `[]int` (multi-select)
- **Index validation**: All chosen indices must be within range of available choices
- **Non-empty choices**: Choices slice cannot be empty
- **Form binding validation**: Invalid form indices return binding errors

**Validation Errors:**
```go
// ❌ Missing Chosen field
type User struct {
    ColorChoices []string // Error: requires ColorChosen
}

// ❌ Wrong Chosen type  
type User struct {
    ColorChoices []string
    ColorChosen  string   // Error: must be int or []int
}

// ❌ Index out of range (during rendering)
user := User{
    ColorChoices: []string{"Red", "Blue"},
    ColorChosen:  5, // Error: index 5 out of range for 2 choices
}

// ❌ Invalid form data (during binding)
formData := map[string][]string{
    "color_chosen": {"5"}, // Error: index 5 out of range for 2 choices
    // or
    "color_chosen": {"invalid"}, // Error: invalid index 'invalid' for field 'color_chosen'
}
```

### Custom Types
Choices can be any type implementing `String()` method:
```go
type Status int
func (s Status) String() string { return "..." }

type User struct {
    StatusChoices []Status
    StatusChosen  int `vee:"type:'select',label:'Status'"`
}
```

## Label Generation

**Default Behavior**: vee automatically generates `<label>` elements for all form fields to improve accessibility and usability.

### Label Text Generation

Labels are generated using this priority order:
1. **Custom label**: Use `label:'Custom Text'` attribute if specified
2. **Human-readable field name**: Convert field name from CamelCase to spaced text

```go
type User struct {
    Name         string // Label: "Name"  
    FirstName    string // Label: "First Name"
    EmailAddress string // Label: "Email Address"  
    IsActive     bool   // Label: "Is Active"
}
```

### Label-Input Association

Labels are properly associated with inputs using the `for` attribute:
```html
<label for="field_id">Field Label</label>
<input type="text" name="field_name" id="field_id" ...>
```

### Customizing Labels

**Custom Label Text:**
```go
Name string `vee:"label:'Full Name'"` 
```

**Custom Label CSS:**
```go
Name string `vee:"required" css:"border-2 border-gray-300 rounded px-3 py-2" labelCss:"font-bold"`
```

**Skip Label Generation:**
```go
Password string `vee:"type:'password',nolabel"`
```

### Multi-Value Field Labels

- **Select dropdowns**: Get a standard `<label>` element
- **Radio/checkbox groups**: Wrapped in `<fieldset><legend>` for semantic grouping

```go
type Form struct {
    ColorChoices []string
    ColorChosen  int    `vee:"type:'select',label:'Favorite Color'"`   // <label>
    
    SizeChoices []string  
    SizeChosen  int       `vee:"type:'radio',label:'Size'"`             // <fieldset><legend>
}
```

## Field Processing

**Default Behavior**: All public struct fields are processed automatically. You don't need to add `vee` tags unless you want to customize field behavior.

```go
type User struct {
    Name  string    // Processed with auto-derived name "name", label "Name"
    Email string    // Processed with auto-derived name "email", label "Email"  
    Age   int       // Processed with auto-derived name "age", label "Age"
}
```

## Skip Fields

Use `vee:"-"` to skip fields during rendering and binding:
```go
type User struct {
    Name     string  // Processed (no tag needed)
    Email    string `vee:"type:'email'"` // Processed with custom type
    Internal string `vee:"-"`           // Skipped
}
```

## Rendering

### Basic Rendering

```go
// Simple rendering - generates form HTML from struct
html, err := vee.Render(user)
```

### Render Options

vee provides flexible rendering options through the `RenderOption` type and helper functions:

#### Form Configuration

**Form Method and Action:**
```go
html, err := vee.Render(user, 
    vee.FormMethodOption("POST"),
    vee.FormActionOption("/submit-user"),
)
```

**Client-Side JavaScript Forms:**
For forms intended for client-side JavaScript handling, use `FormActionScriptOption()`:
```go
html, err := vee.Render(user, vee.FormActionScriptOption())
// Generates: <form> (no method or action attributes)
```

This creates a "pure" form by omitting both `method` and `action` attributes, preventing the browser from navigating away when the form is submitted. This is ideal for:
- AJAX form submissions
- Single-page applications (SPAs) 
- Client-side form validation and processing
- Progressive web apps

The form relies entirely on JavaScript event handlers (like `onsubmit`) for processing.

**Form ID and CSS:**
```go
html, err := vee.Render(user,
    vee.FormIDOption("user-form"),
    vee.FormCSSOption("max-w-md mx-auto p-6 bg-white rounded shadow"),
)
```

#### Default CSS Styling

Apply default CSS classes to all inputs and labels:

```go
html, err := vee.Render(user,
    vee.InputCSSOption("border border-gray-300 rounded px-3 py-2 w-full"),
    vee.LabelCSSOption("block text-sm font-medium text-gray-700 mb-1"),
)
```

#### Combining Options

Multiple render options can be combined:

```go
html, err := vee.Render(user,
    vee.FormIDOption("registration-form"),
    vee.FormMethodOption("POST"),
    vee.FormActionOption("/register"),
    vee.FormCSSOption("space-y-4"),
    vee.InputCSSOption("border border-gray-300 rounded-md px-3 py-2"),
    vee.LabelCSSOption("block text-sm font-medium mb-1"),
)
```

#### Option Priority

Field-specific CSS tags override default options:

```go
type User struct {
    Name  string `css:"border-red-500"`  // Overrides InputCSSOption
    Email string                         // Uses InputCSSOption
}

html, err := vee.Render(user, vee.InputCSSOption("border-gray-300"))
// Name field gets "border-red-500", Email field gets "border-gray-300"
```

### Render Option Reference

| Function | Purpose | Default |
|----------|---------|---------|
| `FormMethodOption(method)` | Sets form HTTP method | "POST" |
| `FormActionOption(action)` | Sets form action URL | "" |
| `FormActionScriptOption()` | Sets form action to "script" for JS handling | - |
| `FormIDOption(id)` | Sets form HTML id | "" |
| `FormCSSOption(css)` | Sets form CSS classes | "" |
| `InputCSSOption(css)` | Default CSS for all inputs | "" |
| `LabelCSSOption(css)` | Default CSS for all labels | "" |

## Example Usage

```go
type User struct {
    Name         string     `vee:"required,label:'Full Name'" css:"border rounded px-3 py-2"`
    Email        string     `vee:"$userEmail,type:'email',required" css:"w-full"`
    Age          int        `vee:"min:18,max:120"`
    Bio          *string    `vee:"placeholder:'Tell us about yourself'" css:"h-24"`
    Website      *string    `vee:"type:'url'"`        // Optional URL field
    Score        *float64   `vee:"min:0,max:100"`     // Optional score field
    Active       bool       `vee:"label:'Account Active'"`
    Birthday     time.Time  `vee:"type:'date'"`
    
    ColorChoices []string   // ["Red", "Blue", "Green"]
    ColorChosen  int        `vee:"type:'select',label:'Favorite Color'"`
}

// Basic rendering
html, err := vee.Render(User{
    ColorChoices: []string{"Red", "Blue", "Green"},
    ColorChosen:  1, // "Blue" selected
})

// Styled rendering with options
html, err := vee.Render(User{
    ColorChoices: []string{"Red", "Blue", "Green"},
    ColorChosen:  1,
}, 
    vee.FormIDOption("user-form"),
    vee.FormActionOption("/users"),
    vee.InputCSSOption("border border-gray-300 rounded px-3 py-2"),
    vee.LabelCSSOption("block font-medium text-gray-700 mb-1"),
)

// Bind form data from HTTP request (recommended)
var user User
err := vee.BindRequest(r, &user) // r is *http.Request

// Or bind from form data directly
err = vee.Bind(r.Form, &user)           // url.Values
err = vee.Bind(formData, &user)         // map[string][]string
```

## Form Data Binding

vee provides two functions for binding HTTP form data to structs:

### BindRequest (Recommended)

```go
func BindRequest(r *http.Request, v any) error
```

**Most convenient approach** - automatically handles form parsing:

```go
// In your HTTP handler
func handleRegistration(w http.ResponseWriter, r *http.Request) {
    var registration AccountRegistration

    err := vee.BindRequest(r, &registration)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Use populated registration struct...
}
```

**Features:**
- Automatically calls `r.ParseForm()`
- Handles both GET query parameters and POST form data
- Returns parsing errors if form parsing fails
- Supports all vee field types and validation

### Bind (Direct)

```go
func Bind(formData any, v any) error
```

**Lower-level approach** for direct form data binding:

```go
// Manual form parsing
if err := r.ParseForm(); err != nil {
    return err
}

var registration AccountRegistration

// Both work:
err := vee.Bind(r.Form, &registration)      // url.Values
err := vee.Bind(formData, &registration)    // map[string][]string
```

**Accepts:**
- `url.Values` (from `r.Form`, `r.PostForm`, or `r.URL.Query()`)
- `map[string][]string` (custom form data)

**Use Cases:**
- Custom form data processing
- Testing with mock data
- Integration with other form parsing libraries

## Implementation Notes

- **Field Processing**: All public struct fields are processed by default - no `vee` tags required unless customizing behavior
- Framework agnostic - works with any `http.Request`
- No nested struct support
- Limited type support for simplicity: `string`, `int`, `int64`, `float64`, `bool`, `time.Time`, `time.Duration` and their pointer equivalents
- **Pointer support**: All base types support pointer variants (`*string`, `*int`, etc.)
- **Pointer rendering**: Nil pointers render with empty/zero values, non-nil render with actual values
- **Pointer binding**: Form data presence creates new pointer with parsed value, absence leaves field nil
- **Multi-value support**: Choices/Chosen convention for select dropdowns, radio groups, and checkbox groups
- **Multi-value validation**: Strict validation of field pairs, types, and index ranges
- Form data binding uses built-in `strconv` package for type conversion
- Invalid numeric values are silently ignored (fields remain unchanged)
- **Boolean checkbox binding**: Presence in form data sets field to `true`, absence sets to `false` (standard checkbox behavior)
- **Time field binding**: Supports `date` (2006-01-02), `time` (15:04), and `datetime-local` (2006-01-02T15:04) formats
- **Duration field binding**: Converts between numeric input and `time.Duration` using configurable units (ms/s/m/h), defaults to seconds