package middleware

import (
	"errors"
	"net/http"
	"strings"
)

type APIPattern struct {
	Path   string
	Method string
}

func (p APIPattern) Key() string {
	return p.Method + ":" + strings.TrimSuffix(p.Path, "/")
}

func parsePatternKey(key string) (*APIPattern, error) {
	sepIdx := strings.Index(key, ":")
	if sepIdx == -1 || sepIdx >= 7 {
		return nil, errors.New("invalid pattern key")
	}
	return &APIPattern{
		Path:   key[sepIdx+1:],
		Method: key[:sepIdx],
	}, nil
}

var routeMatcherIndex map[string]string

//nolint:gochecknoinits
func init() {
	routeMatcherIndex = make(map[string]string)
	for code, patterns := range APIPermissions {
		for _, p := range patterns {
			routeMatcherIndex[p.Key()] = code
		}
	}
}

var APIPermissions = map[string][]APIPattern{
	"BU001": {
		{
			Path:   "/hello",
			Method: http.MethodGet,
		},
	},
}
