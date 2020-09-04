package util

// CopyMap deep copies the given map and returns a new map
func CopyMap(m map[string]string) map[string]string {
	cm := make(map[string]string)
	for k, v := range m {
		cm[k] = v
	}

	return cm
}
