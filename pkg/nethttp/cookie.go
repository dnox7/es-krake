package nethttp

import "github.com/gin-gonic/gin"

func SetCookie(c *gin.Context, cookies []map[string]interface{}, inputMaxAge *int) {
	for _, cookie := range cookies {

		name := cookie["name"].(string)
		value := cookie["value"].(string)
		path := cookie["path"].(string)
		domain := cookie["domain"].(string)
		secure := cookie["secure"].(bool)
		httpOnly := cookie["http_only"].(bool)

		maxAge := *inputMaxAge
		if v, ok := cookie["max_age"].(int); ok && v != 0 {
			maxAge = v
		}

		c.SetCookie(name, value, maxAge, path, domain, secure, httpOnly)
	}
}
