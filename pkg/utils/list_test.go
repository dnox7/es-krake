package utils

import "testing"

func TestIsSubSet(t *testing.T) {
	type testCase[T comparable] struct {
		name     string
		subset   []T
		full     []T
		expected bool
	}

	stringTests := []testCase[string]{
		{
			name:     "valid subset",
			subset:   []string{"read", "write"},
			full:     []string{"read", "write", "delete"},
			expected: true,
		},
		{
			name:     "invalid subset with missing element",
			subset:   []string{"read", "unknown"},
			full:     []string{"read", "write"},
			expected: false,
		},
		{
			name:     "empty subset is always valid",
			subset:   []string{},
			full:     []string{"a", "b"},
			expected: true,
		},
		{
			name:     "identical sets are valid subset",
			subset:   []string{"a", "b"},
			full:     []string{"a", "b"},
			expected: true,
		},
		{
			name:     "subset larger than full set is invalid",
			subset:   []string{"a", "b", "c"},
			full:     []string{"a", "b"},
			expected: false,
		},
	}

	for _, tc := range stringTests {
		t.Run(tc.name, func(t *testing.T) {
			result := IsSubSet(tc.subset, tc.full)
			if result != tc.expected {
				t.Errorf("IsSubSet(%v, %v) = %v; expected %v", tc.subset, tc.full, result, tc.expected)
			}
		})
	}
}
