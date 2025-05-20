package middleware

import (
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

var routeMatcherIndex map[string][]string

func init() {
	routeMatcherIndex = make(map[string][]string)
	for code, patterns := range APIPermissions {
		for _, p := range patterns {
			routeMatcherIndex[p.Key()] = append(routeMatcherIndex[p.Key()], code)
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
