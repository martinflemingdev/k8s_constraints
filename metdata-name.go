// ValidateMetadataName validates the syntax of the metadata.name field in a Kubernetes manifest.
func ValidateMetadataName(name string) error {
	errs := make([]error, 0)

	// Check if the string is empty or exceeds length constraints
	if err := ValidateLength(name, 253); err != nil {
		errs = append(errs, fmt.Errorf("metadata.name must be between 1 and 253 characters: %v", err))
	}

	// Validate DNS Subdomain format
	// if err := ValidateDNSSubdomain(name); err != nil {
	// 	errs = append(errs, err)
	// }

	// If there are errors, join and return them
	if len(errs) > 0 {
		return JoinErrors(errs)
	}

	return nil
}

// ValidateMetadataName validates the metadata.name field in a Kubernetes manifest.
func ValidateMetadataName(name string) error {
	// Regex for DNS label format (no dots, lowercase only)
	namePattern := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)

	if len(name) > 253 {
		return fmt.Errorf("metadata.name exceeds maximum length of 253 characters")
	}
	if !namePattern.MatchString(name) {
		return errors.New("metadata.name must consist of lowercase alphanumeric characters or '-', must start and end with an alphanumeric character, and must not contain '.'")
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
