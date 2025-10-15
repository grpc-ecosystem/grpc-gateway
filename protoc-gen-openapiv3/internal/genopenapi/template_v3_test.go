package genopenapi

import (
	"reflect"
	"sort"
	"testing"
)

// Mock descriptor.Field type to simulate a protobuf field
type MockField struct {
	Name string
}

func (f *MockField) GetName() string {
	return f.Name
}

func Test_generateOneOfCombinations(t *testing.T) {
	t.Run("NoOneOfGroups", func(t *testing.T) {
		oneofGroups := map[string][]MockField{}
		result := generateOneOfCombinations(oneofGroups)

		if len(result) != 1 {
			t.Fatalf("Expected 1 combination, got %d", len(result))
		}
		if len(result[0]) != 0 {
			t.Fatalf("Expected an empty map, got a map with %d elements", len(result[0]))
		}
	})

	t.Run("SingleOneOfGroup", func(t *testing.T) {
		oneofGroups := map[string][]MockField{
			"oneof_group_A": {
				{Name: "field_A1"},
				{Name: "field_A2"},
			},
		}

		result := generateOneOfCombinations(oneofGroups)
		if len(result) != 2 {
			t.Fatalf("Expected 2 combinations, got %d", len(result))
		}

		var foundFieldNames []string
		for _, combination := range result {
			foundFieldNames = append(foundFieldNames, combination["oneof_group_A"].Name)
		}
		sort.Strings(foundFieldNames)
		expectedFieldNames := []string{"field_A1", "field_A2"}

		if !reflect.DeepEqual(foundFieldNames, expectedFieldNames) {
			t.Errorf("Field names do not match. Expected %+v, got %+v", expectedFieldNames, foundFieldNames)
		}
	})
	// This tests the Cartesian product logic.
	t.Run("MultipleOneOfGroups", func(t *testing.T) {
		oneofGroups := map[string][]MockField{
			"oneof_group_A": {
				{Name: "field_A1"},
				{Name: "field_A2"},
			},
			"oneof_group_B": {
				{Name: "field_B1"},
				{Name: "field_B2"},
			},
		}

		result := generateOneOfCombinations(oneofGroups)
		// 2 variants * 2 variants = 4 combinations expected
		if len(result) != 4 {
			t.Fatalf("Expected 4 combinations, got %d", len(result))
		}

		// Check the specific combinations
		expectedCombinations := []map[string]string{
			{"oneof_group_A": "field_A1", "oneof_group_B": "field_B1"},
			{"oneof_group_A": "field_A1", "oneof_group_B": "field_B2"},
			{"oneof_group_A": "field_A2", "oneof_group_B": "field_B1"},
			{"oneof_group_A": "field_A2", "oneof_group_B": "field_B2"},
		}

		// Convert the result to a comparable format and sort for stable comparison.
		foundCombinations := make([]map[string]string, len(result))
		for i, combination := range result {
			foundCombinations[i] = make(map[string]string)
			for k, v := range combination {
				foundCombinations[i][k] = v.Name
			}
		}

		// Sort both slices for consistent comparison
		sort.Slice(foundCombinations, func(i, j int) bool {
			if foundCombinations[i]["oneof_group_A"] != foundCombinations[j]["oneof_group_A"] {
				return foundCombinations[i]["oneof_group_A"] < foundCombinations[j]["oneof_group_A"]
			}
			return foundCombinations[i]["oneof_group_B"] < foundCombinations[j]["oneof_group_B"]
		})

		if !reflect.DeepEqual(foundCombinations, expectedCombinations) {
			t.Errorf("Combinations do not match expected result.\nExpected: %+v\nGot: %+v", expectedCombinations, foundCombinations)
		}
	})

	t.Run("MultipleOneOfGroupsWithDifferentVariantNumbers", func(t *testing.T) {
		oneofGroups := map[string][]MockField{
			"oneof_group_A": {
				{Name: "field_A1"},
				{Name: "field_A2"},
				{Name: "field_A3"},
			},
			"oneof_group_B": {
				{Name: "field_B1"},
				{Name: "field_B2"},
			},
			"oneof_group_C": {
				{Name: "field_C1"},
				{Name: "field_C2"},
				{Name: "field_C3"},
				{Name: "field_C4"},
			},
		}

		result := generateOneOfCombinations(oneofGroups)

		if len(result) != 24 {
			t.Fatalf("Expected 4 combinations, got %d", len(result))
		}

	})
}
