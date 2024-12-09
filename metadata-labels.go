package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

// ValidateMetadataLabels validates the syntax of metadata.labels in a Kubernetes manifest.
func ValidateMetadataLabels(labels map[string]string) error {
	errs := make([]error, 0)

	for key, value := range labels {
		// Validate the label key
		if err := ValidateLabelKey(key); err != nil {
			errs = append(errs, fmt.Errorf("invalid label key '%s': %v", key, err))
		}

		// Validate the label value
		if err := ValidateLabelValue(value); err != nil {
			errs = append(errs, fmt.Errorf("invalid label value for key '%s': %v", key, err))
		}
	}

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateLabelKey validates a label key, which may have an optional prefix.
func ValidateLabelKey(key string) error {
	// A label key may optionally have a prefix (DNS subdomain followed by `/`)
	parts := strings.SplitN(key, "/", 2)

	if len(parts) == 2 {
		// Validate the prefix (must be a valid DNS subdomain)
		prefix := parts[0]
		if err := ValidateDNSSubdomain(prefix); err != nil {
			return fmt.Errorf("invalid prefix: %v", err)
		}

		// Validate the name part (must be a valid DNS label)
		name := parts[1]
		if err := ValidateDNSLabel(name); err != nil {
			return fmt.Errorf("invalid name: %v", err)
		}
	} else if len(parts) == 1 {
		// Validate the key as a DNS label if no prefix is present
		if err := ValidateDNSLabel(parts[0]); err != nil {
			return fmt.Errorf("invalid name: %v", err)
		}
	} else {
		return errors.New("label key must not be empty")
	}

	return nil
}

// ValidateLabelValue validates a label value (must be a valid DNS label or empty).
func ValidateLabelValue(value string) error {
	if value == "" {
		// Empty values are allowed
		return nil
	}

	// Validate value as a DNS label
	if err := ValidateDNSLabel(value); err != nil {
		return fmt.Errorf("invalid value: %v", err)
	}

	return nil
}

// ValidateDNSLabel validates a string against the DNS label format as defined by RFC 1123.
func ValidateDNSLabel(label string) error {
	// DNS label format: Alphanumeric, hyphens allowed, must start/end with alphanumeric.
	// Maximum length of 63 characters.
	labelPattern := regexp.MustCompile(`^[a-zA-Z0-9]([-a-zA-Z0-9]*[a-zA-Z0-9])?$`)
	if len(label) > 63 {
		return fmt.Errorf("label exceeds maximum length of 63 characters")
	}
	if !labelPattern.MatchString(label) {
		return errors.New("label must match DNS label format (alphanumeric, hyphens, max 63 characters, must start and end with alphanumeric)")
	}
	return nil
}

// ValidateDNSSubdomain validates a string against the DNS subdomain format as defined by RFC 1123.
func ValidateDNSSubdomain(subdomain string) error {
	// DNS subdomain format: Lowercase alphanumeric, `-`, `.` allowed.
	// Must start/end with alphanumeric, max 253 characters.
	subdomainPattern := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	if len(subdomain) > 253 {
		return fmt.Errorf("subdomain exceeds maximum length of 253 characters")
	}
	if !subdomainPattern.MatchString(subdomain) {
		return errors.New("subdomain must match DNS subdomain format (lowercase alphanumeric, `-`, `.`, max 253 characters, must start and end with alphanumeric)")
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
	// Test cases for ValidateMetadataLabels
	testLabels := map[string]string{
		"app.kubernetes.io/name": "my-app",            // Valid
		"app.kubernetes.io/role": "frontend",         // Valid
		"app.kubernetes.io/role/extra": "invalid",    // Invalid: key contains extra `/`
		"App.kubernetes.io/Name": "My-App",           // Invalid: uppercase in key and value
		"app.kubernetes.io/name": "",                 // Valid: empty value
		"invalid_key":            "value",            // Invalid: `_` in key
		"example.com/123":        "valid-value",      // Valid: numeric in key and value
	}

	if err := ValidateMetadataLabels(testLabels); err != nil {
		fmt.Printf("Errors: %v\n", err)
	} else {
		fmt.Println("All labels are valid!")
	}
}
