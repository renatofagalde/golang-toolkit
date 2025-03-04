package uri

import (
	"regexp"
	"strings"
)

type URI struct {
}

func (u *URI) ParseURI(method, uri string) map[string]string {
	uuidRegex := regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

	parts := strings.Split(uri, "/")

	for i, part := range parts {
		if uuidRegex.MatchString(part) {
			if i > 0 {
				resourceType := strings.TrimSuffix(parts[i-1], "s")
				parts[i] = "id_" + resourceType
			}
		}
	}

	formattedURI := strings.Join(parts, "/")

	result := map[string]string{
		"action": method,
		"uri":    formattedURI,
	}

	return result
}
