func ValidateAPIGroup(group string) error {
	// Custom regex for API group names
	groupPattern := regexp.MustCompile(`^[A-Za-z0-9]([-.A-Za-z0-9]*[A-Za-z0-9])?$`)

	if len(group) > 253 {
		return fmt.Errorf("API group exceeds maximum length of 253 characters")
	}
	if !groupPattern.MatchString(group) {
		return errors.New("API group must consist of alphanumeric characters, '-', '.', and must start and end with an alphanumeric character")
	}
	return nil
}