package utils

func ToSet[T comparable](items []T) []T {
	set := make(map[T]struct{})
	for _, v := range items {
		set[v] = struct{}{}
	}
	res := []T{}
	for ele := range set {
		res = append(res, ele)
	}
	return res
}

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
