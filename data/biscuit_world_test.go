package data

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWorld_Equals(t *testing.T) {
	// Helper function to create uint64 pointer
	uint64Ptr := func(val uint64) *uint64 {
		return &val
	}

	tests := []struct {
		name     string
		world1   *World
		world2   *World
		expected bool
	}{
		{
			name:     "both nil",
			world1:   nil,
			world2:   nil,
			expected: true,
		},
		{
			name:     "first nil",
			world1:   nil,
			world2:   &World{},
			expected: false,
		},
		{
			name:     "second nil",
			world1:   &World{},
			world2:   nil,
			expected: false,
		},
		{
			name:     "empty worlds",
			world1:   &World{},
			world2:   &World{},
			expected: true,
		},
		{
			name: "identical worlds",
			world1: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1), uint64Ptr(2)}},
				},
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(1)},
				},
				Checks: []ScopedCheck{
					{Checks: []string{"check1", "check2"}, Origin: uint64Ptr(0)},
				},
				Policies: []string{"policy1"},
			},
			world2: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1), uint64Ptr(2)}},
				},
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(1)},
				},
				Checks: []ScopedCheck{
					{Checks: []string{"check1", "check2"}, Origin: uint64Ptr(0)},
				},
				Policies: []string{"policy1"},
			},
			expected: true,
		},
		{
			name: "different facts",
			world1: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1)}},
				},
			},
			world2: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact2"}, Origin: []*uint64{uint64Ptr(1)}},
				},
			},
			expected: false,
		},
		{
			name: "different fact scopes",
			world1: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1)}},
				},
			},
			world2: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(2)}},
				},
			},
			expected: false,
		},
		{
			name: "different number of facts",
			world1: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1)}},
				},
			},
			world2: &World{
				Facts: []ScopedFact{
					{Facts: []string{"fact1"}, Origin: []*uint64{uint64Ptr(1)}},
					{Facts: []string{"fact2"}, Origin: []*uint64{uint64Ptr(2)}},
				},
			},
			expected: false,
		},
		{
			name: "different rules",
			world1: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(1)},
				},
			},
			world2: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule2"}, Origin: uint64Ptr(1)},
				},
			},
			expected: false,
		},
		{
			name: "different rule scopes",
			world1: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(1)},
				},
			},
			world2: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(2)},
				},
			},
			expected: false,
		},
		{
			name: "rule scope nil vs non-nil",
			world1: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: nil},
				},
			},
			world2: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: uint64Ptr(1)},
				},
			},
			expected: false,
		},
		{
			name: "both rule scopes nil",
			world1: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: nil},
				},
			},
			world2: &World{
				Rules: []ScopedRule{
					{Rules: []string{"rule1"}, Origin: nil},
				},
			},
			expected: true,
		},
		{
			name: "different checks",
			world1: &World{
				Checks: []ScopedCheck{
					{Checks: []string{"check1", "check2"}, Origin: uint64Ptr(0)},
				},
			},
			world2: &World{
				Checks: []ScopedCheck{
					{Checks: []string{"check1", "check3"}, Origin: uint64Ptr(0)},
				},
			},
			expected: false,
		},
		{
			name: "different number of checks",
			world1: &World{
				Checks: []ScopedCheck{
					{Checks: []string{"check1"}, Origin: uint64Ptr(0)},
				},
			},
			world2: &World{
				Checks: []ScopedCheck{
					{Checks: []string{"check1"}, Origin: uint64Ptr(0)},
					{Checks: []string{"check2"}, Origin: uint64Ptr(0)},
				},
			},
			expected: false,
		},
		{
			name: "different policies",
			world1: &World{
				Policies: []string{"policy1"},
			},
			world2: &World{
				Policies: []string{"policy2"},
			},
			expected: false,
		},
		{
			name: "different number of policies",
			world1: &World{
				Policies: []string{"policy1"},
			},
			world2: &World{
				Policies: []string{"policy1", "policy2"},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.world1.Equals(tt.world2)
			if result != tt.expected {
				t.Errorf("World.Equals() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestScopedFact_Equals(t *testing.T) {
	uint64Ptr := func(val uint64) *uint64 {
		return &val
	}

	tests := []struct {
		name     string
		fact1    *ScopedFact
		fact2    *ScopedFact
		expected bool
	}{
		{
			name:     "both nil",
			fact1:    nil,
			fact2:    nil,
			expected: true,
		},
		{
			name:     "first nil",
			fact1:    nil,
			fact2:    &ScopedFact{},
			expected: false,
		},
		{
			name:     "identical facts",
			fact1:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1), uint64Ptr(2)}},
			fact2:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1), uint64Ptr(2)}},
			expected: true,
		},
		{
			name:     "different fact strings",
			fact1:    &ScopedFact{Facts: []string{"test1"}, Origin: []*uint64{uint64Ptr(1)}},
			fact2:    &ScopedFact{Facts: []string{"test2"}, Origin: []*uint64{uint64Ptr(1)}},
			expected: false,
		},
		{
			name:     "different scope lengths",
			fact1:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1)}},
			fact2:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1), uint64Ptr(2)}},
			expected: false,
		},
		{
			name:     "different scope values",
			fact1:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1)}},
			fact2:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(2)}},
			expected: false,
		},
		{
			name:     "nil vs non-nil scope element",
			fact1:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{nil}},
			fact2:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{uint64Ptr(1)}},
			expected: false,
		},
		{
			name:     "both nil scope elements",
			fact1:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{nil}},
			fact2:    &ScopedFact{Facts: []string{"test"}, Origin: []*uint64{nil}},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.fact1.Equals(tt.fact2)
			if result != tt.expected {
				t.Errorf("ScopedFact.Equals() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestScopedRule_Equals(t *testing.T) {
	uint64Ptr := func(val uint64) *uint64 {
		return &val
	}

	tests := []struct {
		name     string
		rule1    *ScopedRule
		rule2    *ScopedRule
		expected bool
	}{
		{
			name:     "both nil",
			rule1:    nil,
			rule2:    nil,
			expected: true,
		},
		{
			name:     "first nil",
			rule1:    nil,
			rule2:    &ScopedRule{},
			expected: false,
		},
		{
			name:     "identical rules",
			rule1:    &ScopedRule{Rules: []string{"test"}, Origin: uint64Ptr(1)},
			rule2:    &ScopedRule{Rules: []string{"test"}, Origin: uint64Ptr(1)},
			expected: true,
		},
		{
			name:     "different rule strings",
			rule1:    &ScopedRule{Rules: []string{"test1"}, Origin: uint64Ptr(1)},
			rule2:    &ScopedRule{Rules: []string{"test2"}, Origin: uint64Ptr(1)},
			expected: false,
		},
		{
			name:     "different scope values",
			rule1:    &ScopedRule{Rules: []string{"test"}, Origin: uint64Ptr(1)},
			rule2:    &ScopedRule{Rules: []string{"test"}, Origin: uint64Ptr(2)},
			expected: false,
		},
		{
			name:     "nil vs non-nil scope",
			rule1:    &ScopedRule{Rules: []string{"test"}, Origin: nil},
			rule2:    &ScopedRule{Rules: []string{"test"}, Origin: uint64Ptr(1)},
			expected: false,
		},
		{
			name:     "both nil scopes",
			rule1:    &ScopedRule{Rules: []string{"test"}, Origin: nil},
			rule2:    &ScopedRule{Rules: []string{"test"}, Origin: nil},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rule1.Equals(tt.rule2)
			if result != tt.expected {
				t.Errorf("ScopedRule.Equals() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestWorld_String(t *testing.T) {
	expected := `// Facts:
// origin: 0
right("file1", "read");
right("file1", "write");
right("file2", "read");
// origin: authorizer
resource("file1");

// Checks:
// origin: 1
check if resource($0), operation("read"), right($0, "read");

// Policies:
allow if true;
`

	uint64Ptr := func(val uint64) *uint64 {
		return &val
	}

	world := World{
		Facts: []ScopedFact{
			{
				Facts:  []string{`right("file1", "read")`},
				Origin: []*uint64{uint64Ptr(0)},
			},
			{
				Facts:  []string{`right("file1", "write")`},
				Origin: []*uint64{uint64Ptr(0)},
			},
			{
				Facts:  []string{`right("file2", "read")`},
				Origin: []*uint64{uint64Ptr(0)},
			},
			{
				Facts:  []string{`resource("file1")`},
				Origin: []*uint64{uint64Ptr(OriginAuthorizer)},
			},
		},
		Checks: []ScopedCheck{
			{
				Checks: []string{
					`check if resource($0), operation("read"), right($0, "read")`,
				},
				Origin: uint64Ptr(1),
			},
		},
		Policies: []string{
			`allow if true`,
		},
	}

	result := world.String()
	if result != expected {

		require.Equal(t, expected, result)
		t.Errorf("World.String() = %v, expected %v", result, expected)
	}
}
