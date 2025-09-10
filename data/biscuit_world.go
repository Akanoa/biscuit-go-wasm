package data

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

type ScopedFact struct {
	Facts  []string  `json:"facts"`
	Origin []*uint64 `json:"origin"`
}

type ScopedRule struct {
	Rules  []string `json:"rules"`
	Origin *uint64  `json:"origin"`
}

type ScopedCheck struct {
	Checks []string `json:"checks"`
	Origin *uint64  `json:"origin"`
}

type World struct {
	Facts    []ScopedFact  `json:"facts"`
	Rules    []ScopedRule  `json:"rules"`
	Checks   []ScopedCheck `json:"checks"`
	Policies []string      `json:"policies"`
}

// OriginAuthorizer Rust u64::MAX is 18446744073709551615
const OriginAuthorizer = uint64(18446744073709551615)

func (world *World) Equals(other *World) bool {
	if world == nil && other == nil {
		return true
	}
	if world == nil || other == nil {
		return false
	}

	// Compare Facts
	if len(world.Facts) != len(other.Facts) {
		return false
	}
	for i, fact := range world.Facts {
		if !fact.Equals(&other.Facts[i]) {
			return false
		}
	}

	// Compare Rules
	if len(world.Rules) != len(other.Rules) {
		return false
	}
	for i, rule := range world.Rules {
		if !rule.Equals(&other.Rules[i]) {
			return false
		}
	}

	// Compare Checks
	if len(world.Checks) != len(other.Checks) {
		return false
	}
	for i, check := range world.Checks {
		if !check.Equals(&other.Checks[i]) {
			return false
		}
	}

	// Compare Policies
	if len(world.Policies) != len(other.Policies) {
		return false
	}
	for i, policy := range world.Policies {
		if policy != other.Policies[i] {
			return false
		}
	}

	return true
}

func (world World) FromString(str string) (World, error) {
	err := json.Unmarshal([]byte(str), &world)
	if err != nil {
		return World{}, err
	}
	return world, nil
}

func (fact *ScopedFact) Equals(other *ScopedFact) bool {
	if fact == nil && other == nil {
		return true
	}
	if fact == nil || other == nil {
		return false
	}

	for f := range fact.Facts {
		if fact.Facts[f] != other.Facts[f] {
			return false
		}
	}

	if len(fact.Origin) != len(other.Origin) {
		return false
	}

	for i, scope := range fact.Origin {
		if scope == nil && other.Origin[i] == nil {
			continue
		}
		if scope == nil || other.Origin[i] == nil {
			return false
		}
		if *scope != *other.Origin[i] {
			return false
		}
	}

	return true
}

func (rule *ScopedRule) Equals(other *ScopedRule) bool {
	if rule == nil && other == nil {
		return true
	}
	if rule == nil || other == nil {
		return false
	}

	for r := range rule.Rules {
		if rule.Rules[r] != other.Rules[r] {
			return false
		}
	}

	if rule.Origin == nil && other.Origin == nil {
		return true
	}
	if rule.Origin == nil || other.Origin == nil {
		return false
	}

	return *rule.Origin == *other.Origin
}

func (check *ScopedCheck) Equals(other *ScopedCheck) bool {
	if check == nil && other == nil {
		return true
	}
	if check == nil || other == nil {
		return false
	}
	if len(check.Checks) != len(other.Checks) {
		return false
	}
	for i, check := range check.Checks {
		if check != other.Checks[i] {
			return false
		}
	}
	if check.Origin == nil && other.Origin == nil {
		return true
	}
	if check.Origin == nil || other.Origin == nil {
		return false
	}
	return *check.Origin == *other.Origin
}

func (world *World) String() string {
	var result strings.Builder

	// Facts section
	if len(world.Facts) > 0 {
		result.WriteString("// Facts:\n")

		// Group facts by scope for origin comments
		factsByScope := make(map[string][]string)
		for _, fact := range world.Facts {
			var scopeKey string

			var scopes []string
			for _, scope := range fact.Origin {

				if scope != nil && *scope == OriginAuthorizer {
					scopes = append(scopes, "authorizer")
					continue
				}

				if scope != nil {
					scopes = append(scopes, fmt.Sprintf("%d", *scope))
				} else {
					scopes = append(scopes, "authorizer")
				}
			}
			sort.Strings(scopes)
			scopeKey = strings.Join(scopes, ", ")

			factsByScope[scopeKey] = append(factsByScope[scopeKey], fact.Facts...)
		}

		// Sort scope keys for consistent output, with numeric origins first
		var sortedScopeKeys []string
		for scope := range factsByScope {
			sortedScopeKeys = append(sortedScopeKeys, scope)
		}
		sort.Slice(sortedScopeKeys, func(i, j int) bool {
			if sortedScopeKeys[i] == "" {
				return false
			}
			if sortedScopeKeys[j] == "" {
				return true
			}
			return sortedScopeKeys[i] < sortedScopeKeys[j]
		})

		for _, scopeKey := range sortedScopeKeys {
			facts := factsByScope[scopeKey]
			if len(facts) > 0 {

				result.WriteString(fmt.Sprintf("// origin: %s\n", scopeKey))

				sort.Strings(facts)
				for _, fact := range facts {
					result.WriteString(fmt.Sprintf("%s;\n", fact))
				}
			}
		}
		result.WriteString("\n")
	}

	// Rules section
	if len(world.Rules) > 0 {
		result.WriteString("// Rules:\n")

		// Group rules by scope
		rulesByScope := make(map[string][]string)
		for _, rule := range world.Rules {
			var scopeKey string
			if rule.Origin == nil {
				scopeKey = ""
			} else {
				scopeKey = fmt.Sprintf("%d", *rule.Origin)
			}
			rulesByScope[scopeKey] = append(rulesByScope[scopeKey], rule.Rules...)
		}

		// Sort scope keys, with "authorizer" first
		var sortedScopeKeys []string
		for scope := range rulesByScope {
			sortedScopeKeys = append(sortedScopeKeys, scope)
		}
		sort.Slice(sortedScopeKeys, func(i, j int) bool {
			if sortedScopeKeys[i] == "" {
				return true
			}
			if sortedScopeKeys[j] == "" {
				return false
			}
			return sortedScopeKeys[i] < sortedScopeKeys[j]
		})

		for _, scopeKey := range sortedScopeKeys {
			rules := rulesByScope[scopeKey]
			if len(rules) > 0 {
				if scopeKey == "" {
					result.WriteString("// origin: \n")
				} else {
					result.WriteString(fmt.Sprintf("// origin: %s\n", scopeKey))
				}

				sort.Strings(rules)
				for _, rule := range rules {
					result.WriteString(fmt.Sprintf("%s;\n", rule))
				}
			}
		}
		result.WriteString("\n")
	}

	// Checks section
	if len(world.Checks) > 0 {
		result.WriteString("// Checks:\n")

		// Group rules by scope
		checksByScope := make(map[string][]string)
		for _, rule := range world.Checks {
			var scopeKey string
			if rule.Origin != nil && *rule.Origin == OriginAuthorizer {
				scopeKey = "authorizer"
			} else {
				scopeKey = fmt.Sprintf("%d", *rule.Origin)
			}

			checksByScope[scopeKey] = append(checksByScope[scopeKey], rule.Checks...)
		}

		// Sort scope keys, with "authorizer" first
		var sortedScopeKeys []string
		for scope := range checksByScope {
			sortedScopeKeys = append(sortedScopeKeys, scope)
		}
		sort.Slice(sortedScopeKeys, func(i, j int) bool {
			if sortedScopeKeys[i] == "" {
				return true
			}
			if sortedScopeKeys[j] == "" {
				return false
			}
			return sortedScopeKeys[i] < sortedScopeKeys[j]
		})

		for _, scopeKey := range sortedScopeKeys {
			checks := checksByScope[scopeKey]
			if len(checks) > 0 {

				result.WriteString(fmt.Sprintf("// origin: %s\n", scopeKey))

				for _, rule := range checks {
					result.WriteString(fmt.Sprintf("%s;\n", rule))
				}
			}
		}
		result.WriteString("\n")
	}

	// Policies section
	if len(world.Policies) > 0 {
		result.WriteString("// Policies:\n")

		for _, policy := range world.Policies {
			result.WriteString(fmt.Sprintf("%s;\n", policy))
		}
	}

	return result.String()
}
