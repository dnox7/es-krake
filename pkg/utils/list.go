package utils

func IsSubSet[T comparable](subset, full []T) bool {
	set := make(map[T]struct{}, len(full))
	for _, v := range full {
		set[v] = struct{}{}
	}
	for _, v := range subset {
		if _, ok := set[v]; !ok {
			return false
		}
	}
	return true
}
