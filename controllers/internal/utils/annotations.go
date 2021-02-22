package utils

// MergeAnnotations adds or updates annotations from source to target and returns merge
func MergeAnnotations(target map[string]string, source map[string]string) map[string]string {
	if target == nil {
		target = make(map[string]string)
	}
	if source == nil {
		source = make(map[string]string)
	}
	for k, v := range source {
		if target[k] != v {
			target[k] = v
		}
	}
	return target
}
