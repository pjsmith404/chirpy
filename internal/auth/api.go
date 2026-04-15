package auth

import (
	"fmt"
	"strings"
	"net/http"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")

	if apiKey == "" {
		return "", fmt.Errorf("No API key found in Authorization header")
	}

	apiKey = strings.TrimPrefix(apiKey, "ApiKey ")

	return apiKey, nil
}
