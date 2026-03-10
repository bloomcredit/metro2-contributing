package utils

// IsValidSocialSecurityNumber checks if the provided SSN is valid. It should be a 9-digit number
// that is not all zeros or all nines.
func IsValidSocialSecurityNumber(ssn int) bool {
	// Do not count zero- or 9-filled SSNs.
	if ssn <= 0 || ssn >= 999999999 {
		return false
	}
	return true
}
