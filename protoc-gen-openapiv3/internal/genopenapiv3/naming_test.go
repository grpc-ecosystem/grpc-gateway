package genopenapiv3

import "testing"

func TestNaming(t *testing.T) {
	type expectedNames struct {
		fqn, legacy, simple, pkg string
	}
	messageNameToExpected := map[string]expectedNames{
		".A":     {"A", "A", "A", "A"},
		".a.B.C": {"a.B.C", "aBC", "B.C", "B.C"},
		".a.D.C": {"a.D.C", "aDC", "D.C", "D.C"},
		".a.E.F": {"a.E.F", "aEF", "a.E.F", "a.E.F"},
		".b.E.F": {"b.E.F", "bEF", "b.E.F", "b.E.F"},
		".c.G.H": {"c.G.H", "GH", "H", "G.H"},
	}

	allMessageNames := make([]string, 0, len(messageNameToExpected))
	for msgName := range messageNameToExpected {
		allMessageNames = append(allMessageNames, msgName)
	}

	t.Run("fqn", func(t *testing.T) {
		uniqueNames := resolveNamesFQN(allMessageNames)
		for _, msgName := range allMessageNames {
			expected := messageNameToExpected[msgName].fqn
			actual := uniqueNames[msgName]
			if expected != actual {
				t.Errorf("fqn unique name %q does not match expected name %q", actual, expected)
			}
		}
	})
	t.Run("legacy", func(t *testing.T) {
		uniqueNames := resolveNamesLegacy(allMessageNames)
		for _, msgName := range allMessageNames {
			expected := messageNameToExpected[msgName].legacy
			actual := uniqueNames[msgName]
			if expected != actual {
				t.Errorf("legacy unique name %q does not match expected name %q", actual, expected)
			}
		}
	})
	t.Run("simple", func(t *testing.T) {
		uniqueNames := resolveNamesSimple(allMessageNames)
		for _, msgName := range allMessageNames {
			expected := messageNameToExpected[msgName].simple
			actual := uniqueNames[msgName]
			if expected != actual {
				t.Errorf("simple unique name %q does not match expected name %q", actual, expected)
			}
		}
	})
	t.Run("package", func(t *testing.T) {
		uniqueNames := resolveNamesPackage(allMessageNames)
		for _, msgName := range allMessageNames {
			expected := messageNameToExpected[msgName].pkg
			actual := uniqueNames[msgName]
			if expected != actual {
				t.Errorf("package unique name %q does not match expected name %q", actual, expected)
			}
		}
	})
}

func TestLookupNamingStrategy(t *testing.T) {
	tests := []struct {
		name     string
		strategy string
		wantNil  bool
	}{
		{"fqn", "fqn", false},
		{"FQN uppercase", "FQN", false},
		{"legacy", "legacy", false},
		{"simple", "simple", false},
		{"package", "package", false},
		{"invalid", "invalid", true},
		{"empty", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := LookupNamingStrategy(tt.strategy)
			if (result == nil) != tt.wantNil {
				t.Errorf("LookupNamingStrategy(%q) returned nil=%v, want nil=%v", tt.strategy, result == nil, tt.wantNil)
			}
		})
	}
}

func TestResolveFullyQualifiedNameToOpenAPINames(t *testing.T) {
	fqns := []string{".test.Message", ".test.AnotherMessage"}

	// Test with valid strategy
	result := resolveFullyQualifiedNameToOpenAPINames(fqns, "fqn")
	if result[".test.Message"] != "test.Message" {
		t.Errorf("expected 'test.Message', got %q", result[".test.Message"])
	}

	// Test with invalid strategy (should fall back to fqn)
	result = resolveFullyQualifiedNameToOpenAPINames(fqns, "invalid")
	if result[".test.Message"] != "test.Message" {
		t.Errorf("expected fallback to fqn strategy, got %q", result[".test.Message"])
	}
}
