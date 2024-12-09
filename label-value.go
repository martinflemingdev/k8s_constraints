// ValidateLabelValue validates the value of a Kubernetes label.
// Label values must conform to the DNS label convention:
// - Max length of 63 characters.
// - Alphanumeric, '-' and '.' allowed.
// - Must start and end with an alphanumeric character.
func ValidateLabelValue(value string) error {
	// DNS Label regex for label values
	labelValuePattern := regexp.MustCompile(`^[A-Za-z0-9]([-A-Za-z0-9.]*[A-Za-z0-9])?$`)

	if len(value) > 63 {
		return fmt.Errorf("label value exceeds maximum length of 63 characters")
	}
	if !labelValuePattern.MatchString(value) {
		return errors.New("label value must consist of alphanumeric characters, '-', '.', and must start and end with an alphanumeric character")
	}
	return nil
}
