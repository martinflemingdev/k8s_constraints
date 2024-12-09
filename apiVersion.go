package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateApiVersion validates the syntax of an apiVersion string.
func ValidateApiVersion(apiVersion string) error {
	errs := make([]error, 0)

	// Check if the string is empty
	if apiVersion == "" {
		errs = append(errs, errors.New("apiVersion cannot be empty"))
	}

	// Length check: entire apiVersion should not exceed 63 characters
	if err := ValidateLength(apiVersion, 63); err != nil {
		errs = append(errs, err)
	}

	// Check allowed characters
	if err := ValidateApiVersionAllowedCharacters(apiVersion); err != nil {
		errs = append(errs, err)
	}

	// Check the group/version format
	if err := ValidateGroupVersionFormat(apiVersion); err != nil {
		errs = append(errs, err)
	}

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateLength checks if a string exceeds the maximum allowed length.
func ValidateLength(input string, maxLength int) error {
	if len(input) > maxLength {
		return fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}
	return nil
}

// ValidateApiVersionAllowedCharacters ensures the string only contains valid characters
// for an apiVersion and contains at most one slash (/).
func ValidateApiVersionAllowedCharacters(input string) error {
	// Valid characters: alphanumeric, hyphen (-), and slash (/)
	validChars := regexp.MustCompile(`^[a-zA-Z0-9/-]+$`)
	if !validChars.MatchString(input) {
		return errors.New("input contains invalid characters; only alphanumeric, hyphen (-), and slash (/) are allowed")
	}

	// Ensure at most one slash (/)
	if strings.Count(input, "/") > 1 {
		return errors.New("input contains more than one slash (/); only one slash is allowed")
	}

	return nil
}

// ValidateGroupVersionFormat validates the group/version format.
func ValidateGroupVersionFormat(apiVersion string) error {
	// Split into group and version (e.g., apps/v1)
	parts := strings.Split(apiVersion, "/")
	if len(parts) == 1 {
		// Core API group (e.g., v1)
		if !isValidVersion(parts[0]) {
			return errors.New("core API version is invalid; must match pattern `v\\d+` or `v\\d+(alpha|beta)\\d+`")
		}
	} else if len(parts) == 2 {
		// Non-core API group (e.g., apps/v1)
		group, version := parts[0], parts[1]

		// Validate group using DNS label conventions
		if err := ValidateDNSLabel(group); err != nil {
			return fmt.Errorf("API group is invalid: %v", err)
		}

		// Validate version using regex
		if !isValidVersion(version) {
			return errors.New("API version is invalid; must match pattern `v\\d+` or `v\\d+(alpha|beta)\\d+`")
		}
	} else {
		// Too many slashes in the apiVersion
		return errors.New("apiVersion has an invalid format; expected `group/version` or `version`")
	}
	return nil
}

// isValidVersion checks if the version matches valid Kubernetes version patterns.
func isValidVersion(version string) bool {
	versionPattern := regexp.MustCompile(`^v\d+((alpha|beta)\d+)?$`)
	return versionPattern.MatchString(version)
}

// ValidateDNSLabel validates a string against the DNS label format as defined by RFC 1123.
// Kubernetes uses this for group names, metadata.labels keys/values, and similar fields.
func ValidateDNSLabel(label string) error {
	// DNS label format: Alphanumeric, hyphens, periods allowed, must start/end with alphanumeric.
	// Maximum length of 63 characters.
	labelPattern := regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`)
	if len(label) > 63 {
		return fmt.Errorf("label exceeds maximum length of 63 characters")
	}
	if !labelPattern.MatchString(label) {
		return errors.New("label must match DNS label format (alphanumeric, hyphens, periods, max 63 characters)")
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

func main() {
	// Test cases
	testCases := []string{
		"v1",
		"apps/v1",
		"apps/v1beta1",
		"",
		"Apps/v1",             // Invalid due to case sensitivity
		"apps/v1.1",           // Invalid due to period
		"apps//v1",            // Invalid due to double slashes
		"this-is-a-very-long-api-group-name-that-exceeds-the-limit/v1",
	}

	for _, tc := range testCases {
		fmt.Printf("Testing apiVersion: %s\n", tc)
		if err := ValidateApiVersion(tc); err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Println("Valid!")
		}
	}
}
