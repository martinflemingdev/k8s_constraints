// ValidateDNSSubdomain validates a string against the DNS subdomain format as defined by RFC 1123.
// DNS subdomain format: Lowercase alphanumeric, `-`, `.` allowed.
// Must start/end with alphanumeric, max 253 characters.
func ValidateDNSSubdomain(subdomain string) error {
	// Updated regex to match Kubernetes constraints:
	// [a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*
	subdomainPattern := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$`)

	if len(subdomain) > 253 {
		return fmt.Errorf("subdomain exceeds maximum length of 253 characters")
	}
	if !subdomainPattern.MatchString(subdomain) {
		return errors.New("subdomain must match DNS subdomain format (lowercase alphanumeric, `-`, `.`, max 253 characters, must start and end with alphanumeric)")
	}
	return nil
}

// ValidateLabelOrAnnotationKey validates a label or annotation key based on Kubernetes constraints.
// Keys can have an optional prefix (DNS subdomain) followed by a `/` and a name part.
// The name part must match the regex: ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]
func ValidateLabelOrAnnotationKey(key string) error {
	// Split the key into prefix and name part
	parts := strings.SplitN(key, "/", 2)

	if len(parts) == 2 {
		// Validate the prefix (must be a valid DNS subdomain)
		prefix := parts[0]
		if err := ValidateDNSSubdomain(prefix); err != nil {
			return fmt.Errorf("invalid prefix: %v", err)
		}

		// Validate the name part
		name := parts[1]
		if err := ValidateLabelOrAnnotationNamePart(name); err != nil {
			return fmt.Errorf("invalid name part: %v", err)
		}
	} else if len(parts) == 1 {
		// Validate the key as a name part if no prefix is present
		if err := ValidateLabelOrAnnotationNamePart(parts[0]); err != nil {
			return fmt.Errorf("invalid name part: %v", err)
		}
	} else {
		return errors.New("label or annotation key must not be empty")
	}

	return nil
}

// ValidateLabelOrAnnotationNamePart validates the name part of a label or annotation key.
// It must match the regex: ([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]
func ValidateLabelOrAnnotationNamePart(name string) error {
	// Regex for the name part of a label or annotation key
	namePattern := regexp.MustCompile(`^([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9]$`)

	if len(name) > 63 {
		return fmt.Errorf("name part exceeds maximum length of 63 characters")
	}
	if !namePattern.MatchString(name) {
		return errors.New("name part must consist of alphanumeric characters, '-', '_', or '.', and must start and end with an alphanumeric character")
	}
	return nil
}