package validator

import (
	"fmt"
	"regexp"
)

// Rule is a function that validates the value of an environment variable.
// It returns an error if the value does not satisfy the rule, or nil otherwise.
type Rule func(name, value string) error

// RequireNonEmpty returns a Rule that fails when the value is an empty string.
func RequireNonEmpty() Rule {
	return func(name, value string) error {
		if value == "" {
			return fmt.Errorf("%s: value must not be empty", name)
		}
		return nil
	}
}

// RequireMatch returns a Rule that fails when the value does not match the
// provided regular expression pattern.
func RequireMatch(pattern string) Rule {
	re := regexp.MustCompile(pattern)
	return func(name, value string) error {
		if !re.MatchString(value) {
			return fmt.Errorf("%s: value %q does not match pattern %q", name, value, pattern)
		}
		return nil
	}
}

// RequireOneOf returns a Rule that fails when the value is not in the allowed set.
func RequireOneOf(allowed ...string) Rule {
	set := make(map[string]struct{}, len(allowed))
	for _, a := range allowed {
		set[a] = struct{}{}
	}
	return func(name, value string) error {
		if _, ok := set[value]; !ok {
			return fmt.Errorf("%s: value %q is not one of %v", name, value, allowed)
		}
		return nil
	}
}

// ApplyRules runs each rule against the given name/value pair and returns all
// errors encountered.
func ApplyRules(name, value string, rules []Rule) []error {
	var errs []error
	for _, r := range rules {
		if err := r(name, value); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}
