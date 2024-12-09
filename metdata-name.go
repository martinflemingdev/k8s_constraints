// ValidateMetadataName validates the syntax of the metadata.name field in a Kubernetes manifest.
func ValidateMetadataName(name string) error {
	errs := make([]error, 0)

	// Check if the string is empty or exceeds length constraints
	if err := ValidateLength(name, 253); err != nil {
		errs = append(errs, fmt.Errorf("metadata.name must be between 1 and 253 characters: %v", err))
	}

	// Validate DNS Subdomain format
	if err := ValidateDNSSubdomain(name); err != nil {
		errs = append(errs, err)
	}

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateDNSSubdomain validates a string against the DNS subdomain format as defined by RFC 1123.
// Kubernetes uses this for metadata.name and similar fields.
func ValidateDNSSubdomain(subdomain string) error {
	// DNS subdomain format: Lowercase alphanumeric, `-`, `.` allowed.
	// Must start/end with alphanumeric, max 253 characters.
	subdomainPattern := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	if !subdomainPattern.MatchString(subdomain) {
		return errors.New("subdomain must match DNS subdomain format (lowercase alphanumeric, `-`, `.`, max 253 characters, must start and end with alphanumeric)")
	}
	return nil
}

// ValidateLength checks if a string exceeds the maximum allowed length.
func ValidateLength(input string, maxLength int) error {
	if len(input) == 0 {
		return errors.New("input cannot be empty")
	}
	if len(input) > maxLength {
		return fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}
	return nil
}
