package cors

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

type converter func(string) string

func generateNormalHeaders(c ConfigCors) http.Header {
	header := make(http.Header)
	if c.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.ExposeHeaders) > 0 {
		exposeHeaders := convert(normalize(c.ExposeHeaders), http.CanonicalHeaderKey)
		header.Set("Access-Control-Expose-Headers", strings.Join(exposeHeaders, ","))
	}
	if c.AllowAllOrigins {
		header.Set("Access-Control-Allow-Origin", "*")
	} else {
		header.Set("Vary", "Origin")
	}
	return header
}

func generatePreflightHeaders(c ConfigCors) http.Header {
	header := make(http.Header)
	if c.AllowCredentials {
		header.Set("Access-Control-Allow-Credentials", "true")
	}
	if len(c.AllowMethods) > 0 {
		allowMethods := convert(normalize(c.AllowMethods), strings.ToUpper)
		value := strings.Join(allowMethods, ",")
		header.Set("Access-Control-Allow-Methods", value)
	}
	if len(c.AllowHeaders) > 0 {
		allowHeaders := convert(normalize(c.AllowHeaders), http.CanonicalHeaderKey)
		value := strings.Join(allowHeaders, ",")
		header.Set("Access-Control-Allow-Headers", value)
	}
	if c.MaxAge > time.Duration(0) {
		value := strconv.FormatInt(int64(c.MaxAge/time.Second), 10)
		header.Set("Access-Control-Max-Age", value)
	}
	if c.AllowAllOrigins {
		header.Set("Access-Control-Allow-Origin", "*")
	} else {
		header.Add("Vary", "Origin")
		header.Add("Vary", "Access-Control-Request-Method")
		header.Add("Vary", "Access-Control-Request-Headers")
	}

	return header
}

func normalize(values []string) []string {
	if values == nil {
		return nil
	}

	distinctMap := make(map[string]bool, len(values))
	normalized := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		value = strings.ToLower(value)
		if _, seen := distinctMap[value]; !seen {
			normalized = append(normalized, value)
			distinctMap[value] = true
		}
	}
	return normalized
}

func convert(s []string, c converter) []string {
	var out []string
	for _, i := range s {
		out = append(out, c(i))
	}
	return out
}
