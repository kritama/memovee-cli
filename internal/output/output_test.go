package output

import (
	"bytes"
	"testing"
)

type testResult struct {
	Value string `json:"value"`
}

func TestJSONSuccessIsDeterministicAndDoesNotEscapeHTML(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	renderer := NewRenderer(&stdout, &stderr, true)

	err := renderer.Success(testResult{Value: "<safe>"}, "ignored")
	if err != nil {
		t.Fatalf("Success() error = %v", err)
	}

	want := "{\"schema_version\":\"1\",\"ok\":true," +
		"\"result\":{\"value\":\"<safe>\"}}\n"
	if got := stdout.String(); got != want {
		t.Fatalf("stdout = %q, want %q", got, want)
	}
	if got := stderr.String(); got != "" {
		t.Fatalf("stderr = %q, want empty", got)
	}
}

func TestHumanFailureUsesStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	renderer := NewRenderer(&stdout, &stderr, false)

	err := renderer.Failure(Failure{
		Category: "network",
		ExitCode: 16,
		Message:  "endpoint unavailable",
		Next:     "check the configured origin",
	})
	if err != nil {
		t.Fatalf("Failure() error = %v", err)
	}

	if got := stdout.String(); got != "" {
		t.Fatalf("stdout = %q, want empty", got)
	}
	want := "error (network): endpoint unavailable\n" +
		"next: check the configured origin\n"
	if got := stderr.String(); got != want {
		t.Fatalf("stderr = %q, want %q", got, want)
	}
}
