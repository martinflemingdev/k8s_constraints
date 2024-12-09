package main

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"
)

// ValidateMetadataAnnotations validates the syntax of metadata.annotations in a Kubernetes manifest.
func ValidateMetadataAnnotations(annotations map[string]string) error {
	errs := make([]error, 0)

	for key, value := range annotations {
		// Validate the annotation key
		if err := ValidateLabelKey(key); err != nil {
			errs = append(errs, fmt.Errorf("invalid annotation key '%s': %v", key, err))
		}

		// Validate the annotation value
		if err := ValidateAnnotationValue(value); err != nil {
			errs = append(errs, fmt.Errorf("invalid annotation value for key '%s': %v", key, err))
		}
	}

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateAnnotationValue validates an annotation value (must be a valid UTF-8 string).
func ValidateAnnotationValue(value string) error {
	// Annotation values can be any valid UTF-8 string
	if !utf8.ValidString(value) {
		return errors.New("annotation value must be a valid UTF-8 string")
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

// ValidateLabelKey validates a label or annotation key, which may have an optional prefix.
func ValidateLabelKey(key string) error {
	// A label or annotation key may optionally have a prefix (DNS subdomain followed by `/`)
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
		return errors.New("label or annotation key must not be empty")
	}

	return nil
}

func main() {
	// Test cases for ValidateMetadataAnnotations
	testAnnotations := map[string]string{
		"example.com/description": "This is a valid annotation value", // Valid
		"example.com/role":        "",                                // Valid: empty value
		"example.com/Name":        "Uppercase key, should fail",      // Invalid: uppercase in key
		"invalid_key":             "value",                          // Invalid: `_` in key
		"example.com/utf8":        string([]byte{0xff, 0xfe}),       // Invalid: non-UTF-8 value
		"example.com/another":     "Another valid value",            // Valid
	}

	if err := ValidateMetadataAnnotations(testAnnotations); err != nil {
		fmt.Printf("Errors: %v\n", err)
	} else {
		fmt.Println("All annotations are valid!")
	}
}
