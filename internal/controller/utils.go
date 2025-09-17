package controller

// containsString checks if a slice contains a given string.
func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

// removeString removes a given string from a slice.
func removeString(slice []string, s string) []string {
	for i, item := range slice {
		if item == s {
			slice = append(slice[:i], slice[i+1:]...)
			return slice
		}
	}
	return slice
}
