package utils

// TrimLeadingSlash removes the first slash (if any) from the input path.
func TrimLeadingSlash(path string) string {
	if len(path) > 0 && path[0] == '/' {
		return path[1:]
	}
	return path
}

// EnsureTrailingSlash ensures the input path ends with a slash.
func EnsureTrailingSlash(path string) string {
	if len(path) > 0 && path[len(path)-1] != '/' {
		return path + "/"
	}
	return path
}

// EnsureLeadingSlash ensures that the path starts with a slash.
func EnsureLeadingSlash(path string) string {
	if len(path) > 0 && path[0] != '/' {
		path = "/" + path
	}
	return path
}
