package utils

// GetKeysFromMap is a helper for fetching a slice of strings from a mpa
func GetKeysFromMap(m map[string]string) (keys []string) {
	for k := range m {
		keys = append(keys, k)
	}

	return keys
}
