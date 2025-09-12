# VEE Implementation Plan

This document outlines the iterative implementation approach for the VEE library, starting small and building up functionality in phases.

## Phase 1: Foundation (MVP) [DONE]
1. **Setup project structure** - go.mod, basic package
2. **Tag parser** - Parse simple `vee` tags (name extraction only)
3. **Basic string rendering** - Render `string` fields as `<input type="text">`
4. **Write tests** for tag parsing and basic rendering
5. **Basic binding** - Bind simple string fields from form data

## Phase 2: Form Configuration & Options [DONE]
6. **Form options** - Add RenderOptions struct for form-level configuration
7. **Default CSS styling** - Support default CSS classes for inputs
8. **Form wrapper** - Optional form tag with id, class, method attributes
9. **CSS tag support** - Simple pass-through to class attribute per field
10. **Write tests** for form options and CSS support

## Phase 3: Numeric Types [DONE]
11. **Add numeric types** - int, int64, float64 with `<input type="number">`
12. **Numeric attributes** - min, max, step support
13. **Improve binding** for numeric types
14. **Write tests** for numeric types

## Phase 4: Boolean Types [DONE]
15. **Add bool type** - checkbox rendering
16. **Boolean attributes** - checked state support
17. **Improve binding** for bool types
18. **Write tests** for bool types

## Phase 5: Time Types [DONE]
19. **Add time.Time support** - date/datetime-local/time inputs
20. **Add time.Duration support** - with units field
21. **Time attributes** - min/max date validation
22. **Improve binding** for time types
23. **Write tests** for time types

## Phase 5.5: Round-Trip Testing [DONE]
24. **Add round-trip tests** - Verify render→bind fidelity
25. **Test all supported types** - Ensure end-to-end functionality
26. **Test edge cases** - Zero values, empty forms, mixed types

## Phase 6: Attribute Parsing [DONE]
24. **Add attribute parsing** - required, label, placeholder, help
25. **Extend tag parser** to handle all attribute types
26. **Write tests** for attribute parsing
27. Update SPEC.md

## Phase 7: Basic Validation [DONE]
28. **Basic validation** - required fields, min/max constraints
29. **Validation functions** - Validate struct against rules
30. **Write comprehensive tests**
31. Update SPEC.md

## Phase 8: Default Field Processing [DONE]
32. **Change default behavior** - Process all public fields by default
33. **Update render/bind logic** - Skip only fields with vee:"-"
34. **Update all tests** - Remove unnecessary vee:"" tags
35. **Update SPEC.md** - Document new default behavior

## Phase 9: Pointer Type Support [DONE]
36. **Pointer type support** - Optional field rendering logic
37. **Update render/bind logic** - Handle *string, *int, *bool, etc.
38. **Write tests** for pointer types
39. Update SPEC.md

## Phase 10: Multi-value Support [DONE]
40. **Choices/Chosen pattern** - Support for select/radio/checkbox groups using {Name}Choices + {Name}Chosen convention
41. **Select element rendering** - Single and multiple select support
42. **Radio/checkbox groups** - Grouped input rendering
43. **Write tests** for multi-value support
44. **Update SPEC.md**

## Phase 11: Integration Tests
45. **Integration tests** with real HTTP requests
46. **End-to-end testing** - Full request/response cycle
47. **Real-world scenarios** - Complex forms and data
48. **Write comprehensive tests**
49. **Update SPEC.md**

## Phase 12: Error Handling & Display
50. **Error collection** - Collect and organize validation errors
51. **Error rendering** - Display errors in HTML forms
52. **Error styling** - CSS classes for error states
53. **Write comprehensive tests**
54. **Update SPEC.md**

## Approach

- Start with MVP and iterate
- Write tests for each phase before moving to the next
- Keep implementation simple and focused
- Build confidence through working code at each step