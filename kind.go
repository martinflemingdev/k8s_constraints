package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateKind validates the syntax of the kind field in a Kubernetes manifest.
func ValidateKind(kind string) error {
	errs := make([]error, 0)

	// Check if the string is empty
	if kind == "" {
		errs = append(errs, errors.New("kind cannot be empty"))
	}

	// Length check: kind should not exceed 63 characters
	if err := ValidateLength(kind, 63); err != nil {
		errs = append(errs, err)
	}

	// Check allowed characters (alphanumeric only)
	if err := ValidateAlphanumeric(kind); err != nil {
		errs = append(errs, err)
	}

	// Check if it starts with an uppercase letter
	if err := ValidateStartsWithUppercase(kind); err != nil {
		errs = append(errs, err)
	}

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateAlphanumeric ensures the string contains only alphanumeric characters.
func ValidateAlphanumeric(input string) error {
	alphanumericPattern := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !alphanumericPattern.MatchString(input) {
		return errors.New("input contains invalid characters; only alphanumeric characters are allowed")
	}
	return nil
}

// ValidateStartsWithUppercase ensures the string starts with an uppercase letter.
func ValidateStartsWithUppercase(input string) error {
	if len(input) == 0 {
		return errors.New("input cannot be empty")
	}
	if !strings.HasPrefix(input, strings.ToUpper(string(input[0]))) || !regexp.MustCompile(`^[A-Z]`).MatchString(input) {
		return errors.New("input must start with an uppercase letter")
	}
	return nil
}

// JoinErrors joins multiple error messages into one error.
func JoinErrors(errs []error) error {
	messages := make([]string, len(errs))
	for i, err := range errs {
		messages[i] = err.Error()
	}
	return errors.New(strings.Join(messages, "; "))
}

// ValidateLength checks if a string exceeds the maximum allowed length.
func ValidateLength(input string, maxLength int) error {
	if len(input) > maxLength {
		return fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}
	return nil
}

func main() {
	// Test cases for ValidateKind
	testCases := []string{
		"Pod",              // Valid
		"Service",          // Valid
		"deployment",       // Invalid: does not start with uppercase
		"123Pod",           // Invalid: starts with a number
		"MyCustomResource", // Valid
		"",                 // Invalid: empty
		"thisisaverylongkindnamethatexceedsthemaxlengthallowed", // Invalid: too long
		"Pod-Service",      // Invalid: contains non-alphanumeric characters
	}

	for _, tc := range testCases {
		fmt.Printf("Testing kind: %s\n", tc)
		if err := ValidateKind(tc); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Valid!")
		}
	}
}
