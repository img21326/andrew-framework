package helper

import (
	"net/url"
)

func UrlValuesToStringMap(values url.Values, sep string) map[string]string {
	result := make(map[string]string)

	for key, val := range values {
		var v Collection[string] = val
		result[key] = v.Join(sep, func(s string) string {
			return s
		})
	}

	return result
}
