# CLAUDE.md

## Commonly Used Commands

- `go build` - Build the project
- `go test ./...` - Run all tests
- `go mod tidy` - Clean up dependencies
- `go vet ./...` - Run Go static analysis
- `go fmt ./...` - Format Go code
- `go test -v ./...` - Run all tests verbosely
- `go test -v -run=TestName` - Run a specific test by name

## Code Style
- Follow standard Go formatting conventions
- Group imports: standard library first, then third-party
- Use PascalCase for exported types/methods, camelCase for variables
- Add comments for public API and complex logic
- Place related functionality in logically named files

## Error Handling
- Use custom `Error` type with detailed context
- Include error wrapping with `Unwrap()` method
- Return errors with proper context information (line, position)

## Testing
- Write table-driven tests with clear input/output expectations
- Include detailed error messages (expected vs. actual)
- Test every exported function and error case

## Dependencies
- Minimum Go version: 1.25.0
- External dependencies managed through go modules

## Modernization Notes
- Use `errors.Is()` and `errors.As()` for error checking
- Replace `interface{}` with `any` type alias
- Replace type assertions with type switches where appropriate
- Use generics for type-safe operations
- Implement context cancellation handling for long operations
- Add proper docstring comments for exported functions and types
