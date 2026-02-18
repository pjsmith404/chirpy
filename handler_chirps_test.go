package main

import (
	"testing"
)

func TestCleanBody(t *testing.T) {
	body := "This is a chirp"
	cleanedBody := cleanBody(body)
	if body != cleanedBody {
		t.Errorf(`cleanBody(body) = %q, want %q`, cleanedBody, body)
	}

	body = "This is fornax chirp"
	want := "This is **** chirp"
	cleanedBody = cleanBody(body)
	if cleanedBody != want {
		t.Errorf(`cleanedBody(body) = %q, want %q`, cleanedBody, want)
	}
}
