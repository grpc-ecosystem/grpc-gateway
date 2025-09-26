package httprule

import (
	"fmt"
	"strings"
)

// validation is enabled by default
var (
	validationEnabled = true // Always enabled ---> only disable for exceptional cases
	defaultMaxStatic  = 4    // Conservative limit ---> only increase for exceptional cases
)

// URLPatternValidator validates URL patterns to prevent routing issues
type URLPatternValidator struct {
	MaxStaticSegmentsAtStart int
	AllowConsecutiveStatics  bool
}

// ValidationError represents a URL pattern validation error
type ValidationError struct {
	Message string
	Pattern string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("invalid URL pattern: %s (pattern: %s)", e.Message, e.Pattern)
}

// SetURLValidationEnabled allows disabling validation only for exceptional cases
// All patterns should follow the validation rules
func SetURLValidationEnabled(enabled bool) {
	validationEnabled = enabled
}

// SetMaxStaticSegments allows increasing the limit only for exceptional cases
// Note: The default limit of 4 prevents routing issues ---> only increase if absolutely necessary
func SetMaxStaticSegments(max int) {
	defaultMaxStatic = max
}

// ValidateTemplate validates a URL template ---> always enforced by default
// This prevents routing issues
// Maximum 4 static segments at the beginning
// No consecutive static segments after dynamic segments
func ValidateTemplate(template string) error {
	if !validationEnabled {
		// Only allowed for exceptional cases
		return nil
	}

	validator := &URLPatternValidator{
		MaxStaticSegmentsAtStart: defaultMaxStatic,
		AllowConsecutiveStatics:  false,
	}
	return validator.ValidateURLPattern(template)
}

// ValidateURLPattern validates a URL pattern according to the configured rules
func (v *URLPatternValidator) ValidateURLPattern(pattern string) error {
	if pattern == "" {
		return ValidationError{
			Message: "empty URL template",
			Pattern: pattern,
		}
	}

	if !strings.HasPrefix(pattern, "/") {
		return ValidationError{
			Message: "URL template must start with '/'",
			Pattern: pattern,
		}
	}

	// Remove leading slash and any verb suffix for analysis
	pathOnly := strings.TrimPrefix(pattern, "/")
	if colonIdx := strings.LastIndex(pathOnly, ":"); colonIdx > 0 {
		// Check if this colon is part of a verb (not inside a variable)
		if !strings.Contains(pathOnly[colonIdx:], "{") {
			pathOnly = pathOnly[:colonIdx]
		}
	}

	if pathOnly == "" {
		// Root path "/" is valid
		return nil
	}

	segments := v.parseSegments(pathOnly)
	return v.validateSegments(segments, pattern)
}

// parseSegments properly parses path segments, handling complex variable definitions
func (v *URLPatternValidator) parseSegments(path string) []string {
	var segments []string
	var current strings.Builder
	braceCount := 0

	for _, r := range path {
		if r == '{' {
			braceCount++
		} else if r == '}' {
			braceCount--
		}

		if r == '/' && braceCount == 0 {
			// Only split on '/' when we're not inside braces
			if current.Len() > 0 {
				segments = append(segments, current.String())
				current.Reset()
			}
		} else {
			current.WriteRune(r)
		}
	}

	// Add the last segment
	if current.Len() > 0 {
		segments = append(segments, current.String())
	}

	return segments
}

// validateSegments validates the sequence of path segments
func (v *URLPatternValidator) validateSegments(segments []string, originalPattern string) error {
	staticCount := 0
	foundDynamic := false
	lastWasStatic := false

	for i, segment := range segments {
		isStatic := v.isStaticSegment(segment)

		if isStatic {
			// Count consecutive static segments at the beginning
			if !foundDynamic {
				staticCount++
			}

			// Check for consecutive static segments after dynamic ones
			if foundDynamic && lastWasStatic && !v.AllowConsecutiveStatics {
				return ValidationError{
					Message: fmt.Sprintf("static segment '%s' followed by static segment '%s' at positions %d-%d",
						segments[i-1], segment, i, i+1),
					Pattern: originalPattern,
				}
			}

			lastWasStatic = true
		} else {
			foundDynamic = true
			lastWasStatic = false
		}
	}

	// Check if we exceed the maximum static segments at the start
	if staticCount > v.MaxStaticSegmentsAtStart {
		return ValidationError{
			Message: fmt.Sprintf("more than %d static segments at the beginning: found %d",
				v.MaxStaticSegmentsAtStart, staticCount),
			Pattern: originalPattern,
		}
	}

	return nil
}

// isStaticSegment determines if a segment is static (not a variable or wildcard)
func (v *URLPatternValidator) isStaticSegment(segment string) bool {
	// Handle special wildcards
	if segment == "*" || segment == "**" {
		return false
	}
	// Check if it's a variable ( is enclosed with these brackets ---> {} )
	if strings.HasPrefix(segment, "{") && strings.HasSuffix(segment, "}") {
		return false
	}

	return true
}

// IsValidURLPattern is a convenience function to check if a URL pattern is valid
func IsValidURLPattern(pattern string) bool {
	return ValidateTemplate(pattern) == nil
}

// GetValidationConfig returns the current global validation configuration
func GetValidationConfig() (enabled bool, maxStatic int) {
	return validationEnabled, defaultMaxStatic
}
