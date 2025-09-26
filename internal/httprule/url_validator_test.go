package httprule

import (
	"testing"
)

func TestURLPatternValidator_ValidateURLPattern(t *testing.T) {
	validator := &URLPatternValidator{
		MaxStaticSegmentsAtStart: 4,
		AllowConsecutiveStatics:  false,
	}

	testCases := []struct {
		name        string
		pattern     string
		expectError bool
		errorMsg    string
	}{
		// Valid patterns
		{
			name:        "Root path",
			pattern:     "/",
			expectError: false,
		},
		{
			name:        "Simple dynamic pattern",
			pattern:     "/{id}",
			expectError: false,
		},
		{
			name:        "Valid pattern with 4 static segments",
			pattern:     "/api/compute/v1/instances/{instance_id}",
			expectError: false,
		},
		{
			name:        "Valid pattern with mixed segments",
			pattern:     "/api/{service}/v1/{resource}",
			expectError: false,
		},
		{
			name:        "Valid pattern with wildcards",
			pattern:     "/api/*/v1/{resource}",
			expectError: false,
		},
		{
			name:        "Valid pattern with deep wildcards",
			pattern:     "/api/**/resources",
			expectError: false,
		},
		{
			name:        "Valid pattern with verb",
			pattern:     "/api/compute/v1/{instance}:start",
			expectError: false,
		},
		{
			name:        "Dynamic followed by static",
			pattern:     "/{service}/static/{resource}",
			expectError: false,
		},
		{
			name:        "Exactly 4 static segments at start",
			pattern:     "/a/b/c/d/{dynamic}",
			expectError: false,
		},

		// Invalid patterns - empty or malformed
		{
			name:        "Empty pattern",
			pattern:     "",
			expectError: true,
			errorMsg:    "empty URL template",
		},
		{
			name:        "No leading slash",
			pattern:     "api/compute",
			expectError: true,
			errorMsg:    "URL template must start with '/'",
		},

		// Invalid patterns - too many static segments at start
		{
			name:        "5 static segments at beginning",
			pattern:     "/api/mytenant/v1/idp/types/{type_id}",
			expectError: true,
			errorMsg:    "more than 4 static segments at the beginning: found 5",
		},
		{
			name:        "6 static segments at beginning",
			pattern:     "/a/b/c/d/e/f/{dynamic}",
			expectError: true,
			errorMsg:    "more than 4 static segments at the beginning: found 6",
		},

		// Invalid patterns - consecutive static segments after dynamic
		{
			name:        "Consecutive static after dynamic",
			pattern:     "/api/{service}/static1/static2/{resource}",
			expectError: true,
			errorMsg:    "static segment 'static1' followed by static segment 'static2' at positions 3-4",
		},
		{
			name:        "Multiple consecutive static after dynamic",
			pattern:     "/{service}/a/b/c/{resource}",
			expectError: true,
			errorMsg:    "static segment 'a' followed by static segment 'b' at positions 2-3",
		},

		// Edge cases with verbs
		{
			name:        "Invalid pattern with verb - too many statics",
			pattern:     "/api/mytenant/v1/idp/types/{type}:action",
			expectError: true,
			errorMsg:    "more than 4 static segments at the beginning: found 5",
		},
		{
			name:        "Valid complex pattern with variables and verb",
			pattern:     "/api/{service}/v1/{resource=projects/*/instances/*}:update",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateURLPattern(tc.pattern)
			
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for pattern %q, but got none", tc.pattern)
					return
				}
				if tc.errorMsg != "" && !contains(err.Error(), tc.errorMsg) {
					t.Errorf("Expected error message to contain %q, but got %q", tc.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for pattern %q: %v", tc.pattern, err)
				}
			}
		})
	}
}

func TestURLPatternValidator_AllowConsecutiveStatics(t *testing.T) {
	// Test validator that allows consecutive static segments
	validator := &URLPatternValidator{
		MaxStaticSegmentsAtStart: 4,
		AllowConsecutiveStatics:  true,
	}

	testCases := []struct {
		name        string
		pattern     string
		expectError bool
	}{
		{
			name:        "Consecutive static after dynamic - allowed",
			pattern:     "/api/{service}/static1/static2/{resource}",
			expectError: false,
		},
		{
			name:        "Multiple consecutive static - allowed",
			pattern:     "/{service}/a/b/c/{resource}",
			expectError: false,
		},
		{
			name:        "Still fails on too many static at start",
			pattern:     "/a/b/c/d/e/{dynamic}",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validator.ValidateURLPattern(tc.pattern)
			
			if tc.expectError && err == nil {
				t.Errorf("Expected error for pattern %q, but got none", tc.pattern)
			} else if !tc.expectError && err != nil {
				t.Errorf("Unexpected error for pattern %q: %v", tc.pattern, err)
			}
		})
	}
}

func TestValidateTemplate_GlobalConfiguration(t *testing.T) {
	// Save original configuration
	originalEnabled, originalMaxStatic := GetValidationConfig()
	defer func() {
		SetURLValidationEnabled(originalEnabled)
		SetMaxStaticSegments(originalMaxStatic)
	}()

	// Test with validation disabled (should only be used for exceptional cases)
	SetURLValidationEnabled(false)
	err := ValidateTemplate("/a/b/c/d/e/f/{dynamic}")
	if err != nil {
		t.Errorf("Expected no error when validation disabled (exceptional case), got: %v", err)
	}

	// Test with validation enabled and custom max static (exceptional case)
	SetURLValidationEnabled(true)
	SetMaxStaticSegments(2)
	
	err = ValidateTemplate("/a/b/c/{dynamic}")
	if err == nil {
		t.Error("Expected error with max static segments = 2, but got none")
	}

	err = ValidateTemplate("/a/b/{dynamic}")
	if err != nil {
		t.Errorf("Unexpected error with valid pattern: %v", err)
	}
}

func TestIsValidURLPattern(t *testing.T) {
	testCases := []struct {
		pattern string
		valid   bool
	}{
		{"/api/compute/v1/{instance}", true},
		{"/api/mytenant/v1/idp/types/{type}", false},
		{"", false},
		{"/", true},
		{"/api/{service}/static1/static2/{resource}", false},
	}

	for _, tc := range testCases {
		t.Run(tc.pattern, func(t *testing.T) {
			result := IsValidURLPattern(tc.pattern)
			if result != tc.valid {
				t.Errorf("Expected IsValidURLPattern(%q) = %v, got %v", tc.pattern, tc.valid, result)
			}
		})
	}
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Message: "too many static segments",
		Pattern: "/a/b/c/d/e/{id}",
	}

	expected := "invalid URL pattern: too many static segments (pattern: /a/b/c/d/e/{id})"
	if err.Error() != expected {
		t.Errorf("Expected error message %q, got %q", expected, err.Error())
	}
}

func TestGetValidationConfig(t *testing.T) {
	// Save original configuration
	originalEnabled, originalMaxStatic := GetValidationConfig()
	defer func() {
		SetURLValidationEnabled(originalEnabled)
		SetMaxStaticSegments(originalMaxStatic)
	}()

	// Test setting and getting configuration
	SetURLValidationEnabled(false)
	SetMaxStaticSegments(10)

	enabled, maxStatic := GetValidationConfig()
	if enabled != false || maxStatic != 10 {
		t.Errorf("Expected config (enabled=false, maxStatic=10), got (enabled=%v, maxStatic=%v)", enabled, maxStatic)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(substr) == 0 || (len(s) >= len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Benchmark tests
func BenchmarkValidateTemplate(b *testing.B) {
	patterns := []string{
		"/api/compute/v1/{instance}",
		"/api/mytenant/v1/idp/types/{type}",
		"/api/{service}/v1/{resource}",
		"/{service}/static1/static2/{resource}",
	}

	for i := 0; i < b.N; i++ {
		for _, pattern := range patterns {
			ValidateTemplate(pattern)
		}
	}
}