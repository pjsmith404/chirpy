package auth

import (
	"testing"
	"net/http"
	"strings"
)

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}
	api := strings.Join([]string{"ApiKey", "123456789"}, " ")
	headers.Add("Authorization", api)

	apiKey, err := GetAPIKey(headers)
	if err != nil {
		t.Fatal(err)
	}

	if apiKey != "123456789" {
		t.Errorf(`GetApiKey(headers) = %q, want %q`, apiKey, "123456789")
	}
}

